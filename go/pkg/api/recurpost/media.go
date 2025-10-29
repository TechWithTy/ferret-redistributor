package recurpost

import "context"

// MediaService manages media endpoints.
type MediaService struct{ c *Client }

// Upload uploads media and returns a media object.
func (s *MediaService) Upload(ctx context.Context, in UploadMediaRequest) (*Media, error) {
	// POST /media
	return nil, ErrNotImplemented
}

// Get retrieves media by ID.
func (s *MediaService) Get(ctx context.Context, id string) (*Media, error) {
	// GET /media/{id}
	return nil, ErrNotImplemented
}
