package notion

import (
	"context"
	"fmt"
	"strings"
)

// QueryFirstByTitle returns the first page in a data source whose title property equals the given value.
func (c *Client) QueryFirstByTitle(ctx context.Context, dataSourceID, titlePropName, equals string) (*pageRef, error) {
	ds := strings.TrimSpace(dataSourceID)
	if ds == "" {
		return nil, fmt.Errorf("notion: missing data_source_id")
	}
	prop := strings.TrimSpace(titlePropName)
	if prop == "" {
		return nil, fmt.Errorf("notion: missing title property name")
	}
	req := map[string]any{
		"filter": map[string]any{
			"property": prop,
			"title": map[string]any{
				"equals": equals,
			},
		},
		"page_size": 1,
	}
	var out queryResponse
	_, _, err := c.doJSON(ctx, "POST", "/data_sources/"+ds+"/query", req, &out)
	if err != nil {
		return nil, err
	}
	if len(out.Results) == 0 {
		return nil, nil
	}
	return &out.Results[0], nil
}

// QueryPageRefsByTitle loads all pages in a data source and returns a map from the given title property's plain_text
// value to the page ref (id/url).
//
// This is useful for bulk upserts: do one paginated query, then resolve pages locally instead of querying per item.
func (c *Client) QueryPageRefsByTitle(ctx context.Context, dataSourceID, titlePropName string) (map[string]pageRef, error) {
	ds := strings.TrimSpace(dataSourceID)
	if ds == "" {
		return nil, fmt.Errorf("notion: missing data_source_id")
	}
	prop := strings.TrimSpace(titlePropName)
	if prop == "" {
		return nil, fmt.Errorf("notion: missing title property name")
	}

	out := make(map[string]pageRef, 256)
	startCursor := ""
	for page := 0; page < 10_000; page++ { // safety cap
		req := map[string]any{
			"page_size": 100,
		}
		if startCursor != "" {
			req["start_cursor"] = startCursor
		}

		var res queryResponseFull
		_, _, err := c.doJSON(ctx, "POST", "/data_sources/"+ds+"/query", req, &res)
		if err != nil {
			return nil, err
		}
		for _, p := range res.Results {
			key := titleFromProperties(p.Properties, prop)
			if strings.TrimSpace(key) == "" {
				continue
			}
			out[key] = pageRef{ID: p.ID, URL: p.URL}
		}
		if !res.HasMore || strings.TrimSpace(res.NextCursor) == "" {
			break
		}
		startCursor = res.NextCursor
	}
	return out, nil
}

func titleFromProperties(props map[string]propertyValue, titlePropName string) string {
	if props == nil {
		return ""
	}
	p, ok := props[titlePropName]
	if !ok {
		return ""
	}
	if strings.ToLower(strings.TrimSpace(p.Type)) != "title" {
		return ""
	}
	for _, rt := range p.Title {
		if strings.TrimSpace(rt.PlainText) != "" {
			// Title can be split across segments; join if needed later.
			return rt.PlainText
		}
	}
	return ""
}

// CreatePageInDataSource creates a new page under the given data source.
func (c *Client) CreatePageInDataSource(ctx context.Context, dataSourceID string, properties map[string]any) (*pageRef, error) {
	ds := strings.TrimSpace(dataSourceID)
	if ds == "" {
		return nil, fmt.Errorf("notion: missing data_source_id")
	}
	req := map[string]any{
		"parent": map[string]any{
			"type":           "data_source_id",
			"data_source_id": ds,
		},
		"properties": properties,
	}
	var out pageRef
	_, _, err := c.doJSON(ctx, "POST", "/pages", req, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdatePageProperties updates a page's properties.
func (c *Client) UpdatePageProperties(ctx context.Context, pageID string, properties map[string]any) error {
	pid := strings.TrimSpace(pageID)
	if pid == "" {
		return fmt.Errorf("notion: missing page id")
	}
	req := map[string]any{
		"properties": properties,
	}
	_, _, err := c.doJSON(ctx, "PATCH", "/pages/"+pid, req, nil)
	return err
}
