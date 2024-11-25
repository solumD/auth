package config

import (
	"errors"
	"os"
)

const (
	loggerLevelEnvName = "LOGGER_LEVEL"
)

type loggerConfig struct {
	level string
}

// NewLoggerConfig returns new grpc config
func NewLoggerConfig() (LoggerConfig, error) {
	level := os.Getenv(loggerLevelEnvName)
	if len(level) == 0 {
		return nil, errors.New("logger level not found")
	}

	return &loggerConfig{
		level: level,
	}, nil
}

// Level returns level of logger
func (cfg *loggerConfig) Level() string {
	return cfg.level
}
