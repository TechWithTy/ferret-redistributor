# Postiz Public API Client (Scaffold)

This package provides a typed client for the Postiz Public API.

- Base URL (hosted): `https://api.postiz.com/public/v1`
- Auth: Send API key in header `Authorization: {apiKey}`
- Rate limit: 30 requests/hour
- UI term "channel" is called "integration" in the API

Implemented endpoints

- GET `/integrations` — list added channels
- GET `/find-slot/{id}` — next available slot for a channel
- POST `/upload` — upload a file (multipart/form-data)
- POST `/upload-from-url` — upload a file from URL
- GET `/posts` — list posts (startDate, endDate, customer)
- POST `/posts` — create/update posts (type: draft|schedule|now)
- DELETE `/posts/{id}` — delete a post
- POST `/generate-video` — generate videos with AI
- POST `/video/function` — AI utility (e.g., load voices)

Usage

```go
cli := postiz.NewClient(
  postiz.WithBaseURL("https://api.postiz.com/public/v1"),
  postiz.WithAPIKey(os.Getenv("POSTIZ_API_KEY")),
)
ctx := context.Background()
ints, err := cli.Integrations.List(ctx)
```

Tests

- See `pkg/api/postiz/tests` — tests use a local mock server and never call real API.
- Set `POSTIZ_API_KEY` in your root `.env` (auto-loaded) or environment to run tests.

