package postiz

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	baseURL   string
	apiKey    string
	http      *http.Client
	userAgent string

	Integrations *IntegrationsService
	Slots        *SlotsService
	Upload       *UploadService
	Posts        *PostsService
	Video        *VideoService
}

type Option func(*Client)

func WithBaseURL(u string) Option          { return func(c *Client) { c.baseURL = strings.TrimRight(u, "/") } }
func WithAPIKey(k string) Option           { return func(c *Client) { c.apiKey = k } }
func WithHTTPClient(h *http.Client) Option { return func(c *Client) { c.http = h } }
func WithUserAgent(ua string) Option       { return func(c *Client) { c.userAgent = ua } }

func NewClient(opts ...Option) *Client {
	c := &Client{
		baseURL:   "https://api.postiz.com/public/v1",
		http:      &http.Client{Timeout: 30 * time.Second},
		userAgent: "SocialScale-Postiz/0.1",
	}
	for _, o := range opts {
		o(c)
	}
	c.Integrations = &IntegrationsService{c}
	c.Slots = &SlotsService{c}
	c.Upload = &UploadService{c}
	c.Posts = &PostsService{c}
	c.Video = &VideoService{c}
	return c
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	if c.baseURL == "" {
		return nil, fmt.Errorf("empty base url")
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	req.Header.Set("Accept", "application/json")
	if body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", c.apiKey)
	}
	return req, nil
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		var e struct{ Code, Message string }
		_ = json.Unmarshal(b, &e)
		return &APIError{StatusCode: resp.StatusCode, Code: e.Code, Message: e.Message}
	}
	if out != nil {
		if err := json.Unmarshal(b, out); err != nil {
			return err
		}
	}
	return nil
}
