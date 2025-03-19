package zapx

import (
	"fmt"
	"go.uber.org/zap"
)

type loggerImpl struct {
	lg *zap.Logger
	al *zap.AtomicLevel
}

func NewLogger(l *zap.Logger, al *zap.AtomicLevel) LoggerX {
	return &loggerImpl{
		lg: l,
		al: al,
	}
}

func (l *loggerImpl) Debug(msg string, fields ...Field) {
	l.lg.Debug(msg, fields...)
}

func (l *loggerImpl) Debugf(format string, v ...interface{}) {
	l.lg.Debug(fmt.Sprintf(format, v...))
}

func (l *loggerImpl) Debugw(msg string, keysAndValues ...interface{}) {
	l.lg.Sugar().Debugw(msg, keysAndValues...)
}

func (l *loggerImpl) Info(msg string, fields ...Field) {
	l.lg.Info(msg, fields...)
}

func (l *loggerImpl) Infof(format string, v ...interface{}) {
	l.lg.Info(fmt.Sprintf(format, v...))
}

func (l *loggerImpl) Infow(msg string, keysAndValues ...interface{}) {
	l.lg.Sugar().Infow(msg, keysAndValues...)
}

func (l *loggerImpl) Warn(msg string, fields ...Field) {
	l.lg.Warn(msg, fields...)
}

func (l *loggerImpl) Warnf(format string, v ...interface{}) {
	l.lg.Warn(fmt.Sprintf(format, v...))
}

func (l *loggerImpl) Warnw(msg string, keysAndValues ...interface{}) {
	l.lg.Sugar().Warnw(msg, keysAndValues...)
}

func (l *loggerImpl) Error(msg string, fields ...Field) {
	l.lg.Error(msg, fields...)
}

func (l *loggerImpl) Errorf(format string, v ...interface{}) {
	l.lg.Error(fmt.Sprintf(format, v...))
}

func (l *loggerImpl) Errorw(msg string, keysAndValues ...interface{}) {
	l.lg.Sugar().Errorw(msg, keysAndValues...)
}

func (l *loggerImpl) Panic(msg string, fields ...Field) {
	l.lg.Panic(msg, fields...)
}

func (l *loggerImpl) Panicf(format string, v ...interface{}) {
	l.lg.Panic(fmt.Sprintf(format, v...))
}

func (l *loggerImpl) Panicw(msg string, keysAndValues ...interface{}) {
	l.lg.Sugar().Panicw(msg, keysAndValues...)
}

func (l *loggerImpl) Fatal(msg string, fields ...Field) {
	l.lg.Fatal(msg, fields...)
}

func (l *loggerImpl) Fatalf(format string, v ...interface{}) {
	l.lg.Fatal(fmt.Sprintf(format, v...))
}

func (l *loggerImpl) Fatalw(msg string, keysAndValues ...interface{}) {
	l.lg.Sugar().Fatalw(msg, keysAndValues...)
}

func (l *loggerImpl) Write(p []byte) (n int, err error) {
	l.lg.Info(string(p))
	return len(p), nil
}

func (l *loggerImpl) WithValues(keysAndValues ...interface{}) LoggerX {
	newLogger := l.lg.With(handleFields(l.lg, keysAndValues)...)
	return NewLogger(newLogger, l.al)
}

func (l *loggerImpl) WithName(name string) LoggerX {
	named := l.lg.Named(name)
	return NewLogger(named, l.al)
}

func (l *loggerImpl) With(field ...Field) LoggerX {
	newLogger := l.lg.With(field...)
	return NewLogger(newLogger, l.al)
}

func (l *loggerImpl) Flush() {
	l.lg.Sync()
}

func (l *loggerImpl) SetLevel(level Level) {
	if l.al != nil {
		l.al.SetLevel(level)
	}
}

func handleFields(l *zap.Logger, args []interface{}, additional ...zap.Field) []zap.Field {
	if len(args) == 0 {
		// fast-return if we have no suggared fields.
		return additional
	}

	fields := make([]zap.Field, 0, len(args)/2+len(additional))
	for i := 0; i < len(args); {
		if _, ok := args[i].(zap.Field); ok {
			l.DPanic("strongly-typed Zap Field passed to logr", zap.Any("zap field", args[i]))
			break
		}
		if i == len(args)-1 {
			l.DPanic("odd number of arguments passed as key-value pairs for logging", zap.Any("ignored key", args[i]))

			break
		}
		key, val := args[i], args[i+1]
		keyStr, isString := key.(string)
		if !isString {
			l.DPanic(
				"non-string key argument passed to logging, ignoring all later arguments",
				zap.Any("invalid key", key),
			)

			break
		}

		fields = append(fields, zap.Any(keyStr, val))
		i += 2
	}

	return append(fields, additional...)
}
