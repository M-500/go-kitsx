package eventbus

import (
	"context"
	"encoding/json"
	"fmt"
)

type KafkaRecord struct {
	Topic     string
	Key       string
	Headers   map[string]string
	Value     []byte
	Partition int32
	Offset    int64
}

type KafkaConsumeRequest struct {
	Topic         string
	ConsumerGroup string
}

type KafkaProducer interface {
	Publish(ctx context.Context, record KafkaRecord) error
}

type KafkaConsumer interface {
	Subscribe(ctx context.Context, req KafkaConsumeRequest, handler func(context.Context, KafkaRecord) error) error
}

type KafkaTransport struct {
	producer KafkaProducer
	consumer KafkaConsumer
}

func NewKafkaTransport(producer KafkaProducer, consumer KafkaConsumer) *KafkaTransport {
	return &KafkaTransport{
		producer: producer,
		consumer: consumer,
	}
}

func (t *KafkaTransport) Name() string { return TransportKafka }

func (t *KafkaTransport) Publish(ctx context.Context, msg *Message) error {
	if t.producer == nil {
		return fmt.Errorf("eventbus: kafka producer is nil")
	}
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return t.producer.Publish(ctx, KafkaRecord{
		Topic:   msg.Topic,
		Key:     msg.Key,
		Headers: cloneStringMap(msg.Headers),
		Value:   raw,
	})
}

func (t *KafkaTransport) Subscribe(ctx context.Context, req SubscribeRequest) (Subscription, error) {
	if t.consumer == nil {
		return nil, fmt.Errorf("eventbus: kafka consumer is nil")
	}

	subCtx, cancel := context.WithCancel(ctx)
	done := make(chan error, 1)
	go func() {
		done <- t.consumer.Subscribe(subCtx, KafkaConsumeRequest{
			Topic:         req.Topic,
			ConsumerGroup: req.Options.ConsumerGroup,
		}, func(msgCtx context.Context, record KafkaRecord) error {
			var msg Message
			if err := json.Unmarshal(record.Value, &msg); err != nil {
				return err
			}
			if msg.Topic == "" {
				msg.Topic = record.Topic
			}
			if !matches(req.Topic, msg.Topic, req.Options.MatchMode) {
				return nil
			}
			return req.FinalHandler(msgCtx, msg.Clone())
		})
	}()

	return &subscriptionFunc{fn: func() error {
		cancel()
		select {
		case err := <-done:
			return err
		default:
			return nil
		}
	}}, nil
}

func (t *KafkaTransport) Close(ctx context.Context) error {
	return nil
}

func cloneStringMap(in map[string]string) map[string]string {
	if len(in) == 0 {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
