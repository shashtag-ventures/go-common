package httputil

import (
	"log/slog"
	"net/http"
	"time"
)

// ResponseWriterWrapper captures status code and response size.
// It also implements http.Flusher for streaming support.
type ResponseWriterWrapper struct {
	http.ResponseWriter
	StatusCode int
	Size       int64
}

// WriteHeader captures the status code before sending it to the underlying ResponseWriter.
func (w *ResponseWriterWrapper) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// Write captures the number of bytes written to the response.
func (w *ResponseWriterWrapper) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.Size += int64(size)
	return size, err
}

// Flush ensures that the wrapped ResponseWriter flushes when requested.
// This is critical for SSE and streaming responses.
func (w *ResponseWriterWrapper) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// RetryRoundTripper is a custom transport that retries requests on transient errors (502/503)
type RetryRoundTripper struct {
	Base       http.RoundTripper
	MaxRetries int
	Logger     *slog.Logger
}

// RoundTrip executes a single HTTP transaction with exponential backoff retries on transient errors.
func (rt *RetryRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= rt.MaxRetries; i++ {
		if i > 0 {
			// Exponential backoff
			waitTime := time.Duration(i*800) * time.Millisecond
			if rt.Logger != nil {
				rt.Logger.Info("Retrying HTTP request due to transient error", 
					"attempt", i, 
					"url", req.URL.String(), 
					"wait_ms", waitTime.Milliseconds(),
				)
			}
			time.Sleep(waitTime)
		}

		resp, err = rt.Base.RoundTrip(req)

		// Success path (non-transient error codes)
		if err == nil && resp.StatusCode != http.StatusServiceUnavailable && resp.StatusCode != http.StatusBadGateway {
			return resp, nil
		}

		// Transient error path (502/503)
		if err == nil && (resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusBadGateway) {
			if i < rt.MaxRetries {
				resp.Body.Close()
				continue
			}
			return resp, nil
		}

		// Network error path (e.g. connection reset)
		if err != nil && i < rt.MaxRetries {
			continue
		}
		
		break
	}

	return resp, err
}
