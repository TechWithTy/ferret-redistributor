package postiz

import "time"

// Integration represents a Postiz channel (integration).
type Integration struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Identifier string        `json:"identifier"`
	Picture    string        `json:"picture"`
	Disabled   bool          `json:"disabled"`
	Profile    string        `json:"profile"`
	Customer   *CustomerInfo `json:"customer"`
}

// CustomerInfo describes the customer that owns an integration.
type CustomerInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SlotResponse is returned by the find-slot endpoint.
type SlotResponse struct {
	Date time.Time `json:"date"`
}

// FileAsset represents an uploaded file.
type FileAsset struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	Path           string    `json:"path"`
	OrganizationID string    `json:"organizationId"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// PostsResponse returns posts with metadata.
type PostsResponse struct {
	Posts []Post `json:"posts"`
}

// Post represents a scheduled/published item.
type Post struct {
	ID          string            `json:"id"`
	Content     string            `json:"content"`
	PublishDate time.Time         `json:"publishDate"`
	ReleaseURL  string            `json:"releaseURL"`
	State       string            `json:"state"`
	Integration IntegrationDigest `json:"integration"`
}

// IntegrationDigest is a summary inside post payloads.
type IntegrationDigest struct {
	ID                 string `json:"id"`
	ProviderIdentifier string `json:"providerIdentifier"`
	Name               string `json:"name"`
	Picture            string `json:"picture"`
}

// Tag represents a Postiz tag.
type Tag struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// Media represents media attachments inside a post.
type Media struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

// CreatePostsResponse is returned after creating/updating posts.
type CreatePostsResponse struct {
	PostID      string `json:"postId"`
	Integration string `json:"integration"`
}

// DeletePostResponse is returned when deleting a post.
type DeletePostResponse struct {
	ID string `json:"id"`
}

// VideoAsset is returned after generating videos.
type VideoAsset struct {
	ID   string `json:"id"`
	Path string `json:"path"`
}

// VideoFunctionResponse contains helper data such as voices.
type VideoFunctionResponse struct {
	Voices []VideoVoice `json:"voices"`
}

// VideoVoice describes a text-to-speech voice.
type VideoVoice struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
