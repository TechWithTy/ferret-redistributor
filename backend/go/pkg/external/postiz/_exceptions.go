package postiz

import "errors"

var (
	// ErrMissingAPIKey indicates the client configuration lacks an API key.
	ErrMissingAPIKey = errors.New("postiz: API key is required")

	// ErrMissingIntegrationID indicates an integration identifier was not supplied.
	ErrMissingIntegrationID = errors.New("postiz: integration id is required")

	// ErrInvalidDateRange indicates ListPosts was called without valid dates.
	ErrInvalidDateRange = errors.New("postiz: start and end dates are required")

	// ErrEndBeforeStart indicates the end date precedes the start date.
	ErrEndBeforeStart = errors.New("postiz: endDate must be after startDate")

	// ErrMissingPostPayload indicates CreateOrUpdatePosts was called with incomplete data.
	ErrMissingPostPayload = errors.New("postiz: posts payload is incomplete")

	// ErrMissingPostID indicates a delete call without an id.
	ErrMissingPostID = errors.New("postiz: post id is required")

	// ErrMissingFilename indicates UploadFile was called without a filename.
	ErrMissingFilename = errors.New("postiz: filename is required")

	// ErrMissingURL indicates UploadFromURL was called without a URL.
	ErrMissingURL = errors.New("postiz: url is required")

	// ErrMissingVideoType indicates Video generation lacks a type.
	ErrMissingVideoType = errors.New("postiz: video type is required")

	// ErrMissingVideoOutput indicates Video generation lacks an output orientation.
	ErrMissingVideoOutput = errors.New("postiz: video output is required")

	// ErrMissingFunctionName indicates a helper invocation without a function.
	ErrMissingFunctionName = errors.New("postiz: function name is required")

	// ErrMissingIdentifier indicates helper invocation without identifier.
	ErrMissingIdentifier = errors.New("postiz: identifier is required")
)
