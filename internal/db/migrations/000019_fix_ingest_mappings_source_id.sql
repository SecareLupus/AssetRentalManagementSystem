-- Migration: 000019_fix_ingest_mappings_source_id.sql
-- Description: Makes source_id optional in ingest_mappings since it's redundant with endpoint_id

ALTER TABLE ingest_mappings ALTER COLUMN source_id DROP NOT NULL;
