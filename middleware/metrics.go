package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// httpRequestsTotal is a Prometheus counter that tracks the total number of HTTP requests.
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)
	// httpRequestDuration is a Prometheus histogram that tracks the duration of HTTP requests.
	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
	// httpResponseSize is a Prometheus histogram that tracks the size of HTTP responses.
	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Size of HTTP responses.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path", "status"},
	)
)

// MetricsMiddleware collects HTTP request metrics (total requests, duration, response size).
// It wraps the next http.Handler in the chain.
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now() // Record the start time of the request.

		// Wrap the ResponseWriter to capture status code and response size.
		lw := &loggingResponseWriter{ResponseWriter: w}
		next.ServeHTTP(lw, r) // Serve the actual request.

		duration := time.Since(start).Seconds() // Calculate request duration.
		status := strconv.Itoa(lw.statusCode)   // Get the HTTP status code as a string.

		// Increment the total request counter.
		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		// Observe the request duration.
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
		// Observe the response size.
		httpResponseSize.WithLabelValues(r.Method, r.URL.Path, status).Observe(float64(lw.size))
	})
}

// loggingResponseWriter is a wrapper around http.ResponseWriter to capture the status code and response size.
type loggingResponseWriter struct {
	http.ResponseWriter     // Embedded to satisfy the http.ResponseWriter interface.
	statusCode          int // Stores the HTTP status code written.
	size                int // Stores the size of the response body written.
}

// WriteHeader captures the status code before calling the underlying WriteHeader.
func (lw *loggingResponseWriter) WriteHeader(code int) {
	lw.statusCode = code
	lw.ResponseWriter.WriteHeader(code)
}

// Write captures the size of the response body before calling the underlying Write.
func (lw *loggingResponseWriter) Write(data []byte) (int, error) {
	size, err := lw.ResponseWriter.Write(data)
	lw.size += size
	return size, err
}
