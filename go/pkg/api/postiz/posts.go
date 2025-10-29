package postiz

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type PostsService struct{ c *Client }

func (s *PostsService) List(ctx context.Context, in PostsListRequest) (*PostsListResponse, error) {
	// Build query
	path := "/posts?startDate=" + in.StartDate + "&endDate=" + in.EndDate
	if in.Customer != "" {
		path += "&customer=" + in.Customer
	}
	req, err := s.c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	var out PostsListResponse
	if err := s.c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (s *PostsService) CreateOrUpdate(ctx context.Context, in CreateUpdatePostRequest) ([]CreateUpdateResult, error) {
	b, _ := json.Marshal(in)
	req, err := s.c.newRequest(ctx, http.MethodPost, "/posts", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var out []CreateUpdateResult
	if err := s.c.do(req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *PostsService) Delete(ctx context.Context, id string) (*DeletePostResponse, error) {
	req, err := s.c.newRequest(ctx, http.MethodDelete, "/posts/"+id, nil)
	if err != nil {
		return nil, err
	}
	var out DeletePostResponse
	if err := s.c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
