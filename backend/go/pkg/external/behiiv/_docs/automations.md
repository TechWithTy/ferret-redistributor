# LinkedIn Posting Automations (with GitHub Actions)

This guide shows how to post to LinkedIn using Social Scale and how to wire automation in GitHub Actions.

- Client: `pkg/external/linkedin.go`
- Guidelines schema: `pkg/external/linkedin/_docs/guidlines.json`

## Prerequisites

- LinkedIn Developer App with `w_member_social` scope approved
- Member access token for the target account (the member that will post)
- Your app allowed to create posts for that member

## Environment

- `LINKEDIN_ACCESS_TOKEN` – Bearer token for the posting member
- Optional (from `config.json`):
  - `base_url`, `does_meta_og_image_has_relative_path` for OG image resolution

Update `config.json` to include LinkedIn:

```json
{
  "socials": [
    "facebook",
    "linkedin",
    "mastodon",
    "twitter",
    "thread"
  ]
}
```

## How Posting Works

Social Scale’s LinkedIn integration:
- Fetches the OG image from your article (handles relative URLs if configured)
- Initializes an image upload with LinkedIn REST Images API and uploads the bytes
- Creates a post with commentary, landing page URL, and the uploaded image as thumbnail

See: `pkg/external/linkedin.go` for the full flow and `getOGImageURL` for image discovery.

## GitHub Actions: Scheduled Posting

Use a scheduled workflow to run Social Scale periodically and post new feed items.

```yaml
name: Social Post
on:
  schedule: [{ cron: '*/30 * * * *' }]
  workflow_dispatch:

jobs:
  post:
    runs-on: ubuntu-latest
    env:
      LINKEDIN_ACCESS_TOKEN: ${{ secrets.LINKEDIN_ACCESS_TOKEN }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.21.x' }
      - run: go build -o bin/ferret ./cmd/ferret
      - run: ./bin/ferret
```

Store `LINKEDIN_ACCESS_TOKEN` in repository secrets.

## Content Guidelines

Use `pkg/external/linkedin/_docs/guidlines.json` as your reference for:
- Formats: post, carousel, article, reel
- Visual and caption guidelines (hook, structure, CTA)
- Hashtag strategy and posting cadence

## Interactions and Webhooks

- LinkedIn does not provide general-purpose public webhooks for post comments to third-party apps. For engagement workflows, consider:
  - Periodic analytics checks via LinkedIn APIs (not implemented here)
  - Manual moderation/notification tooling
- If you need automated replies, design a separate job to poll recent posts and comments (ensure API access and compliance).

## Troubleshooting

- 401/403 errors: verify token scope (`w_member_social`) and that the token belongs to the posting member
- Image upload failing: confirm OG image URL is publicly accessible and valid (PNG/JPEG)
- Post creation 4xx: ensure `contentLandingPage` and `article` fields meet LinkedIn API requirements

## Security

- Keep `LINKEDIN_ACCESS_TOKEN` in GitHub Secrets
- Rotate tokens periodically and avoid logging request/response bodies with tokens

