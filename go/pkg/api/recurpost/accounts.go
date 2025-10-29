package recurpost

import "context"

// AccountsService manages connected accounts/profiles.
type AccountsService struct{ c *Client }

// List returns connected accounts.
func (s *AccountsService) List(ctx context.Context) ([]Account, error) {
	// GET /accounts
	return nil, ErrNotImplemented
}

// Get returns a single account by ID.
func (s *AccountsService) Get(ctx context.Context, id string) (*Account, error) {
	// GET /accounts/{id}
	return nil, ErrNotImplemented
}
