package worker

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/shashtag-ventures/go-common/middleware"
)

// Task is a function that represents a background job.
type Task func(ctx context.Context) error

// SafeGo executes a task in a new goroutine with panic recovery and standard logging.
// It uses the provided context for the task.
func SafeGo(ctx context.Context, name string, task Task) {
	go func() {
		// 1. Setup metadata for logging
		start := time.Now()
		logger := middleware.GetLoggerFromContext(ctx).With("task", name)
		
		// 2. Panic Recovery
		defer func() {
			if r := recover(); r != nil {
				logger.Error("TASK PANIC RECOVERED",
					"error", r,
					"duration_ms", time.Since(start).Milliseconds(),
					"stack", string(debug.Stack()),
				)
			}
		}()

		logger.Info("Task started")

		// 3. Execute Task
		err := task(ctx)

		// 4. Final Logging
		duration := time.Since(start)
		if err != nil {
			logger.Error("Task failed", 
				"error", err, 
				"duration_ms", duration.Milliseconds(),
			)
		} else {
			logger.Info("Task completed successfully", 
				"duration_ms", duration.Milliseconds(),
			)
		}
	}()
}

// RunSync executes a task synchronously with panic recovery.
// Useful for one-off CLI tasks or migrations.
func RunSync(ctx context.Context, name string, task Task) error {
	start := time.Now()
	logger := middleware.GetLoggerFromContext(ctx).With("task", name)

	defer func() {
		if r := recover(); r != nil {
			logger.Error("TASK PANIC RECOVERED",
				"error", r,
				"duration_ms", time.Since(start).Milliseconds(),
				"stack", string(debug.Stack()),
			)
		}
	}()

	logger.Info("Task started")
	err := task(ctx)
	
	duration := time.Since(start)
	if err != nil {
		logger.Error("Task failed", "error", err, "duration_ms", duration.Milliseconds())
		return fmt.Errorf("task %s failed: %w", name, err)
	}

	logger.Info("Task completed successfully", "duration_ms", duration.Milliseconds())
	return nil
}
