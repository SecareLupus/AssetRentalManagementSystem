-- Migration: 000015_unique_constraints.sql
-- Description: Adds unique constraints to support idempotent upserts for the ingestion engine

-- 1. ItemTypes should have unique codes (SKUs)
ALTER TABLE item_types ADD CONSTRAINT item_types_code_key UNIQUE (code);

-- 2. Companies should have unique names
ALTER TABLE companies ADD CONSTRAINT companies_name_key UNIQUE (name);

-- 3. Places should have unique names (at least within a site, but for now globally)
ALTER TABLE places ADD CONSTRAINT places_name_key UNIQUE (name);

-- 4. People should have unique names (fallback)
-- Note: In a real system, you'd use email or employee ID
ALTER TABLE people ADD CONSTRAINT people_names_key UNIQUE (given_name, family_name);

-- 5. Assets already have unique constraints on asset_tag and serial_number from convergence? 
-- Let's double check initial schema.
-- Actually, ON CONFLICT requires a specific index/constraint.
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'assets_asset_tag_key') THEN
        ALTER TABLE assets ADD CONSTRAINT assets_asset_tag_key UNIQUE (asset_tag);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'assets_serial_number_key') THEN
        ALTER TABLE assets ADD CONSTRAINT assets_serial_number_key UNIQUE (serial_number);
    END IF;
END $$;
