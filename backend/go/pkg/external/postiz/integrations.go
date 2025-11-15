package postiz

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// ListIntegrations retrieves all connected channels.
func (c *Client) ListIntegrations(ctx context.Context) ([]Integration, error) {
	var res []Integration
	if err := c.doJSON(ctx, http.MethodGet, "/integrations", nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// FindNextSlot returns the next available publishing slot for an integration.
func (c *Client) FindNextSlot(ctx context.Context, integrationID string) (time.Time, error) {
	if integrationID == "" {
		return time.Time{}, ErrMissingIntegrationID
	}
	var res SlotResponse
	path := fmt.Sprintf("/find-slot/%s", integrationID)
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &res); err != nil {
		return time.Time{}, err
	}
	return res.Date, nil
}
