# Calendar (Go)

Helpers for reading and updating scheduled posts in Postgres.

- `FetchAndClaimDuePosts(ctx, db, within, limit)` — atomic claim (status `scheduled` -> `processing`) using `FOR UPDATE SKIP LOCKED`.
- `FetchScheduledPostsWithin(ctx, db, start, end)` — read-only range fetch.
- `UpdatePostStatus(ctx, db, id, status, externalID, publishedAt, metadata)` — set status and fields after posting.

## Indexes
Run: `psql "$DATABASE_URL" -f python/calendar/migrations/001_indexes.sql`

