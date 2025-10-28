package youtube

import "errors"

var (
    ErrUnauthorized = errors.New("youtube: unauthorized or invalid credentials")
    ErrForbidden    = errors.New("youtube: insufficient permissions or quota")
    ErrRateLimited  = errors.New("youtube: rate limited")
    ErrValidation   = errors.New("youtube: validation error")
    ErrNotFound     = errors.New("youtube: resource not found")
    ErrServer       = errors.New("youtube: server error")
    ErrUnsupported  = errors.New("youtube: feature unsupported (requires OAuth/partner access)")
)

