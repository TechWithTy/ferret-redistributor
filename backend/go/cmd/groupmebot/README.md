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
- Notion logging (best-effort, non-blocking):
  - `NOTION_API_KEY` (or `NOTION_TOKEN` / `NOTION_KEY`)
  - `NOTION_DATA_SOURCE_ID_BOT_MESSAGE_LOGS`
  - `NOTION_DATA_SOURCE_ID_BOTS` (to link `Bot` relation)
  - `NOTION_DATA_SOURCE_ID_GROUPS` (to link `Group` relation)

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

## Scheduled posts (Send Now + recurring)

This repo reuses the existing Notion **🚁 Bot Message Logs** database as a lightweight “outbox”:

- To **send now**: create/edit a row with:
  - `Bot` relation set (points to the bot in the Bots DB)
  - `Message Text` set
  - check `Send Now`
- To **schedule recurring**:
  - set `Message Text`
  - set `Bot`
  - check `Schedule Enabled`
  - set `Recurrence` to one of: `Daily`, `Every 7 days`, `Every 2 weeks`, `Monthly`
  - set `Next Run At` to the first date+time you want it to send

### Run the scheduler (for n8n Cron)

Call this endpoint (protected by `GROUPME_WEBHOOK_TOKEN`):

```bash
curl -i -X POST "http://localhost:8081/groupme/schedule/run?token=$GROUPME_WEBHOOK_TOKEN"
```

In n8n: **Cron** node → **HTTP Request** node (POST) to the same URL (your tunnel/public host if needed).

Suggested Cron frequencies:
- `Send Now` usage: every 1–5 minutes
- Recurrence usage: every 1–15 minutes (depending on how strict you need the send time)

