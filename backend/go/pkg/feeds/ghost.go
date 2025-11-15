package feeds

import (
    "context"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "path"
    "strings"
    "time"

    "github.com/mmcdole/gofeed"
)

// GhostFetcher fetches and parses RSS feeds from Ghost.org blogs.
type GhostFetcher struct {
    // HTTP client to use for requests. If nil, a default client is used.
    Client *http.Client
    // Request timeout if Client is nil.
    Timeout time.Duration
    // Optional custom User-Agent header.
    UserAgent string
}

// Fetch retrieves and parses the RSS feed for a Ghost site base URL.
// Typically Ghost exposes an RSS feed at "<base>/rss/".
func (g *GhostFetcher) Fetch(ctx context.Context, baseURL string) (*gofeed.Feed, string, error) {
    rssURL, err := ghostRSSURL(baseURL)
    if err != nil {
        return nil, "", fmt.Errorf("invalid base URL: %w", err)
    }

    client := g.Client
    if client == nil {
        timeout := g.Timeout
        if timeout == 0 {
            timeout = 30 * time.Second
        }
        client = &http.Client{Timeout: timeout}
    }

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, rssURL, nil)
    if err != nil {
        return nil, rssURL, fmt.Errorf("create request: %w", err)
    }
    ua := g.UserAgent
    if strings.TrimSpace(ua) == "" {
        ua = "SocialScaleRSSFetcher/1.0 (+https://github.com/bitesinbyte/ferret)"
    }
    req.Header.Set("User-Agent", ua)
    req.Header.Set("Accept", "application/rss+xml, application/xml, text/xml; q=0.9, */*; q=0.8")

    resp, err := client.Do(req)
    if err != nil {
        return nil, rssURL, fmt.Errorf("fetch rss: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
        return nil, rssURL, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
    }

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, rssURL, fmt.Errorf("read rss body: %w", err)
    }

    parser := gofeed.NewParser()
    feed, err := parser.ParseString(string(data))
    if err != nil {
        return nil, rssURL, fmt.Errorf("parse rss: %w", err)
    }

    // Basic normalization of URLs inside the feed where applicable
    if feed.Link != "" {
        feed.Link = canonicalURL(feed.Link)
    }
    for _, item := range feed.Items {
        if item.Link != "" {
            item.Link = canonicalURL(item.Link)
        }
        if item.GUID != "" {
            item.GUID = strings.TrimSpace(item.GUID)
        }
    }

    return feed, rssURL, nil
}

// ghostRSSURL builds the canonical RSS URL for a Ghost site ("/rss/").
func ghostRSSURL(base string) (string, error) {
    base = strings.TrimSpace(base)
    if base == "" {
        return "", fmt.Errorf("empty base URL")
    }
    u, err := url.Parse(base)
    if err != nil {
        return "", err
    }
    if u.Scheme == "" {
        u.Scheme = "https"
    }
    // Ensure path joins with trailing rss/
    p := u.Path
    if !strings.HasSuffix(p, "/") {
        p += "/"
    }
    u.Path = path.Join(p, "rss") + "/"
    return u.String(), nil
}

// canonicalURL trims and normalizes host case; leaves path/query as-is.
func canonicalURL(s string) string {
    s = strings.TrimSpace(s)
    u, err := url.Parse(s)
    if err != nil || u.Scheme == "" || u.Host == "" {
        return s
    }
    u.Host = strings.ToLower(u.Host)
    return u.String()
}

