package postiz

import (
	"context"
	"net/http"
)

type SlotsService struct{ c *Client }

func (s *SlotsService) Find(ctx context.Context, id string) (*FindSlotResponse, error) {
	req, err := s.c.newRequest(ctx, http.MethodGet, "/find-slot/"+id, nil)
	if err != nil {
		return nil, err
	}
	var out FindSlotResponse
	if err := s.c.do(req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
