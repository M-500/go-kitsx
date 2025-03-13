package logx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

type Level = zapcore.Level

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	PanicLevel = zapcore.PanicLevel
	FatalLevel = zapcore.FatalLevel
)

var std = NewLogger(os.Stdout, InfoLevel)

type logger struct {
	l       *zap.Logger
	al      *zap.AtomicLevel
	rotaCfg *RotateConfig
}

func (l *logger) Debug(msg string, fields ...zap.Field) {
	l.l.Debug(msg, fields...)
}
func (l *logger) Info(msg string, fields ...zap.Field) {
	l.l.Info(msg, fields...)
}

func (l *logger) Warn(msg string, fields ...zap.Field) {
	l.l.Warn(msg, fields...)
}

func (l *logger) Error(msg string, fields ...zap.Field) {
	l.l.Error(msg, fields...)
}

func (l *logger) Panic(msg string, fields ...zap.Field) {
	l.l.Panic(msg, fields...)
}

func (l *logger) Fatal(msg string, fields ...zap.Field) {
	l.l.Fatal(msg, fields...)
}

func (l *logger) Sync() {
	_ = l.l.Sync()
}

func NewLogger(out io.Writer, level Level, zapOpts ...zap.Option) LoggerX {
	if out == nil {
		out = os.Stdout
	}
	al := zap.NewAtomicLevel()
	al.SetLevel(level)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(out),
		al,
	)
	l := zap.New(core, zapOpts...)
	return &logger{l: l, al: &al}
}

func Info(msg string, fields ...zap.Field) {
	std.Info(msg, fields...)
}

func Debug(msg string, fields ...zap.Field) {
	std.Debug(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	std.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	std.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	std.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	std.Fatal(msg, fields...)
}

func Sync() {
	std.Sync()
}
