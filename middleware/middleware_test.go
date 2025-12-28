package middleware_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shashtag-ventures/go-common/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCorsMiddleware(t *testing.T) {
	cfg := middleware.CorsConfig{
		AllowedOrigins: []string{"http://example.com"},
		AllowedMethods: []string{"GET", "POST"},
	}
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	t.Run("Allowed Origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		req.Header.Set("Origin", "http://example.com")
		req.Header.Set("Access-Control-Request-Method", "GET")
		rr := httptest.NewRecorder()

		middleware.CorsMiddleware(cfg, nextHandler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Equal(t, "http://example.com", rr.Header().Get("Access-Control-Allow-Origin"))
	})

	t.Run("Disallowed Origin", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodOptions, "/", nil)
		req.Header.Set("Origin", "http://malicious.com")
		req.Header.Set("Access-Control-Request-Method", "GET")
		rr := httptest.NewRecorder()

		middleware.CorsMiddleware(cfg, nextHandler).ServeHTTP(rr, req)

		assert.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"))
	})
}

func TestRecovery(t *testing.T) {
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	t.Run("Recover from Panic", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		middleware.Recovery()(panicHandler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), http.StatusText(http.StatusInternalServerError))
	})

	t.Run("Recover from Panic with RequestID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-req-id")
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		middleware.Recovery()(panicHandler).ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestRateLimitMiddleware(t *testing.T) {
	cfg := middleware.RateLimitConfig{
		Enabled: true,
		Limit:   2,
		Window:  1,
	}
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.RateLimitMiddleware(cfg)(nextHandler)

	t.Run("Under Limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "1.2.3.4:1234"
		
		rr1 := httptest.NewRecorder()
		handler.ServeHTTP(rr1, req)
		assert.Equal(t, http.StatusOK, rr1.Code)

		rr2 := httptest.NewRecorder()
		handler.ServeHTTP(rr2, req)
		assert.Equal(t, http.StatusOK, rr2.Code)
	})

	t.Run("Over Limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "5.6.7.8:1234"
		
		handler.ServeHTTP(httptest.NewRecorder(), req) // 1
		handler.ServeHTTP(httptest.NewRecorder(), req) // 2
		
		rr3 := httptest.NewRecorder()
		handler.ServeHTTP(rr3, req) // 3 - Should fail
		assert.Equal(t, http.StatusTooManyRequests, rr3.Code)
	})

	t.Run("Disabled", func(t *testing.T) {
		disabledCfg := middleware.RateLimitConfig{Enabled: false}
		disabledHandler := middleware.RateLimitMiddleware(disabledCfg)(nextHandler)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()
		
		disabledHandler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestRequestLogger(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/error" {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	handler := middleware.RequestLogger()(nextHandler)

	t.Run("Health Check - Skip Logging", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Standard Request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Request with Body and Scrubbing", func(t *testing.T) {
		body := map[string]string{"password": "secret-password", "name": "john"}
		bodyBytes, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	t.Run("Error Request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})

	t.Run("Request with Breadcrumbs and Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/error", nil)
		state := &middleware.LogState{Breadcrumbs: []string{"start", "middle"}}
		ctx := context.WithValue(req.Context(), middleware.LogStateKey, state)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()
		
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}

func TestMetricsMiddleware(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("measured"))
	})

	handler := middleware.MetricsMiddleware(nextHandler)

	t.Run("Measure Request", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/metrics-test", nil)
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})
}

func TestLogState_EnrichLogger(t *testing.T) {
	t.Run("Enrich with state and attrs", func(t *testing.T) {
		state := &middleware.LogState{}
		ctx := context.WithValue(context.Background(), middleware.LogStateKey, state)
		
		ctx = middleware.EnrichLogger(ctx, slog.String("test-key", "test-val"))
		
		assert.Equal(t, "test-val", state.Fields["test-key"])
	})
}

func TestAddBreadcrumb(t *testing.T) {
	t.Run("Add breadcrumb", func(t *testing.T) {
		state := &middleware.LogState{}
		ctx := context.WithValue(context.Background(), middleware.LogStateKey, state)
		
		middleware.AddBreadcrumb(ctx, "step 1")
		
		assert.Contains(t, state.Breadcrumbs, "step 1")
	})
}
