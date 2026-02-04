-- Create build_specs table
CREATE TABLE build_specs (
    id BIGSERIAL PRIMARY KEY,
    version VARCHAR(64) NOT NULL UNIQUE,
    hardware_config JSONB,
    software_config JSONB,
    firmware_url TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Add current_build_spec_id to assets
ALTER TABLE assets ADD COLUMN current_build_spec_id BIGINT;
ALTER TABLE assets ADD CONSTRAINT fk_asset_build_spec FOREIGN KEY (current_build_spec_id) REFERENCES build_specs(id);

-- Create provision_actions table for audit log
CREATE TABLE provision_actions (
    id BIGSERIAL PRIMARY KEY,
    asset_id BIGINT NOT NULL,
    build_spec_id BIGINT,
    status VARCHAR(32) NOT NULL,
    performed_by VARCHAR(255) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT fk_pa_asset FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    CONSTRAINT fk_pa_build_spec FOREIGN KEY (build_spec_id) REFERENCES build_specs(id)
);

CREATE INDEX idx_pa_asset_id ON provision_actions(asset_id);
