package postiz

import (
	"context"
	"net/http"
)

type IntegrationsService struct{ c *Client }

func (s *IntegrationsService) List(ctx context.Context) ([]Integration, error) {
	req, err := s.c.newRequest(ctx, http.MethodGet, "/integrations", nil)
	if err != nil {
		return nil, err
	}
	var out []Integration
	if err := s.c.do(req, &out); err != nil {
		return nil, err
	}
	return out, nil
}
