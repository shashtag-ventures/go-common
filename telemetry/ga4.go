package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	ga4Endpoint = "https://www.google-analytics.com/mp/collect"
)

// GA4Tracker implements the Tracker interface for Google Analytics 4 (Measurement Protocol).
type GA4Tracker struct {
	measurementID string
	apiSecret     string
	httpClient    *http.Client
}

// NewGA4Tracker creates a new GA4Tracker instance.
func NewGA4Tracker(measurementID, apiSecret string) *GA4Tracker {
	return &GA4Tracker{
		measurementID: measurementID,
		apiSecret:     apiSecret,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type ga4Payload struct {
	ClientID string      `json:"client_id"`
	Events   []ga4Event `json:"events"`
}

type ga4Event struct {
	Name   string                 `json:"name"`
	Params map[string]interface{} `json:"params"`
}

func (g *GA4Tracker) Track(ctx context.Context, event Event) error {
	payload := ga4Payload{
		ClientID: event.UserID,
		Events: []ga4Event{
			{
				Name:   event.Name,
				Params: event.Properties,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s?measurement_id=%s&api_secret=%s", ga4Endpoint, g.measurementID, g.apiSecret)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("ga4 request failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (g *GA4Tracker) Identify(ctx context.Context, userID string, traits map[string]interface{}) error {
	// GA4 doesn't have a direct "Identify" call like PostHog.
	// Typically, you send a user_id parameter with events.
	return nil
}

func (g *GA4Tracker) Flush() error {
	return nil
}

func (g *GA4Tracker) Close() error {
	return nil
}
