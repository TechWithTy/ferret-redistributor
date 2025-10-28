-- AI generation, variants, experiments, and outcomes for auto‑optimization

BEGIN;

-- AI generations (record prompts and outputs)
CREATE TABLE IF NOT EXISTS ai_generations (
  id             TEXT PRIMARY KEY,
  org_id         TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  user_id        TEXT REFERENCES users(id) ON DELETE SET NULL,
  model          TEXT NOT NULL, -- e.g., gpt-4o-mini
  prompt         TEXT NOT NULL,
  parameters     JSONB NOT NULL DEFAULT '{}'::jsonb,
  output_text    TEXT,
  output_json    JSONB,
  content_item_id TEXT REFERENCES content_items(id) ON DELETE SET NULL,
  status         TEXT NOT NULL DEFAULT 'succeeded', -- queued, running, succeeded, failed
  error_message  TEXT,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Variants generated per generation to test different copy/media
CREATE TABLE IF NOT EXISTS ai_variants (
  id             TEXT PRIMARY KEY,
  generation_id  TEXT NOT NULL REFERENCES ai_generations(id) ON DELETE CASCADE,
  content_item_id TEXT REFERENCES content_items(id) ON DELETE SET NULL,
  variant_index  INT NOT NULL,
  payload        JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(generation_id, variant_index)
);

-- Experiments linking variants to scheduled posts and measuring outcomes
CREATE TABLE IF NOT EXISTS experiments (
  id             TEXT PRIMARY KEY,
  org_id         TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  name           TEXT NOT NULL,
  hypothesis     TEXT,
  started_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  ended_at       TIMESTAMPTZ,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS experiment_arms (
  id             TEXT PRIMARY KEY,
  experiment_id  TEXT NOT NULL REFERENCES experiments(id) ON DELETE CASCADE,
  variant_id     TEXT REFERENCES ai_variants(id) ON DELETE SET NULL,
  weight         REAL NOT NULL DEFAULT 1.0,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Outcomes collected from platform analytics per scheduled post
CREATE TABLE IF NOT EXISTS post_outcomes (
  id              TEXT PRIMARY KEY,
  scheduled_post_id TEXT NOT NULL REFERENCES scheduled_posts(id) ON DELETE CASCADE,
  platform        TEXT NOT NULL,
  external_id     TEXT, -- published id
  impressions     BIGINT NOT NULL DEFAULT 0,
  reach           BIGINT NOT NULL DEFAULT 0,
  likes           BIGINT NOT NULL DEFAULT 0,
  comments        BIGINT NOT NULL DEFAULT 0,
  shares          BIGINT NOT NULL DEFAULT 0,
  clicks          BIGINT NOT NULL DEFAULT 0,
  saves           BIGINT NOT NULL DEFAULT 0,
  conversions     BIGINT NOT NULL DEFAULT 0,
  collected_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  metadata        JSONB NOT NULL DEFAULT '{}'::jsonb
);

-- Link experiments to scheduled posts (multi‑arm bandits/A-B tests)
CREATE TABLE IF NOT EXISTS scheduled_post_arms (
  scheduled_post_id TEXT NOT NULL REFERENCES scheduled_posts(id) ON DELETE CASCADE,
  arm_id            TEXT NOT NULL REFERENCES experiment_arms(id) ON DELETE CASCADE,
  PRIMARY KEY (scheduled_post_id, arm_id)
);

COMMIT;

