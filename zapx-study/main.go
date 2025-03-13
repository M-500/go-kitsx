package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newCustomLogger() (*zap.Logger, error) {
	cfg := zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:     "time",
			LevelKey:    "level",
			NameKey:     "logger",
			CallerKey:   "", // 不记录日志调用位置
			FunctionKey: zapcore.OmitKey,
			MessageKey:  "message",
			LineEnding:  zapcore.DefaultLineEnding,
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			EncodeTime:  zapcore.RFC3339TimeEncoder,
			//EncodeDuration: zapcore.SecondsDurationEncoder,
			//EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout", "test.log"},
		ErrorOutputPaths: []string{"error.log"},
	}
	return cfg.Build()
}

func main() {
	logger, _ := newCustomLogger()
	defer logger.Sync()

	// 增加一个 skip 选项，触发 zap 内部 error，将错误输出到 error.log
	logger = logger.WithOptions(zap.AddCallerSkip(100))

	logger.Info("Info msg")
	logger.Error("Error msg")
}
