# Postiz Publisher CLI

Moves claimed scheduler rows into Postiz queues so the hosted scheduler can post
on our behalf. It only touches rows mapped to a Postiz integration and rewrites
`go/due_posts.json` with the remaining items so the legacy `poster` CLI can
handle everything else.

## Usage

```bash
POSTIZ_API_KEY=pk_xxx \
POSTIZ_INTEGRATIONS='{"linkedin":"cm4ean69r0003w8w1cdomox9n"}' \
DATABASE_URL=postgres://... \
go run ./cmd/postizpublisher --input go/due_posts.json
```

### Environment

| Variable              | Required | Description                                      |
| --------------------- | -------- | ------------------------------------------------ |
| `POSTIZ_API_KEY`      | Yes      | Public API key copied from Postiz settings.      |
| `POSTIZ_BASE_URL`     | No       | Override for self-hosted Postiz deployments.     |
| `POSTIZ_INTEGRATIONS` | No       | JSON map of `platform -> integration ID`.        |

You can also provide `--integration-map path/to/map.json` instead of
`POSTIZ_INTEGRATIONS`. Keys are case-insensitive platform names (`linkedin`,
`facebook`, etc.). Set `{"*": "<integration-id>"}` to declare a default.

### Behavior

- Skips rows without an integration mapping (they remain in `due_posts.json`).
- Uploads media via `UploadFromURL` when `media_urls`/`image_urls` metadata is present.
- Marks records as `published` locally with Postiz metadata and external IDs.
- Rewrites the input JSON (unless `--keep-input` is set) so downstream CLIs only see unprocessed rows.



