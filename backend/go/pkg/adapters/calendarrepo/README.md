# Calendar Repository (Go)

Small adapter to insert rows into `scheduled_posts`.

- `repo.go`: `SchedulePost` and `BulkSchedule` insert with `status='scheduled'` and timestamps.

## Expected Schema
- Table: `scheduled_posts`
  - id (text/uuid), campaign_id, content_id, platform (text), caption (text), hashtags (text), scheduled_at (timestamptz), status (text), metadata (json), created_at, updated_at
- See `python/calendar/models.py` for the Python ORM that defines the same shape.
- Recommended indexes (apply via SQL): `python/calendar/migrations/001_indexes.sql`.

