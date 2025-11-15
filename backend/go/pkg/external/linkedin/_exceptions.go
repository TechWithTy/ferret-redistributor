package linkedin

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
)

var (
    ErrUnauthorized = errors.New("linkedin: unauthorized or invalid token")
    ErrForbidden    = errors.New("linkedin: insufficient permissions or scope")
    ErrRateLimited  = errors.New("linkedin: rate limited")
    ErrValidation   = errors.New("linkedin: validation error")
    ErrNotFound     = errors.New("linkedin: resource not found")
    ErrServer       = errors.New("linkedin: server error")
)

// APIError captures LinkedIn REST error payload fields when available.
type APIError struct {
    Message          string `json:"message"`
    Status           int    `json:"status"`
    Code             string `json:"code"`
    ServiceErrorCode int    `json:"serviceErrorCode"`
    Raw              string `json:"-"`
}

func (e APIError) Error() string {
    if e.Message != "" {
        return fmt.Sprintf("linkedin: %s (code=%s status=%d service=%d)", e.Message, e.Code, e.Status, e.ServiceErrorCode)
    }
    if e.Raw != "" {
        return "linkedin: " + e.Raw
    }
    return "linkedin: unknown error"
}

// MapHTTPError maps an HTTP response to a sentinel error, preserving APIError details when present.
func MapHTTPError(resp *http.Response) error {
    if resp == nil {
        return ErrServer
    }
    status := resp.StatusCode
    body, _ := io.ReadAll(resp.Body)
    // Best-effort unmarshal
    var env APIError
    _ = json.Unmarshal(body, &env)
    env.Status = status
    env.Raw = string(body)

    switch status {
    case http.StatusUnauthorized:
        return ErrUnauthorized
    case http.StatusForbidden:
        return ErrForbidden
    case http.StatusTooManyRequests:
        return ErrRateLimited
    case http.StatusNotFound:
        return ErrNotFound
    }

    if status >= 500 {
        return ErrServer
    }
    if status >= 400 {
        // Provide detailed APIError when available
        if env.Message != "" || env.Code != "" || env.ServiceErrorCode != 0 {
            return env
        }
        return ErrValidation
    }
    return nil
}
