# GroupMe Bot (Local Harness)

This is a minimal local harness for testing GroupMe bot callbacks and bot posting **without Docker**.

## Requirements

- Go (see [`go/docs/deploy.md`](../../docs/deploy.md))
- A GroupMe bot created in your target group (you’ll need the bot id)
- A tunnel so GroupMe can reach your machine (ngrok or cloudflared)

## Environment

Create a local `.env` at the repository root (or export env vars in your shell).

Required:
- `GROUPME_WEBHOOK_TOKEN`: shared secret for inbound webhook auth (query `?token=` or header `X-Webhook-Token`)
- `GROUPME_BOT_ID`: used to post messages back to GroupMe

Optional:
- `GROUPME_PORT`: default `8081`
- `GROUPME_BASE_URL`: default `https://api.groupme.com/v3`

See [`env.example`](../../../env.example) (rename to `.env` for local use).

## Run

From `backend/go`:

```bash
go run ./cmd/groupmebot
```

## Tunnel and callback URL

- Start the harness locally on port 8081.
- Expose it publicly:
  - `ngrok http 8081`
  - or `cloudflared tunnel --url http://localhost:8081`
- Set your GroupMe bot callback URL to:
  - `https://<tunnel-host>/webhooks/groupme?token=<GROUPME_WEBHOOK_TOKEN>`

## Smoke tests

### Synthetic webhook (local)

```bash
curl -i -X POST "http://localhost:8081/webhooks/groupme?token=$GROUPME_WEBHOOK_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"id":"m1","group_id":"g1","sender_id":"u1","sender_type":"user","name":"Test","text":"!ping","system":false,"created_at":0}'
```

### Manual outbound post (real GroupMe)

```bash
curl -i -X POST "http://localhost:8081/groupme/send?token=$GROUPME_WEBHOOK_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"text":"hello from local harness"}'
```

## Notes

- This harness intentionally uses a **shared secret** because GroupMe bot callbacks are not consistently documented as providing a verifiable signature header.
- The webhook handler ignores `sender_type=bot` (and `system=true`) to avoid reply loops.

