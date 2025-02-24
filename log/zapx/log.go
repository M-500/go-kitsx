package zapx

import (
	"go.uber.org/zap"
	"sync"
)

type loggerX struct {
	opt *options
	l   *zap.Logger      // l is the underlying zap logger.
	al  *zap.AtomicLevel // al is the atomic level for the logger.
	mu  sync.Mutex
}

var std = NewLoggerX()

func NewLoggerX(opts ...Option) LoggerX {
	logger := &loggerX{opt: initOptions(opts...)}

	// 注册一个钩子函数，用于在日志输出前添加前缀
	logger.l.WithOptions(zap.Hooks(logger.opt.hooks...))
	return logger
}

func SetOptions(opts ...Option) {
	std.SetOptions(opts...)
}

func (l *loggerX) SetOptions(opts ...Option) {
	l.mu.Lock()
	defer l.mu.Unlock()

	for _, opt := range opts {
		opt(l.opt)
	}
}

func (l *loggerX) Info(args ...any) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerX) InfoWithFields(msg string, args ...Field) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerX) Infof(format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerX) Debug(args ...any) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerX) DebugWithFields(msg string, args ...Field) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerX) Debugf(format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerX) Warn(args ...any) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerX) WarnWithFields(msg string, args ...Field) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerX) Warnf(format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerX) Error(args ...any) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerX) ErrorWithFields(msg string, args ...Field) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerX) Errorf(format string, args ...any) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerX) Fatal(args ...any) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerX) FatalWithFields(msg string, args ...Field) {
	//TODO implement me
	panic("implement me")
}

func (l *loggerX) Fatalf(format string, args ...any) {
	//TODO implement me
	panic("implement me")
}
