-- Migration: 000015_unique_constraints.sql
-- Description: Adds unique constraints to support idempotent upserts for the ingestion engine

DO $$ 
BEGIN 
    -- 1. ItemTypes should have unique codes (SKUs)
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'item_types_code_key') THEN
        ALTER TABLE item_types ADD CONSTRAINT item_types_code_key UNIQUE (code);
    END IF;

    -- 2. Companies should have unique names
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'companies_name_key') THEN
        ALTER TABLE companies ADD CONSTRAINT companies_name_key UNIQUE (name);
    END IF;

    -- 3. Places cleanup: De-duplicate names before adding constraint
    -- Append ID to names that are duplicated to satisfy uniqueness
    UPDATE places p
    SET name = p.name || ' (ID ' || p.id || ')'
    WHERE p.id IN (
        SELECT id FROM (
            SELECT id, row_number() OVER(PARTITION BY name ORDER BY id) as rn
            FROM places
        ) t WHERE t.rn > 1
    );

    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'places_name_key') THEN
        ALTER TABLE places ADD CONSTRAINT places_name_key UNIQUE (name);
    END IF;

    -- 4. People should have unique names (fallback)
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'people_names_key') THEN
        ALTER TABLE people ADD CONSTRAINT people_names_key UNIQUE (given_name, family_name);
    END IF;

    -- 5. Assets already have unique constraints on asset_tag and serial_number?
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'assets_asset_tag_key') THEN
        ALTER TABLE assets ADD CONSTRAINT assets_asset_tag_key UNIQUE (asset_tag);
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'assets_serial_number_key') THEN
        ALTER TABLE assets ADD CONSTRAINT assets_serial_number_key UNIQUE (serial_number);
    END IF;
END $$;
