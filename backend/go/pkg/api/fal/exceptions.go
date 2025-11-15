package fal

import "errors"

var (
	// ErrEmptyPath indicates the caller did not provide a Fal path.
	ErrEmptyPath = errors.New("fal: path is required")
	// ErrMissingCredentials indicates no API key was configured on the client.
	ErrMissingCredentials = errors.New("fal: api key is required")
	// ErrCallIncomplete is returned when the workflow fails to reach COMPLETED.
	ErrCallIncomplete = errors.New("fal: request did not complete successfully")
)
