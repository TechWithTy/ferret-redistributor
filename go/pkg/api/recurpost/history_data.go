package recurpost

import (
    "bytes"
    "context"
    "encoding/json"
)

// HistoryService exposes /api/history_data
type HistoryService struct{ c *Client }

// Data fetches history data for a user/social account.
func (s *HistoryService) Data(ctx context.Context, in HistoryDataRequest) (*HistoryDataResponse, error) {
    if in.EmailID == "" || in.PassKey == "" {
        return nil, &APIError{StatusCode: 400, Message: "emailid and pass_key are required"}
    }
    payload, err := json.Marshal(in)
    if err != nil {
        return nil, err
    }
    req, err := s.c.newRequest(ctx, "POST", "/api/history_data", bytes.NewReader(payload))
    if err != nil {
        return nil, err
    }
    // Decode into map so we don't constrain schema
    var raw map[string]any
    if err := s.c.do(req, &raw); err != nil {
        return nil, err
    }
    return &HistoryDataResponse{Data: raw}, nil
}

