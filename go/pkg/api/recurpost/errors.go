package recurpost

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrNotImplemented indicates a scaffold method needs implementation.
	ErrNotImplemented = errors.New("not implemented")
)

// APIError represents an HTTP error returned by the remote API.
type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Code != "" {
		return fmt.Sprintf("api error %d %s: %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("api error %d: %s", e.StatusCode, e.Message)
}

func IsNotFound(err error) bool {
	var ae *APIError
	return errors.As(err, &ae) && ae.StatusCode == http.StatusNotFound
}

func IsUnauthorized(err error) bool {
	var ae *APIError
	return errors.As(err, &ae) && (ae.StatusCode == http.StatusUnauthorized || ae.StatusCode == http.StatusForbidden)
}
