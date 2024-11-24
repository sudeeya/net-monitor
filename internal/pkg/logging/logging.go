// Package logging provides a way to create logger.
package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Possible log levels.
const (
	Info  = "INFO"
	Error = "ERROR"
	Fatal = "FATAL"
)

// NewLogger returns a logger.
// The function sets the log level and adds a file for writing logs.
// If [logFile] is an empty string, logs will be output to standard out only.
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

	if logFile != "" {
		cfg.OutputPaths = append(cfg.OutputPaths, logFile)
	}

	return cfg.Build()
}
