-- Admin analytics: product events & lightweight KPI snapshots

BEGIN;

CREATE TABLE IF NOT EXISTS product_events (
  id           BIGSERIAL PRIMARY KEY,
  org_id       TEXT REFERENCES organizations(id) ON DELETE SET NULL,
  user_id      TEXT REFERENCES users(id) ON DELETE SET NULL,
  session_id   TEXT,
  source       TEXT, -- web, api, worker
  name         TEXT NOT NULL, -- event name
  properties   JSONB NOT NULL DEFAULT '{}'::jsonb,
  occurred_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_product_events_time ON product_events(occurred_at);
CREATE INDEX IF NOT EXISTS idx_product_events_name ON product_events(name);
CREATE INDEX IF NOT EXISTS idx_product_events_org ON product_events(org_id);

-- Daily KPI snapshots (optional backfill)
CREATE TABLE IF NOT EXISTS admin_kpis (
  day        DATE NOT NULL,
  metric     TEXT NOT NULL, -- dau, wau, mau, new_orgs, churned_orgs, mrr_cents
  value      DOUBLE PRECISION NOT NULL,
  meta       JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  PRIMARY KEY(day, metric)
);

COMMIT;

