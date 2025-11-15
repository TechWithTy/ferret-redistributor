package postiz

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// ListPosts fetches posts between the provided dates.
func (c *Client) ListPosts(ctx context.Context, params ListPostsParams) (PostsResponse, error) {
	if params.StartDate.IsZero() || params.EndDate.IsZero() {
		return PostsResponse{}, ErrInvalidDateRange
	}
	if params.EndDate.Before(params.StartDate) {
		return PostsResponse{}, ErrEndBeforeStart
	}
	values := url.Values{}
	values.Set("startDate", params.StartDate.UTC().Format(time.RFC3339))
	values.Set("endDate", params.EndDate.UTC().Format(time.RFC3339))
	if params.Customer != "" {
		values.Set("customer", params.Customer)
	}

	path := "/posts?" + values.Encode()
	var res PostsResponse
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &res); err != nil {
		return PostsResponse{}, err
	}
	return res, nil
}

// CreateOrUpdatePosts creates, schedules, or updates posts.
func (c *Client) CreateOrUpdatePosts(ctx context.Context, payload CreatePostsRequest) ([]CreatePostsResponse, error) {
	if payload.Type == "" {
		return nil, ErrMissingPostPayload
	}
	if payload.Date.IsZero() {
		return nil, ErrMissingPostPayload
	}
	if len(payload.Posts) == 0 {
		return nil, ErrMissingPostPayload
	}
	var res []CreatePostsResponse
	if err := c.doJSON(ctx, http.MethodPost, "/posts", payload, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// DeletePost removes a post by ID.
func (c *Client) DeletePost(ctx context.Context, id string) (DeletePostResponse, error) {
	if id == "" {
		return DeletePostResponse{}, ErrMissingPostID
	}
	path := fmt.Sprintf("/posts/%s", id)
	var res DeletePostResponse
	if err := c.doJSON(ctx, http.MethodDelete, path, nil, &res); err != nil {
		return DeletePostResponse{}, err
	}
	return res, nil
}
