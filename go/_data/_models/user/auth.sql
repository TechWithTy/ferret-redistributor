-- User authentication: identities for email, phone, LinkedIn, Meta; sessions; audits

BEGIN;

-- Identities represent login handles for a user across providers.
-- provider: email | phone | linkedin | meta
CREATE TABLE IF NOT EXISTS auth_identities (
  id             TEXT PRIMARY KEY,
  user_id        TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  provider       TEXT NOT NULL,
  identifier     TEXT NOT NULL,   -- email, phone E.164, or provider user id
  secret_hash    TEXT,            -- e.g., bcrypt hash for email/password; NULL for OAuth
  oauth_data     JSONB NOT NULL DEFAULT '{}'::jsonb, -- tokens/claims meta for OAuth providers
  verified_at    TIMESTAMPTZ,     -- when identifier was verified (email/phone)
  is_primary     BOOLEAN NOT NULL DEFAULT FALSE,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  UNIQUE(provider, identifier)
);

CREATE INDEX IF NOT EXISTS idx_auth_identities_user ON auth_identities(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_identities_provider ON auth_identities(provider);
-- Ensure only one primary identity per (user, provider)
CREATE UNIQUE INDEX IF NOT EXISTS ux_auth_identities_primary_per_provider
  ON auth_identities(user_id, provider)
  WHERE is_primary = TRUE;

-- Email verification tokens (one-time use)
CREATE TABLE IF NOT EXISTS auth_email_verifications (
  id           TEXT PRIMARY KEY,
  user_id      TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  email        TEXT NOT NULL,
  token        TEXT NOT NULL UNIQUE,
  expires_at   TIMESTAMPTZ NOT NULL,
  used_at      TIMESTAMPTZ,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_auth_email_verifications_user ON auth_email_verifications(user_id);

-- Password reset tokens (one-time use)
CREATE TABLE IF NOT EXISTS auth_password_resets (
  id           TEXT PRIMARY KEY,
  user_id      TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  email        TEXT NOT NULL,
  token        TEXT NOT NULL UNIQUE,
  expires_at   TIMESTAMPTZ NOT NULL,
  used_at      TIMESTAMPTZ,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_auth_password_resets_user ON auth_password_resets(user_id);

-- Phone verification codes (short-lived OTP)
CREATE TABLE IF NOT EXISTS auth_phone_codes (
  id           TEXT PRIMARY KEY,
  user_id      TEXT REFERENCES users(id) ON DELETE CASCADE,
  phone        TEXT NOT NULL,        -- E.164
  code_hash    TEXT NOT NULL,        -- store hash of the OTP, not plaintext
  expires_at   TIMESTAMPTZ NOT NULL,
  attempts     SMALLINT NOT NULL DEFAULT 0,
  used_at      TIMESTAMPTZ,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_auth_phone_codes_phone ON auth_phone_codes(phone);

-- Sessions (server-managed) using opaque tokens; token itself should be hashed before storage
CREATE TABLE IF NOT EXISTS auth_sessions (
  id            TEXT PRIMARY KEY,     -- session id (random)
  user_id       TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash    TEXT NOT NULL,        -- hash of the bearer token
  created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  last_seen_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  expires_at    TIMESTAMPTZ NOT NULL,
  ip_address    TEXT,
  user_agent    TEXT,
  revoked_at    TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_auth_sessions_user ON auth_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_auth_sessions_expires ON auth_sessions(expires_at);

-- Login audit trail
CREATE TABLE IF NOT EXISTS auth_logins (
  id           BIGSERIAL PRIMARY KEY,
  user_id      TEXT REFERENCES users(id) ON DELETE SET NULL,
  provider     TEXT NOT NULL,
  identifier   TEXT,
  ip_address   TEXT,
  user_agent   TEXT,
  success      BOOLEAN NOT NULL,
  error        TEXT,
  occurred_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_auth_logins_user_time ON auth_logins(user_id, occurred_at DESC);

COMMIT;
