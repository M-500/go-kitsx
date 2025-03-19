package zapx

import (
	"context"
	"go.uber.org/zap"
)

var (
	glbStd = New(NewOptionX())
)

func New(opt *OptionsX, opts ...zap.Option) LoggerX {
	if opt == nil {
		opt = NewOptionX()
	}
	if opt.DisableCaller {
		opts = append(opts, zap.AddCallerSkip(1))
	}
	if opt.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(zap.DPanicLevel))
	}
	logger := zap.New(opt.BuildCore(), opts...)
	return &loggerImpl{
		lg: logger,
		al: nil,
	}
}

func Debug(msg string, fields ...Field) {
	glbStd.Debug(msg, fields...)
}

func Debugf(format string, v ...interface{}) {
	glbStd.Debugf(format, v...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	glbStd.Debugw(msg, keysAndValues...)
}

func Info(msg string, fields ...Field) {
	glbStd.Info(msg, fields...)
}

func Infof(format string, v ...interface{}) {
	glbStd.Infof(format, v...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	glbStd.Infow(msg, keysAndValues...)
}

func Warn(msg string, fields ...Field) {
	glbStd.Warn(msg, fields...)
}

func Warnf(format string, v ...interface{}) {
	glbStd.Warnf(format, v...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	glbStd.Warnw(msg, keysAndValues...)
}

func Error(msg string, fields ...Field) {
	glbStd.Error(msg, fields...)
}

func Errorf(format string, v ...interface{}) {
	glbStd.Errorf(format, v...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	glbStd.Errorw(msg, keysAndValues...)
}

func Panic(msg string, fields ...Field) {
	glbStd.Panic(msg, fields...)
}

func Panicf(format string, v ...interface{}) {
	glbStd.Panicf(format, v...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	glbStd.Panicw(msg, keysAndValues...)
}

func Fatal(msg string, fields ...Field) {
	glbStd.Fatal(msg, fields...)
}

func Fatalf(format string, v ...interface{}) {
	glbStd.Fatalf(format, v...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	glbStd.Fatalw(msg, keysAndValues...)
}

func Write(p []byte) (n int, err error) {

	return glbStd.Write(p)
}

func WithValues(keysAndValues ...interface{}) LoggerX {

	return glbStd.WithValues(keysAndValues...)
}

func WithName(name string) LoggerX {
	return glbStd.WithName(name)
}

func With(field ...Field) LoggerX {
	return glbStd.With(field...)
}

func WithContext(ctx context.Context) context.Context {
	return glbStd.WithContext(ctx)
}

func Flush() {
	glbStd.Flush()
}

func SetLevel(level Level) {
	glbStd.SetLevel(level)
}
