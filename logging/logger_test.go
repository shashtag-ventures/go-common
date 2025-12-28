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

	t.Run("Sensitive key masking", func(t *testing.T) {
		cfg := logging.Config{
			Level: slog.LevelInfo,
			Env:   "production", // Use production to get JSON and verify masking easily
		}

		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		logger := logging.New(cfg)
		logger.Info("Login attempt", "password", "123456", "token", "secret-token")

		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = oldStdout

		assert.Contains(t, string(out), "[MASKED]")
		assert.NotContains(t, string(out), "123456")
		assert.NotContains(t, string(out), "secret-token")
	})

	t.Run("Production JSON Handler", func(t *testing.T) {
		cfg := logging.Config{
			Env: "prod",
		}
		logger := logging.New(cfg)
		assert.NotNil(t, logger)
		// Smoke test for handler type indirectly
	})

	t.Run("Version metadata", func(t *testing.T) {
		cfg := logging.Config{
			Version: "1.2.3",
			Env:     "prod",
		}
		
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		logger := logging.New(cfg)
		logger.Info("ver test")

		w.Close()
		out, _ := io.ReadAll(r)
		os.Stdout = oldStdout
		
		assert.Contains(t, string(out), "1.2.3")
		assert.Contains(t, string(out), "version")
	})
}