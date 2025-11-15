# Postiz Go SDK

This package wraps the [Postiz Public API](https://api.postiz.com/public/v1) so Go
services can mirror the workflows currently offered in the Postiz NodeJS SDK and
custom n8n nodes. It supports:

- Listing integrations (channels) and checking the next available slot.
- Uploading media (binary files or existing URLs).
- Listing, creating, scheduling, and deleting posts.
- Generating AI-powered videos and loading helper metadata (voices, etc.).

## Configuration

```bash
export POSTIZ_API_KEY="pk_xxx"
# Optional if you self-host Postiz
export POSTIZ_BASE_URL="https://<NEXT_PUBLIC_BACKEND_URL>/public/v1"
```

```go
cfg := postiz.NewConfigFromEnv()
client, err := postiz.New(cfg)
```

Every request automatically attaches `Authorization: {apiKey}` and respects the
30-requests-per-hour public API threshold noted in the Postiz docs.

## Usage Examples

```go
integrations, _ := client.ListIntegrations(ctx)
slot, _ := client.FindNextSlot(ctx, integrations[0].ID)

media, _ := client.UploadFromURL(ctx, "https://uploads.gitroom.com/example.png")

payload := postiz.CreateOrUpdatePostsRequest{
    Type:      "schedule",
    ShortLink: true,
    Date:      slot,
    Posts: []postiz.PostDraft{
        {
            Integration: postiz.IntegrationRef{ID: integrations[0].ID},
            Value: []postiz.PostValue{
                {Content: "Hello Postiz!", Image: []postiz.Media{{ID: media.ID, Path: media.Path}}},
            },
        },
    },
}
client.CreateOrUpdatePosts(ctx, payload)
```

See `types.go` for the rest of the supported payloads and helpers.

## Testing

Integration-style tests live under `_tests` to mirror the structure used by the
other Social Scale SDKs. Run them with:

```bash
cd go
go test ./pkg/external/postiz/_tests
```

