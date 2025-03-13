package cachex

import "context"

type CacheX[K comparable, V any] interface {
	Set(ctx context.Context, k K, v V) error
	Get(ctx context.Context, k K) (V, error)
	Del(ctx context.Context, k K) error
	Keys(ctx context.Context) ([]K, error)
}
