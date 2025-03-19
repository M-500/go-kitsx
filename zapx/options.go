package zapx

import (
	"go.uber.org/zap/zapcore"
)

type Option func(*OptionsX)

func NewOptionX(opts ...Option) *OptionsX {
	return &OptionsX{}
}

type OptionsX struct {
	OutputPaths       []string   `json:"output-paths"       mapstructure:"output-paths"`       // Support output to multiple outputs, separated by commas. Supports outputting to standard output (UWP) and files.
	ErrorOutputPaths  []string   `json:"error-output-paths" mapstructure:"error-output-paths"` // 错误日志输出路径
	Level             Level      `json:"level"              mapstructure:"level"`              // The log level: Debug, Info, Warn, Error, Panic, Fatal。
	Format            FormatType `json:"format"             mapstructure:"format"`             // 编码器类型
	DisableCaller     bool       `json:"disable-caller"     mapstructure:"disable-caller"`     // Whether to enable caller, if enabled, will display the file, function, and line number where the call log is located in the log.
	DisableStacktrace bool       `json:"disable-stacktrace" mapstructure:"disable-stacktrace"` // Is printing stack information prohibited at Panic and above levels.
	EnableColor       bool       `json:"enable-color"       mapstructure:"enable-color"`       // Whether to enable color output
	Development       bool       `json:"development"        mapstructure:"development"`        // Is it a development mode. If it is in development mode, stack tracing will be performed on DPanicLevel
	Name              string     `json:"name"               mapstructure:"name"`               // logger Name
}

func (o *OptionsX) BuildCore() zapcore.Core {
	encodeLevel := zapcore.CapitalLevelEncoder
	// when output to local path, with color is forbidden
	if o.Format == ConsoleFormat && o.EnableColor {
		encodeLevel = zapcore.CapitalColorLevelEncoder
	}
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
	var encoder zapcore.Encoder
	if o.Format == ConsoleFormat {
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // use JSON encoder in production
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig) // use console encoder in development
	}
	return zapcore.NewCore(encoder, nil, o.Level)
}

func WithOutputPaths(outputPaths []string) Option {
	return func(o *OptionsX) {
		o.OutputPaths = outputPaths
	}
}

func WithErrorOutputPaths(errorOutputPaths []string) Option {
	return func(o *OptionsX) {
		o.ErrorOutputPaths = errorOutputPaths
	}
}

func WithLevel(level Level) Option {
	return func(OptionsX *OptionsX) {
		OptionsX.Level = level
	}
}

func WithEnableColor(enableColor bool) Option {
	return func(OptionsX *OptionsX) {
		OptionsX.EnableColor = enableColor
	}
}

func WithDisableCaller(disableCaller bool) Option {
	return func(OptionsX *OptionsX) {
		OptionsX.DisableCaller = disableCaller
	}
}

func WithDevelopment(development bool) Option {
	return func(OptionsX *OptionsX) {
		OptionsX.Development = development
	}
}

func WithFormat(format FormatType) Option {
	return func(OptionsX *OptionsX) {
		OptionsX.Format = format
	}
}
