package instagram

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
)

// MediaBasic holds quick, near-realtime counters available via fields query.
type MediaBasic struct {
    ID            string `json:"id"`
    LikeCount     int    `json:"like_count"`
    CommentsCount int    `json:"comments_count"`
    PlayCount     int    `json:"play_count"`
    SaveCount     int    `json:"save_count"`
}

// GetMediaBasic fetches basic counters for a media using fields query.
// fields: id,like_count,comments_count,play_count,save_count
func (c *Client) GetMediaBasic(ctx context.Context, mediaID string) (*MediaBasic, error) {
    if mediaID == "" { return nil, ErrValidation }
    fields := "id,like_count,comments_count,play_count,save_count"
    endpoint := fmt.Sprintf("%s/%s?fields=%s&access_token=%s", c.cfg.BaseURL, url.PathEscape(mediaID), url.QueryEscape(fields), url.QueryEscape(c.cfg.AccessToken))
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
    if err != nil { return nil, err }
    resp, err := c.httpClient.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return nil, mapHTTPError(resp) }
    body, err := io.ReadAll(resp.Body)
    if err != nil { return nil, err }
    var mb MediaBasic
    if err := json.Unmarshal(body, &mb); err != nil { return nil, err }
    return &mb, nil
}

// MediaInsight represents a single metric insight for a media.
type MediaInsight struct {
    Name   string  `json:"name"`
    Period string  `json:"period"`
    Value  float64 `json:"value"`
}

// GetMediaInsights fetches lifetime insights for a media. Metrics vary by media type.
// Example metrics: impressions, reach, saved, plays, likes, comments, video_views.
func (c *Client) GetMediaInsights(ctx context.Context, mediaID string, metrics []string) ([]MediaInsight, error) {
    if mediaID == "" || len(metrics) == 0 { return nil, ErrValidation }
    metricParam := strings.Join(metrics, ",")
    endpoint := fmt.Sprintf("%s/%s/insights?metric=%s&access_token=%s", c.cfg.BaseURL, url.PathEscape(mediaID), url.QueryEscape(metricParam), url.QueryEscape(c.cfg.AccessToken))
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
    if err != nil { return nil, err }
    resp, err := c.httpClient.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return nil, mapHTTPError(resp) }
    body, err := io.ReadAll(resp.Body)
    if err != nil { return nil, err }
    var payload struct {
        Data []struct {
            Name   string `json:"name"`
            Period string `json:"period"`
            Values []struct{
                Value any `json:"value"`
            } `json:"values"`
        } `json:"data"`
    }
    if err := json.Unmarshal(body, &payload); err != nil { return nil, err }
    out := make([]MediaInsight, 0, len(payload.Data))
    for _, d := range payload.Data {
        var val float64
        if len(d.Values) > 0 {
            switch v := d.Values[0].Value.(type) {
            case float64:
                val = v
            case int:
                val = float64(v)
            case map[string]any:
                if num, ok := v["value"].(float64); ok { val = num }
            }
        }
        out = append(out, MediaInsight{Name: d.Name, Period: d.Period, Value: val})
    }
    return out, nil
}

// UserInsight represents a user-level insight metric over a period.
type UserInsight struct {
    Name  string  `json:"name"`
    Value float64 `json:"value"`
    EndTime string `json:"end_time"`
}

// GetUserInsights fetches account insights for given metrics and period.
// period: day, week, days_28, month, lifetime (dependent on metric)
// If since/until are provided (unix seconds as strings), include them to bound range.
func (c *Client) GetUserInsights(ctx context.Context, metrics []string, period string, since, until string) ([]UserInsight, error) {
    if len(metrics) == 0 || period == "" { return nil, ErrValidation }
    vals := url.Values{}
    vals.Set("metric", strings.Join(metrics, ","))
    vals.Set("period", period)
    vals.Set("access_token", c.cfg.AccessToken)
    if since != "" { vals.Set("since", since) }
    if until != "" { vals.Set("until", until) }
    endpoint := fmt.Sprintf("%s/%s/insights?%s", c.cfg.BaseURL, c.cfg.IGUserID, vals.Encode())
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
    if err != nil { return nil, err }
    resp, err := c.httpClient.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return nil, mapHTTPError(resp) }
    body, err := io.ReadAll(resp.Body)
    if err != nil { return nil, err }
    var payload struct {
        Data []struct {
            Name string `json:"name"`
            Values []struct{ Value float64 `json:"value"`; EndTime string `json:"end_time"` } `json:"values"`
        } `json:"data"`
    }
    if err := json.Unmarshal(body, &payload); err != nil { return nil, err }
    out := make([]UserInsight, 0)
    for _, d := range payload.Data {
        for _, v := range d.Values {
            out = append(out, UserInsight{Name: d.Name, Value: v.Value, EndTime: v.EndTime})
        }
    }
    return out, nil
}

