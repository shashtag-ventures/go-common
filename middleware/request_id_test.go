package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shashtag-ventures/go-common/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRequestIDMiddleware(t *testing.T) {
	// Create a simple handler that checks for the request ID in the context.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := middleware.GetRequestIDFromContext(r.Context())
		assert.NotEmpty(t, requestID, "Request ID should be in the context")

		logger := middleware.GetLoggerFromContext(r.Context())
		// This is harder to test without inspecting logs, but we can check it's not nil
		assert.NotNil(t, logger)

		w.WriteHeader(http.StatusOK)
	})

	// Wrap the handler with the middleware.
	testHandler := middleware.RequestIDMiddleware(handler)

	// Create a test request and response recorder.
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Serve the HTTP request.
	testHandler.ServeHTTP(rr, req)

	// Assert that the response status is OK.
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert that the X-Request-ID header is present in the response.
	assert.NotEmpty(t, rr.Header().Get("X-Request-ID"), "X-Request-ID header should be set")
}
