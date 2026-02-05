-- Initial schema for Rental Management System
-- Aligned with Schema.org standards

-- 1. item_types table (ProductModel)
CREATE TABLE item_types (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(64) UNIQUE NOT NULL, -- Maps to schema.org/sku
    name VARCHAR(191) NOT NULL,       -- Maps to schema.org/name
    kind VARCHAR(32) NOT NULL,       -- 'serialized', 'fungible', 'kit' (Maps to schema.org/category)
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    schema_org JSONB,                -- Standard JSON-LD representation
    metadata JSONB,                  -- Internal attributes
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. assets table (IndividualProduct)
CREATE TABLE assets (
    id BIGSERIAL PRIMARY KEY,
    item_type_id BIGINT NOT NULL,
    asset_tag VARCHAR(128) UNIQUE,   -- Maps to schema.org/identifier
    serial_number VARCHAR(128),      -- Maps to schema.org/serialNumber
    status VARCHAR(32) NOT NULL DEFAULT 'available',
    location VARCHAR(191),
    assigned_to VARCHAR(191),
    mesh_node_id VARCHAR(191),
    wireguard_hostname VARCHAR(191),
    schema_org JSONB,                -- Standard JSON-LD representation
    metadata JSONB,                  -- Internal attributes
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_assets_item_type FOREIGN KEY (item_type_id) REFERENCES item_types(id)
);

-- 3. rent_actions table (RentAction)
CREATE TABLE rent_actions (
    id BIGSERIAL PRIMARY KEY,
    requester_ref VARCHAR(191) NOT NULL, -- Maps to schema.org/agent
    created_by_ref VARCHAR(191) NOT NULL,
    approved_by_ref VARCHAR(191),
    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    priority VARCHAR(32) NOT NULL DEFAULT 'normal',
    start_time TIMESTAMP WITH TIME ZONE NOT NULL, -- Maps to schema.org/startTime
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,   -- Maps to schema.org/endTime
    is_asap BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT,                -- Maps to schema.org/description
    external_source VARCHAR(191),
    external_ref VARCHAR(191),
    schema_org JSONB,                -- Standard JSON-LD representation
    metadata JSONB,                  -- Internal attributes
    approved_at TIMESTAMP WITH TIME ZONE,
    rejected_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 4. rent_action_items table (Line Items for RentAction)
CREATE TABLE rent_action_items (
    id BIGSERIAL PRIMARY KEY,
    rent_action_id BIGINT NOT NULL,
    item_kind VARCHAR(32) NOT NULL, -- 'item_type', 'kit_template'
    item_id BIGINT NOT NULL,        -- polymorphic link
    requested_quantity INT NOT NULL,
    allocated_quantity INT NOT NULL DEFAULT 0,
    notes TEXT,
    metadata JSONB,
    CONSTRAINT fk_rai_rent_action FOREIGN KEY (rent_action_id) REFERENCES rent_actions(id)
);

-- Indices for performance and windows
CREATE INDEX idx_assets_item_type_id ON assets(item_type_id);
CREATE INDEX idx_assets_status ON assets(status);
CREATE INDEX idx_rent_actions_status ON rent_actions(status);
CREATE INDEX idx_rent_actions_window ON rent_actions(start_time, end_time);
CREATE INDEX idx_rai_rent_action_id ON rent_action_items(rent_action_id);
