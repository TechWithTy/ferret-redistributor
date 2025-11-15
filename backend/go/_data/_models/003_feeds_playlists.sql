-- Add support for RSS feeds and YouTube playlists
-- This migration adds tables for managing user content sources

BEGIN;

-- RSS Feeds
CREATE TABLE IF NOT EXISTS rss_feeds (
  id              TEXT PRIMARY KEY,
  org_id          TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  user_id         TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name            TEXT NOT NULL,
  url             TEXT NOT NULL,
  last_fetched_at TIMESTAMPTZ,
  last_error      TEXT,
  refresh_interval_minutes INTEGER NOT NULL DEFAULT 1440, -- Default: 24 hours
  is_active       BOOLEAN NOT NULL DEFAULT TRUE,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(org_id, url)
);

-- RSS Feed Items (cached content)
CREATE TABLE IF NOT EXISTS rss_feed_items (
  id              TEXT PRIMARY KEY,
  feed_id         TEXT NOT NULL REFERENCES rss_feeds(id) ON DELETE CASCADE,
  title           TEXT NOT NULL,
  url             TEXT NOT NULL,
  content         TEXT,
  published_at    TIMESTAMPTZ NOT NULL,
  processed_at    TIMESTAMPTZ,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(feed_id, url)
);

-- YouTube Playlists
CREATE TABLE IF NOT EXISTS youtube_playlists (
  id              TEXT PRIMARY KEY,
  org_id          TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  user_id         TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  playlist_id     TEXT NOT NULL, -- YouTube's playlist ID
  title           TEXT NOT NULL,
  description     TEXT,
  channel_id      TEXT,
  channel_title   TEXT,
  thumbnail_url   TEXT,
  last_synced_at  TIMESTAMPTZ,
  is_active       BOOLEAN NOT NULL DEFAULT TRUE,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(org_id, playlist_id)
);

-- YouTube Playlist Items (cached content)
CREATE TABLE IF NOT EXISTS youtube_playlist_items (
  id              TEXT PRIMARY KEY,
  playlist_id     TEXT NOT NULL REFERENCES youtube_playlists(id) ON DELETE CASCADE,
  video_id        TEXT NOT NULL,
  title           TEXT NOT NULL,
  description     TEXT,
  thumbnail_url   TEXT,
  published_at    TIMESTAMPTZ NOT NULL,
  processed_at    TIMESTAMPTZ,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(playlist_id, video_id)
);

-- Admin management for content sources
CREATE TABLE IF NOT EXISTS admin_content_sources (
  id              TEXT PRIMARY KEY,
  org_id          TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  created_by      TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  source_type     TEXT NOT NULL, -- 'rss' or 'youtube_playlist'
  source_id       TEXT NOT NULL, -- ID from rss_feeds or youtube_playlists
  is_approved     BOOLEAN NOT NULL DEFAULT FALSE,
  approved_by     TEXT REFERENCES users(id) ON DELETE SET NULL,
  approved_at     TIMESTAMPTZ,
  notes           TEXT,
  created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(source_type, source_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_rss_feeds_org_user ON rss_feeds(org_id, user_id);
CREATE INDEX IF NOT EXISTS idx_rss_feed_items_feed_published ON rss_feed_items(feed_id, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_youtube_playlists_org_user ON youtube_playlists(org_id, user_id);
CREATE INDEX IF NOT EXISTS idx_youtube_playlist_items_playlist_published ON youtube_playlist_items(playlist_id, published_at DESC);
CREATE INDEX IF NOT EXISTS idx_admin_content_sources_org_approved ON admin_content_sources(org_id, is_approved);

-- Add permissions to admin role if RBAC is enabled
DO $$
BEGIN
  -- Check if the rbac.roles table exists
  IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'roles') THEN
    -- Add permissions for RSS feeds
    INSERT INTO role_permissions (role_id, resource_type, action, created_at, updated_at)
    SELECT id, 'rss_feeds', 'create', NOW(), NOW() FROM roles WHERE name = 'admin'
    ON CONFLICT DO NOTHING;

    INSERT INTO role_permissions (role_id, resource_type, action, created_at, updated_at)
    SELECT id, 'rss_feeds', 'read', NOW(), NOW() FROM roles WHERE name = 'admin'
    ON CONFLICT DO NOTHING;

    INSERT INTO role_permissions (role_id, resource_type, action, created_at, updated_at)
    SELECT id, 'rss_feeds', 'update', NOW(), NOW() FROM roles WHERE name = 'admin'
    ON CONFLICT DO NOTHING;

    INSERT INTO role_permissions (role_id, resource_type, action, created_at, updated_at)
    SELECT id, 'rss_feeds', 'delete', NOW(), NOW() FROM roles WHERE name = 'admin'
    ON CONFLICT DO NOTHING;

    -- Add permissions for YouTube playlists
    INSERT INTO role_permissions (role_id, resource_type, action, created_at, updated_at)
    SELECT id, 'youtube_playlists', 'create', NOW(), NOW() FROM roles WHERE name = 'admin'
    ON CONFLICT DO NOTHING;

    INSERT INTO role_permissions (role_id, resource_type, action, created_at, updated_at)
    SELECT id, 'youtube_playlists', 'read', NOW(), NOW() FROM roles WHERE name = 'admin'
    ON CONFLICT DO NOTHING;

    INSERT INTO role_permissions (role_id, resource_type, action, created_at, updated_at)
    SELECT id, 'youtube_playlists', 'update', NOW(), NOW() FROM roles WHERE name = 'admin'
    ON CONFLICT DO NOTHING;

    INSERT INTO role_permissions (role_id, resource_type, action, created_at, updated_at)
    SELECT id, 'youtube_playlists', 'delete', NOW(), NOW() FROM roles WHERE name = 'admin'
    ON CONFLICT DO NOTHING;
  END IF;
END $$;

COMMIT;
