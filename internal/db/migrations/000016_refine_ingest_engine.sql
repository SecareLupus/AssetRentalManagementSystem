-- Migration: 000016_refine_ingest_engine.sql
-- Description: Refines ingestion engine to support multiple endpoints per source and better auth

-- 1. Update ingest_sources for better auth and base URL
ALTER TABLE ingest_sources RENAME COLUMN api_url TO base_url;
ALTER TABLE ingest_sources ADD COLUMN IF NOT EXISTS auth_endpoint TEXT;
ALTER TABLE ingest_sources ADD COLUMN IF NOT EXISTS last_token TEXT;
ALTER TABLE ingest_sources ADD COLUMN IF NOT EXISTS token_expiry TIMESTAMP WITH TIME ZONE;
ALTER TABLE ingest_sources DROP COLUMN IF EXISTS target_model; -- Moved to mapping/endpoint level

-- 2. Create ingest_endpoints table
CREATE TABLE IF NOT EXISTS ingest_endpoints (
    id SERIAL PRIMARY KEY,
    source_id INTEGER NOT NULL REFERENCES ingest_sources(id) ON DELETE CASCADE,
    path TEXT NOT NULL,
    method TEXT NOT NULL DEFAULT 'GET',
    request_body JSONB,
    resp_strategy TEXT NOT NULL DEFAULT 'auto', -- 'single', 'list', 'auto'
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    
    last_sync_at TIMESTAMP WITH TIME ZONE,
    last_success_at TIMESTAMP WITH TIME ZONE,
    last_etag TEXT,
    last_payload_hash TEXT,
    
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 3. Refactor ingest_mappings to link to endpoints and have target_model
-- First, if there are existing mappings, we might need to migrate them.
-- But since this is a refinement phase, we'll recreate for simplicity or careful migration.
ALTER TABLE ingest_mappings ADD COLUMN IF NOT EXISTS endpoint_id INTEGER REFERENCES ingest_endpoints(id) ON DELETE CASCADE;
ALTER TABLE ingest_mappings ADD COLUMN IF NOT EXISTS target_model TEXT NOT NULL DEFAULT 'asset';

-- If we have old data, try to link it to a default endpoint (we'll fix this in the worker/handlers)
-- For a fresh start, you'd just drop and recreate.

-- 4. Clean up old mappings that don't have endpoint_id if we want strictness later
-- ALTER TABLE ingest_mappings DROP COLUMN source_id; -- We'll keep it for now but it's redundant if endpoint_id is set
