package zapx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type FormatType string

const (
	flagLevel             = "log.level"
	flagDisableCaller     = "log.disable-caller"
	flagDisableStacktrace = "log.disable-stacktrace"
	flagFormat            = "log.format"
	flagEnableColor       = "log.enable-color"
	flagOutputPaths       = "log.output-paths"
	flagErrorOutputPaths  = "log.error-output-paths"
	flagDevelopment       = "log.development"
	flagName              = "log.name"

	ConsoleFormat FormatType = "console"
	JsonFormat    FormatType = "json"
)

type Option func(*Options)

type Options struct {
	OutputPaths       []string      `json:"output-paths"       mapstructure:"output-paths"`       // Support output to multiple outputs, separated by commas. Supports outputting to standard output (UWP) and files.
	ErrorOutputPaths  []string      `json:"error-output-paths" mapstructure:"error-output-paths"` // 错误日志输出路径
	Level             zapcore.Level `json:"level"              mapstructure:"level"`              // The log level: Debug, Info, Warn, Error, Panic, Fatal。
	Format            FormatType    `json:"format"             mapstructure:"format"`             // 编码器类型
	DisableCaller     bool          `json:"disable-caller"     mapstructure:"disable-caller"`     // Whether to enable caller, if enabled, will display the file, function, and line number where the call log is located in the log.
	DisableStacktrace bool          `json:"disable-stacktrace" mapstructure:"disable-stacktrace"` // Is printing stack information prohibited at Panic and above levels.
	EnableColor       bool          `json:"enable-color"       mapstructure:"enable-color"`       // Whether to enable color output
	Development       bool          `json:"development"        mapstructure:"development"`        // Is it a development mode. If it is in development mode, stack tracing will be performed on DPanicLevel
	Name              string        `json:"name"               mapstructure:"name"`               // logger Name
}

func WithOutputPaths(outputPaths []string) Option {
	return func(o *Options) {
		o.OutputPaths = outputPaths
	}
}

func WithErrorOutputPaths(errorOutputPaths []string) Option {
	return func(o *Options) {
		o.ErrorOutputPaths = errorOutputPaths
	}
}

func WithLevel(level zapcore.Level) Option {
	return func(options *Options) {
		options.Level = level
	}
}

func WithEnableColor(enableColor bool) Option {
	return func(options *Options) {
		options.EnableColor = enableColor
	}
}

func WithDisableCaller(disableCaller bool) Option {
	return func(options *Options) {
		options.DisableCaller = disableCaller
	}
}

func WithDevelopment(development bool) Option {
	return func(options *Options) {
		options.Development = development
	}
}

func WithFormat(format FormatType) Option {
	return func(options *Options) {
		options.Format = format
	}
}

// NewOptions creates an Options object with default parameters.
func NewOptions(opts ...Option) *Options {
	res := &Options{
		Level:             DebugLevel,
		DisableCaller:     true,
		DisableStacktrace: true,
		Format:            ConsoleFormat,
		EnableColor:       true,
		Development:       true,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

func (o *Options) Build() error {
	//var zapLevel zapcore.Level
	//if err := zapLevel.UnmarshalText([]byte(o.Level)); err != nil {
	//	zapLevel = zapcore.InfoLevel
	//}
	encodeLevel := zapcore.CapitalLevelEncoder
	if o.Format == ConsoleFormat && o.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}

	zc := &zap.Config{
		Level:             zap.NewAtomicLevelAt(o.Level),
		Development:       o.Development,
		DisableCaller:     o.DisableCaller,
		DisableStacktrace: o.DisableStacktrace,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: string(o.Format),
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "timestamp",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    encodeLevel,
			EncodeTime:     timeEncoder,
			EncodeDuration: milliSecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
		},
		OutputPaths:      o.OutputPaths,
		ErrorOutputPaths: o.ErrorOutputPaths,
	}
	logger, err := zc.Build(zap.AddStacktrace(zapcore.PanicLevel))
	if err != nil {
		return err
	}
	zap.RedirectStdLog(logger.Named(o.Name))
	zap.ReplaceGlobals(logger)

	return nil
}
