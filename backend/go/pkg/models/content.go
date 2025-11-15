package models

import (
	"time"

	"github.com/google/uuid"
)

// Content represents a piece of content in the system
type Content struct {
	ID           string    `json:"id"`
	OrgID        string    `json:"org_id"`
	Title        string    `json:	itle"`
	Slug         string    `json:"slug"`
	ContentType  string    `json:"content_type"` // article, video, podcast, etc.
	Status       string    `json:"status"`       // draft, published, archived
	Content      string    `json:"content"`
	Metadata     JSONMap   `json:"metadata"`
	PublishedAt  time.Time `json:"published_at,omitempty"`
	AuthorID     string    `json:"author_id"`
	Featured     bool      `json:"featured"`
	ViewCount    int64     `json:"view_count"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ContentVersion represents a version of a content piece
type ContentVersion struct {
	ID        string    `json:"id"`
	ContentID string    `json:"content_id"`
	Version   int       `json:"version"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  string    `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ContentCollection groups related content
type ContentCollection struct {
	ID          string    `json:"id"`
	OrgID       string    `json:"org_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Slug        string    `json:"slug"`
	Featured    bool      `json:"featured"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CollectionItem represents an item in a content collection
type CollectionItem struct {
	ID           string    `json:"id"`
	CollectionID string    `json:"collection_id"`
	ContentID    string    `json:"content_id"`
	Position     int       `json:"position"`
	CreatedAt    time.Time `json:"created_at"`
}

// NewContent creates a new content piece with default values
func NewContent(orgID, title, contentType, authorID string) *Content {
	now := time.Now().UTC()
	return &Content{
		ID:          uuid.New().String(),
		OrgID:       orgID,
		Title:       title,
		Slug:        generateSlug(title),
		ContentType: contentType,
		Status:      "draft",
		AuthorID:    authorID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// NewVersion creates a new version of the content
func (c *Content) NewVersion(authorID string) *ContentVersion {
	return &ContentVersion{
		ID:        uuid.New().String(),
		ContentID: c.ID,
		Version:   getNextVersionNumber(c.ID),
		Title:     c.Title,
		Content:   c.Content,
		AuthorID:  authorID,
		CreatedAt: time.Now().UTC(),
	}
}

// Publish marks the content as published
func (c *Content) Publish() {
	c.Status = "published"
	if c.PublishedAt.IsZero() {
		c.PublishedAt = time.Now().UTC()
	}
	c.UpdatedAt = time.Now().UTC()
}

// helper functions (implement these based on your needs)
func generateSlug(title string) string {
	// TODO: Implement slug generation
	return "generated-slug-from-title"
}

func getNextVersionNumber(contentID string) int {
	// TODO: Implement version number generation
	return 1
}
