package config

import (
	"encoding/hex"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRequiredEnv(t *testing.T) {
	key := "TEST_REQUIRED_ENV"
	val := "test-value"

	// Test missing
	os.Unsetenv(key)
	_, err := GetRequiredEnv(key)
	assert.Error(t, err)

	// Test present
	os.Setenv(key, val)
	defer os.Unsetenv(key)
	got, err := GetRequiredEnv(key)
	assert.NoError(t, err)
	assert.Equal(t, val, got)
}

func TestGetEnvAsInt(t *testing.T) {
	key := "TEST_INT_ENV"
	
	// Test default
	os.Unsetenv(key)
	assert.Equal(t, 10, GetEnvAsInt(key, 10))

	// Test valid
	os.Setenv(key, "42")
	defer os.Unsetenv(key)
	assert.Equal(t, 42, GetEnvAsInt(key, 10))

	// Test invalid
	os.Setenv(key, "not-an-int")
	assert.Equal(t, 10, GetEnvAsInt(key, 10))
}

func TestGetEnvAsBool(t *testing.T) {
	key := "TEST_BOOL_ENV"

	// Test default
	os.Unsetenv(key)
	assert.True(t, GetEnvAsBool(key, true))

	// Test valid true
	os.Setenv(key, "true")
	assert.True(t, GetEnvAsBool(key, false))

	// Test valid false
	os.Setenv(key, "false")
	assert.False(t, GetEnvAsBool(key, true))

	// Test invalid
	os.Setenv(key, "maybe")
	assert.True(t, GetEnvAsBool(key, true))
}

func TestGetEnvAsSlogLevel(t *testing.T) {
	key := "TEST_LOG_LEVEL"
	defaultLevel := slog.LevelInfo

	// Test default
	os.Unsetenv(key)
	assert.Equal(t, defaultLevel, GetEnvAsSlogLevel(key, defaultLevel))

	// Test valid (Debug = -4)
	os.Setenv(key, "-4")
	assert.Equal(t, slog.LevelDebug, GetEnvAsSlogLevel(key, defaultLevel))

	// Test invalid
	os.Setenv(key, "not-a-level")
	assert.Equal(t, defaultLevel, GetEnvAsSlogLevel(key, defaultLevel))
}

func TestGetEnvAsHexBytes(t *testing.T) {
	key := "TEST_HEX_ENV"

	// 1. Test missing
	os.Unsetenv(key)
	_, err := GetEnvAsHexBytes(key)
	assert.Error(t, err)

	// 2. Test valid hex (64 chars / 32 bytes)
	rawBytes := make([]byte, 32)
	for i := range rawBytes {
		rawBytes[i] = byte(i)
	}
	hexVal := hex.EncodeToString(rawBytes)
	os.Setenv(key, hexVal)
	
	got, err := GetEnvAsHexBytes(key)
	assert.NoError(t, err)
	assert.Equal(t, rawBytes, got)

	// 3. Test raw string (not hex, or odd length)
	rawStr := "plain-text-key"
	os.Setenv(key, rawStr)
	got, err = GetEnvAsHexBytes(key)
	assert.NoError(t, err)
	assert.Equal(t, []byte(rawStr), got)

	// 4. Test even length but invalid hex characters
	invalidHex := "zzzzzzzz" // even length but not hex
	os.Setenv(key, invalidHex)
	got, err = GetEnvAsHexBytes(key)
	assert.NoError(t, err)
	assert.Equal(t, []byte(invalidHex), got)
	
	os.Unsetenv(key)
}
