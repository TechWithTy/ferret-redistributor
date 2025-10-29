package recurpost

import (
	"bytes"
	"context"
	"encoding/json"
)

// PublishingService exposes /api/post_content
type PublishingService struct{ c *Client }

// Post publishes content directly or schedules it.
func (s *PublishingService) Post(ctx context.Context, in PostContentRequest) (*PostContentResponse, error) {
	if in.EmailID == "" || in.PassKey == "" || in.ID == "" || in.Message == "" {
		return nil, &APIError{StatusCode: 400, Message: "emailid, pass_key, id and message are required"}
	}
	payload, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	req, err := s.c.newRequest(ctx, "POST", "/api/post_content", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	var out PostContentResponse
	if err := s.c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
