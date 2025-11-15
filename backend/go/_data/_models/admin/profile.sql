-- Admin/staff profiles for internal support & operations

BEGIN;

CREATE TABLE IF NOT EXISTS admin_profiles (
  id           TEXT PRIMARY KEY,
  user_id      TEXT UNIQUE REFERENCES users(id) ON DELETE SET NULL,
  kind         TEXT NOT NULL DEFAULT 'agent', -- admin | agent | support
  skills       TEXT[] DEFAULT ARRAY[]::TEXT[],
  notes        TEXT,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;

