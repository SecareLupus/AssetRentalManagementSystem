-- Migration: 000014_ingest_engine.sql
-- Description: Supports harvesting data from external REST APIs

CREATE TABLE IF NOT EXISTS ingest_sources (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    target_model TEXT NOT NULL, -- 'item_type', 'asset', 'company', 'person', 'place'
    api_url TEXT NOT NULL,
    auth_type TEXT NOT NULL DEFAULT 'none', -- 'none', 'bearer'
    auth_credentials JSONB, -- store tokens or keys
    sync_interval_seconds INTEGER NOT NULL DEFAULT 3600,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    
    last_sync_at TIMESTAMP WITH TIME ZONE,
    last_success_at TIMESTAMP WITH TIME ZONE,
    last_status TEXT,
    last_error TEXT,
    next_sync_at TIMESTAMP WITH TIME ZONE,
    
    last_etag TEXT,
    last_payload_hash TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS ingest_mappings (
    id SERIAL PRIMARY KEY,
    source_id INTEGER NOT NULL REFERENCES ingest_sources(id) ON DELETE CASCADE,
    json_path TEXT NOT NULL, -- e.g. '$.sku'
    target_field TEXT NOT NULL, -- e.g. 'code'
    is_identity BOOLEAN NOT NULL DEFAULT FALSE, -- used for UPSERT
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Index for background worker
CREATE INDEX idx_ingest_sources_next_sync ON ingest_sources(next_sync_at) WHERE is_active = TRUE;
