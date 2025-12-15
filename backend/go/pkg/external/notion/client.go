package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.notion.com/v1"
	// As of 2025-09-03 Notion introduces multi-source databases; write operations should use data_source_id.
	// See: https://developers.notion.com/docs/upgrade-guide-2025-09-03
	notionVersion = "2025-09-03"
)

// Client is a minimal Notion HTTP client for data source query + page create/update.
type Client struct {
	baseURL string
	apiKey  string
	hc      *http.Client
}

type Config struct {
	APIKey      string
	BaseURL     string
	HTTPTimeout time.Duration
}

func New(cfg Config) (*Client, error) {
	key := strings.TrimSpace(cfg.APIKey)
	if key == "" {
		return nil, fmt.Errorf("notion: missing api key")
	}
	base := strings.TrimSpace(cfg.BaseURL)
	if base == "" {
		base = defaultBaseURL
	}
	to := cfg.HTTPTimeout
	if to <= 0 {
		to = 20 * time.Second
	}
	return &Client{
		baseURL: strings.TrimRight(base, "/"),
		apiKey:  key,
		hc:      &http.Client{Timeout: to},
	}, nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, reqBody any, out any) (*http.Response, []byte, error) {
	var body io.Reader
	if reqBody != nil {
		b, err := json.Marshal(reqBody)
		if err != nil {
			return nil, nil, err
		}
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Notion-Version", notionVersion)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.hc.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	raw, _ := io.ReadAll(res.Body)
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		msg := strings.TrimSpace(string(raw))
		if len(msg) > 2000 {
			msg = msg[:2000] + "…"
		}
		return res, raw, fmt.Errorf("notion: %s %s -> %s: %s", method, path, res.Status, msg)
	}
	if out != nil {
		if err := json.Unmarshal(raw, out); err != nil {
			return res, raw, err
		}
	}
	return res, raw, nil
}
