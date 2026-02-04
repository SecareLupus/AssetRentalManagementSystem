-- Add fleet lifecycle fields to assets table
ALTER TABLE assets ADD COLUMN build_spec_version VARCHAR(64);
ALTER TABLE assets ADD COLUMN provisioning_status VARCHAR(32) DEFAULT 'unprovisioned';
ALTER TABLE assets ADD COLUMN firmware_version VARCHAR(64);
ALTER TABLE assets ADD COLUMN hostname VARCHAR(191);
ALTER TABLE assets ADD COLUMN mesh_central_id VARCHAR(191);
ALTER TABLE assets ADD COLUMN last_inspection_at TIMESTAMP WITH TIME ZONE;

-- Create maintenance_logs table
CREATE TABLE maintenance_logs (
    id BIGSERIAL PRIMARY KEY,
    asset_id BIGINT NOT NULL,
    action_type VARCHAR(32) NOT NULL,
    notes TEXT,
    performed_by VARCHAR(191) NOT NULL,
    test_bits JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_ml_asset FOREIGN KEY (asset_id) REFERENCES assets(id)
);

CREATE INDEX idx_ml_asset_id ON maintenance_logs(asset_id);
