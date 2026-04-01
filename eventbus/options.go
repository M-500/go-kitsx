package eventbus

import "time"

type MatchMode string

const (
	MatchExact   MatchMode = "exact"
	MatchPrefix  MatchMode = "prefix"
	MatchPattern MatchMode = "pattern"
)

const (
	TransportLocal   = "local"
	TransportChannel = "channel"
	TransportKafka   = "kafka"
)

type RetryPolicy struct {
	MaxAttempts     int
	InitialInterval time.Duration
	MaxInterval     time.Duration
	Multiplier      float64
}

func (r RetryPolicy) normalized() RetryPolicy {
	if r.MaxAttempts <= 0 {
		r.MaxAttempts = 1
	}
	if r.InitialInterval <= 0 {
		r.InitialInterval = 200 * time.Millisecond
	}
	if r.MaxInterval <= 0 {
		r.MaxInterval = 3 * time.Second
	}
	if r.Multiplier < 1 {
		r.Multiplier = 2
	}
	return r
}

type SubscribeOptions struct {
	Transport     string
	ConsumerGroup string
	Concurrency   int
	Buffer        int
	Timeout       time.Duration
	RetryPolicy   RetryPolicy
	MatchMode     MatchMode
}

func (o SubscribeOptions) normalized() SubscribeOptions {
	if o.Transport == "" {
		o.Transport = TransportLocal
	}
	if o.Concurrency <= 0 {
		o.Concurrency = 1
	}
	if o.Buffer <= 0 {
		o.Buffer = 64
	}
	if o.MatchMode == "" {
		o.MatchMode = MatchExact
	}
	o.RetryPolicy = o.RetryPolicy.normalized()
	return o
}

type PublishOptions struct {
	Transport string
}

func (o PublishOptions) normalized() PublishOptions {
	if o.Transport == "" {
		o.Transport = TransportLocal
	}
	return o
}
