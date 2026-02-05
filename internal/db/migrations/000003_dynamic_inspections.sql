-- Rename mesh_central_id to remote_management_id for abstraction
ALTER TABLE assets RENAME COLUMN mesh_central_id TO remote_management_id;

-- Add supported_features to item_types for per-category lifecycle toggles
ALTER TABLE item_types ADD COLUMN supported_features JSONB DEFAULT '{}';

-- Dynamic Inspection System Tables
CREATE TABLE inspection_templates (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE inspection_fields (
    id BIGSERIAL PRIMARY KEY,
    template_id BIGINT NOT NULL,
    label VARCHAR(255) NOT NULL,
    field_type VARCHAR(32) NOT NULL, -- boolean, text, image
    required BOOLEAN DEFAULT TRUE,
    display_order INT DEFAULT 0,
    CONSTRAINT fk_if_template FOREIGN KEY (template_id) REFERENCES inspection_templates(id) ON DELETE CASCADE
);

CREATE TABLE item_type_inspections (
    item_type_id BIGINT NOT NULL,
    template_id BIGINT NOT NULL,
    PRIMARY KEY (item_type_id, template_id),
    CONSTRAINT fk_iti_item_type FOREIGN KEY (item_type_id) REFERENCES item_types(id) ON DELETE CASCADE,
    CONSTRAINT fk_iti_template FOREIGN KEY (template_id) REFERENCES inspection_templates(id) ON DELETE CASCADE
);

CREATE TABLE inspection_submissions (
    id BIGSERIAL PRIMARY KEY,
    asset_id BIGINT NOT NULL,
    template_id BIGINT NOT NULL,
    performed_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_is_asset FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    CONSTRAINT fk_is_template FOREIGN KEY (template_id) REFERENCES inspection_templates(id) ON DELETE CASCADE
);

CREATE TABLE inspection_responses (
    id BIGSERIAL PRIMARY KEY,
    submission_id BIGINT NOT NULL,
    field_id BIGINT NOT NULL,
    response_value TEXT, -- Stores boolean (true/false), text, or image URL/path
    CONSTRAINT fk_ir_submission FOREIGN KEY (submission_id) REFERENCES inspection_submissions(id) ON DELETE CASCADE,
    CONSTRAINT fk_ir_field FOREIGN KEY (field_id) REFERENCES inspection_fields(id) ON DELETE CASCADE
);

CREATE INDEX idx_iti_item_type ON item_type_inspections(item_type_id);
CREATE INDEX idx_is_asset ON inspection_submissions(asset_id);
