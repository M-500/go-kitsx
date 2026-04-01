package eventbus

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	ErrBusClosed          = errors.New("eventbus: bus is closed")
	ErrTransportNotFound  = errors.New("eventbus: transport not found")
	ErrTransportDuplicate = errors.New("eventbus: transport already registered")
)

type Bus struct {
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	transports  map[string]Driver
	middleware  []Middleware
	metrics     metricsRegistry
	closed      atomic.Bool
	defaultMode string
}

func NewBus() *Bus {
	ctx, cancel := context.WithCancel(context.Background())
	b := &Bus{
		ctx:         ctx,
		cancel:      cancel,
		transports:  make(map[string]Driver),
		defaultMode: TransportLocal,
	}
	_ = b.RegisterTransport(TransportLocal, NewLocalTransport())
	_ = b.RegisterTransport(TransportChannel, NewChannelTransport(ChannelConfig{}))
	return b
}

func (b *Bus) Use(mw ...Middleware) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.middleware = append(b.middleware, mw...)
}

func (b *Bus) RegisterTransport(name string, driver Driver) error {
	if name == "" {
		return fmt.Errorf("eventbus: empty transport name")
	}
	if driver == nil {
		return fmt.Errorf("eventbus: nil transport")
	}

	b.mu.Lock()
	defer b.mu.Unlock()
	if _, ok := b.transports[name]; ok {
		return fmt.Errorf("%w: %s", ErrTransportDuplicate, name)
	}
	b.transports[name] = driver
	return nil
}

func (b *Bus) Publish(ctx context.Context, msg *Message, options PublishOptions) error {
	if msg == nil {
		return fmt.Errorf("eventbus: nil message")
	}
	if b.closed.Load() {
		return ErrBusClosed
	}

	options = options.normalized()
	driver, err := b.transport(options.Transport)
	if err != nil {
		b.metrics.publishFailed.Add(1)
		return err
	}

	b.metrics.published.Add(1)
	if err := driver.Publish(withFallbackContext(ctx, b.ctx), msg.Clone()); err != nil {
		b.metrics.publishFailed.Add(1)
		return err
	}
	return nil
}

func (b *Bus) PublishJSON(ctx context.Context, topic string, payload any, options PublishOptions) error {
	msg, err := NewJSONMessage(topic, payload)
	if err != nil {
		return err
	}
	return b.Publish(ctx, msg, options)
}

func (b *Bus) Subscribe(topic string, handler Handler, options SubscribeOptions) (Subscription, error) {
	if topic == "" {
		return nil, fmt.Errorf("eventbus: empty topic")
	}
	if handler == nil {
		return nil, fmt.Errorf("eventbus: nil handler")
	}
	if b.closed.Load() {
		return nil, ErrBusClosed
	}

	options = options.normalized()
	driver, err := b.transport(options.Transport)
	if err != nil {
		return nil, err
	}

	b.mu.RLock()
	mws := append([]Middleware(nil), b.middleware...)
	b.mu.RUnlock()

	final := handler
	final = withTimeout(final, options.Timeout)
	final = Chain([]Middleware{
		withRecover(&b.metrics),
		withRetry(&b.metrics, options.RetryPolicy),
	}, final)
	final = Chain(mws, final)

	return driver.Subscribe(b.ctx, SubscribeRequest{
		Topic:         topic,
		Options:       options,
		FinalHandler:  wrapMetrics(&b.metrics, final),
		Middlewares:   mws,
		SharedMetrics: &b.metrics,
	})
}

func (b *Bus) Metrics() Metrics {
	return b.metrics.snapshot()
}

func (b *Bus) Close(ctx context.Context) error {
	if !b.closed.CompareAndSwap(false, true) {
		return nil
	}
	b.cancel()

	b.mu.RLock()
	transports := make([]Driver, 0, len(b.transports))
	for _, t := range b.transports {
		transports = append(transports, t)
	}
	b.mu.RUnlock()

	var errs []error
	for _, t := range transports {
		if err := t.Close(withFallbackContext(ctx, context.Background())); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", t.Name(), err))
		}
	}
	return errors.Join(errs...)
}

func (b *Bus) transport(name string) (Driver, error) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	driver, ok := b.transports[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrTransportNotFound, name)
	}
	return driver, nil
}

func wrapMetrics(metrics *metricsRegistry, next Handler) Handler {
	return func(ctx context.Context, msg *Message) error {
		err := next(ctx, msg)
		if err != nil {
			metrics.handlerFailed.Add(1)
			return err
		}
		metrics.delivered.Add(1)
		return nil
	}
}

func withFallbackContext(ctx context.Context, fallback context.Context) context.Context {
	if ctx != nil {
		return ctx
	}
	return fallback
}
