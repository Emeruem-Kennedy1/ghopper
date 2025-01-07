package logging

import (
    "go.uber.org/zap"
    "go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.Logger, error) {
    config := zap.NewProductionConfig()

    config.EncoderConfig.TimeKey = "timestamp"
    config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    config.OutputPaths = []string{"stdout", "app.log"}
    config.EncoderConfig.StacktraceKey = "stacktrace"
    config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

    return config.Build()
}