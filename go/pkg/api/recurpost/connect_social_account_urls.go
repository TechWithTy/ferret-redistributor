package recurpost

import (
	"bytes"
	"context"
	"encoding/json"
)

// SocialConnectService exposes /api/connect_social_account_urls
type SocialConnectService struct{ c *Client }

// GetURLs posts emailid and pass_key and returns provider connect URLs.
func (s *SocialConnectService) GetURLs(ctx context.Context, in ConnectSocialAccountURLsRequest) (*ConnectURLsResponse, error) {
	if in.EmailID == "" || in.PassKey == "" {
		return nil, &APIError{StatusCode: 400, Message: "emailid and pass_key are required"}
	}
	payload, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	req, err := s.c.newRequest(ctx, "POST", "/api/connect_social_account_urls", bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	var out ConnectURLsResponse
	if err := s.c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
