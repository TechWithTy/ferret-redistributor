# Postiz Automation Playbook

This guide mirrors the Playbook format used by the other Social Scale API SDKs and
illustrates how to integrate the Postiz Public API with GitHub Actions, n8n, or
custom Go workers.

## 1. Capabilities

- List integrations (channels) and their customers.
- Query the next publishing slot for a given integration.
- Upload media from a buffer or an existing URL.
- List, create, update, or delete posts (draft, schedule, now).
- Generate AI videos via the Postiz modal API or fetch helper metadata
  (available voices, formats, etc.).

## 2. Environment Variables

| Variable            | Description                                          |
| ------------------- | ---------------------------------------------------- |
| `POSTIZ_API_KEY`      | Public API key copied from Postiz settings.        |
| `POSTIZ_BASE_URL`     | Optional – set when self-hosting Postiz.           |
| `POSTIZ_INTEGRATIONS` | Optional JSON map like `{"linkedin":"integration"}` |

Add them to `.env` locally and to GitHub Secrets/Variables when using Actions.

## 3. GitHub Actions Example

```yaml
name: Postiz Scheduler
on:
  schedule: [{ cron: "*/15 * * * *" }]
  workflow_dispatch:

jobs:
  publish:
    runs-on: ubuntu-latest
    env:
      POSTIZ_API_KEY: ${{ secrets.POSTIZ_API_KEY }}
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.21.x"
      - run: go test ./go/pkg/external/postiz/...
      - run: go run ./go/cmd/postizpublisher --input go/due_posts.json
```

`cmd/postizpublisher` can call the SDK to pull queued Social Scale posts, find the
next slot, and publish in batch while respecting the 30-requests-per-hour
limit.

## 4. n8n / NodeJS Nodes

- Use the same API key and base URL.
- Inputs map 1:1 to the OpenAPI spec at `go/pkg/external/postiz/openapi.yaml`.
- Remember that the UI term “channel” is named `integration` in the API.

## 5. Self-hosted Tips

- Base URL: `https://{NEXT_PUBLIC_BACKEND_URL}/public/v1`
- Access the modal generator locally via `http://localhost:5000/modal/dark/all`
  when building AI video presets.
- If you proxy requests through Social Scale, ensure the `Authorization` header is
  passed through untouched.

## 6. Error Handling & Retries

- 401/403 -> check API key or workspace permissions.
- 429 -> respect the documented rate limit; use exponential backoff.
- 5xx -> transient; retry with jitter.
- SDK maps validation issues to sentinel errors (see `_exceptions.go`) before
  hitting the network.


