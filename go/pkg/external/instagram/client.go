package instagram

import (
    "context"
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "net/url"
    "strconv"
    "strings"
    "time"
)

type Client struct {
    httpClient *http.Client
    cfg        Config
}

// MediaItem represents a minimal media object returned by hashtag media endpoints.
type MediaItem struct {
    ID               string `json:"id"`
    Caption          string `json:"caption"`
    MediaType        string `json:"media_type"`
    MediaURL         string `json:"media_url"`
    Permalink        string `json:"permalink"`
    CommentsCount    int    `json:"comments_count"`
    LikeCount        int    `json:"like_count"`
    Timestamp        string `json:"timestamp"`
    MediaProductType string `json:"media_product_type"`
}

func New(cfg Config) *Client {
    timeout := cfg.HTTPTimeout
    if timeout <= 0 {
        timeout = 30 * time.Second
    }
    return &Client{
        httpClient: &http.Client{Timeout: timeout},
        cfg:        cfg,
    }
}

func (c *Client) PostFeedImage(ctx context.Context, imageURL, caption string) (string, error) {
    if imageURL == "" {
        return "", ErrValidation
    }
    cr, err := c.createMedia(ctx, createMediaParams{Caption: caption, ImageURL: imageURL})
    if err != nil {
        return "", err
    }
    return c.publish(ctx, publishParams{CreationID: cr.ID})
}

func (c *Client) IGUserID() string { return c.cfg.IGUserID }

func (c *Client) PostFeedVideo(ctx context.Context, videoURL, caption string) (string, error) {
    if videoURL == "" {
        return "", ErrValidation
    }
    cr, err := c.createMedia(ctx, createMediaParams{Caption: caption, VideoURL: videoURL})
    if err != nil {
        return "", err
    }
    if err := c.waitForContainerReady(ctx, cr.ID); err != nil {
        return "", err
    }
    return c.publish(ctx, publishParams{CreationID: cr.ID})
}

func (c *Client) PostCarousel(ctx context.Context, mediaURLs []string, caption string) (string, error) {
    if len(mediaURLs) == 0 {
        return "", ErrValidation
    }
    children := make([]string, 0, len(mediaURLs))
    for _, m := range mediaURLs {
        if m == "" {
            return "", ErrValidation
        }
        var child createMediaParams
        if isVideoURL(m) {
            child = createMediaParams{VideoURL: m, IsCarousel: true}
        } else {
            child = createMediaParams{ImageURL: m, IsCarousel: true}
        }
        cr, err := c.createMedia(ctx, child)
        if err != nil {
            return "", err
        }
        children = append(children, cr.ID)
    }
    cr, err := c.createMedia(ctx, createMediaParams{Caption: caption, IsCarousel: true, Children: children})
    if err != nil {
        return "", err
    }
    return c.publish(ctx, publishParams{CreationID: cr.ID})
}

func (c *Client) PostReel(ctx context.Context, videoURL, caption string) (string, error) {
    if videoURL == "" {
        return "", ErrValidation
    }
    cr, err := c.createMedia(ctx, createMediaParams{Caption: caption, VideoURL: videoURL, IsReel: true})
    if err != nil {
        return "", err
    }
    if err := c.waitForContainerReady(ctx, cr.ID); err != nil {
        return "", err
    }
    return c.publish(ctx, publishParams{CreationID: cr.ID})
}

func (c *Client) PostStoryImage(ctx context.Context, imageURL, caption string) (string, error) {
    if imageURL == "" {
        return "", ErrValidation
    }
    cr, err := c.createMedia(ctx, createMediaParams{Caption: caption, ImageURL: imageURL, IsStory: true})
    if err != nil {
        return "", err
    }
    return c.publish(ctx, publishParams{CreationID: cr.ID})
}

func (c *Client) PostStoryVideo(ctx context.Context, videoURL, caption string) (string, error) {
    if videoURL == "" {
        return "", ErrValidation
    }
    cr, err := c.createMedia(ctx, createMediaParams{Caption: caption, VideoURL: videoURL, IsStory: true})
    if err != nil {
        return "", err
    }
    if err := c.waitForContainerReady(ctx, cr.ID); err != nil {
        return "", err
    }
    return c.publish(ctx, publishParams{CreationID: cr.ID})
}

func (c *Client) createMedia(ctx context.Context, p createMediaParams) (*creationResponse, error) {
    endpoint := c.cfg.BaseURL + "/" + c.cfg.IGUserID + "/media"
    form := c.buildCreateMediaForm(p)
    if p.IsCarousel && p.ImageURL == "" && p.VideoURL == "" && len(p.Children) == 0 {
        return nil, ErrValidation
    }
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, mapHTTPError(resp)
    }
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    var cr creationResponse
    if err := json.Unmarshal(body, &cr); err != nil {
        return nil, err
    }
    if cr.ID == "" {
        return nil, errors.New("instagram: empty creation id")
    }
    return &cr, nil
}

func (c *Client) buildCreateMediaForm(p createMediaParams) url.Values {
    form := url.Values{}
    form.Set("access_token", c.cfg.AccessToken)
    if p.Caption != "" {
        form.Set("caption", p.Caption)
    }
    if p.ImageURL != "" {
        form.Set("image_url", p.ImageURL)
    }
    if p.VideoURL != "" {
        form.Set("video_url", p.VideoURL)
    }
    if p.ThumbOffsetSeconds > 0 {
        form.Set("thumb_offset", strconv.Itoa(p.ThumbOffsetSeconds))
    }
    if p.DisableComments {
        form.Set("disable_comments", "true")
    }
    if p.CoverURL != "" {
        form.Set("cover_url", p.CoverURL)
    }
    if p.IsCarousel && len(p.Children) > 0 {
        form.Set("children", strings.Join(p.Children, ","))
        form.Set("media_type", "CAROUSEL")
    }
    if p.IsReel {
        form.Set("media_type", "REELS")
    }
    if p.IsStory {
        form.Set("media_type", "STORIES")
    }
    return form
}

func (c *Client) waitForContainerReady(ctx context.Context, creationID string) error {
    if creationID == "" {
        return ErrValidation
    }
    deadline := time.Now().Add(2 * time.Minute)
    backoff := 2 * time.Second
    for {
        if time.Now().After(deadline) {
            return errors.New("instagram: processing timeout")
        }
        status, err := c.getContainerStatus(ctx, creationID)
        if err != nil {
            if errors.Is(err, ErrServer) || errors.Is(err, ErrRateLimited) {
                time.Sleep(backoff)
                if backoff < 10*time.Second { backoff *= 2 }
                continue
            }
            return err
        }
        switch strings.ToUpper(status) {
        case "FINISHED", "PUBLISHED":
            return nil
        case "ERROR", "FAILED":
            return errors.New("instagram: video processing failed")
        default:
            time.Sleep(backoff)
            if backoff < 10*time.Second { backoff *= 2 }
        }
    }
}

func (c *Client) getContainerStatus(ctx context.Context, creationID string) (string, error) {
    endpoint := c.cfg.BaseURL + "/" + creationID + "?fields=status_code&access_token=" + url.QueryEscape(c.cfg.AccessToken)
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
    if err != nil {
        return "", err
    }
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return "", mapHTTPError(resp)
    }
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    var cs containerStatusResponse
    if err := json.Unmarshal(body, &cs); err != nil {
        return "", err
    }
    return cs.StatusCode, nil
}

func (c *Client) publish(ctx context.Context, p publishParams) (string, error) {
    endpoint := c.cfg.BaseURL + "/" + c.cfg.IGUserID + "/media_publish"
    form := url.Values{}
    form.Set("access_token", c.cfg.AccessToken)
    form.Set("creation_id", p.CreationID)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
    if err != nil {
        return "", err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return "", mapHTTPError(resp)
    }
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }
    var pr publishResponse
    if err := json.Unmarshal(body, &pr); err != nil {
        return "", err
    }
    if pr.ID == "" {
        return "", errors.New("instagram: empty published id")
    }
    return pr.ID, nil
}

func mapHTTPError(resp *http.Response) error {
    if resp == nil {
        return ErrServer
    }
    if resp.StatusCode == http.StatusUnauthorized {
        return ErrUnauthorized
    }
    if resp.StatusCode == http.StatusForbidden {
        return ErrForbidden
    }
    if resp.StatusCode == http.StatusTooManyRequests {
        return ErrRateLimited
    }
    body, _ := io.ReadAll(resp.Body)
    var ge graphErrorResponse
    if err := json.Unmarshal(body, &ge); err == nil && ge.Error.Message != "" {
        switch resp.StatusCode {
        case http.StatusBadRequest:
            return ErrValidation
        case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
            return ErrServer
        default:
            return errors.New(ge.Error.Message)
        }
    }
    switch {
    case resp.StatusCode >= 500:
        return ErrServer
    case resp.StatusCode >= 400:
        return ErrValidation
    default:
        return ErrServer
    }
}

func isVideoURL(u string) bool {
    u = strings.ToLower(u)
    return strings.HasSuffix(u, ".mp4") || strings.HasSuffix(u, ".mov") || strings.HasSuffix(u, ".mkv") || strings.Contains(u, "video")
}

// SearchHashtag resolves a hashtag name to its Graph ID (or returns an error if not found).
func (c *Client) SearchHashtag(ctx context.Context, hashtag string) (string, error) {
    hashtag = strings.TrimPrefix(strings.TrimSpace(hashtag), "#")
    if hashtag == "" { return "", ErrValidation }
    endpoint := c.cfg.BaseURL + "/ig_hashtag_search?user_id=" + url.QueryEscape(c.cfg.IGUserID) + "&q=" + url.QueryEscape(hashtag)
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
    if err != nil { return "", err }
    resp, err := c.httpClient.Do(req)
    if err != nil { return "", err }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 { return "", mapHTTPError(resp) }
    body, err := io.ReadAll(resp.Body)
    if err != nil { return "", err }
    var payload struct{ Data []struct{ ID string `json:"id"` } `json:"data"` }
    if err := json.Unmarshal(body, &payload); err != nil { return "", err }
    if len(payload.Data) == 0 || payload.Data[0].ID == "" { return "", ErrNotFound }
    return payload.Data[0].ID, nil
}

// GetHashtagRecentMedia returns up to limit media items for a hashtag, optionally filtered by min likes/comments.
func (c *Client) GetHashtagRecentMedia(ctx context.Context, hashtagID string, limit int, minLikes, minComments int) ([]MediaItem, error) {
    if hashtagID == "" { return nil, ErrValidation }
    fields := "id,caption,media_type,media_url,permalink,comments_count,like_count,timestamp,media_product_type"
    endpoint := c.cfg.BaseURL + "/" + url.PathEscape(hashtagID) + "/recent_media?user_id=" + url.QueryEscape(c.cfg.IGUserID) + "&fields=" + url.QueryEscape(fields) + "&limit=100"
    out := make([]MediaItem, 0, min(100, max(1, limit)))
    next := endpoint
    remaining := limit
    for remaining > 0 && next != "" {
        req, err := http.NewRequestWithContext(ctx, http.MethodGet, next, nil)
        if err != nil { return nil, err }
        resp, err := c.httpClient.Do(req)
        if err != nil { return nil, err }
        if resp.StatusCode < 200 || resp.StatusCode >= 300 {
            _ = resp.Body.Close()
            return nil, mapHTTPError(resp)
        }
        body, err := io.ReadAll(resp.Body)
        _ = resp.Body.Close()
        if err != nil { return nil, err }
        var payload struct{
            Data   []MediaItem `json:"data"`
            Paging struct{
                Cursors struct{
                    After string `json:"after"`
                } `json:"cursors"`
                Next string `json:"next"`
            } `json:"paging"`
        }
        if err := json.Unmarshal(body, &payload); err != nil { return nil, err }
        for _, m := range payload.Data {
            if m.LikeCount >= minLikes && m.CommentsCount >= minComments {
                out = append(out, m)
                remaining--
                if remaining <= 0 { break }
            }
        }
        if remaining <= 0 { break }
        if payload.Paging.Next != "" {
            next = payload.Paging.Next
        } else if payload.Paging.Cursors.After != "" {
            next = endpoint + "&after=" + url.QueryEscape(payload.Paging.Cursors.After)
        } else {
            next = ""
        }
    }
    return out, nil
}

func min(a, b int) int { if a < b { return a }; return b }
func max(a, b int) int { if a > b { return a }; return b }

// PostFirstComment posts a comment on the authenticated user's media.
func (c *Client) PostFirstComment(ctx context.Context, mediaID, message string) error {
    if mediaID == "" || message == "" {
        return ErrValidation
    }
    endpoint := c.cfg.BaseURL + "/" + mediaID + "/comments"
    form := url.Values{}
    form.Set("access_token", c.cfg.AccessToken)
    form.Set("message", message)
    req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return mapHTTPError(resp)
    }
    return nil
}

type Comment struct {
    ID       string
    Text     string
    FromID   string
    Username string
}

// ListComments retrieves comments on a media with minimal fields.
func (c *Client) ListComments(ctx context.Context, mediaID string) ([]Comment, error) {
    if mediaID == "" {
        return nil, ErrValidation
    }
    endpoint := c.cfg.BaseURL + "/" + mediaID + "/comments?fields=id,text,from&access_token=" + url.QueryEscape(c.cfg.AccessToken)
    req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
    if err != nil {
        return nil, err
    }
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, mapHTTPError(resp)
    }
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    var payload struct {
        Data []struct {
            ID   string `json:"id"`
            Text string `json:"text"`
            From struct {
                ID       string `json:"id"`
                Username string `json:"username"`
            } `json:"from"`
        } `json:"data"`
    }
    if err := json.Unmarshal(body, &payload); err != nil {
        return nil, err
    }
    out := make([]Comment, 0, len(payload.Data))
    for _, d := range payload.Data {
        out = append(out, Comment{ID: d.ID, Text: d.Text, FromID: d.From.ID, Username: d.From.Username})
    }
    return out, nil
}

// SendDM attempts to send a direct message to a user. For Instagram, this may
// require messaging permissions and Page-scoped IDs; thus this method returns
// ErrUnsupportedFeature unless messaging is explicitly enabled/configured.
func (c *Client) SendDM(ctx context.Context, userID, message string) error {
    _ = userID
    _ = message
    return ErrUnsupportedFeature
}

// TrendingCriteria defines how to fetch trending content via hashtags.
type TrendingCriteria struct {
    Hashtags        []string
    LimitPerHashtag int
    MinLikes        int
    MinComments     int
}

// TrendingItem represents a media item with a simple engagement score and source hashtag.
type TrendingItem struct {
    Hashtag string     `json:"hashtag"`
    Media   MediaItem  `json:"media"`
    Score   int        `json:"score"` // like_count + comments_count
}

// GetTrendingContent aggregates recent media for multiple hashtags and returns items
// sorted by a simple score (likes + comments). Caller may do further scoring.
func (c *Client) GetTrendingContent(ctx context.Context, criteria TrendingCriteria) ([]TrendingItem, error) {
    if len(criteria.Hashtags) == 0 {
        return nil, ErrValidation
    }
    limit := criteria.LimitPerHashtag
    if limit <= 0 { limit = 10 }
    minLikes := criteria.MinLikes
    minComments := criteria.MinComments

    out := make([]TrendingItem, 0, len(criteria.Hashtags)*limit)
    for _, raw := range criteria.Hashtags {
        tag := strings.TrimPrefix(strings.TrimSpace(raw), "#")
        if tag == "" { continue }
        id, err := c.SearchHashtag(ctx, tag)
        if err != nil { continue }
        media, err := c.GetHashtagRecentMedia(ctx, id, limit, minLikes, minComments)
        if err != nil { continue }
        for _, m := range media {
            score := m.LikeCount + m.CommentsCount
            out = append(out, TrendingItem{Hashtag: tag, Media: m, Score: score})
        }
    }
    // simple sort by score desc
    sort.Slice(out, func(i, j int) bool { return out[i].Score > out[j].Score })
    return out, nil
}
