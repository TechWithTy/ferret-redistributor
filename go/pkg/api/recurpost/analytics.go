package recurpost

import "context"

// AnalyticsService provides analytics endpoints.
type AnalyticsService struct{ c *Client }

type PostStats struct {
	PostID      string `json:"post_id"`
	Clicks      int    `json:"clicks"`
	Likes       int    `json:"likes"`
	Shares      int    `json:"shares"`
	Impressions int    `json:"impressions"`
}

// PostAnalytics returns basic metrics for a post.
func (s *AnalyticsService) PostAnalytics(ctx context.Context, id string) (*PostStats, error) {
	// GET /analytics/posts/{id}
	return nil, ErrNotImplemented
}
