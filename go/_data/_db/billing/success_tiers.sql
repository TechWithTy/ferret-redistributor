-- Success-based tiers: thresholds and enrollment for performance-based pricing

BEGIN;

-- Time series of success metrics per org (inputs to performance billing)
CREATE TABLE IF NOT EXISTS success_metrics (
  id            BIGSERIAL PRIMARY KEY,
  org_id        TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  key           TEXT NOT NULL, -- conversions, leads, revenue
  dimension     TEXT,          -- optional dimension (campaign:xyz, platform:instagram)
  period_start  TIMESTAMPTZ NOT NULL,
  period_end    TIMESTAMPTZ NOT NULL,
  value         DOUBLE PRECISION NOT NULL,
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(org_id, key, dimension, period_start, period_end)
);

-- Tiers define thresholds and fee percent for a metric key
CREATE TABLE IF NOT EXISTS success_tiers (
  id           TEXT PRIMARY KEY,
  name         TEXT NOT NULL,
  metric_key   TEXT NOT NULL,
  threshold    DOUBLE PRECISION NOT NULL,
  fee_percent  NUMERIC(5,2) NOT NULL DEFAULT 0.0, -- percentage applied to value or revenue
  effective_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  is_active    BOOLEAN NOT NULL DEFAULT TRUE
);

-- Enrollment enables success billing for an org
CREATE TABLE IF NOT EXISTS org_success_enrollments (
  org_id      TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  tier_id     TEXT NOT NULL REFERENCES success_tiers(id),
  status      TEXT NOT NULL DEFAULT 'active', -- active, paused, canceled
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY(org_id, tier_id)
);

COMMIT;

