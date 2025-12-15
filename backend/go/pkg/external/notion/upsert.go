package notion

import (
	"context"
	"fmt"
	"strings"
)

type UpsertResult struct {
	Page    pageRef
	Created bool
}

// UpsertByTitle finds a page by a title property value; if found it patches properties, otherwise it creates a new page.
// This is intended for idempotent writes (logs/metrics).
func (c *Client) UpsertByTitle(ctx context.Context, dataSourceID, titlePropName, titleValue string, properties map[string]any) (*UpsertResult, error) {
	ds := strings.TrimSpace(dataSourceID)
	if ds == "" {
		return nil, fmt.Errorf("notion: missing data_source_id")
	}
	prop := strings.TrimSpace(titlePropName)
	if prop == "" {
		return nil, fmt.Errorf("notion: missing title property name")
	}
	val := strings.TrimSpace(titleValue)
	if val == "" {
		return nil, fmt.Errorf("notion: missing title value")
	}

	existing, err := c.QueryFirstByTitle(ctx, ds, prop, val)
	if err != nil {
		return nil, err
	}
	if existing == nil || strings.TrimSpace(existing.ID) == "" {
		created, err := c.CreatePageInDataSource(ctx, ds, properties)
		if err != nil {
			return nil, err
		}
		return &UpsertResult{Page: *created, Created: true}, nil
	}
	if err := c.UpdatePageProperties(ctx, existing.ID, properties); err != nil {
		return nil, err
	}
	return &UpsertResult{Page: *existing, Created: false}, nil
}


