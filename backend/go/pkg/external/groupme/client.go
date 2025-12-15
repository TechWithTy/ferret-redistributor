package groupme

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.groupme.com/v3"

// Config controls how the GroupMe client is created.
type Config struct {
	BaseURL     string
	HTTPTimeout time.Duration
}

// Client can call a limited subset of GroupMe's API needed for bots.
type Client struct {
	baseURL string
	hc      *http.Client
}

// New creates a GroupMe client.
func New(cfg Config) *Client {
	base := strings.TrimSpace(cfg.BaseURL)
	if base == "" {
		base = defaultBaseURL
	}
	to := cfg.HTTPTimeout
	if to <= 0 {
		to = 15 * time.Second
	}
	return &Client{
		baseURL: strings.TrimRight(base, "/"),
		hc: &http.Client{
			Timeout: to,
		},
	}
}

// PostBotMessage posts a message using a bot id (no access token required).
//
// Endpoint: POST {baseURL}/bots/post
// Body: {"bot_id":"...","text":"..."}
func (c *Client) PostBotMessage(ctx context.Context, botID, text string) error {
	if c == nil || c.hc == nil {
		return errors.New("groupme: nil client")
	}
	botID = strings.TrimSpace(botID)
	if botID == "" {
		return errors.New("groupme: missing bot id")
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return errors.New("groupme: missing text")
	}

	body, _ := json.Marshal(map[string]string{
		"bot_id": botID,
		"text":   text,
	})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/bots/post", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.hc.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("groupme: non-2xx response: %s", res.Status)
	}
	return nil
}
