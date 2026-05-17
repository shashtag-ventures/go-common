package httputil

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockRoundTripper struct {
	responses []*http.Response
	errs      []error
	attempts  int
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	idx := m.attempts
	m.attempts++
	if idx < len(m.errs) && m.errs[idx] != nil {
		return nil, m.errs[idx]
	}
	
	var resp *http.Response
	if idx < len(m.responses) && m.responses[idx] != nil {
		resp = m.responses[idx]
	} else {
		resp = &http.Response{StatusCode: http.StatusOK}
	}

	if resp.Body == nil {
		resp.Body = io.NopCloser(strings.NewReader(""))
	}
	return resp, nil
}

func TestRetryRoundTripper(t *testing.T) {
	t.Run("Success on first attempt", func(t *testing.T) {
		mock := &mockRoundTripper{
			responses: []*http.Response{{StatusCode: http.StatusOK}},
		}
		rt := &RetryRoundTripper{
			Base:       mock,
			MaxRetries: 3,
		}

		req := httptest.NewRequest("GET", "http://example.com", nil)
		resp, err := rt.RoundTrip(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 1, mock.attempts)
	})

	t.Run("Retry on 503 and eventual success", func(t *testing.T) {
		mock := &mockRoundTripper{
			responses: []*http.Response{
				{StatusCode: http.StatusServiceUnavailable},
				{StatusCode: http.StatusServiceUnavailable},
				{StatusCode: http.StatusOK},
			},
		}
		rt := &RetryRoundTripper{
			Base:       mock,
			MaxRetries: 3,
			Backoff:    func(i int) time.Duration { return 0 },
		}

		req := httptest.NewRequest("GET", "http://example.com", nil)
		resp, err := rt.RoundTrip(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 3, mock.attempts)
	})

	t.Run("Max retries exceeded on 503", func(t *testing.T) {
		mock := &mockRoundTripper{
			responses: []*http.Response{
				{StatusCode: http.StatusServiceUnavailable},
				{StatusCode: http.StatusServiceUnavailable},
				{StatusCode: http.StatusServiceUnavailable},
				{StatusCode: http.StatusServiceUnavailable},
			},
		}
		rt := &RetryRoundTripper{
			Base:       mock,
			MaxRetries: 2,
			Backoff:    func(i int) time.Duration { return 0 },
		}

		req := httptest.NewRequest("GET", "http://example.com", nil)
		resp, err := rt.RoundTrip(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)
		assert.Equal(t, 3, mock.attempts) // 1 initial + 2 retries
	})

	t.Run("Retry on network error", func(t *testing.T) {
		mock := &mockRoundTripper{
			errs: []error{errors.New("network error"), nil},
			responses: []*http.Response{
				nil,
				{StatusCode: http.StatusOK},
			},
		}
		rt := &RetryRoundTripper{
			Base:       mock,
			MaxRetries: 3,
			Backoff:    func(i int) time.Duration { return 0 },
		}

		req := httptest.NewRequest("GET", "http://example.com", nil)
		resp, err := rt.RoundTrip(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, 2, mock.attempts)
	})

	t.Run("POST body is reset on 503 retry", func(t *testing.T) {
		bodyContent := `{"jsonrpc":"2.0","method":"initialize"}`
		var capturedBodies []string

		mock := &mockRoundTripper{
			responses: []*http.Response{
				{StatusCode: http.StatusServiceUnavailable},
				{StatusCode: http.StatusOK},
			},
		}

		// Override the mock to capture bodies
		originalRT := mock
		capturingRT := roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Body != nil {
				b, _ := io.ReadAll(req.Body)
				capturedBodies = append(capturedBodies, string(b))
			}
			return originalRT.RoundTrip(req)
		})

		rt := &RetryRoundTripper{
			Base:       capturingRT,
			MaxRetries: 3,
			Backoff:    func(i int) time.Duration { return 0 },
		}

		req, _ := http.NewRequest("POST", "http://example.com/mcp", strings.NewReader(bodyContent))
		req.Header.Set("Content-Type", "application/json")
		resp, err := rt.RoundTrip(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		// Both attempts should have received the full body
		assert.Len(t, capturedBodies, 2)
		assert.Equal(t, bodyContent, capturedBodies[0], "first attempt should have full body")
		assert.Equal(t, bodyContent, capturedBodies[1], "retry should have full body (reset via GetBody)")
	})
}

// roundTripFunc adapts a function to http.RoundTripper
type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
