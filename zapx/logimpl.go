package zapx

import (
	"context"
	"go.uber.org/zap"
)

type loggerImpl struct {
	lg *zap.Logger
	al *zap.AtomicLevel
}

func (l loggerImpl) Debug(msg string, fields ...Field) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Debugf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Debugw(msg string, keysAndValues ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Info(msg string, fields ...Field) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Infof(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Infow(msg string, keysAndValues ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Warn(msg string, fields ...Field) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Warnf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Warnw(msg string, keysAndValues ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Error(msg string, fields ...Field) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Errorf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Errorw(msg string, keysAndValues ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Panic(msg string, fields ...Field) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Panicf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Panicw(msg string, keysAndValues ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Fatal(msg string, fields ...Field) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Fatalf(format string, v ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Fatalw(msg string, keysAndValues ...interface{}) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Write(p []byte) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) WithValues(keysAndValues ...interface{}) LoggerX {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) WithName(name string) LoggerX {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) With(field ...Field) LoggerX {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) WithContext(ctx context.Context) context.Context {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) Flush() {
	//TODO implement me
	panic("implement me")
}

func (l loggerImpl) SetLevel(level Level) {
	//TODO implement me
	panic("implement me")
}
