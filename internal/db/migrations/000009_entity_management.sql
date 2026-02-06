-- Migration 000009: Entity Management
-- Companies, Contacts, Sites, Locations, Events

-- 1. companies table
CREATE TABLE companies (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(191) NOT NULL,
    legal_name VARCHAR(191),
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. contacts table
CREATE TABLE contacts (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES companies(id),
    first_name VARCHAR(64) NOT NULL,
    last_name VARCHAR(64) NOT NULL,
    email VARCHAR(191),
    phone VARCHAR(32),
    role VARCHAR(64),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. sites table (Facilities)
CREATE TABLE sites (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies(id),
    name VARCHAR(191) NOT NULL,
    address_street TEXT,
    address_city VARCHAR(128),
    address_state VARCHAR(64),
    address_zip VARCHAR(20),
    address_country VARCHAR(64),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 4. locations table (Nested physical spaces)
CREATE TABLE locations (
    id BIGSERIAL PRIMARY KEY,
    site_id BIGINT NOT NULL REFERENCES sites(id),
    parent_id BIGINT REFERENCES locations(id),
    name VARCHAR(191) NOT NULL,
    location_type VARCHAR(64), -- 'room', 'floor', 'zone'
    presumed_asset_needs JSONB, -- Default needs for this location
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 5. events table (Windows/Projects)
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies(id),
    name VARCHAR(191) NOT NULL,
    description TEXT,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'confirmed', -- 'assumed', 'confirmed', 'cancelled'
    parent_event_id BIGINT REFERENCES events(id),
    recurrence_rule TEXT, -- iCal-like rule string
    last_confirmed_at TIMESTAMP WITH TIME ZONE,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 6. event_asset_needs table
CREATE TABLE event_asset_needs (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL REFERENCES events(id),
    item_type_id BIGINT NOT NULL REFERENCES item_types(id),
    quantity INT NOT NULL DEFAULT 1,
    is_assumed BOOLEAN NOT NULL DEFAULT FALSE,
    location_id BIGINT REFERENCES locations(id),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indices
CREATE INDEX idx_contacts_company ON contacts(company_id);
CREATE INDEX idx_sites_company ON sites(company_id);
CREATE INDEX idx_locations_site ON locations(site_id);
CREATE INDEX idx_locations_parent ON locations(parent_id);
CREATE INDEX idx_events_company ON events(company_id);
CREATE INDEX idx_events_window ON events(start_time, end_time);
CREATE INDEX idx_ean_event ON event_asset_needs(event_id);
