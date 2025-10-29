package recurpost

import (
    "context"
    "net/url"
)

// SocialAccountsService exposes /api/social_account_list
type SocialAccountsService struct{ c *Client }

// List returns the connected social accounts for a user (POST with query param emailid).
func (s *SocialAccountsService) List(ctx context.Context, in SocialAccountListRequest) (*SocialAccountListResponse, error) {
    if in.EmailID == "" {
        return nil, &APIError{StatusCode: 400, Message: "emailid is required"}
    }
    q := url.Values{}
    q.Set("emailid", in.EmailID)
    p := "/api/social_account_list?" + q.Encode()
    req, err := s.c.newRequest(ctx, "POST", p, nil)
    if err != nil {
        return nil, err
    }
    var out SocialAccountListResponse
    if err := s.c.do(req, &out); err != nil {
        return nil, err
    }
    return &out, nil
}

