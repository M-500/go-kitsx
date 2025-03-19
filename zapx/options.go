package zapx

import (
	"go.uber.org/zap/zapcore"
	"io"
	"os"
)

type Option func(*OptionsX)

func NewOptionX(opts ...Option) *OptionsX {
	res := &OptionsX{
		Out:               os.Stdout,
		Level:             DebugLevel,
		Format:            ConsoleFormat,
		DisableCaller:     false,
		DisableStacktrace: false,
		EnableColor:       true,
		Development:       true,
		Name:              "Default",
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

type OptionsX struct {
	Out               io.Writer
	Level             Level      `json:"level"              mapstructure:"level"`              // The log level: Debug, Info, Warn, Error, Panic, Fatal。
	Format            FormatType `json:"format"             mapstructure:"format"`             // 编码器类型
	DisableCaller     bool       `json:"disable-caller"     mapstructure:"disable-caller"`     // Whether to enable caller, if enabled, will display the file, function, and line number where the call log is located in the log.
	DisableStacktrace bool       `json:"disable-stacktrace" mapstructure:"disable-stacktrace"` // Is printing stack information prohibited at Panic and above levels.
	EnableColor       bool       `json:"enable-color"       mapstructure:"enable-color"`       // Whether to enable color output
	Development       bool       `json:"development"        mapstructure:"development"`        // Is it a development mode. If it is in development mode, stack tracing will be performed on DPanicLevel
	Name              string     `json:"name"               mapstructure:"name"`               // logger Name
}

func (o *OptionsX) BuildCore() zapcore.Core {
	if o.Out == nil {
		o.Out = os.Stdout
	}
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
	return zapcore.NewCore(encoder, zapcore.AddSync(o.Out), o.Level)
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
