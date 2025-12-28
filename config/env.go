package config

import (
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Load loads the .env file if it exists.
func Load() {
	_ = godotenv.Load()
}

// GetRequiredEnv retrieves an environment variable or returns an error if not set.
func GetRequiredEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("environment variable %s is not set", key)
	}
	return val, nil
}

// GetEnvAsInt retrieves an environment variable as an integer or returns a default value.
func GetEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// GetEnvAsBool retrieves an environment variable as a boolean or returns a default value.
func GetEnvAsBool(name string, defaultVal bool) bool {
	valueStr := os.Getenv(name)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultVal
}

// GetEnvAsSlogLevel retrieves an environment variable as an slog.Level or returns a default value.
func GetEnvAsSlogLevel(name string, defaultVal slog.Level) slog.Level {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return slog.Level(value)
	}
	return defaultVal
}

// GetEnvAsHexBytes retrieves an environment variable as hex-decoded bytes.
// If the variable is not hex-encoded but matches the target length, it returns the raw bytes.
func GetEnvAsHexBytes(name string) ([]byte, error) {
	val, err := GetRequiredEnv(name)
	if err != nil {
		return nil, err
	}

	// If it looks like a hex string (even length and common key lengths like 32, 64)
	if len(val) % 2 == 0 {
		keyBytes, err := hex.DecodeString(val)
		if err == nil {
			return keyBytes, nil
		}
	}

	return []byte(val), nil
}
