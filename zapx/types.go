package zapx

import (
	"context"
	"go.uber.org/zap/zapcore"
)

// Defines common log fields.
const (
	KeyRequestID   string = "requestID"
	KeyUsername    string = "username"
	KeyWatcherName string = "watcher"
)

// Field is an alias for the field structure in the underlying log frame.
type Field = zapcore.Field

// Level is an alias for the level structure in the underlying log frame.
type Level = zapcore.Level

var (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	PanicLevel = zapcore.PanicLevel
	FatalLevel = zapcore.FatalLevel
)

type LoggerX interface {
	//SetLevel(level Level)
	Debug(msg string, fields ...Field)
	Debugf(format string, v ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Info(msg string, fields ...Field)
	Infof(format string, v ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warn(msg string, fields ...Field)
	Warnf(format string, v ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Error(msg string, fields ...Field)
	Errorf(format string, v ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panic(msg string, fields ...Field)
	Panicf(format string, v ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatal(msg string, fields ...Field)
	Fatalf(format string, v ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	Write(p []byte) (n int, err error)

	// WithValues adds some key-value pairs of context to a LoggerX.
	// See Info for documentation on how key/value pairs work.
	//可以返回一个携带指定 key-value 的子 Logger，供后面使用
	WithValues(keysAndValues ...interface{}) LoggerX

	// WithName adds a new element to the LoggerX's name.
	// Successive calls with WithName continue to append
	// suffixes to the LoggerX's name.  It's strongly recommended
	// that name segments contain only letters, digits, and hyphens
	// (see the package documentation for more information).
	WithName(name string) LoggerX

	With(field ...Field) LoggerX

	// WithContext returns a copy of context in which the log value is set.
	WithContext(ctx context.Context) context.Context

	// Flush calls the underlying Core's Sync method, flushing any buffered
	// log entries. Applications should take care to call Sync before exiting.
	Flush()
}
