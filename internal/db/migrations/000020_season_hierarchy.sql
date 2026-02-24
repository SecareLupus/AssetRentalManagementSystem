-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE show_companies (
    id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(company_id)
);

CREATE TABLE seasons (
    id SERIAL PRIMARY KEY,
    show_company_id INTEGER NOT NULL REFERENCES show_companies(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE rings (
    id SERIAL PRIMARY KEY,
    show_company_id INTEGER NOT NULL REFERENCES show_companies(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE shows (
    id SERIAL PRIMARY KEY,
    season_id INTEGER NOT NULL REFERENCES seasons(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    location_id INTEGER REFERENCES places(id) ON DELETE SET NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE show_rings (
    id SERIAL PRIMARY KEY,
    show_id INTEGER NOT NULL REFERENCES shows(id) ON DELETE CASCADE,
    ring_id INTEGER NOT NULL REFERENCES rings(id) ON DELETE CASCADE,
    UNIQUE(show_id, ring_id)
);

CREATE TABLE ring_loadout_items (
    id SERIAL PRIMARY KEY,
    show_ring_id INTEGER NOT NULL REFERENCES show_rings(id) ON DELETE CASCADE,
    item_type_id INTEGER NOT NULL REFERENCES item_types(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 1,
    UNIQUE(show_ring_id, item_type_id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE IF EXISTS ring_loadout_items;
DROP TABLE IF EXISTS show_rings;
DROP TABLE IF EXISTS shows;
DROP TABLE IF EXISTS rings;
DROP TABLE IF EXISTS seasons;
DROP TABLE IF EXISTS show_companies;
