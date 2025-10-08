package logging_test

import (
	"io"
	"log/slog"
	"os"
	"testing"

	"github.com/shashtag-ventures/go-common/logging"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("Logger is created with INFO level", func(t *testing.T) {
		cfg := logging.Config{
			Level: slog.LevelInfo,
		}

		// Capture stdout to verify log output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		logger := logging.New(cfg)
		assert.NotNil(t, logger)

		// Log a message at INFO level
		logger.Info("Test info message")

		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = oldStdout // Restore stdout

		// Verify that the INFO message is present in the output
		assert.Contains(t, string(out), "INF")
		assert.Contains(t, string(out), "Test info message")
	})

	t.Run("Logger is created with DEBUG level", func(t *testing.T) {
		cfg := logging.Config{
			Level: slog.LevelDebug,
		}

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		logger := logging.New(cfg)
		assert.NotNil(t, logger)

		// Log a message at DEBUG level
		logger.Debug("Test debug message")
		// Log a message at INFO level (should also appear)
		logger.Info("Test info message at debug level")

		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = oldStdout

		// Verify that the DEBUG message is present in the output
		assert.Contains(t, string(out), "DBG")
		assert.Contains(t, string(out), "Test debug message")
		assert.Contains(t, string(out), "INF")
		assert.Contains(t, string(out), "Test info message at debug level")
	})

	t.Run("Logger is created with ERROR level", func(t *testing.T) {
		cfg := logging.Config{
			Level: slog.LevelError,
		}

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		logger := logging.New(cfg)
		assert.NotNil(t, logger)

		// Log a message at INFO level (should NOT appear)
		logger.Info("Test info message at error level")
		// Log a message at ERROR level (should appear)
		logger.Error("Test error message")

		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = oldStdout

		// Verify that the ERROR message is present and INFO is not
		assert.NotContains(t, string(out), "level=INFO")
		assert.NotContains(t, string(out), "msg=\"Test info message at error level\"")
		assert.Contains(t, string(out), "ERR")
		assert.Contains(t, string(out), "Test error message")
	})
}
