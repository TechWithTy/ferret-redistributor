package recurpost

import (
    "bytes"
    "context"
    "encoding/json"
)

// UserLoginService exposes the /api/user_login route.
type UserLoginService struct{ c *Client }

// Login performs a POST /api/user_login?emailid=...
func (s *UserLoginService) Login(ctx context.Context, req UserLoginRequest) (*UserLoginResponse, error) {
    if req.EmailID == "" || req.PassKey == "" {
        return nil, &APIError{StatusCode: 400, Message: "emailid and pass_key are required"}
    }
    payload, err := json.Marshal(req)
    if err != nil {
        return nil, err
    }
    httpReq, err := s.c.newRequest(ctx, "POST", "/api/user_login", bytes.NewReader(payload))
    if err != nil {
        return nil, err
    }
    var out UserLoginResponse
    if err := s.c.do(httpReq, &out); err != nil {
        return nil, err
    }
    return &out, nil
}
