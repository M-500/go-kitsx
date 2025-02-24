package zapx

import "go.uber.org/zap/zapcore"

type LoggerX interface {
	Info(args ...any)
	InfoWithFields(msg string, args ...Field)
	Infof(format string, args ...any)

	Debug(args ...any)
	DebugWithFields(msg string, args ...Field)
	Debugf(format string, args ...any)

	Warn(args ...any)
	WarnWithFields(msg string, args ...Field)
	Warnf(format string, args ...any)

	Error(args ...any)
	ErrorWithFields(msg string, args ...Field)
	Errorf(format string, args ...any)

	Fatal(args ...any)
	FatalWithFields(msg string, args ...Field)
	Fatalf(format string, args ...any)

	SetOptions(opts ...Option)
}

type AlertHook interface {
	Alert(entry zapcore.Entry) bool
}
