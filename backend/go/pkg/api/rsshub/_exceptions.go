package rsshub

import (
	"errors"
	"fmt"
)

var (
	// ErrMissingRoutePath occurs when FeedRequest.Path is empty.
	ErrMissingRoutePath = errors.New("rsshub: route path is required")
	// ErrMissingURL occurs when a URL parameter is required but missing.
	ErrMissingURL = errors.New("rsshub: url is required")
)

// APIError represents non-2xx responses from RSSHub.
type APIError struct {
	StatusCode int
	Body       []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("rsshub: api error status=%d body=%s", e.StatusCode, string(e.Body))
}
