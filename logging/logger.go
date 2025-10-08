package logging

import (
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
)

// Config holds the configuration for the logger.
type Config struct {
	Env   string
	Level slog.Level
}

// New creates a new slog.Logger instance.
func New(cfg Config) *slog.Logger {
	var handler slog.Handler
	if strings.ToLower(cfg.Env) == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: cfg.Level,
		})
	} else {
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level: cfg.Level,
		})
	}
	return slog.New(handler)
}
