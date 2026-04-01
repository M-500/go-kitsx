package eventbus

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func ExampleBus_local() {
	bus := NewBus()
	defer func() { _ = bus.Close(context.Background()) }()

	done := make(chan struct{})
	_, _ = bus.Subscribe("user.created", func(ctx context.Context, msg *Message) error {
		var payload struct {
			Name string `json:"name"`
		}
		_ = msg.DecodeJSON(&payload)
		fmt.Println(payload.Name)
		close(done)
		return nil
	}, SubscribeOptions{
		Transport:   TransportLocal,
		Concurrency: 1,
	})

	_ = bus.PublishJSON(context.Background(), "user.created", map[string]string{"name": "alice"}, PublishOptions{
		Transport: TransportLocal,
	})
	<-done
	// Output:
	// alice
}

func ExampleBus_channel() {
	bus := NewBus()
	defer func() { _ = bus.Close(context.Background()) }()

	done := make(chan struct{})
	_, _ = bus.Subscribe("order.*", func(ctx context.Context, msg *Message) error {
		fmt.Println(msg.Topic)
		close(done)
		return nil
	}, SubscribeOptions{
		Transport: TransportChannel,
		MatchMode: MatchPattern,
	})

	_ = bus.PublishJSON(context.Background(), "order.created", map[string]any{"id": 1}, PublishOptions{
		Transport: TransportChannel,
	})
	<-done
	// Output:
	// order.created
}

func TestLocalTransportRetryAndMetrics(t *testing.T) {
	bus := NewBus()
	defer func() { _ = bus.Close(context.Background()) }()

	var mu sync.Mutex
	attempts := 0
	done := make(chan struct{})

	_, err := bus.Subscribe("task.run", func(ctx context.Context, msg *Message) error {
		mu.Lock()
		defer mu.Unlock()
		attempts++
		if attempts < 3 {
			return fmt.Errorf("temporary")
		}
		close(done)
		return nil
	}, SubscribeOptions{
		Transport: TransportLocal,
		RetryPolicy: RetryPolicy{
			MaxAttempts:     3,
			InitialInterval: 10 * time.Millisecond,
			MaxInterval:     10 * time.Millisecond,
		},
	})
	if err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	if err := bus.PublishJSON(context.Background(), "task.run", map[string]any{"ok": true}, PublishOptions{
		Transport: TransportLocal,
	}); err != nil {
		t.Fatalf("publish: %v", err)
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("handler was not delivered")
	}

	metrics := bus.Metrics()
	if metrics.Delivered != 1 {
		t.Fatalf("expected delivered=1 got %d", metrics.Delivered)
	}
	if metrics.Retries != 2 {
		t.Fatalf("expected retries=2 got %d", metrics.Retries)
	}
}

func TestKafkaTransportWithMock(t *testing.T) {
	producer := &mockKafkaProducer{}
	consumer := &mockKafkaConsumer{}
	bus := NewBus()
	if err := bus.RegisterTransport(TransportKafka, NewKafkaTransport(producer, consumer)); err != nil {
		t.Fatalf("register kafka: %v", err)
	}
	defer func() { _ = bus.Close(context.Background()) }()

	done := make(chan struct{})
	_, err := bus.Subscribe("audit.*", func(ctx context.Context, msg *Message) error {
		close(done)
		return nil
	}, SubscribeOptions{
		Transport:     TransportKafka,
		ConsumerGroup: "audit-service",
		MatchMode:     MatchPattern,
	})
	if err != nil {
		t.Fatalf("subscribe kafka: %v", err)
	}

	if err := bus.PublishJSON(context.Background(), "audit.created", map[string]string{"id": "1"}, PublishOptions{
		Transport: TransportKafka,
	}); err != nil {
		t.Fatalf("publish kafka: %v", err)
	}

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		if consumer.ready() {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	consumer.emit(producer.last)
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("kafka message not consumed")
	}
}

type mockKafkaProducer struct {
	last KafkaRecord
}

func (m *mockKafkaProducer) Publish(ctx context.Context, record KafkaRecord) error {
	m.last = record
	return nil
}

type mockKafkaConsumer struct {
	mu      sync.RWMutex
	handler func(context.Context, KafkaRecord) error
}

func (m *mockKafkaConsumer) Subscribe(ctx context.Context, req KafkaConsumeRequest, handler func(context.Context, KafkaRecord) error) error {
	m.mu.Lock()
	m.handler = handler
	m.mu.Unlock()
	<-ctx.Done()
	return nil
}

func (m *mockKafkaConsumer) emit(record KafkaRecord) {
	m.mu.RLock()
	handler := m.handler
	m.mu.RUnlock()
	if handler != nil {
		_ = handler(context.Background(), record)
	}
}

func (m *mockKafkaConsumer) ready() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.handler != nil
}
