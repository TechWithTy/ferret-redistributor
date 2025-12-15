package notion

import (
	"context"
	"fmt"
	"strings"
)

// GetPage fetches a Notion page object (including properties).
func (c *Client) GetPage(ctx context.Context, pageID string) (*PageObject, error) {
	pid := strings.TrimSpace(pageID)
	if pid == "" {
		return nil, fmt.Errorf("notion: missing page id")
	}
	var out PageObject
	_, _, err := c.doJSON(ctx, "GET", "/pages/"+pid, nil, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}


