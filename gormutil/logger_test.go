package gormutil_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/shashtag-ventures/go-common/gormutil"
	"github.com/shashtag-ventures/go-common/middleware"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGormLogger(t *testing.T) {
	var buf bytes.Buffer
	h := slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := &gormutil.GormLogger{
		SlowThreshold: 100 * time.Millisecond,
	}
	
	// Create a context with the logger
	ctx := context.WithValue(context.Background(), middleware.LoggerContextKey, slog.New(h))

	t.Run("LogMode", func(t *testing.T) {
		assert.Equal(t, l, l.LogMode(0))
	})

	t.Run("Info", func(t *testing.T) {
		buf.Reset()
		l.Info(ctx, "test info", "foo", "bar")
		assert.Contains(t, buf.String(), "test info")
	})

	t.Run("Warn", func(t *testing.T) {
		buf.Reset()
		l.Warn(ctx, "test warn", "foo", "bar")
		assert.Contains(t, buf.String(), "test warn")
	})

	t.Run("Error", func(t *testing.T) {
		buf.Reset()
		l.Error(ctx, "test error", "foo", "bar")
		assert.Contains(t, buf.String(), "test error")
	})

	t.Run("Trace - Normal", func(t *testing.T) {
		buf.Reset()
		l.Trace(ctx, time.Now(), func() (string, int64) {
			return "SELECT 1", 1
		}, nil)
		assert.Contains(t, buf.String(), "SQL QUERY")
	})

	t.Run("Trace - Error", func(t *testing.T) {
		buf.Reset()
		l.Trace(ctx, time.Now(), func() (string, int64) {
			return "SELECT 1", 0
		}, errors.New("db error"))
		assert.Contains(t, buf.String(), "SQL ERROR")
		assert.Contains(t, buf.String(), "db error")
	})

	t.Run("Trace - Slow", func(t *testing.T) {
		buf.Reset()
		l.Trace(ctx, time.Now().Add(-200*time.Millisecond), func() (string, int64) {
			return "SELECT pg_sleep(1)", 1
		}, nil)
		assert.Contains(t, buf.String(), "SLOW SQL")
	})

	t.Run("Trace - RecordNotFound (No Error Log)", func(t *testing.T) {
		buf.Reset()
		l.Trace(ctx, time.Now(), func() (string, int64) {
			return "SELECT * FROM users WHERE id = 999", 0
		}, gorm.ErrRecordNotFound)
		assert.NotContains(t, buf.String(), "SQL ERROR")
	})
}
