package eventbus

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

type LocalTransport struct {
	mu      sync.RWMutex
	subs    map[uint64]*localSubscription
	nextID  atomic.Uint64
	closed  atomic.Bool
	closeWg sync.WaitGroup
}

func NewLocalTransport() *LocalTransport {
	return &LocalTransport{
		subs: make(map[uint64]*localSubscription),
	}
}

func (t *LocalTransport) Name() string { return TransportLocal }

func (t *LocalTransport) Publish(ctx context.Context, msg *Message) error {
	if t.closed.Load() {
		return ErrBusClosed
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
		if err := sub.enqueue(ctx, msg.Clone()); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (t *LocalTransport) Subscribe(ctx context.Context, req SubscribeRequest) (Subscription, error) {
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
		t.closeWg.Add(1)
		go func() {
			defer t.closeWg.Done()
			sub.worker(subCtx)
		}()
	}

	t.mu.Lock()
	t.subs[sub.id] = sub
	t.mu.Unlock()
	return sub.closer(func() {
		cancel()
		t.mu.Lock()
		delete(t.subs, sub.id)
		t.mu.Unlock()
	}), nil
}

func (t *LocalTransport) Close(ctx context.Context) error {
	if !t.closed.CompareAndSwap(false, true) {
		return nil
	}

	t.mu.Lock()
	for _, sub := range t.subs {
		sub.cancel()
	}
	t.subs = make(map[uint64]*localSubscription)
	t.mu.Unlock()

	done := make(chan struct{})
	go func() {
		defer close(done)
		t.closeWg.Wait()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}

type localSubscription struct {
	id      uint64
	topic   string
	options SubscribeOptions
	handler Handler
	queue   chan *Message
	cancel  context.CancelFunc
	once    sync.Once
}

func (s *localSubscription) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-s.queue:
			if msg == nil {
				continue
			}
			_ = s.handler(ctx, msg)
		}
	}
}

func (s *localSubscription) enqueue(ctx context.Context, msg *Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.queue <- msg:
		return nil
	default:
		select {
		case <-ctx.Done():
			return ctx.Err()
		case s.queue <- msg:
			return nil
		}
	}
}

func (s *localSubscription) Close() error {
	s.once.Do(func() {
		s.cancel()
	})
	return nil
}

func (s *localSubscription) closer(cleanup func()) Subscription {
	return &subscriptionFunc{fn: func() error {
		s.once.Do(cleanup)
		return nil
	}}
}

type subscriptionFunc struct {
	once sync.Once
	fn   func() error
	err  error
}

func (s *subscriptionFunc) Close() error {
	s.once.Do(func() {
		if s.fn != nil {
			s.err = s.fn()
		}
	})
	return s.err
}

func (s *localSubscription) String() string {
	return fmt.Sprintf("localSubscription{id=%d,topic=%s}", s.id, s.topic)
}
