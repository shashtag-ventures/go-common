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
	Env     string
	Service string
	Version string
	Level   slog.Level
}

// New creates a new slog.Logger instance with global metadata and sensitive data masking.
func New(cfg Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: cfg.Level,
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
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      cfg.Level,
			TimeFormat: "15:04:05",
		})
	}
	
	// Inject Global Metadata (The "Identity" of the service)
	logger := slog.New(handler).With(
		slog.String("service", cfg.Service),
		slog.String("env", cfg.Env),
	)
	
	if cfg.Version != "" {
		logger = logger.With(slog.String("version", cfg.Version))
	}
	
	return logger
}
