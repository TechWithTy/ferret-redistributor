package behiiv

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strconv"
)

type AggregateStats struct {
    Email struct {
        Recipients   int `json:"recipients"`
        Delivered    int `json:"delivered"`
        Opens        int `json:"opens"`
        UniqueOpens  int `json:"unique_opens"`
        OpenRate     int `json:"open_rate"`
        Clicks       int `json:"clicks"`
        UniqueClicks int `json:"unique_clicks"`
        ClickRate    int `json:"click_rate"`
        Unsubscribes int `json:"unsubscribes"`
        SpamReports  int `json:"spam_reports"`
    } `json:"email"`
    Web struct {
        Views  int `json:"views"`
        Clicks int `json:"clicks"`
    } `json:"web"`
    Clicks []struct {
        URL         string `json:"url"`
        Email struct {
            Clicks           int `json:"clicks"`
            UniqueClicks     int `json:"unique_clicks"`
            ClickThroughRate int `json:"click_through_rate"`
        } `json:"email"`
        Web struct {
            Clicks           int `json:"clicks"`
            UniqueClicks     int `json:"unique_clicks"`
            ClickThroughRate int `json:"click_through_rate"`
        } `json:"web"`
        TotalClicks           int `json:"total_clicks"`
        TotalUniqueClicks     int `json:"total_unique_clicks"`
        TotalClickThroughRate int `json:"total_click_through_rate"`
    } `json:"clicks"`
}

type AggregateStatsResponse struct {
    Data struct {
        Stats AggregateStats `json:"stats"`
    } `json:"data"`
}

func (c *Client) GetAggregateStats(ctx context.Context, publicationID string) (*AggregateStatsResponse, error) {
    if publicationID == "" { return nil, fmt.Errorf("publicationID required") }
    endpoint := fmt.Sprintf("%s/%s/publications/%s/posts/aggregate_stats", c.cfg.BaseURL, c.cfg.Version, url.PathEscape(publicationID))
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
    if err != nil { return nil, err }
    req.Header.Set("Authorization", "Bearer "+c.cfg.Token)
    resp, err := c.httpClient.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return nil, fmt.Errorf("beehiiv http %d", resp.StatusCode) }
    body, err := io.ReadAll(resp.Body)
    if err != nil { return nil, err }
    var ar AggregateStatsResponse
    if err := json.Unmarshal(body, &ar); err != nil { return nil, err }
    return &ar, nil
}

// Post minimal representation from Beehiiv posts list
type Post struct {
    ID          string   `json:"id"`
    Title       string   `json:"title"`
    Subtitle    string   `json:"subtitle"`
    Authors     []string `json:"authors"`
    Created     int64    `json:"created"`
    Status      string   `json:"status"`
    SplitTested bool     `json:"split_tested"`

    SubjectLine            string   `json:"subject_line"`
    PreviewText            string   `json:"preview_text"`
    Slug                   string   `json:"slug"`
    ThumbnailURL           string   `json:"thumbnail_url"`
    WebURL                 string   `json:"web_url"`
    Audience               string   `json:"audience"`
    Platform               string   `json:"platform"`
    ContentTags            []string `json:"content_tags"`
    HiddenFromFeed         bool     `json:"hidden_from_feed"`
    PublishDate            int64    `json:"publish_date"`
    DisplayedDate          int64    `json:"displayed_date"`
    MetaDefaultDescription string   `json:"meta_default_description"`
    MetaDefaultTitle       string   `json:"meta_default_title"`

    Content struct {
        Free struct {
            Web   string `json:"web"`
            Email string `json:"email"`
            RSS   string `json:"rss"`
        } `json:"free"`
        Premium struct {
            Web   string `json:"web"`
            Email string `json:"email"`
        } `json:"premium"`
    } `json:"content"`

    Stats struct {
        Email struct {
            Recipients   int `json:"recipients"`
            Delivered    int `json:"delivered"`
            Opens        int `json:"opens"`
            UniqueOpens  int `json:"unique_opens"`
            OpenRate     int `json:"open_rate"`
            Clicks       int `json:"clicks"`
            UniqueClicks int `json:"unique_clicks"`
            ClickRate    int `json:"click_rate"`
            Unsubscribes int `json:"unsubscribes"`
            SpamReports  int `json:"spam_reports"`
        } `json:"email"`
        Web struct {
            Views  int `json:"views"`
            Clicks int `json:"clicks"`
        } `json:"web"`
        Clicks []struct {
            URL   string `json:"url"`
            Email struct {
                Clicks           int `json:"clicks"`
                UniqueClicks     int `json:"unique_clicks"`
                ClickThroughRate int `json:"click_through_rate"`
            } `json:"email"`
            Web struct {
                Clicks           int `json:"clicks"`
                UniqueClicks     int `json:"unique_clicks"`
                ClickThroughRate int `json:"click_through_rate"`
            } `json:"web"`
            TotalClicks           int `json:"total_clicks"`
            TotalUniqueClicks     int `json:"total_unique_clicks"`
            TotalClickThroughRate int `json:"total_click_through_rate"`
        } `json:"clicks"`
    } `json:"stats"`
}

type ListPostsResponse struct {
    Data         []Post `json:"data"`
    Limit        int    `json:"limit"`
    Page         int    `json:"page"`
    TotalResults int    `json:"total_results"`
    TotalPages   int    `json:"total_pages"`
    Links        struct {
        Next string `json:"next"`
        Prev string `json:"prev"`
    } `json:"links"`
}

func (c *Client) ListPosts(ctx context.Context, publicationID string, page, limit int) (*ListPostsResponse, error) {
    if publicationID == "" { return nil, fmt.Errorf("publicationID required") }
    v := url.Values{}
    if page > 0 { v.Set("page", strconv.Itoa(page)) }
    if limit > 0 { v.Set("limit", strconv.Itoa(limit)) }
    endpoint := fmt.Sprintf("%s/%s/publications/%s/posts", c.cfg.BaseURL, c.cfg.Version, url.PathEscape(publicationID))
    if qs := v.Encode(); qs != "" { endpoint += "?" + qs }
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
    if err != nil { return nil, err }
    req.Header.Set("Authorization", "Bearer "+c.cfg.Token)
    resp, err := c.httpClient.Do(req)
    if err != nil { return nil, err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return nil, fmt.Errorf("beehiiv http %d", resp.StatusCode) }
    body, err := io.ReadAll(resp.Body)
    if err != nil { return nil, err }
    var lr ListPostsResponse
    if err := json.Unmarshal(body, &lr); err != nil { return nil, err }
    return &lr, nil
}
