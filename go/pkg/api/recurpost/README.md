# RecurPost API (Scaffold)

This package provides a typed client scaffold for integrating with a RecurPost‑style social scheduling API. It includes:

- A lightweight API spec (OpenAPI skeleton) in `openapi.yaml`
- Consolidated request models in `requests.go`
- Consolidated response models in `responses.go`
- Centralized error definitions in `errors.go`
- Per‑route service files: `auth.go`, `accounts.go`, `posts.go`, `schedules.go`, `media.go`, `analytics.go`

Notes:
- Endpoints and shapes are inferred/generic placeholders. Replace with real fields and paths as needed.
- Methods currently return `ErrNotImplemented`. Wire them to HTTP in `client.go` when the real API is available.

## Routes (Planned)

- Auth: `/oauth/token`, `/oauth/refresh`
- Accounts: `/accounts`, `/accounts/{id}`
- Posts: `/posts`, `/posts/{id}`
- Schedules: `/schedules`, `/schedules/{id}`
- Media: `/media`, `/media/{id}`
- Analytics: `/analytics/posts/{id}`

## Usage (Scaffold)

```go
cli := recurpost.NewClient(
    recurpost.WithBaseURL("https://api.recurpost.example"),
    recurpost.WithToken("<token>"),
)

// Replace with real request types/fields
resp, err := cli.Posts.Create(ctx, recurpost.CreatePostRequest{ 
    Text: "Hello World",
})
if err != nil { /* handle */ }
_ = resp
```

