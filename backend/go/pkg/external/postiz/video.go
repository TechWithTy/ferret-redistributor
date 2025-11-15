package postiz

import (
	"context"
	"net/http"
)

// GenerateVideo triggers AI video creation.
func (c *Client) GenerateVideo(ctx context.Context, payload VideoRequest) ([]VideoAsset, error) {
	if payload.Type == "" {
		return nil, ErrMissingVideoType
	}
	if payload.Output == "" {
		return nil, ErrMissingVideoOutput
	}
	var res []VideoAsset
	if err := c.doJSON(ctx, http.MethodPost, "/generate-video", payload, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// LoadVideoVoices fetches helper data such as available voices.
func (c *Client) LoadVideoVoices(ctx context.Context, payload VideoFunctionRequest) (VideoFunctionResponse, error) {
	if payload.FunctionName == "" {
		return VideoFunctionResponse{}, ErrMissingFunctionName
	}
	if payload.Identifier == "" {
		return VideoFunctionResponse{}, ErrMissingIdentifier
	}
	var res VideoFunctionResponse
	if err := c.doJSON(ctx, http.MethodPost, "/video/function", payload, &res); err != nil {
		return VideoFunctionResponse{}, err
	}
	return res, nil
}
