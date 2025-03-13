package zapx

import (
	"io"
)

type LevelEnablerFunc func(Level) bool

type TeeOption struct {
	Out io.Writer
	LevelEnablerFunc
}

func NewTee(tees []TeeOption, opts *Options) LoggerX {
	//if opts == nil {
	//	opts = NewOptions()
	//}
	//var cores []zapcore.Core
	//for _, tee := range tees {
	//	cfg := zap.NewProductionEncoderConfig()
	//	cfg.EncodeTime = zapcore.RFC3339TimeEncoder
	//	core := zapcore.NewCore(
	//		zapcore.NewJSONEncoder(cfg),
	//		zapcore.AddSync(tee.Out),
	//		zap.LevelEnablerFunc(tee.LevelEnablerFunc),
	//	)
	//	cores = append(cores, core)
	//}
	//logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))
	////return &Logger{l: zap.New(zapcore.NewTee(cores...), opts...)}
	//return &zapLogger{
	//	lg: logger,
	//	al: nil,
	//}
	return nil
}
