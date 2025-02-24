package zapx

import (
	"go.uber.org/zap/zapcore"
	"io"
)

type options struct {
	output io.Writer
	prefix string // prefix is the prefix to add to all log messages.
	hooks  []func(entry zapcore.Entry) error
}

type Option func(*options)

func RegHook(fn func(entry zapcore.Entry) error) Option {
	return func(o *options) {
		o.hooks = append(o.hooks, fn)
	}
}

func WithOutput(output io.Writer) Option {
	return func(o *options) {
		o.output = output
	}
}

func WithPrefix(prefix string) Option {
	return func(o *options) {
		o.prefix = prefix
	}
}

func initOptions(opts ...Option) (o *options) {
	o = &options{}
	for _, opt := range opts {
		opt(o)
	}

	if o.output == nil {
		o.output = io.Discard
	}

	return
}
