package recurpost

import (
    "bytes"
    "context"
    "encoding/json"
)

// SocialAccountsService exposes /api/social_account_list
type SocialAccountsService struct{ c *Client }

// List returns the connected social accounts for a user (POST with JSON body emailid+pass_key).
func (s *SocialAccountsService) List(ctx context.Context, in SocialAccountListRequest) (*SocialAccountListResponse, error) {
    if in.EmailID == "" || in.PassKey == "" {
        return nil, &APIError{StatusCode: 400, Message: "emailid and pass_key are required"}
    }
    payload, err := json.Marshal(in)
    if err != nil {
        return nil, err
    }
    req, err := s.c.newRequest(ctx, "POST", "/api/social_account_list", bytes.NewReader(payload))
    if err != nil {
        return nil, err
    }
	var out SocialAccountListResponse
	if err := s.c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
