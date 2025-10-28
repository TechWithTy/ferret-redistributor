package youtube

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
)

// VideoStatistics basic counts from YouTube Data API v3 videos.list?part=statistics
type VideoStatistics struct {
    ViewCount    int64 `json:"viewCount"`
    LikeCount    int64 `json:"likeCount"`
    CommentCount int64 `json:"commentCount"`
}

// GetVideoStatistics fetches basic counters for a video ID using API key.
func (c *Client) GetVideoStatistics(ctx context.Context, videoID string) (*VideoStatistics, error) {
    if videoID == "" || c.cfg.APIKey == "" { return nil, ErrValidation }
    u := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=statistics&id=%s&key=%s", url.QueryEscape(videoID), url.QueryEscape(c.cfg.APIKey))
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
    if err != nil { return nil, err }
    resp, err := c.httpClient.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    switch resp.StatusCode {
    case http.StatusUnauthorized:
        return nil, ErrUnauthorized
    case http.StatusForbidden:
        return nil, ErrForbidden
    case http.StatusTooManyRequests:
        return nil, ErrRateLimited
    }
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return nil, ErrServer }
    body, err := io.ReadAll(resp.Body)
    if err != nil { return nil, err }
    var payload struct {
        Items []struct{
            Statistics struct{
                ViewCount    string `json:"viewCount"`
                LikeCount    string `json:"likeCount"`
                CommentCount string `json:"commentCount"`
            } `json:"statistics"`
        } `json:"items"`
    }
    if err := json.Unmarshal(body, &payload); err != nil { return nil, err }
    if len(payload.Items) == 0 { return nil, ErrNotFound }
    s := payload.Items[0].Statistics
    vs := &VideoStatistics{ViewCount: atoi64(s.ViewCount), LikeCount: atoi64(s.LikeCount), CommentCount: atoi64(s.CommentCount)}
    return vs, nil
}

func atoi64(s string) int64 { var n int64; for _, r := range s { if r >= '0' && r <= '9' { n = n*10 + int64(r-'0') } }; return n }

