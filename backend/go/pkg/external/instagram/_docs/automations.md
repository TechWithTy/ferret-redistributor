# Instagram Posting and Webhook Automations (with GitHub Actions)

This guide explains how to:
- Post to Instagram (feed, reels, carousel, stories) via Social Scale
- Set up Instagram Graph Webhooks to react to comments/mentions/messages
- Wire environment in GitHub Actions so automation runs on schedule or on push

> Note: Webhooks require a publicly reachable HTTPS URL. GitHub Actions cannot directly receive webhooks; use a hosted service (Render, Fly.io, Railway, Cloud Run, etc.) for your webhook endpoint, and use Actions to build/deploy it and manage secrets.

## 1) Prerequisites

- Instagram Professional account connected to a Facebook Page
- Meta app with required permissions:
  - instagram_basic, pages_show_list, pages_read_engagement, instagram_content_publish
  - For Messaging webhooks: `instagram_manage_messages` and enable IG Messaging for the Page
- Long‑lived user access token and IG User ID

## 2) Environment Variables

Add these to your `.env` locally and to GitHub Secrets/Variables for workflows:

- IG_ACCESS_TOKEN
- IG_USER_ID
- IG_GRAPH_VERSION (optional, default v19.0)
- IG_FIRST_COMMENT_TEXT (optional)
- IG_TRIGGER_WORDS (optional, e.g. `READY,TEMPLATE`)
- IG_MESSAGING_ENABLED (optional, set `true` when configured)
- IG_APP_SECRET (for webhook signature validation)
- IG_VERIFY_TOKEN (for webhook subscription verification)

## 3) Enable Posting in Social Scale

1. Add `"instagram"` to your `config.json` `socials` list.
2. Ensure env vars from step 2 are set.
3. Social Scale calls `external.Instagram` (pkg/external/instagram.go) which:
   - Resolves OG image from your article link
   - Posts a feed image with caption `Title + HashTags`
   - Optionally posts a first comment using `IG_FIRST_COMMENT_TEXT`

### Reels, Carousel, Stories (Programmatic)

Use the client directly:

```go
cfg := instagram.NewFromEnv()
cli := instagram.New(cfg)
id, err := cli.PostReel(ctx, videoURL, caption) // polls processing then publishes
```

## 4) Webhooks Setup (Comments/Mentions/Messages)

Social Scale includes a webhook handler: `pkg/external/instagram/webhooks.go`.

It supports:
- GET verification (hub.challenge)
- POST deliveries with `X-Hub-Signature-256` validation
- Routing for fields: `comments`, `mentions`, `messages`, `likes`

### Minimal Server Example

```go
package main

import (
  "net/http"
  "os"
  ig "github.com/bitesinbyte/ferret/pkg/external/instagram"
  wk "github.com/bitesinbyte/ferret/pkg/external/instagram/wokers"
  gen "github.com/bitesinbyte/ferret/pkg/ai/generator"
)

func main() {
  handler := ig.WebhookHandler{
    AppSecret:   os.Getenv("IG_APP_SECRET"),
    VerifyToken: os.Getenv("IG_VERIFY_TOKEN"),
    OnComment: func(ctx context.Context, c ig.CommentChange) error {
      // Trigger → DM
      cli := ig.New(ig.NewFromEnv())
      worker := wk.DMWorker{IG: cli, Matcher: wk.TriggerMatcher{Triggers: splitEnv("IG_TRIGGER_WORDS")}, Generator: gen.NoopGenerator{}}
      // Process only the delivered comment: quick path
      if _, ok := worker.Matcher.Match(c.Text); ok {
        // Attempt DM (requires IG Messaging perms), else no‑op if unsupported
        _ = cli.SendDM(ctx, c.FromID, "Thanks! Check your inbox soon.")
      }
      return nil
    },
  }
  http.Handle("/webhooks/instagram", handler)
  _ = http.ListenAndServe(":8080", nil)
}
```

> Host this service behind HTTPS and use the public URL in your Meta app’s Webhooks configuration for the `instagram` object.

### Configure Meta App

1. In Meta App Dashboard → Add Product → Webhooks.
2. Subscribe to `instagram` object; fields: `comments`, `mentions`, `messages`.
3. Set Callback URL to `https://your-host/webhooks/instagram` and Verify Token to `IG_VERIFY_TOKEN`.
4. Save; Meta will perform GET verification (challenge response is handled).

## 5) GitHub Actions Automation

Use Actions to:
- Build/Deploy your webhook server (container or binary) to your hosting platform
- Run posting on a schedule (e.g., to process RSS posts)

### Example Workflow Snippet (Posting)

```yaml
name: Social Post
on:
  schedule: [{ cron: '*/30 * * * *' }]
  workflow_dispatch:

jobs:
  post:
    runs-on: ubuntu-latest
    env:
      IG_ACCESS_TOKEN: ${{ secrets.IG_ACCESS_TOKEN }}
      IG_USER_ID: ${{ secrets.IG_USER_ID }}
      IG_GRAPH_VERSION: v19.0
      IG_FIRST_COMMENT_TEXT: ${{ vars.IG_FIRST_COMMENT_TEXT }}
      IG_TRIGGER_WORDS: ${{ vars.IG_TRIGGER_WORDS }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.21.x' }
      - run: go build -o bin/ferret ./cmd/ferret
      - run: ./bin/ferret
```

### Example Workflow Snippet (Build & Push Container)

```yaml
name: Deploy Webhook Server
on:
  push:
    paths:
      - 'cmd/webhooks/**'
      - 'pkg/external/instagram/webhooks.go'

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: docker/setup-buildx-action@v3
      - uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: ghcr.io/${{ github.repository_owner }}/social-scale-webhooks:latest
```

> Deploy the built image to your hosting (Render/Fly/Railway/Cloud Run). Store `IG_APP_SECRET` and `IG_VERIFY_TOKEN` as environment variables there.

## 6) Security & Tips

- Always validate `X-Hub-Signature-256` using `IG_APP_SECRET` (implemented in handler).
- Rate limits: the client maps 429 to `ErrRateLimited`. Consider retry/backoff.
- Messaging requires additional review and setup; until enabled, DM sending returns `ErrUnsupportedFeature`.

