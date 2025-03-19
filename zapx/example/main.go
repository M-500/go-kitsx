package main

import (
	"go-kitsx/zapx"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	//zapx.Info("hello world")
	tees := []zapx.TeeOption{
		{
			LevelEnablerFunc: zapx.DebugLevel.Enabled,
			Encoder: zapx.NewOptionX(
				zapx.WithFormat(zapx.JsonFormat)).BuildEncoder(),
			Syncer: zapcore.AddSync(os.Stdout),
		},
		{
			LevelEnablerFunc: zapx.InfoLevel.Enabled,
			Encoder: zapx.NewOptionX(
				zapx.WithFormat(zapx.JsonFormat),
				zapx.WithDevelopment(false)).BuildEncoder(),
			Syncer: zapx.NewRotateBySizeWriter(&zapx.RotateConfig{
				Filename: "./info.log",
				MaxSize:  1,
			}),
		},
		{
			LevelEnablerFunc: zapx.ErrorLevel.Enabled,
			Encoder: zapx.NewOptionX(
				zapx.WithFormat(zapx.JsonFormat),
				zapx.WithDevelopment(false)).BuildEncoder(),
			Syncer: zapx.NewRotateBySizeWriter(&zapx.RotateConfig{
				Filename: "./errors.log",
				MaxSize:  1,
			}),
		},
	}

	lg := zapx.NewTee(tees)
	defer lg.Flush()
	zapx.ReplaceDefault(lg)
	for {
		zapx.Info("hello world", zapx.String("key", "value"))
		zapx.Warn("hello world", zapx.String("key", "value"))
		zapx.Error("hello world", zapx.String("key", "value"))
	}

}
