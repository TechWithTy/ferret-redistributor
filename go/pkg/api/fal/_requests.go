package fal

// RequestOptions configures how the initial Fal request is sent.
type RequestOptions struct {
	// Method defaults to POST because the queue endpoints generally require it.
	Method string
	// Payload is JSON-encoded and sent as the request body.
	Payload any
}


