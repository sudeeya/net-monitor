package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	Info  = "INFO"
	Error = "ERROR"
	Fatal = "FATAL"
)

func NewLogger(logLevel, logFile string) (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()

	switch logLevel {
	case Info:
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case Error:
		cfg.Level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case Fatal:
		cfg.Level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		return nil, fmt.Errorf("unknown log level: %s", logLevel)
	}

	cfg.OutputPaths = append(cfg.OutputPaths, logFile)

	return cfg.Build()
}
