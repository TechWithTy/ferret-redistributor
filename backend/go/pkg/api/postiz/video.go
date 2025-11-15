package postiz

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type VideoService struct{ c *Client }

func (s *VideoService) Generate(ctx context.Context, in GenerateVideoRequest) ([]GeneratedVideo, error) {
	b, _ := json.Marshal(in)
	req, err := s.c.newRequest(ctx, http.MethodPost, "/generate-video", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var out []GeneratedVideo
	if err := s.c.do(req, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (s *VideoService) Function(ctx context.Context, in VideoFunctionRequest) (*VideoFunctionResponse, error) {
	b, _ := json.Marshal(in)
	req, err := s.c.newRequest(ctx, http.MethodPost, "/video/function", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	var out VideoFunctionResponse
	if err := s.c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
