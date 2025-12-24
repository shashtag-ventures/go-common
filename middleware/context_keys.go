package middleware

import (
	"context"
	"log/slog"
	"sync"
)

// CtxKey is a custom type for context keys to avoid collisions.
type CtxKey string

const (
	RequestIDKey    CtxKey = "requestID"
	LoggerContextKey CtxKey = "logger"
	UserContextKey   CtxKey = "user"
	LogStateKey      CtxKey = "logState"
)

// LogState holds mutable metadata gathered during the request lifecycle.
type LogState struct {
	mu          sync.RWMutex
	Fields      map[string]any
	Breadcrumbs []string // NEW: Track the "story" of the request
	UserID      uint
}

func (s *LogState) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.Fields == nil {
		s.Fields = make(map[string]any)
	}
	s.Fields[key] = value
}

func (s *LogState) AddBreadcrumb(msg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Breadcrumbs = append(s.Breadcrumbs, msg)
}

func (s *LogState) SetUser(id uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.UserID = id
}

// GetLoggerFromContext retrieves the logger from the context.
func GetLoggerFromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(LoggerContextKey).(*slog.Logger); ok {
		return logger
	}
	return slog.Default()
}

// EnrichLogger adds attributes to the logger AND the mutable log state.
func EnrichLogger(ctx context.Context, attrs ...slog.Attr) context.Context {
	logger := GetLoggerFromContext(ctx)
	for _, attr := range attrs {
		logger = logger.With(attr)
	}
	if state, ok := ctx.Value(LogStateKey).(*LogState); ok {
		for _, attr := range attrs {
			state.Set(attr.Key, attr.Value.Any())
		}
	}
	return context.WithValue(ctx, LoggerContextKey, logger)
}

// AddBreadcrumb logs a milestone in the current request flow.
func AddBreadcrumb(ctx context.Context, msg string) {
	if state, ok := ctx.Value(LogStateKey).(*LogState); ok {
		state.AddBreadcrumb(msg)
	}
}

func GetRequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
