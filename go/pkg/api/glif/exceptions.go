package glif

import (
	"errors"
	"fmt"
)

var (
	// ErrMissingWorkflowID is returned when no workflow ID is supplied.
	ErrMissingWorkflowID = errors.New("glif: workflow id is required")
	// ErrConflictingInputModes is returned when both positional and named inputs are provided together.
	ErrConflictingInputModes = errors.New("glif: provide either positional OR named inputs, not both")
	// ErrWorkflowFailed indicates the workflow responded with an error payload.
	ErrWorkflowFailed = errors.New("glif: workflow returned an error")
	// ErrWorkflowNotFound signals that the requested workflow does not exist or is not public.
	ErrWorkflowNotFound = errors.New("glif: workflow not found")
)

// APIError represents non-2xx HTTP responses from the Glif API.
type APIError struct {
	StatusCode int
	Body       []byte
}

func (e *APIError) Error() string {
	return fmt.Sprintf("glif: api error status=%d body=%s", e.StatusCode, string(e.Body))
}

