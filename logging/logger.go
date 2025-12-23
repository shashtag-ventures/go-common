package logging

import (
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/lmittmann/tint"
)

// Config holds the configuration for the logger.
type Config struct {
	Env   string
	Level slog.Level
}

// New creates a new slog.Logger instance with environment-aware formatting and sensitive data masking.
func New(cfg Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: false, // Cleaner logs: Removed file and line numbers
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// 1. Mask sensitive keys
			key := strings.ToLower(a.Key)
			if key == "password" || key == "token" || key == "access_token" || key == "refresh_token" || key == "secret" {
				return slog.String(a.Key, "[MASKED]")
			}
			
			// 2. Consistent timestamp format
			if a.Key == slog.TimeKey {
				return slog.String(a.Key, a.Value.Time().Format(time.RFC3339))
			}
			
			return a
		},
	}

	var handler slog.Handler
	if strings.ToLower(cfg.Env) == "production" || strings.ToLower(cfg.Env) == "prod" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// Beautiful colored output for local development
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      cfg.Level,
			TimeFormat: "15:04:05",
		})
	}
	
	return slog.New(handler)
}