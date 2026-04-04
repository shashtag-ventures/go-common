package telemetry

import (
	"context"

	"github.com/posthog/posthog-go"
)

// PostHogTracker implements the Tracker interface for PostHog.
type PostHogTracker struct {
	client posthog.Client
}

// NewPostHogTracker creates a new PostHogTracker instance.
func NewPostHogTracker(apiKey string, host string) (*PostHogTracker, error) {
	client, err := posthog.NewWithConfig(apiKey, posthog.Config{
		Endpoint: host,
	})
	if err != nil {
		return nil, err
	}
	return &PostHogTracker{client: client}, nil
}

func (p *PostHogTracker) Track(ctx context.Context, event Event) error {
	properties := posthog.NewProperties()
	for k, v := range event.Properties {
		properties.Set(k, v)
	}

	return p.client.Enqueue(posthog.Capture{
		DistinctId: event.UserID,
		Event:      event.Name,
		Properties: properties,
	})
}

func (p *PostHogTracker) Identify(ctx context.Context, userID string, traits map[string]interface{}) error {
	properties := posthog.NewProperties()
	for k, v := range traits {
		properties.Set(k, v)
	}

	return p.client.Enqueue(posthog.Identify{
		DistinctId: userID,
		Properties: properties,
	})
}

func (p *PostHogTracker) Flush() error {
	return p.client.Close()
}

func (p *PostHogTracker) Close() error {
	return p.client.Close()
}
