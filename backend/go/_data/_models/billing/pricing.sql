-- Pricing plans, subscriptions, and basic usage limits (MVP)

BEGIN;

CREATE TABLE IF NOT EXISTS pricing_plans (
  id             TEXT PRIMARY KEY,
  name           TEXT NOT NULL,
  tier           TEXT NOT NULL, -- free, pro, business, enterprise
  price_cents    INTEGER NOT NULL DEFAULT 0,
  currency       TEXT NOT NULL DEFAULT 'USD',
  interval       TEXT NOT NULL DEFAULT 'monthly', -- monthly, annual
  is_active      BOOLEAN NOT NULL DEFAULT TRUE,
  limits         JSONB NOT NULL DEFAULT '{}'::jsonb, -- e.g., {"posts_per_month":100, "seats":3}
  external_id    TEXT, -- e.g., Stripe price id
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(name, interval)
);

CREATE TABLE IF NOT EXISTS org_subscriptions (
  id                   TEXT PRIMARY KEY,
  org_id               TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  plan_id              TEXT NOT NULL REFERENCES pricing_plans(id),
  status               TEXT NOT NULL DEFAULT 'active', -- trialing, active, past_due, canceled
  seats                INTEGER NOT NULL DEFAULT 1,
  trial_ends_at        TIMESTAMPTZ,
  current_period_start TIMESTAMPTZ,
  current_period_end   TIMESTAMPTZ,
  cancel_at            TIMESTAMPTZ,
  canceled_at          TIMESTAMPTZ,
  metadata             JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_org ON org_subscriptions(org_id);
CREATE INDEX IF NOT EXISTS idx_subscriptions_status ON org_subscriptions(status);

COMMIT;

