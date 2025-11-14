package postiz

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Client provides typed helpers for the Postiz Public API.
type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// APIError represents an error returned by the Postiz API.
type APIError struct {
	Status  int
	Message string
	Body    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("postiz: status=%d message=%s", e.Status, e.Message)
	}
	return fmt.Sprintf("postiz: status=%d", e.Status)
}

// New creates a Client from the provided Config.
func New(cfg Config) (*Client, error) {
	if cfg.APIKey == "" {
		return nil, ErrMissingAPIKey
	}
	return &Client{
		httpClient: &http.Client{Timeout: cfg.normalizedTimeout()},
		baseURL:    strings.TrimRight(cfg.normalizedBaseURL(), "/"),
		apiKey:     cfg.APIKey,
	}, nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.apiKey)
	return req, nil
}

func (c *Client) do(req *http.Request, out any) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		apiErr := &APIError{
			Status: resp.StatusCode,
			Body:   string(bodyBytes),
		}
		var payload map[string]any
		if err := json.Unmarshal(bodyBytes, &payload); err == nil {
			if msg, ok := payload["message"].(string); ok {
				apiErr.Message = msg
			}
		}
		return apiErr
	}

	if out == nil {
		_, err = io.Copy(io.Discard, resp.Body)
		return err
	}

	return json.NewDecoder(resp.Body).Decode(out)
}

func (c *Client) doJSON(ctx context.Context, method, path string, payload any, out any) error {
	var (
		body io.Reader
	)
	if payload != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(payload); err != nil {
			return err
		}
		body = buf
	}
	req, err := c.newRequest(ctx, method, path, body)
	if err != nil {
		return err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.do(req, out)
}
