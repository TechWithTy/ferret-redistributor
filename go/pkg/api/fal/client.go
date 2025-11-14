package fal

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

	"github.com/azer/logger"
)

const (
	defaultQueueBaseURL = "https://queue.fal.run"
	defaultPollInterval = 500 * time.Millisecond
)

var log = logger.New("fal-sdk")

// Option allows customizing the Fal client.
type Option func(*Client)

// WithHTTPClient overrides the http.Client used for outgoing calls.
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		if h != nil {
			c.httpClient = h
		}
	}
}

// WithQueueBaseURL overrides the queue base URL (useful for tests).
func WithQueueBaseURL(base string) Option {
	return func(c *Client) {
		if base != "" {
			c.queueBaseURL = strings.TrimRight(base, "/")
		}
	}
}

// WithPollInterval customizes how often Poll* helpers check the status endpoint.
func WithPollInterval(interval time.Duration) Option {
	return func(c *Client) {
		if interval > 0 {
			c.pollInterval = interval
		}
	}
}

// Client represents a Fal API client.
type Client struct {
	apiKey       string
	queueBaseURL string
	httpClient   *http.Client
	pollInterval time.Duration
}

// NewClient builds a Fal client with the provided API key.
func NewClient(apiKey string, opts ...Option) *Client {
	c := &Client{
		apiKey:       apiKey,
		queueBaseURL: defaultQueueBaseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		pollInterval: defaultPollInterval,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Call submits an initial request to Fal (e.g. "fal-ai/image-tool").
func (c *Client) Call(ctx context.Context, path string, options RequestOptions) (*Call, error) {
	if c.apiKey == "" {
		return nil, ErrMissingCredentials
	}
	if path == "" {
		return nil, ErrEmptyPath
	}

	method := options.Method
	if method == "" {
		method = http.MethodPost
	}

	fullURL, err := url.JoinPath(c.queueBaseURL, path)
	if err != nil {
		return nil, fmt.Errorf("fal: build url: %w", err)
	}

	var body io.Reader
	if options.Payload != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(options.Payload); err != nil {
			return nil, fmt.Errorf("fal: marshal payload: %w", err)
		}
		body = buf
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("fal: create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Key %s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fal: initial request failed: %w", err)
	}
	defer resp.Body.Close()

	payload, err := decodeResponse(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Info("fal initial response received", logger.Attrs{
		"request_id": payload.RequestID,
		"path":       path,
	})

	return &Call{
		client:      c,
		requestID:   payload.RequestID,
		statusURL:   payload.StatusURL,
		responseURL: payload.ResponseURL,
		cancelURL:   payload.CancelURL,
	}, nil
}

// Call represents a queued Fal request.
type Call struct {
	client      *Client
	requestID   string
	statusURL   string
	responseURL string
	cancelURL   string
}

// CheckStatus fetches the latest task status.
func (call *Call) CheckStatus(ctx context.Context) (*Response, error) {
	return call.doGet(ctx, call.statusURL)
}

// FetchResponse pulls the final response from Fal once the run is done.
func (call *Call) FetchResponse(ctx context.Context) (*Response, error) {
	return call.doGet(ctx, call.responseURL)
}

// Cancel aborts the queued request.
func (call *Call) Cancel(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, call.urlOrDefault(call.cancelURL), nil)
	if err != nil {
		return fmt.Errorf("fal: cancel request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Key %s", call.client.apiKey))

	resp, err := call.client.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("fal: cancel request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fal: cancel request unexpected status %d", resp.StatusCode)
	}

	log.Info("fal request cancelled", logger.Attrs{"request_id": call.requestID})
	return nil
}

// PollUntilCompletion keeps checking the status URL until the call finishes or the context cancels.
func (c *Client) PollUntilCompletion(ctx context.Context, call *Call) (*Response, error) {
	ticker := time.NewTicker(c.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			statusResp, err := call.CheckStatus(ctx)
			if err != nil {
				return nil, err
			}

			switch strings.ToUpper(statusResp.Status) {
			case "COMPLETED":
				log.Info("fal request completed", logger.Attrs{"request_id": call.requestID})
				return call.FetchResponse(ctx)
			case "FAILED":
				log.Error("fal request failed", logger.Attrs{"request_id": call.requestID})
				return nil, ErrCallIncomplete
			default:
				// optional: keep logging
			}
		}
	}
}

// PollWithProgress streams intermediate status payloads to the supplied channel.
func (c *Client) PollWithProgress(ctx context.Context, call *Call, progressCh chan<- *Response) (*Response, error) {
	ticker := time.NewTicker(c.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			statusResp, err := call.CheckStatus(ctx)
			if err != nil {
				return nil, err
			}

			select {
			case progressCh <- statusResp:
			default:
			}

			switch strings.ToUpper(statusResp.Status) {
			case "COMPLETED":
				return call.FetchResponse(ctx)
			case "FAILED":
				return nil, ErrCallIncomplete
			}
		}
	}
}

func (call *Call) doGet(ctx context.Context, rawURL string) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, call.urlOrDefault(rawURL), nil)
	if err != nil {
		return nil, fmt.Errorf("fal: build request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Key %s", call.client.apiKey))

	resp, err := call.client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fal: request failed: %w", err)
	}
	defer resp.Body.Close()

	return decodeResponse(resp.Body)
}

func (call *Call) urlOrDefault(raw string) string {
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return raw
	}
	if raw == "" {
		return call.client.queueBaseURL
	}
	u, err := url.JoinPath(call.client.queueBaseURL, raw)
	if err != nil {
		return raw
	}
	return u
}

func decodeResponse(r io.Reader) (*Response, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("fal: read response: %w", err)
	}

	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("fal: decode response: %w", err)
	}
	return &resp, nil
}
