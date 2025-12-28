package worker_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/shashtag-ventures/go-common/worker"
	"github.com/stretchr/testify/assert"
)

func TestSafeGo(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		
		worker.SafeGo(context.Background(), "success-test", func(ctx context.Context) error {
			defer wg.Done()
			return nil
		})
		
		wg.Wait()
	})

	t.Run("Error", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		
		worker.SafeGo(context.Background(), "error-test", func(ctx context.Context) error {
			defer wg.Done()
			return errors.New("failed")
		})
		
		wg.Wait()
	})

	t.Run("Panic Recovery", func(t *testing.T) {
		var wg sync.WaitGroup
		wg.Add(1)
		
		worker.SafeGo(context.Background(), "panic-test", func(ctx context.Context) error {
			defer wg.Done()
			panic("boom")
		})
		
		wg.Wait()
	})
}

func TestRunSync(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		err := worker.RunSync(ctx, "sync-success", func(ctx context.Context) error {
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("Error", func(t *testing.T) {
		err := worker.RunSync(ctx, "sync-error", func(ctx context.Context) error {
			return errors.New("failed")
		})
		assert.Error(t, err)
	})

	t.Run("Panic", func(t *testing.T) {
		// RunSync will return nil if it panics because we recover and don't return the panic as error (optional design choice)
		// but let's just ensure it doesn't crash the test runner.
		_ = worker.RunSync(ctx, "sync-panic", func(ctx context.Context) error {
			panic("boom")
		})
	})
}
