package recurpost

import (
	"context"
)

// AuthService handles OAuth/token flows.
type AuthService struct{ c *Client }

// Exchange exchanges credentials for an access token.
func (s *AuthService) Exchange(ctx context.Context, req TokenRequest) (*TokenResponse, error) {
	// Wire HTTP call to /oauth/token here
	return nil, ErrNotImplemented
}

// Refresh refreshes the access token.
func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	// Wire HTTP call to /oauth/token with refresh_token grant here
	return nil, ErrNotImplemented
}
