# GroupMe → Notion Sync (Groups + Bots)

This command pulls **Groups** and **Bots** from GroupMe using a user access token, then upserts them into your Notion **Groups** and **Bots** data sources.

It uses `data_source_id` and sends `Notion-Version: 2025-09-03` (multi-source databases support).

References:
- Notion upgrade guide: [`https://developers.notion.com/docs/upgrade-guide-2025-09-03`](https://developers.notion.com/docs/upgrade-guide-2025-09-03)
- Notion changelog: [`https://developers.notion.com/page/changelog`](https://developers.notion.com/page/changelog)

## Required env vars

- `GROUPME_ACCESS_TOKEN` (or `GROUP_ME_MASTER_TOKEN`)
- `NOTION_API_KEY`
- `NOTION_DATA_SOURCE_ID_GROUPS`
- `NOTION_DATA_SOURCE_ID_BOTS`

## Run (dry-run)

From `backend/go`:

```bash
go run ./cmd/groupmesync
```

## Run (apply changes)

```bash
go run ./cmd/groupmesync --apply
```

## Update engagement metrics (optional)

This aggregates Notion **Bot Message Logs** into **Group Engagement Metrics**.

Additional required env vars:
- `NOTION_DATA_SOURCE_ID_BOT_MESSAGE_LOGS`
- `NOTION_DATA_SOURCE_ID_GROUP_ENGAGEMENT_METRICS`

Dry-run:

```bash
go run ./cmd/groupmesync --metrics
```

Apply:

```bash
go run ./cmd/groupmesync --metrics --apply
```

Metrics only (skip GroupMe API calls):

```bash
go run ./cmd/groupmesync --metrics-only --apply
```

## Notes

- Dry-run prints intended creates/updates without writing to Notion.
- Apply mode will create/update pages in Notion. Keep your Notion token secure.


