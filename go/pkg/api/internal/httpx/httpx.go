package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

// Doer captures http.Client.Do for easier testing/mocking.
type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

// RetryPolicy controls simple retry behavior.
type RetryPolicy struct {
	MaxRetries int
	Backoff    func(attempt int) time.Duration
	RetryOn    map[int]bool // HTTP status codes to retry on (e.g., 429, 500)
}

func defaultBackoff(attempt int) time.Duration {
	return time.Duration(250*(attempt+1)) * time.Millisecond
}

// DoWithRetry executes an HTTP request with basic retry capability.
func DoWithRetry(client Doer, req *http.Request, rp *RetryPolicy) (*http.Response, error) {
	if rp == nil {
		rp = &RetryPolicy{MaxRetries: 0, Backoff: defaultBackoff, RetryOn: map[int]bool{429: true, 500: true, 502: true, 503: true}}
	}
	var resp *http.Response
	var err error
	for attempt := 0; attempt <= rp.MaxRetries; attempt++ {
		resp, err = client.Do(req)
		if err != nil {
			if attempt == rp.MaxRetries {
				return nil, err
			}
		} else if !rp.RetryOn[resp.StatusCode] {
			return resp, nil
		}
		// retry on next iteration
		time.Sleep(rp.Backoff(attempt))
	}
	return resp, err
}

// NewJSONRequest builds a JSON request with headers.
func NewJSONRequest(ctx context.Context, method, url string, payload any, headers map[string]string) (*http.Request, error) {
	var body io.Reader
	if payload != nil {
		b, _ := json.Marshal(payload)
		body = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req, nil
}
