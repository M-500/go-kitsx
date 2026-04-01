package eventbus

import (
	"context"
	"fmt"
	"time"
)

type Handler func(context.Context, *Message) error

type Middleware func(Handler) Handler

func Chain(middlewares []Middleware, final Handler) Handler {
	if len(middlewares) == 0 {
		return final
	}
	wrapped := final
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}
	return wrapped
}

func withTimeout(next Handler, timeout time.Duration) Handler {
	if timeout <= 0 {
		return next
	}
	return func(ctx context.Context, msg *Message) error {
		childCtx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return next(childCtx, msg)
	}
}

func withRecover(metrics *metricsRegistry) Middleware {
	return func(next Handler) Handler {
		return func(ctx context.Context, msg *Message) (err error) {
			defer func() {
				if r := recover(); r != nil {
					metrics.panics.Add(1)
					err = fmt.Errorf("eventbus: handler panic topic=%s id=%s: %v", msg.Topic, msg.ID, r)
				}
			}()
			return next(ctx, msg)
		}
	}
}

func withRetry(metrics *metricsRegistry, policy RetryPolicy) Middleware {
	policy = policy.normalized()
	return func(next Handler) Handler {
		return func(ctx context.Context, msg *Message) error {
			var lastErr error
			delay := policy.InitialInterval

			for attempt := 1; attempt <= policy.MaxAttempts; attempt++ {
				if attempt > 1 {
					metrics.retries.Add(1)
					timer := time.NewTimer(delay)
					select {
					case <-ctx.Done():
						timer.Stop()
						return ctx.Err()
					case <-timer.C:
					}

					nextDelay := time.Duration(float64(delay) * policy.Multiplier)
					if nextDelay > policy.MaxInterval {
						nextDelay = policy.MaxInterval
					}
					delay = nextDelay
				}

				if err := next(ctx, msg); err != nil {
					lastErr = err
					continue
				}
				return nil
			}
			return lastErr
		}
	}
}
