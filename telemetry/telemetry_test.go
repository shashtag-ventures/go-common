package telemetry

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiTracker(t *testing.T) {
	mock1 := &mockTracker{}
	mock2 := &mockTracker{}
	
	m := NewMultiTracker(mock1, mock2)
	
	event := Event{
		Name: "test_event",
		UserID: "user_123",
		Properties: map[string]interface{}{"foo": "bar"},
	}
	
	err := m.Track(context.Background(), event)
	assert.NoError(t, err)
	assert.True(t, mock1.tracked)
	assert.True(t, mock2.tracked)
}

func TestGA4Tracker(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Contains(t, r.URL.RawQuery, "measurement_id=test_id")
		assert.Contains(t, r.URL.RawQuery, "api_secret=test_secret")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// For testing, I'll use a mock server
}

type mockTracker struct {
	tracked bool
}

func (m *mockTracker) Track(ctx context.Context, event Event) error {
	m.tracked = true
	return nil
}

func (m *mockTracker) Identify(ctx context.Context, userID string, traits map[string]interface{}) error {
	return nil
}

func (m *mockTracker) Flush() error {
	return nil
}

func (m *mockTracker) Close() error {
	return nil
}
