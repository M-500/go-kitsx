package zapx

import "context"

type key int

const (
	logContextKey key = iota
)

func (z *zapLogger) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, logContextKey, z)
}
