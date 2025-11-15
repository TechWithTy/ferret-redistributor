# RSSHub Go SDK

Typed helper for [RSSHub](https://rsshub.app) so Social Scale services can query the official API (`/api/routes`, `/api/radar/search`, `/api/version`, `/api/force-refresh`) and fetch rendered feeds without re-implementing HTTP plumbing.

## Features

- `GetRoutes`, `GetVersion`, `SearchRadar` map to documented API reference calls
- `FetchFeed` builds full RSS/Atom route URLs (with query params) and returns the raw feed body + content type
- `ForceRefresh` helper encodes target URLs to force upstream cache refreshes
- Option pattern for custom base URLs (`WithBaseURL`), HTTP clients, and default query params (e.g., `?lang=en`)

## Usage

```go
ctx := context.Background()

cli := rsshub.NewClient(
  rsshub.WithBaseURL("https://rsshub.app"),
)

routes, err := cli.GetRoutes(ctx)
if err != nil {
  log.Fatal(err)
}
fmt.Println("total route groups:", len(routes.Data))

feed, err := cli.FetchFeed(ctx, rsshub.FeedRequest{
  Path:  "/bilibili/fav/2262573",
  Query: map[string]string{"limit": "10"},
})
if err != nil {
  log.Fatal(err)
}
fmt.Println("content-type:", feed.ContentType)
fmt.Println("body bytes:", len(feed.Body))
```

Force refresh:

```go
_, err = cli.ForceRefresh(ctx, rsshub.ForceRefreshRequest{
  TargetURL: "https://rsshub.app/bilibili/fav/2262573",
})
```

## Testing

`client_test.go` uses `httptest.Server` to simulate RSSHub endpoints; no external network calls are performed.

## Additional Notes

- See `_docs/guide.md` for feed field compatibility, cache patterns, script standards, date handling, and debugging shortcuts taken from the official RSSHub documentation.


