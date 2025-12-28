package middleware

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"net/http"
)

// responseBuffer is a custom ResponseWriter that captures the response body.
type responseBuffer struct {
	http.ResponseWriter
	body   *bytes.Buffer
	status int
}

func (w *responseBuffer) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *responseBuffer) WriteHeader(statusCode int) {
	w.status = statusCode
}

// ETag is a middleware that generates an ETag for the response and handles conditional requests.
// It buffers the response to calculate the hash, so use with caution for large responses.
func ETag(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Only calculate ETag for GET and HEAD requests
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			next.ServeHTTP(w, r)
			return
		}

		buf := &bytes.Buffer{}
		rb := &responseBuffer{
			ResponseWriter: w,
			body:           buf,
			status:         http.StatusOK, // Default status
		}

		next.ServeHTTP(rb, r)

		// If the handler returned a non-success code (or 304/204), pass through
		// We only ETag 200 OK responses usually
		if rb.status != http.StatusOK {
			w.WriteHeader(rb.status)
			w.Write(buf.Bytes())
			return
		}

		// Calculate ETag
		bodyBytes := buf.Bytes()
		etag := fmt.Sprintf("\"%x\"", sha256.Sum256(bodyBytes))

		// Check If-None-Match
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}

		// Set ETag header
		w.Header().Set("ETag", etag)
		// We add "no-cache" to force revalidation with the server using the ETag
		w.Header().Set("Cache-Control", "public, no-cache")

		// Write the buffered response
		w.WriteHeader(rb.status)
		w.Write(bodyBytes)
	})
}
