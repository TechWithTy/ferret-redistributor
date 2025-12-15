package notion

import (
	"context"
	"fmt"
	"strings"
)

// QueryAllPages runs a Notion data source query and returns all matching pages (paginated).
// req is the raw Notion query body (filter/sorts/page_size/start_cursor).
func (c *Client) QueryAllPages(ctx context.Context, dataSourceID string, req map[string]any) ([]pageObject, error) {
	ds := strings.TrimSpace(dataSourceID)
	if ds == "" {
		return nil, fmt.Errorf("notion: missing data_source_id")
	}
	if req == nil {
		req = map[string]any{}
	}

	// ensure page_size is set (max 100)
	if _, ok := req["page_size"]; !ok {
		req["page_size"] = 100
	}

	var out []pageObject
	startCursor := ""
	for page := 0; page < 10_000; page++ { // safety cap
		if startCursor != "" {
			req["start_cursor"] = startCursor
		} else {
			delete(req, "start_cursor")
		}

		var res queryResponseFull
		_, _, err := c.doJSON(ctx, "POST", "/data_sources/"+ds+"/query", req, &res)
		if err != nil {
			return nil, err
		}
		out = append(out, res.Results...)
		if !res.HasMore || strings.TrimSpace(res.NextCursor) == "" {
			break
		}
		startCursor = res.NextCursor
	}
	return out, nil
}


