-- Credit wallets and transactions for AI, posting, etc. (MVP)

BEGIN;

CREATE TABLE IF NOT EXISTS credit_products (
  id           TEXT PRIMARY KEY,
  name         TEXT NOT NULL,
  description  TEXT,
  unit         TEXT NOT NULL DEFAULT 'credit',
  bundle_size  INTEGER NOT NULL DEFAULT 100,
  price_cents  INTEGER NOT NULL DEFAULT 0,
  currency     TEXT NOT NULL DEFAULT 'USD',
  is_active    BOOLEAN NOT NULL DEFAULT TRUE,
  external_id  TEXT, -- Stripe price id, etc
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS credit_wallets (
  id         TEXT PRIMARY KEY,
  org_id     TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  balance    BIGINT NOT NULL DEFAULT 0,
  reserved   BIGINT NOT NULL DEFAULT 0,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(org_id)
);

CREATE TABLE IF NOT EXISTS credit_transactions (
  id          TEXT PRIMARY KEY,
  wallet_id   TEXT NOT NULL REFERENCES credit_wallets(id) ON DELETE CASCADE,
  amount      BIGINT NOT NULL, -- positive for add, negative for spend
  kind        TEXT NOT NULL,   -- purchase, spend, adjust, expire, refund
  ref_id      TEXT,            -- related entity (order id, post id)
  metadata    JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_credit_tx_wallet_time ON credit_transactions(wallet_id, created_at DESC);

COMMIT;

