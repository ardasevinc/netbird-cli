package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	DefaultTimeout  = 20 * time.Second
	MinTimeout      = 1 * time.Second
	MaxTimeout      = 5 * time.Minute
	DefaultLogLevel = "error"
)

func ParseTimeout(value string) (time.Duration, error) {
	timeout, err := time.ParseDuration(strings.TrimSpace(value))
	if err != nil {
		return 0, fmt.Errorf("timeout must be a duration between %s and %s: %w", MinTimeout, MaxTimeout, err)
	}
	if timeout < MinTimeout || timeout > MaxTimeout {
		return 0, fmt.Errorf("timeout must be between %s and %s", MinTimeout, MaxTimeout)
	}
	return timeout, nil
}

func TimeoutFromEnvironment() (time.Duration, bool, error) {
	value := os.Getenv("NB_TIMEOUT")
	if value == "" {
		return DefaultTimeout, false, nil
	}
	timeout, err := ParseTimeout(value)
	if err != nil {
		return 0, true, fmt.Errorf("NB_TIMEOUT: %w", err)
	}
	return timeout, true, nil
}

func ParseLogLevel(value string) (string, error) {
	level := strings.ToLower(strings.TrimSpace(value))
	switch level {
	case "debug", "info", "warn", "error":
		return level, nil
	default:
		return "", errors.New("log level must be one of debug, info, warn, or error")
	}
}

func LogLevelFromEnvironment() (string, bool, error) {
	value := os.Getenv("NB_LOG_LEVEL")
	if value == "" {
		return DefaultLogLevel, false, nil
	}
	level, err := ParseLogLevel(value)
	if err != nil {
		return "", true, fmt.Errorf("NB_LOG_LEVEL: %w", err)
	}
	return level, true, nil
}
