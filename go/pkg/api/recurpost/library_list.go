package recurpost

import (
	"bytes"
	"context"
	"encoding/json"
)

// LibrariesService exposes /api/library_list
type LibrariesService struct{ c *Client }

// List returns list of libraries for a user using emailid and pass_key in body.
func (s *LibrariesService) List(ctx context.Context, in LibraryListRequest) (*LibraryListResponse, error) {
	if in.EmailID == "" || in.PassKey == "" {
		return nil, &APIError{StatusCode: 400, Message: "emailid and pass_key are required"}
	}
	payload, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	req, err := s.c.newRequest(ctx, "POST", "/api/library_list", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	var out LibraryListResponse
	if err := s.c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddContent adds a post/content to a library.
func (s *LibrariesService) AddContent(ctx context.Context, in AddContentInLibraryRequest) (*AddContentInLibraryResponse, error) {
	if in.EmailID == "" || in.PassKey == "" || in.ID == "" || in.Message == "" {
		return nil, &APIError{StatusCode: 400, Message: "emailid, pass_key, id and message are required"}
	}
	payload, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	req, err := s.c.newRequest(ctx, "POST", "/api/add_content_in_library", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	var out AddContentInLibraryResponse
	if err := s.c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
