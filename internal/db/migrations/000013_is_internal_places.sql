-- Add is_internal column to places table
ALTER TABLE places ADD COLUMN is_internal BOOLEAN NOT NULL DEFAULT FALSE;
