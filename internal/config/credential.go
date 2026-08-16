package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func ResolveCredential(ref string) (string, error) {
	if ref == "" {
		return "", nil
	}
	switch {
	case strings.HasPrefix(ref, "env:"):
		name := strings.TrimPrefix(ref, "env:")
		value, ok := os.LookupEnv(name)
		if !ok || value == "" {
			return "", fmt.Errorf("credential environment variable %q is unset", name)
		}
		return value, nil
	case strings.HasPrefix(ref, "file:"):
		path := strings.TrimPrefix(ref, "file:")
		info, err := os.Stat(path)
		if err != nil {
			return "", fmt.Errorf("credential file: %w", err)
		}
		if info.Mode().Perm()&0o077 != 0 {
			return "", errors.New("credential file is accessible by group or other users")
		}
		value, err := os.ReadFile(path) // #nosec G304 -- the credential file path is an explicit local profile reference.
		if err != nil {
			return "", fmt.Errorf("read credential file: %w", err)
		}
		return strings.TrimSpace(string(value)), nil
	default:
		return "", errors.New("unsupported credential reference")
	}
}
