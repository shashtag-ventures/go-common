package telemetry

import (
	"context"
)

// Event represents a telemetry event with its properties.
type Event struct {
	Name       string                 `json:"name"`
	UserID     string                 `json:"user_id"`
	Properties map[string]interface{} `json:"properties"`
}

// Tracker defines the interface for sending telemetry events.
type Tracker interface {
	// Track sends an event to the analytics provider.
	Track(ctx context.Context, event Event) error
	// Identify associates a user with specific traits.
	Identify(ctx context.Context, userID string, traits map[string]interface{}) error
	// Flush ensures all pending events are sent.
	Flush() error
	// Close releases any resources held by the tracker.
	Close() error
}

// MultiTracker wraps multiple trackers and sends events to all of them.
type MultiTracker struct {
	trackers []Tracker
}

func NewMultiTracker(trackers ...Tracker) *MultiTracker {
	return &MultiTracker{trackers: trackers}
}

func (m *MultiTracker) Track(ctx context.Context, event Event) error {
	for _, t := range m.trackers {
		_ = t.Track(ctx, event)
	}
	return nil
}

func (m *MultiTracker) Identify(ctx context.Context, userID string, traits map[string]interface{}) error {
	for _, t := range m.trackers {
		_ = t.Identify(ctx, userID, traits)
	}
	return nil
}

func (m *MultiTracker) Flush() error {
	for _, t := range m.trackers {
		_ = t.Flush()
	}
	return nil
}

func (m *MultiTracker) Close() error {
	for _, t := range m.trackers {
		_ = t.Close()
	}
	return nil
}
