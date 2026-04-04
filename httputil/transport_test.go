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
}
