-- Migration 000011: Schema.org Convergence (Places & People)

-- 1. places table (Unified Site/Location)
CREATE TABLE places (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(191) NOT NULL,
    description TEXT,
    contained_in_place_id BIGINT REFERENCES places(id),
    owner_id BIGINT REFERENCES companies(id), -- Top level org owning the place
    category VARCHAR(64), -- 'site', 'room', 'zone', 'floor', etc.
    address JSONB, -- PostalAddress structure
    presumed_demands JSONB, -- Combined with presumed_asset_needs
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. people table (Standalone Person)
CREATE TABLE people (
    id BIGSERIAL PRIMARY KEY,
    given_name VARCHAR(64) NOT NULL,
    family_name VARCHAR(64) NOT NULL,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 3. contact_points table (Email/Phone entries for People)
CREATE TABLE contact_points (
    id BIGSERIAL PRIMARY KEY,
    person_id BIGINT NOT NULL REFERENCES people(id),
    email VARCHAR(191),
    phone VARCHAR(32),
    contact_type VARCHAR(64), -- 'technical', 'billing', 'emergency'
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 4. organization_roles table (Connecting Persons to Organizations/Companies)
CREATE TABLE organization_roles (
    id BIGSERIAL PRIMARY KEY,
    person_id BIGINT NOT NULL REFERENCES people(id),
    organization_id BIGINT NOT NULL REFERENCES companies(id),
    role_name VARCHAR(64) NOT NULL, -- 'Contact', 'Manager', 'Technician'
    start_date TIMESTAMP WITH TIME ZONE,
    end_date TIMESTAMP WITH TIME ZONE,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Temporary columns for migration mapping
ALTER TABLE places ADD COLUMN tmp_old_site_id BIGINT;
ALTER TABLE places ADD COLUMN tmp_old_location_id BIGINT;
ALTER TABLE people ADD COLUMN tmp_old_contact_id BIGINT;

-- Port Data: Sites -> Places
INSERT INTO places (name, owner_id, category, address, metadata, created_at, updated_at, tmp_old_site_id)
SELECT 
    name, 
    company_id, 
    'site', 
    jsonb_build_object(
        'street_address', address_street,
        'address_locality', address_city,
        'address_region', address_state,
        'postal_code', address_zip,
        'address_country', address_country
    ),
    metadata,
    created_at,
    updated_at,
    id
FROM sites;

-- Port Data: Locations -> Places (Top-level locations within sites)
INSERT INTO places (name, contained_in_place_id, category, presumed_demands, metadata, created_at, updated_at, tmp_old_location_id)
SELECT 
    l.name,
    p.id,
    l.location_type,
    l.presumed_asset_needs,
    l.metadata,
    l.created_at,
    l.updated_at,
    l.id
FROM locations l
JOIN places p ON p.tmp_old_site_id = l.site_id
WHERE l.parent_id IS NULL;

-- Port Data: Locations -> Places (Nested locations - assume max 5 levels for simplicity in this migration)
-- Level 2
INSERT INTO places (name, contained_in_place_id, category, presumed_demands, metadata, created_at, updated_at, tmp_old_location_id)
SELECT 
    l.name,
    p.id,
    l.location_type,
    l.presumed_asset_needs,
    l.metadata,
    l.created_at,
    l.updated_at,
    l.id
FROM locations l
JOIN places p ON p.tmp_old_location_id = l.parent_id
WHERE l.parent_id IS NOT NULL 
AND NOT EXISTS (SELECT 1 FROM places WHERE tmp_old_location_id = l.id);

-- Port Data: Contacts -> People & Roles
INSERT INTO people (given_name, family_name, metadata, created_at, updated_at, tmp_old_contact_id)
SELECT first_name, last_name, metadata, created_at, updated_at, id FROM contacts;

INSERT INTO contact_points (person_id, email, phone, contact_type, created_at, updated_at)
SELECT p.id, c.email, c.phone, 'general', c.created_at, c.updated_at
FROM contacts c
JOIN people p ON p.tmp_old_contact_id = c.id;

INSERT INTO organization_roles (person_id, organization_id, role_name, created_at, updated_at)
SELECT p.id, c.company_id, c.role, c.created_at, c.updated_at
FROM contacts c
JOIN people p ON p.tmp_old_contact_id = c.id;

-- Update other tables
ALTER TABLE assets ADD COLUMN place_id BIGINT REFERENCES places(id);
-- Best effort backfill for assets if we had exact site/location links. 
-- Currently assets only have 'location' string. This will remain for now or be manual.

ALTER TABLE event_asset_needs ADD COLUMN place_id BIGINT REFERENCES places(id);
UPDATE event_asset_needs ean
SET place_id = p.id
FROM places p
WHERE p.tmp_old_location_id = ean.location_id;

-- Clean up temporary columns
ALTER TABLE places DROP COLUMN tmp_old_site_id;
ALTER TABLE places DROP COLUMN tmp_old_location_id;
ALTER TABLE people DROP COLUMN tmp_old_contact_id;

-- Indices
CREATE INDEX idx_places_owner ON places(owner_id);
CREATE INDEX idx_places_parent ON places(contained_in_place_id);
CREATE INDEX idx_people_names ON people(family_name, given_name);
CREATE INDEX idx_org_roles_person ON organization_roles(person_id);
CREATE INDEX idx_org_roles_org ON organization_roles(organization_id);
CREATE INDEX idx_contact_points_person ON contact_points(person_id);
