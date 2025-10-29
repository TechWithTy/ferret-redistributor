package recurpost

import (
	"context"
)

// PostsService manages post endpoints.
type PostsService struct{ c *Client }

// Create creates a post.
func (s *PostsService) Create(ctx context.Context, in CreatePostRequest) (*Post, error) {
	// POST /posts
	return nil, ErrNotImplemented
}

// Update updates a post by ID.
func (s *PostsService) Update(ctx context.Context, id string, in UpdatePostRequest) (*Post, error) {
	// PATCH /posts/{id}
	return nil, ErrNotImplemented
}

// Get retrieves a post.
func (s *PostsService) Get(ctx context.Context, id string) (*Post, error) {
	// GET /posts/{id}
	return nil, ErrNotImplemented
}

// Delete removes a post.
func (s *PostsService) Delete(ctx context.Context, id string) error {
	// DELETE /posts/{id}
	return ErrNotImplemented
}

// List returns a paginated list of posts.
func (s *PostsService) List(ctx context.Context, in ListPostsRequest) (*PostList, error) {
	// GET /posts
	return nil, ErrNotImplemented
}
