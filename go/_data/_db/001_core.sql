-- PostgreSQL core schema
-- Organization, users/teams, content, campaigns, scheduling (compatible with Go code)

BEGIN;

-- Extensions (optional)
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Organizations
CREATE TABLE IF NOT EXISTS organizations (
  id            TEXT PRIMARY KEY,
  name          TEXT NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Users
CREATE TABLE IF NOT EXISTS users (
  id            TEXT PRIMARY KEY,
  org_id        TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  email         TEXT NOT NULL UNIQUE,
  display_name  TEXT,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Teams and membership
CREATE TABLE IF NOT EXISTS teams (
  id            TEXT PRIMARY KEY,
  org_id        TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  name          TEXT NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(org_id, name)
);

CREATE TABLE IF NOT EXISTS team_members (
  team_id       TEXT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
  user_id       TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role          TEXT NOT NULL DEFAULT 'member', -- owner, admin, editor, viewer
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY(team_id, user_id)
);

-- Social accounts (per platform connection)
CREATE TABLE IF NOT EXISTS social_accounts (
  id               TEXT PRIMARY KEY,
  org_id           TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  team_id          TEXT REFERENCES teams(id) ON DELETE SET NULL,
  platform         TEXT NOT NULL, -- instagram, linkedin, twitter, facebook, youtube, behiiv
  handle           TEXT,
  external_id      TEXT,          -- page/channel/profile id
  auth_kind        TEXT,          -- api_key, oauth2, etc
  auth_meta        JSONB NOT NULL DEFAULT '{}'::jsonb, -- token metadata, secrets stored elsewhere
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(org_id, platform, external_id)
);

-- Content library
CREATE TABLE IF NOT EXISTS content_items (
  id              TEXT PRIMARY KEY,
  org_id          TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  title           TEXT,
  body            TEXT,
  canonical_url   TEXT,
  media_url       TEXT,
  metadata        JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Campaigns
CREATE TABLE IF NOT EXISTS campaigns (
  id              TEXT PRIMARY KEY,
  org_id          TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  name            TEXT NOT NULL,
  description     TEXT,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(org_id, name)
);

-- Scheduled posts: structure expected by Go calendar package
CREATE TABLE IF NOT EXISTS scheduled_posts (
  id              TEXT PRIMARY KEY,
  org_id          TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  campaign_id     TEXT REFERENCES campaigns(id) ON DELETE SET NULL,
  content_id      TEXT REFERENCES content_items(id) ON DELETE SET NULL,
  social_account_id TEXT REFERENCES social_accounts(id) ON DELETE SET NULL,
  platform        TEXT NOT NULL,
  caption         TEXT,
  hashtags        TEXT,
  scheduled_at    TIMESTAMPTZ NOT NULL,
  status          TEXT NOT NULL DEFAULT 'scheduled', -- scheduled, processing, published, failed, canceled
  external_id     TEXT,
  published_at    TIMESTAMPTZ,
  metadata        JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Basic indexes used by queries/workers
CREATE INDEX IF NOT EXISTS idx_scheduled_posts_status_scheduled_at
  ON scheduled_posts(status, scheduled_at);
CREATE INDEX IF NOT EXISTS idx_scheduled_posts_campaign
  ON scheduled_posts(campaign_id);
CREATE INDEX IF NOT EXISTS idx_scheduled_posts_content
  ON scheduled_posts(content_id);

COMMIT;
