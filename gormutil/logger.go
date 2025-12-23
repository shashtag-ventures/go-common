package gormutil

import (
	"context"
	"errors"
	"time"

	"github.com/shashtag-ventures/go-common/middleware"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormLogger is a custom GORM logger that uses the contextual slog logger.
type GormLogger struct {
	SlowThreshold time.Duration
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	middleware.GetLoggerFromContext(ctx).Info(msg, "data", data)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	middleware.GetLoggerFromContext(ctx).Warn(msg, "data", data)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	middleware.GetLoggerFromContext(ctx).Error(msg, "data", data)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	log := middleware.GetLoggerFromContext(ctx).With(
		"elapsed", elapsed.String(),
		"rows", rows,
	)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error("SQL ERROR", "sql", sql, "error", err)
		return
	}

	if l.SlowThreshold > 0 && elapsed > l.SlowThreshold {
		log.Warn("SLOW SQL", "sql", sql)
		return
	}

	// For normal queries, we can log at Debug level so they don't spam Production
	log.Debug("SQL QUERY", "sql", sql)
}
