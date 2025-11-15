package instagram

import "errors"

var (
    ErrUnauthorized       = errors.New("instagram: unauthorized or invalid token")
    ErrForbidden          = errors.New("instagram: insufficient permissions")
    ErrRateLimited        = errors.New("instagram: rate limited")
    ErrValidation         = errors.New("instagram: validation error")
    ErrNotFound           = errors.New("instagram: resource not found")
    ErrServer             = errors.New("instagram: server error")
    ErrUnsupportedFeature = errors.New("instagram: feature unsupported in current API/permissions")
)
