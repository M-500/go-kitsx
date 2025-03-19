package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func main() {
	// 创建不同的文件输出目标
	debugFile, err := os.Create("debug.log")
	if err != nil {
		panic(err)
	}
	defer debugFile.Close()

	infoFile, err := os.Create("info.log")
	if err != nil {
		panic(err)
	}
	defer infoFile.Close()

	warnFile, err := os.Create("warn.log")
	if err != nil {
		panic(err)
	}
	defer warnFile.Close()

	errorFile, err := os.Create("error.log")
	if err != nil {
		panic(err)
	}
	defer errorFile.Close()

	// 创建编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()

	// 创建不同级别的 Core
	debugCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // JSON 格式的编码器
		zapcore.AddSync(debugFile),            // 输出到 debug.log 文件
		zapcore.DebugLevel,                    // 仅输出 Debug 级别及以上日志
	)

	infoCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // JSON 格式的编码器
		zapcore.AddSync(infoFile),             // 输出到 info.log 文件
		zapcore.InfoLevel,                     // 仅输出 Info 级别及以上日志
	)

	warnCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // JSON 格式的编码器
		zapcore.AddSync(warnFile),             // 输出到 warn.log 文件
		zapcore.WarnLevel,                     // 仅输出 Warn 级别及以上日志
	)

	errorCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // JSON 格式的编码器
		zapcore.AddSync(errorFile),            // 输出到 error.log 文件
		zapcore.ErrorLevel,                    // 仅输出 Error 级别及以上日志
	)

	// 将多个 Core 合并
	logger := zap.New(zapcore.NewTee(debugCore, infoCore, warnCore, errorCore))

	// 使用 logger 输出不同级别的日志
	logger.Debug("This is a debug message")  // 只会写入 debug.log
	logger.Info("This is an info message")   // 只会写入 info.log
	logger.Warn("This is a warn message")    // 只会写入 warn.log
	logger.Error("This is an error message") // 只会写入 error.log

	// 注意：不同的日志级别不会输出低于当前配置级别的日志

	zap.NewDevelopment()
}
