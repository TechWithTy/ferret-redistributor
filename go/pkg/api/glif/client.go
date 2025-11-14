package glif

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
	// DefaultSimpleAPIBaseURL points to the hosted Simple API endpoint.
	DefaultSimpleAPIBaseURL = "https://simple-api.glif.app"
	// DefaultAPIBaseURL is the root for the REST endpoints (e.g. /api/glifs).
	DefaultAPIBaseURL = "https://glif.app/api"
)

// Option customizes the Glif client.
type Option func(*Client)

// WithHTTPClient overrides the default http.Client (30s timeout).
func WithHTTPClient(h *http.Client) Option {
	return func(c *Client) {
		if h != nil {
			c.httpClient = h
		}
	}
}

// WithSimpleAPIBaseURL sets a custom Simple API base URL.
func WithSimpleAPIBaseURL(base string) Option {
	return func(c *Client) {
		if base != "" {
			c.simpleBaseURL = base
		}
	}
}

// WithAPIBaseURL sets a custom REST API base URL.
func WithAPIBaseURL(base string) Option {
	return func(c *Client) {
		if base != "" {
			c.apiBaseURL = base
		}
	}
}

// Client is a typed helper for talking to the Glif APIs.
type Client struct {
	httpClient    *http.Client
	token         string
	simpleBaseURL string
	apiBaseURL    string
}

// NewClient constructs a Glif client with the given API token.
func NewClient(token string, opts ...Option) *Client {
	c := &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		simpleBaseURL: DefaultSimpleAPIBaseURL,
		apiBaseURL:    DefaultAPIBaseURL,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// RunWorkflow executes a workflow via the Simple API.
func (c *Client) RunWorkflow(ctx context.Context, req RunWorkflowRequest) (*RunWorkflowResponse, error) {
	if req.WorkflowID == "" && !req.UsePathID {
		return nil, ErrMissingWorkflowID
	}

	if len(req.Inputs) > 0 && len(req.NamedInputs) > 0 {
		return nil, ErrConflictingInputModes
	}

	payload := make(map[string]any)
	if !req.UsePathID && req.WorkflowID != "" {
		payload["id"] = req.WorkflowID
	}
	if len(req.NamedInputs) > 0 {
		payload["inputs"] = req.NamedInputs
	} else if len(req.Inputs) > 0 {
		payload["inputs"] = req.Inputs
	}
	if req.Visibility != "" {
		payload["visibility"] = req.Visibility
	}

	var body io.Reader
	if len(payload) > 0 {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(payload); err != nil {
			return nil, fmt.Errorf("glif: encode request: %w", err)
		}
		body = buf
	}

	endpoint := strings.TrimRight(c.simpleBaseURL, "/")
	if req.UsePathID {
		if req.WorkflowID == "" {
			return nil, ErrMissingWorkflowID
		}
		endpoint = fmt.Sprintf("%s/%s", endpoint, req.WorkflowID)
	}

	if req.Strict {
		withStrict, err := addQuery(endpoint, "strict", "1")
		if err != nil {
			return nil, err
		}
		endpoint = withStrict
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("glif: create request: %w", err)
	}
	if body != nil {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		b, _ := io.ReadAll(resp.Body)
		return nil, &APIError{StatusCode: resp.StatusCode, Body: b}
	}

	var result RunWorkflowResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("glif: decode workflow response: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("%w: %s", ErrWorkflowFailed, result.Error)
	}

	return &result, nil
}

// GetWorkflow retrieves a workflow definition by ID using /api/glifs?id=<id>.
func (c *Client) GetWorkflow(ctx context.Context, workflowID string) (*Workflow, error) {
	if workflowID == "" {
		return nil, ErrMissingWorkflowID
	}

	u, err := url.Parse(strings.TrimRight(c.apiBaseURL, "/") + "/glifs")
	if err != nil {
		return nil, fmt.Errorf("glif: parse api base url: %w", err)
	}

	q := u.Query()
	q.Set("id", workflowID)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("glif: create workflow request: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		b, _ := io.ReadAll(resp.Body)
		return nil, &APIError{StatusCode: resp.StatusCode, Body: b}
	}

	var workflows []Workflow
	if err := json.NewDecoder(resp.Body).Decode(&workflows); err != nil {
		return nil, fmt.Errorf("glif: decode workflow: %w", err)
	}

	if len(workflows) == 0 {
		return nil, ErrWorkflowNotFound
	}

	return &workflows[0], nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	if c.token != "" && req.Header.Get("Authorization") == "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("glif: send request: %w", err)
	}

	return resp, nil
}

func addQuery(rawURL, key, value string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("glif: parse url: %w", err)
	}

	q := u.Query()
	q.Set(key, value)
	u.RawQuery = q.Encode()
	return u.String(), nil
}
