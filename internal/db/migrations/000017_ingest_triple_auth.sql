-- Migration: 000017_ingest_triple_auth.sql
-- Description: Supports triple-auth model (Login, Verify, Refresh)

ALTER TABLE ingest_sources ADD COLUMN IF NOT EXISTS verify_endpoint TEXT;
