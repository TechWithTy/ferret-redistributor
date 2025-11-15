package linkedin

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
)

// PostStats contains commonly used counters for a LinkedIn post.
type PostStats struct {
    Impressions       int `json:"impressions"`
    UniqueImpressions int `json:"uniqueImpressions"`
    Reactions         int `json:"reactions"`
    Comments          int `json:"comments"`
    Shares            int `json:"shares"`
}

// GetPostStatistics tries the REST posts statistics endpoint, then falls back to v2 socialActions.
// Accepts either a full URN (e.g., urn:li:post:123) or a numeric/string id (assumed post id).
func (c *Client) GetPostStatistics(ctx context.Context, postIDOrURN string) (*PostStats, error) {
    urn := postIDOrURN
    if !strings.HasPrefix(urn, "urn:li:") {
        urn = "urn:li:post:" + urn
    }
    if ps, err := c.getRestPostStatistics(ctx, urn); err == nil {
        return ps, nil
    }
    return c.getV2SocialActions(ctx, urn)
}

func (c *Client) getRestPostStatistics(ctx context.Context, urn string) (*PostStats, error) {
    endpoint := "https://api.linkedin.com/rest/posts/" + url.PathEscape(urn) + "/statistics"
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
    if err != nil { return nil, err }
    req.Header.Set("Authorization", "Bearer "+c.cfg.AccessToken)
    req.Header.Set("LinkedIn-Version", "202401")
    req.Header.Set("X-Restli-Protocol-Version", "2.0.0")
    resp, err := c.httpClient.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return nil, MapHTTPError(resp) }
    body, err := io.ReadAll(resp.Body)
    if err != nil { return nil, err }
    var data struct {
        TotalImpressionStatistics *struct {
            ImpressionCount       int `json:"impressionCount"`
            UniqueImpressionsCount int `json:"uniqueImpressionsCount"`
        } `json:"totalImpressionStatistics"`
        TotalReactionStatistics *struct {
            Count int `json:"count"`
        } `json:"totalReactionStatistics"`
        TotalCommentStatistics *struct {
            CommentCount int `json:"commentCount"`
        } `json:"totalCommentStatistics"`
        TotalShareStatistics *struct {
            ShareCount int `json:"shareCount"`
        } `json:"totalShareStatistics"`
    }
    if err := json.Unmarshal(body, &data); err != nil { return nil, err }
    out := PostStats{}
    if data.TotalImpressionStatistics != nil {
        out.Impressions = data.TotalImpressionStatistics.ImpressionCount
        out.UniqueImpressions = data.TotalImpressionStatistics.UniqueImpressionsCount
    }
    if data.TotalReactionStatistics != nil { out.Reactions = data.TotalReactionStatistics.Count }
    if data.TotalCommentStatistics != nil { out.Comments = data.TotalCommentStatistics.CommentCount }
    if data.TotalShareStatistics != nil { out.Shares = data.TotalShareStatistics.ShareCount }
    return &out, nil
}

func (c *Client) getV2SocialActions(ctx context.Context, urn string) (*PostStats, error) {
    // Encode URN for URL path
    endpoint := "https://api.linkedin.com/v2/socialActions/" + url.PathEscape(urn)
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
    if err != nil { return nil, err }
    req.Header.Set("Authorization", "Bearer "+c.cfg.AccessToken)
    resp, err := c.httpClient.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return nil, MapHTTPError(resp) }
    body, err := io.ReadAll(resp.Body)
    if err != nil { return nil, err }
    var data map[string]any
    if err := json.Unmarshal(body, &data); err != nil { return nil, err }
    // Best-effort extraction of counts
    ps := &PostStats{}
    if ls, ok := data["likesSummary"].(map[string]any); ok {
        if v, ok := ls["totalLikes"].(float64); ok { ps.Reactions = int(v) }
    }
    if cs, ok := data["commentsSummary"].(map[string]any); ok {
        if v, ok := cs["totalFirstLevelComments"].(float64); ok { ps.Comments = int(v) }
    }
    // Shares and impressions may not be present in this endpoint
    return ps, nil
}

// Debug helper to format errors; not exported.
func formatHTTPError(status int, body []byte) error { return fmt.Errorf("linkedin: http %d: %s", status, string(body)) }

