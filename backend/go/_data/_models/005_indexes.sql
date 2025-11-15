-- Additional indexes and constraints for performance and data quality

BEGIN;

-- Scheduled posts operational indexes
CREATE INDEX IF NOT EXISTS idx_scheduled_posts_status_time
  ON scheduled_posts(status, scheduled_at);
CREATE INDEX IF NOT EXISTS idx_scheduled_posts_org_time
  ON scheduled_posts(org_id, scheduled_at);

-- Content lookups
CREATE INDEX IF NOT EXISTS idx_content_items_org
  ON content_items(org_id);

-- Social accounts lookups
CREATE INDEX IF NOT EXISTS idx_social_accounts_org_platform
  ON social_accounts(org_id, platform);

-- Outcomes lookups
CREATE INDEX IF NOT EXISTS idx_post_outcomes_post_time
  ON post_outcomes(scheduled_post_id, collected_at);

COMMIT;

