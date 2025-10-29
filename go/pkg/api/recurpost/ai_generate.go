package recurpost

import (
	"bytes"
	"context"
	"encoding/json"
)

// AIService exposes AI endpoints.
type AIService struct{ c *Client }

// GenerateContent generates text content with AI.
func (s *AIService) GenerateContent(ctx context.Context, in GenerateContentWithAIRequest) (*GenerateContentResponse, error) {
	if in.EmailID == "" || in.PassKey == "" || in.PromptText == "" {
		return nil, &APIError{StatusCode: 400, Message: "emailid, pass_key and prompt_text are required"}
	}
	payload, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	req, err := s.c.newRequest(ctx, "POST", "/api/generate_content_with_ai", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	var out GenerateContentResponse
	if err := s.c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GenerateImage generates an image with AI.
func (s *AIService) GenerateImage(ctx context.Context, in GenerateImageWithAIRequest) (*GenerateImageResponse, error) {
	if in.EmailID == "" || in.PassKey == "" || in.PromptText == "" {
		return nil, &APIError{StatusCode: 400, Message: "emailid, pass_key and prompt_text are required"}
	}
	payload, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	req, err := s.c.newRequest(ctx, "POST", "/api/generate_image_with_ai", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	var out GenerateImageResponse
	if err := s.c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
