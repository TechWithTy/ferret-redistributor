-- Migration: Calendar performance indexes
-- Purpose: Speed up queries that fetch scheduled posts within a time window
--          and filter by status='scheduled'.
--
-- Apply with any Postgres client:
--   psql "$DATABASE_URL" -f python/calendar/migrations/001_indexes.sql
-- Or using a migration tool of your choice.

BEGIN;

-- Composite index on (status, scheduled_at) to support queries like:
-- SELECT ... FROM scheduled_posts WHERE status='scheduled' AND scheduled_at BETWEEN ... ORDER BY scheduled_at;
CREATE INDEX IF NOT EXISTS idx_scheduled_posts_status_sched_at
  ON scheduled_posts (status, scheduled_at);

-- Partial index focusing only on scheduled posts ordered by time.
-- This is very effective if most rows are not 'scheduled'.
CREATE INDEX IF NOT EXISTS idx_scheduled_posts_sched_at_partial
  ON scheduled_posts (scheduled_at)
  WHERE status = 'scheduled';

-- Optional: speed up lookups by campaign when viewing a campaign calendar.
CREATE INDEX IF NOT EXISTS idx_scheduled_posts_campaign_sched_at
  ON scheduled_posts (campaign_id, scheduled_at);

COMMIT;

