package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

const (
	LevelDebug = "debug"
	LevelWarn  = "warn"
	LevelInfo  = "info"
	LevelError = "error"
)

func init() {
	logger = getLogger(zapcore.DebugLevel)
}

func getLogger(logLevel zapcore.Level) *zap.Logger {
	var err error
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Level = zap.NewAtomicLevelAt(logLevel)
	cfg.Encoding = "console"
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	logger, err = cfg.Build()
	if err != nil {
		panic(fmt.Sprintf("Can`t init logger: %s", err.Error()))
	}
	return logger
}

func Info(text string, fields ...zap.Field) {
	logger.Info(text, fields...)
}

func Error(text string, fields ...zap.Field) {
	logger.Error(text, fields...)
}

func Warn(text string, fields ...zap.Field) {
	logger.Warn(text, fields...)
}

func Fatal(text string, fields ...zap.Field) {
	logger.Fatal(text, fields...)
}
