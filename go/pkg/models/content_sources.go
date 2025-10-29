package models

import "time"

// ContentSourceType defines the type of content source
type ContentSourceType string

const (
	ContentSourceTypeRSS          ContentSourceType = "rss"
	ContentSourceTypeYouTube      ContentSourceType = "youtube"
	ContentSourceTypeYouTubeVideo ContentSourceType = "youtube_video"
	ContentSourceTypeYouTubePlaylist ContentSourceType = "youtube_playlist"
)

// ContentSource represents a source of content like an RSS feed or YouTube channel
type ContentSource struct {
	ID          string          `json:"id"`
	OrgID       string          `json:"org_id"`
	TeamID      *string         `json:"team_id,omitempty"`
	Name        string          `json:
ame"`
	Type        ContentSourceType `json:"type"`
	URL         string          `json:"url"`
	Description string          `json:"description"`
	IsActive    bool            `json:"is_active"`
	LastFetched *time.Time      `json:"last_fetched,omitempty"`
	Metadata    map[string]any  `json:"metadata"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// ContentSourceItem represents an item from a content source (e.g., an RSS feed item or YouTube video)
type ContentSourceItem struct {
	ID              string         `json:"id"`
	SourceID        string         `json:"source_id"`
	OrgID           string         `json:"org_id"`
	ExternalID      string         `json:"external_id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	URL             string         `json:"url"`
	ThumbnailURL    *string        `json:"thumbnail_url,omitempty"`
	PublishedAt     *time.Time     `json:"published_at,omitempty"`
	Author          *string        `json:"author,omitempty"`
	DurationSeconds *int           `json:"duration_seconds,omitempty"` // For videos
	ViewCount       *int64         `json:"view_count,omitempty"`       // For videos
	Metadata        map[string]any `json:"metadata"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// YouTubePlaylist represents a YouTube playlist
type YouTubePlaylist struct {
	ID              string         `json:"id"`
	OrgID           string         `json:"org_id"`
	SourceID        string         `json:"source_id"` // Reference to ContentSource
	ExternalID      string         `json:"external_id"`
	Title           string         `json:"title"`
	Description     string         `json:"description"`
	ThumbnailURL    *string        `json:"thumbnail_url,omitempty"`
	ItemCount       int            `json:"item_count"`
	IsActive        bool           `json:"is_active"`
	LastSyncedAt    *time.Time     `json:"last_synced_at,omitempty"`
	Metadata        map[string]any `json:"metadata"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// ContentSourceSync represents a sync operation for a content source
type ContentSourceSync struct {
	ID            string     `json:"id"`
	SourceID      string     `json:"source_id"`
	OrgID         string     `json:"org_id"`
	Status        string     `json:"status"` // "pending", "in_progress", "completed", "failed"
	ItemsFetched  int        `json:"items_fetched"`
	ItemsImported int        `json:"items_imported"`
	Error         *string    `json:"error,omitempty"`
	StartedAt     *time.Time `json:"started_at,omitempty"`
	CompletedAt   *time.Time `json:"completed_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
