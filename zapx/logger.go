package zapx

import "go.uber.org/zap"

func New(opt *OptionsX, opts ...zap.Option) LoggerX {
	if opt == nil {
		opt = NewOptionX()
	}
	logger := zap.New(opt.BuildCore(), opts...)
	return &loggerImpl{
		lg: logger,
		al: nil,
	}
}
