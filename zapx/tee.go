package zapx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LevelEnablerFunc func(Level) bool

type TeeOption struct {
	LevelEnablerFunc
	Encoder zapcore.Encoder
	Syncer  zapcore.WriteSyncer
}

func NewTee(tees []TeeOption, options ...zap.Option) LoggerX {
	var cores []zapcore.Core
	for _, tee := range tees {
		core := zapcore.NewCore(tee.Encoder, tee.Syncer, zap.LevelEnablerFunc(tee.LevelEnablerFunc))
		cores = append(cores, core)
	}
	logger := zap.New(zapcore.NewTee(cores...), options...)
	return &loggerImpl{
		lg: logger,
		al: nil,
	}
}
