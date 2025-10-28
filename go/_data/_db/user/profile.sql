-- User profiles and preferences to personalize content

BEGIN;

CREATE TABLE IF NOT EXISTS user_profiles (
  id             TEXT PRIMARY KEY,
  user_id        TEXT NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
  org_id         TEXT REFERENCES organizations(id) ON DELETE SET NULL,
  display_name   TEXT,
  avatar_url     TEXT,
  bio            TEXT,
  timezone       TEXT,
  locale         TEXT,
  notification_prefs JSONB NOT NULL DEFAULT '{}'::jsonb, -- {email:true,sms:false}
  content_prefs  JSONB NOT NULL DEFAULT '{}'::jsonb, -- {hashtags:[..], keywords:[..]}
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;

