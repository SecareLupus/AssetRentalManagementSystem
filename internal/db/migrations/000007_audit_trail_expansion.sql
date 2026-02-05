-- Migration: Audit Trail Expansion
-- Adds CreatedByUserID and UpdatedByUserID to item_types, assets, and rent_actions.

ALTER TABLE item_types ADD COLUMN created_by_user_id BIGINT;
ALTER TABLE item_types ADD COLUMN updated_by_user_id BIGINT;

ALTER TABLE assets ADD COLUMN updated_by_user_id BIGINT;

ALTER TABLE rent_actions ADD COLUMN updated_by_user_id BIGINT;

-- Foreign key constraints (optional, depending on if we want hard links)
-- ALTER TABLE item_types ADD CONSTRAINT fk_item_types_created_by FOREIGN KEY (created_by_user_id) REFERENCES users(id);
-- ALTER TABLE item_types ADD CONSTRAINT fk_item_types_updated_by FOREIGN KEY (updated_by_user_id) REFERENCES users(id);
-- ... etc
