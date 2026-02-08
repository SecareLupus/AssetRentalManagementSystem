-- Migration: 000017_dual_token_auth.sql
-- Description: Add support for refresh tokens in the ingestion engine

ALTER TABLE ingest_sources ADD COLUMN IF NOT EXISTS refresh_endpoint TEXT;
ALTER TABLE ingest_sources ADD COLUMN IF NOT EXISTS refresh_token TEXT;
