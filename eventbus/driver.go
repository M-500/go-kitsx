package eventbus

import "context"

type Subscription interface {
	Close() error
}

type Driver interface {
	Name() string
	Publish(ctx context.Context, msg *Message) error
	Subscribe(ctx context.Context, req SubscribeRequest) (Subscription, error)
	Close(ctx context.Context) error
}

type SubscribeRequest struct {
	Topic         string
	Options       SubscribeOptions
	FinalHandler  Handler
	Middlewares   []Middleware
	SharedMetrics *metricsRegistry
}
