package tracing_test

import (
	"testing"

	"github.com/shashtag-ventures/go-common/tracing"
	"github.com/stretchr/testify/assert"
)

func TestInitTracer(t *testing.T) {
	t.Run("Initialize and cleanup", func(t *testing.T) {
		cleanup := tracing.InitTracer("test-service")
		assert.NotNil(t, cleanup)
		
		// Ensure cleanup doesn't panic
		assert.NotPanics(t, func() {
			cleanup()
		})
	})
}
