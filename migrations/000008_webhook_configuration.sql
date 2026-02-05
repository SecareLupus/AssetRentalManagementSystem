-- 000008_webhook_configuration.sql
CREATE TABLE IF NOT EXISTS webhooks (
    id BIGSERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    secret TEXT,
    enabled_events JSONB NOT NULL DEFAULT '[]', -- Array of event_type strings
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_webhooks_active ON webhooks(is_active) WHERE is_active = TRUE;
