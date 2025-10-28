-- Marketplace to exchange or buy posts among members

BEGIN;

CREATE TABLE IF NOT EXISTS marketplace_posts (
  id             TEXT PRIMARY KEY,
  org_id         TEXT NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
  seller_user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  content_item_id TEXT REFERENCES content_items(id) ON DELETE SET NULL,
  title          TEXT NOT NULL,
  description    TEXT,
  price_cents    INTEGER NOT NULL DEFAULT 0,
  currency       TEXT NOT NULL DEFAULT 'USD',
  status         TEXT NOT NULL DEFAULT 'active', -- active, sold, withdrawn
  metadata       JSONB NOT NULL DEFAULT '{}'::jsonb,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS marketplace_transactions (
  id                TEXT PRIMARY KEY,
  post_id           TEXT NOT NULL REFERENCES marketplace_posts(id) ON DELETE CASCADE,
  buyer_user_id     TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  amount_cents      INTEGER NOT NULL,
  currency          TEXT NOT NULL,
  status            TEXT NOT NULL DEFAULT 'completed', -- pending, completed, refunded, canceled
  created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

COMMIT;

