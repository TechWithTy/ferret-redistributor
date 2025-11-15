package rsshub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL points to the public RSSHub instance.
	DefaultBaseURL = "https://rsshub.app"
)

// Option customizes the client.
type Option func(*Client)

// WithHTTPClient overrides the default http.Client (timeout 30s).
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		if h != nil {
			c.httpClient = h
		}
	}
}

// WithBaseURL overrides the RSSHub base URL (useful for self-hosted instances).
func WithBaseURL(base string) Option {
	return func(c *Client) {
		if base != "" {
			c.baseURL = strings.TrimRight(base, "/")
		}
	}
}

// Client interacts with RSSHub's API and feed routes.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient constructs a Client with optional overrides.
func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL: DefaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// GetRoutes fetches metadata for all RSSHub routes.
func (c *Client) GetRoutes(ctx context.Context) (*RoutesResponse, error) {
	var resp RoutesResponse
	if err := c.getJSON(ctx, "/api/routes", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// GetVersion retrieves the running RSSHub version.
func (c *Client) GetVersion(ctx context.Context) (*VersionResponse, error) {
	var resp VersionResponse
	if err := c.getJSON(ctx, "/api/version", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SearchRadar hits /api/radar/search with the provided URL.
func (c *Client) SearchRadar(ctx context.Context, req RadarSearchRequest) (*RadarSearchResponse, error) {
	if req.URL == "" {
		return nil, ErrMissingURL
	}

	u, err := url.Parse(c.baseURL + "/api/radar/search")
	if err != nil {
		return nil, fmt.Errorf("rsshub: parse radar url: %w", err)
	}
	q := u.Query()
	q.Set("url", req.URL)
	u.RawQuery = q.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("rsshub: create radar request: %w", err)
	}

	resp, err := c.do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RadarSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("rsshub: decode radar response: %w", err)
	}

	return &result, nil
}

// ForceRefresh triggers RSSHub to refresh a cached feed.
func (c *Client) ForceRefresh(ctx context.Context, req ForceRefreshRequest) (*ForceRefreshResponse, error) {
	if req.TargetURL == "" {
		return nil, ErrMissingURL
	}

	refreshURL := fmt.Sprintf("%s/api/force-refresh/%s", c.baseURL, url.QueryEscape(req.TargetURL))
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, refreshURL, nil)
	if err != nil {
		return nil, fmt.Errorf("rsshub: create refresh request: %w", err)
	}

	resp, err := c.do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ForceRefreshResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("rsshub: decode refresh response: %w", err)
	}

	return &result, nil
}

// FetchFeed requests the rendered feed for a given route.
func (c *Client) FetchFeed(ctx context.Context, req FeedRequest) (*FeedResult, error) {
	if strings.TrimSpace(req.Path) == "" {
		return nil, ErrMissingRoutePath
	}

	fullURL := c.baseURL + "/" + strings.TrimLeft(req.Path, "/")

	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, fmt.Errorf("rsshub: parse feed url: %w", err)
	}

	q := u.Query()
	for k, v := range req.Query {
		q.Set(k, v)
	}
	if req.Format != "" {
		q.Set("format", req.Format)
	}
	if len(q) > 0 {
		u.RawQuery = q.Encode()
	}
	fullURL = u.String()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("rsshub: create feed request: %w", err)
	}

	resp, err := c.do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("rsshub: read feed response: %w", err)
	}

	return &FeedResult{
		ContentType: resp.Header.Get("Content-Type"),
		Body:        body,
	}, nil
}

func (c *Client) getJSON(ctx context.Context, path string, target any) error {
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return fmt.Errorf("rsshub: create request: %w", err)
	}

	resp, err := c.do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("rsshub: decode response: %w", err)
	}
	return nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("rsshub: send request: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4<<10))
		_ = resp.Body.Close()
		return nil, &APIError{StatusCode: resp.StatusCode, Body: body}
	}

	return resp, nil
}

// EncodeRoutesResponse pretty prints a RoutesResponse (mainly for debugging).
func EncodeRoutesResponse(routes *RoutesResponse) ([]byte, error) {
	if routes == nil {
		return nil, nil
	}
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(routes); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// FetchDebugJSON fetches ctx.set('json', obj) payloads via ?format=debug.json.
func (c *Client) FetchDebugJSON(ctx context.Context, path string, query map[string]string) ([]byte, error) {
	res, err := c.FetchFeed(ctx, FeedRequest{
		Path:   path,
		Query:  cloneQuery(query),
		Format: "debug.json",
	})
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

// FetchItemDebugHTML fetches ?format={index}.debug.html output for previewing descriptions.
func (c *Client) FetchItemDebugHTML(ctx context.Context, path string, itemIndex int, query map[string]string) ([]byte, error) {
	format := fmt.Sprintf("%d.debug.html", itemIndex)
	res, err := c.FetchFeed(ctx, FeedRequest{
		Path:   path,
		Query:  cloneQuery(query),
		Format: format,
	})
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func cloneQuery(src map[string]string) map[string]string {
	if len(src) == 0 {
		return nil
	}
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
