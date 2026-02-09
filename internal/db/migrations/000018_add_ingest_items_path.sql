-- Add items_path to ingest_endpoints
ALTER TABLE ingest_endpoints ADD COLUMN items_path TEXT DEFAULT '$';
