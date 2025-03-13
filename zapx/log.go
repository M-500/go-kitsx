package zapx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

var (
	glbStd = New(NewOptions())
	mu     sync.Mutex
)

// NewLogger creates a new logr.Logger using the given Zap Logger to log.
func NewLogger(l *zap.Logger, al *zap.AtomicLevel) LoggerX {
	return &zapLogger{
		lg: l,
		al: al,
	}
}
func New(opts *Options) LoggerX {
	if opts == nil {
		opts = NewOptions()
	}
	//var zapLevel zapcore.Level
	//if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
	//	zapLevel = zapcore.InfoLevel
	//}
	al := zap.NewAtomicLevelAt(opts.Level)
	encodeLevel := zapcore.CapitalLevelEncoder //  默认大写
	// when output to local path, with color is forbidden
	if opts.Format == ConsoleFormat && opts.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// set encoder
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "line_num",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    encodeLevel,
		EncodeTime:     timeEncoder,
		EncodeDuration: milliSecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	loggerConfig := &zap.Config{
		Level:             zap.NewAtomicLevelAt(opts.Level),
		Development:       opts.Development,
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         string(opts.Format),
		EncoderConfig:    encoderConfig,
		OutputPaths:      opts.OutputPaths,
		ErrorOutputPaths: opts.ErrorOutputPaths,
	}

	var err error
	//l, err := loggerConfig.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(1))
	l, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	logger := &zapLogger{
		lg: l.Named(opts.Name),
		al: &al,
	}
	zap.RedirectStdLog(l)
	return logger
}

type zapLogger struct {
	lg *zap.Logger
	al *zap.AtomicLevel
}

func (z *zapLogger) Info(msg string, fields ...Field) {
	z.lg.Info(msg, fields...)
}

func (z *zapLogger) Infof(format string, v ...interface{}) {
	z.lg.Sugar().Infof(format, v...)
}

func (z *zapLogger) Infow(msg string, keysAndValues ...interface{}) {
	z.lg.Sugar().Infow(msg, keysAndValues...)
}

func (z *zapLogger) Debug(msg string, fields ...Field) {
	z.lg.Debug(msg, fields...)
}

func (z *zapLogger) Debugf(format string, v ...interface{}) {
	z.lg.Sugar().Debugf(format, v...)
}

func (z *zapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	z.lg.Sugar().Debugw(msg, keysAndValues...)
}

func (z *zapLogger) Warn(msg string, fields ...Field) {
	z.lg.Warn(msg, fields...)
}

func (z *zapLogger) Warnf(format string, v ...interface{}) {
	z.lg.Sugar().Warnf(format, v...)
}

func (z *zapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	z.lg.Sugar().Warnw(msg, keysAndValues...)
}

func (z *zapLogger) Error(msg string, fields ...Field) {
	z.lg.Error(msg, fields...)
}

func (z *zapLogger) Errorf(format string, v ...interface{}) {
	z.lg.Sugar().Errorf(format, v...)
}

func (z *zapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	z.lg.Sugar().Errorw(msg, keysAndValues...)
}

func (z *zapLogger) Panic(msg string, fields ...Field) {
	z.lg.Panic(msg, fields...)
}

func (z *zapLogger) Panicf(format string, v ...interface{}) {
	z.lg.Sugar().Panicf(format, v...)
}

func (z *zapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	z.lg.Sugar().Panicw(msg, keysAndValues...)
}

func (z *zapLogger) Fatal(msg string, fields ...Field) {
	z.lg.Fatal(msg, fields...)
}

func (z *zapLogger) Fatalf(format string, v ...interface{}) {
	z.lg.Sugar().Fatalf(format, v...)
}

func (z *zapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	z.lg.Sugar().Fatalw(msg, keysAndValues...)
}

func (z *zapLogger) Write(p []byte) (n int, err error) {
	z.lg.Info(string(p))
	return len(p), nil
}

func (z *zapLogger) WithValues(keysAndValues ...interface{}) LoggerX {
	newLogger := z.lg.With(handleFields(z.lg, keysAndValues)...)
	return NewLogger(newLogger, z.al)
}

func (z *zapLogger) With(field ...Field) LoggerX {
	newLogger := z.lg.With(field...)
	return NewLogger(newLogger, z.al)
}

func (z *zapLogger) WithName(name string) LoggerX {
	named := z.lg.Named(name)
	return NewLogger(named, z.al)
}

func (z *zapLogger) Flush() {
	_ = z.lg.Sync()
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

// var std = New()
//
//	func New() LoggerX {
//		al := zap.NewAtomicLevelAt(DebugLevel)
//		cfg := zap.NewProductionEncoderConfig()
//		cfg.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
//		cfg.LineEnding = zapcore.DefaultLineEnding
//		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
//		cfg.ConsoleSeparator = " | "
//
//		core := zapcore.NewCore(
//			zapcore.NewConsoleEncoder(cfg), // 编码器，用来定义日志的输出格式
//			zapcore.AddSync(os.Stdout),     // 指定输出位置
//			al,                             // 日志级别
//		)
//		return &logger{
//			lg: zap.New(core),
//			al: &al,
//		}
//	}
//
//	func NewLogger(level Level) LoggerX {
//		al := zap.NewAtomicLevelAt(level)
//		return &logger{
//			lg: nil,
//			al: &al,
//		}
//	}

func ReplaceDefault(l LoggerX) { glbStd = l }

func Default() LoggerX { return glbStd }

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
func WithName(s string) LoggerX { return glbStd.WithName(s) }

func Flush() {
	glbStd.Flush()
}
