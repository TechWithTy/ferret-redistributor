-- Ideal Customer Profile (ICP) and personalization settings per organization

BEGIN;

-- Core ICP profile per org (one primary profile recommended)
CREATE TABLE IF NOT EXISTS icp_profiles (
  id               TEXT PRIMARY KEY,
  org_id           TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  name             TEXT NOT NULL DEFAULT 'Default',
  timezone         TEXT,                     -- e.g., America/Los_Angeles
  languages        TEXT[] DEFAULT ARRAY['en'],
  region           TEXT,                     -- e.g., US, EU, APAC
  industry         TEXT,
  company_size     TEXT,                     -- e.g., solo, smb, mid, enterprise
  stage            TEXT,                     -- e.g., prelaunch, growth, scale
  goals            JSONB NOT NULL DEFAULT '{}'::jsonb,  -- {"primary": ["brand", "leads"], ...}
  pains            JSONB NOT NULL DEFAULT '{}'::jsonb,  -- {"primary": ["time", "tools"], ...}
  brand_voice      JSONB NOT NULL DEFAULT '{}'::jsonb,  -- {"tone":"friendly", "style":"concise", "persona":"operator"}
  guidelines       TEXT,                     -- freeform brand/style/compliance guidance
  compliance       JSONB NOT NULL DEFAULT '{}'::jsonb,  -- {"banned_topics":["..."], "disclaimers":["..."]}
  audience         JSONB NOT NULL DEFAULT '{}'::jsonb,  -- {"segments":[{"name":"...","interests":["..."]}]}
  content_pillars  JSONB NOT NULL DEFAULT '{}'::jsonb, -- {"pillars":["education","case_study","promo"]}
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(org_id, name)
);

-- Platform-specific preferences for publishing
CREATE TABLE IF NOT EXISTS icp_platform_prefs (
  id               TEXT PRIMARY KEY,
  profile_id       TEXT NOT NULL REFERENCES icp_profiles(id) ON DELETE CASCADE,
  platform         TEXT NOT NULL,              -- instagram, linkedin, twitter, youtube, facebook
  enabled          BOOLEAN NOT NULL DEFAULT TRUE,
  cadence_per_week INTEGER NOT NULL DEFAULT 3,
  post_times       JSONB NOT NULL DEFAULT '[]'::jsonb,  -- ["09:00", "13:00"] local times
  hashtags         JSONB NOT NULL DEFAULT '[]'::jsonb,  -- ["ai", "realestate"]
  keywords         JSONB NOT NULL DEFAULT '[]'::jsonb,  -- ["automation", "seller leads"]
  mentions         JSONB NOT NULL DEFAULT '[]'::jsonb,  -- ["@brand"]
  link_policy      JSONB NOT NULL DEFAULT '{}'::jsonb,  -- {"include": true, "utm": {"source":"..."}}
  media_policy     JSONB NOT NULL DEFAULT '{}'::jsonb,  -- {"ratio":"1:1","video_pref":true}
  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(profile_id, platform)
);

-- Competitors and references for tone/positioning
CREATE TABLE IF NOT EXISTS icp_competitors (
  id           TEXT PRIMARY KEY,
  profile_id   TEXT NOT NULL REFERENCES icp_profiles(id) ON DELETE CASCADE,
  platform     TEXT,
  handle       TEXT,
  url          TEXT,
  notes        TEXT,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(profile_id, platform, handle)
);

-- Constraints: banned topics/phrases or required disclaimers
CREATE TABLE IF NOT EXISTS icp_constraints (
  id           TEXT PRIMARY KEY,
  profile_id   TEXT NOT NULL REFERENCES icp_profiles(id) ON DELETE CASCADE,
  kind         TEXT NOT NULL,   -- banned_topic | banned_phrase | required_disclaimer
  value        TEXT NOT NULL,
  notes        TEXT,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;

