package zapx

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"strings"
	"time"
)

type RotateConfig struct {
	// 共用配置
	Filename string // 完整文件名
	MaxAge   int    // 保留旧日志文件的最大天数

	// 按时间轮转配置
	RotationTime time.Duration // 日志文件轮转时间

	// 按大小轮转配置
	MaxSize    int  // 日志文件最大大小（MB）
	MaxBackups int  // 保留日志文件的最大数量
	Compress   bool // 是否对日志文件进行压缩归档
	LocalTime  bool // 是否使用本地时间，默认 UTC 时间
}

func newRotateCfg() *RotateConfig {
	return &RotateConfig{
		Filename:     "info.log",
		MaxAge:       30, // 日志保留 30 天
		RotationTime: 24 * time.Hour,
		MaxSize:      100, // 100M
		MaxBackups:   1000,
		Compress:     true,
		LocalTime:    true,
	}
}

func NewRotateBySizeWriter(cfg *RotateConfig) zapcore.WriteSyncer {
	if cfg == nil {
		cfg = newRotateCfg()
	}
	lumberJackLogger := &lumberjack.Logger{
		Filename:   cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  cfg.LocalTime,
		Compress:   cfg.Compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func NewRotateByTime(cfg *RotateConfig) zapcore.WriteSyncer {
	opts := []rotatelogs.Option{
		rotatelogs.WithMaxAge(time.Duration(cfg.MaxAge) * time.Hour * 24),
		rotatelogs.WithRotationTime(cfg.RotationTime),
		rotatelogs.WithLinkName(cfg.Filename),
	}
	if !cfg.LocalTime {
		rotatelogs.WithClock(rotatelogs.UTC)
	}
	filename := strings.SplitN(cfg.Filename, ".", 2)
	l, _ := rotatelogs.New(
		filename[0]+".%Y-%m-%d-%H-%M-%S."+filename[1],
		opts...,
	)
	return zapcore.AddSync(l)
}
