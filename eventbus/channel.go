package eventbus

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type ChannelConfig struct {
	Buffer int
}

type ChannelTransport struct {
	input   chan *Message
	mu      sync.RWMutex
	subs    map[uint64]*localSubscription
	nextID  atomic.Uint64
	closed  atomic.Bool
	runCtx  context.Context
	cancel  context.CancelFunc
	runWg   sync.WaitGroup
	workerW sync.WaitGroup
}

func NewChannelTransport(cfg ChannelConfig) *ChannelTransport {
	if cfg.Buffer <= 0 {
		cfg.Buffer = 256
	}
	ctx, cancel := context.WithCancel(context.Background())
	t := &ChannelTransport{
		input:  make(chan *Message, cfg.Buffer),
		subs:   make(map[uint64]*localSubscription),
		runCtx: ctx,
		cancel: cancel,
	}
	t.runWg.Add(1)
	go t.dispatch()
	return t
}

func (t *ChannelTransport) Name() string { return TransportChannel }

func (t *ChannelTransport) Publish(ctx context.Context, msg *Message) error {
	if t.closed.Load() {
		return ErrBusClosed
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.runCtx.Done():
		return ErrBusClosed
	case t.input <- msg.Clone():
		return nil
	}
}

func (t *ChannelTransport) Subscribe(ctx context.Context, req SubscribeRequest) (Subscription, error) {
	if t.closed.Load() {
		return nil, ErrBusClosed
	}

	subCtx, cancel := context.WithCancel(ctx)
	sub := &localSubscription{
		id:      t.nextID.Add(1),
		topic:   req.Topic,
		options: req.Options,
		handler: req.FinalHandler,
		queue:   make(chan *Message, req.Options.Buffer),
		cancel:  cancel,
	}

	for i := 0; i < req.Options.Concurrency; i++ {
		t.workerW.Add(1)
		go func() {
			defer t.workerW.Done()
			sub.worker(subCtx)
		}()
	}

	t.mu.Lock()
	t.subs[sub.id] = sub
	t.mu.Unlock()

	return &subscriptionFunc{fn: func() error {
		sub.once.Do(cancel)
		t.mu.Lock()
		delete(t.subs, sub.id)
		t.mu.Unlock()
		return nil
	}}, nil
}

func (t *ChannelTransport) Close(ctx context.Context) error {
	if !t.closed.CompareAndSwap(false, true) {
		return nil
	}
	t.cancel()

	t.mu.Lock()
	for _, sub := range t.subs {
		sub.cancel()
	}
	t.subs = make(map[uint64]*localSubscription)
	t.mu.Unlock()

	done := make(chan struct{})
	go func() {
		defer close(done)
		t.runWg.Wait()
		t.workerW.Wait()
	}()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

func (t *ChannelTransport) dispatch() {
	defer t.runWg.Done()
	for {
		select {
		case <-t.runCtx.Done():
			return
		case msg := <-t.input:
			if msg == nil {
				continue
			}
			t.mu.RLock()
			snapshot := make([]*localSubscription, 0, len(t.subs))
			for _, sub := range t.subs {
				if matches(sub.topic, msg.Topic, sub.options.MatchMode) {
					snapshot = append(snapshot, sub)
				}
			}
			t.mu.RUnlock()

			var errs []error
			for _, sub := range snapshot {
				if err := sub.enqueue(t.runCtx, msg.Clone()); err != nil {
					errs = append(errs, err)
				}
			}
			_ = errors.Join(errs...)
		}
	}
}
