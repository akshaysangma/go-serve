// Package logging contains logic to initialize zap Logger
package logging

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(levelStr, encoding string) (*zap.Logger, error) {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid log level '%s', defaulting to 'info'. Error: %v\n", levelStr, err)
		level = zap.InfoLevel
	}

	var config zap.Config
	switch encoding {
	case "json":
		config = zap.NewProductionConfig()
	case "console":
		config = zap.NewDevelopmentConfig()
		// Use colored levels in console output for easier visual scanning.
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config = zap.NewDevelopmentConfig()
	}
	config.Level = zap.NewAtomicLevelAt(level)
	config.Encoding = encoding
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true

	// AddCallerSkip(1) ensures the correct caller (e.g., handler function) is shown in logs, not this InitLogger function.
	logger, err := config.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	return logger, nil
}
