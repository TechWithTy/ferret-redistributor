package postiz

import "time"

// UploadFromURLRequest represents the payload for upload-from-url.
type UploadFromURLRequest struct {
	URL string `json:"url"`
}

// ListPostsParams configures the posts query.
type ListPostsParams struct {
	StartDate time.Time
	EndDate   time.Time
	Customer  string
}

// PostValue is the per-channel content.
type PostValue struct {
	Content string  `json:"content"`
	ID      string  `json:"id,omitempty"`
	Image   []Media `json:"image,omitempty"`
}

// IntegrationRef references a channel by ID.
type IntegrationRef struct {
	ID string `json:"id"`
}

// PostDraft describes posts submitted for creation/update.
type PostDraft struct {
	Integration IntegrationRef    `json:"integration"`
	Value       []PostValue       `json:"value"`
	Group       string            `json:"group,omitempty"`
	Settings    map[string]any    `json:"settings,omitempty"`
	Tags        []Tag             `json:"tags,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	Extras      map[string]any    `json:"extras,omitempty"`
}

// CreatePostsRequest payload for POST /posts.
type CreatePostsRequest struct {
	Type      string      `json:"type"`
	Order     string      `json:"order,omitempty"`
	ShortLink bool        `json:"shortLink"`
	Interval  *int        `json:"inter,omitempty"`
	Date      time.Time   `json:"date"`
	Tags      []Tag       `json:"tags,omitempty"`
	Posts     []PostDraft `json:"posts"`
}

// VideoRequest represents payloads for AI video generation.
type VideoRequest struct {
	Type         string         `json:"type"`
	Output       string         `json:"output"`
	CustomParams map[string]any `json:"customParams,omitempty"`
}

// VideoFunctionRequest loads server-side video helpers.
type VideoFunctionRequest struct {
	FunctionName string `json:"functionName"`
	Identifier   string `json:"identifier"`
}

