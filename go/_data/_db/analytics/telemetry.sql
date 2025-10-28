-- Optional application metrics table for fallback/local tracking

BEGIN;

CREATE TABLE IF NOT EXISTS app_metrics (
  id           BIGSERIAL PRIMARY KEY,
  name         TEXT NOT NULL,
  value        DOUBLE PRECISION NOT NULL,
  attributes   JSONB NOT NULL DEFAULT '{}'::jsonb,
  recorded_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_app_metrics_name_time ON app_metrics(name, recorded_at);

COMMIT;

