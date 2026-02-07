-- Migration 000012: Logistics Modernization (Reservations & Demands)

-- 1. Rename rent_actions to rental_reservations
ALTER TABLE rent_actions RENAME TO rental_reservations;
ALTER TABLE rental_reservations RENAME COLUMN status TO reservation_status;

-- 2. Add new columns to rental_reservations for Schema.org alignment
ALTER TABLE rental_reservations ADD COLUMN reservation_name VARCHAR(191);
ALTER TABLE rental_reservations ADD COLUMN under_name_id BIGINT REFERENCES people(id);
ALTER TABLE rental_reservations ADD COLUMN booking_time TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE rental_reservations ADD COLUMN provider_id BIGINT REFERENCES companies(id);

-- Update existing data: map status to schema.org values
-- pending -> ReservationPending
-- approved -> ReservationConfirmed
-- rejected -> ReservationCancelled
-- cancelled -> ReservationCancelled
-- fulfilled -> ReservationFulfilled
UPDATE rental_reservations SET reservation_status = 'ReservationPending' WHERE reservation_status = 'pending';
UPDATE rental_reservations SET reservation_status = 'ReservationConfirmed' WHERE reservation_status = 'approved';
UPDATE rental_reservations SET reservation_status = 'ReservationCancelled' WHERE reservation_status IN ('rejected', 'cancelled');
UPDATE rental_reservations SET reservation_status = 'ReservationFulfilled' WHERE reservation_status = 'fulfilled';
UPDATE rental_reservations SET booking_time = created_at;

-- 3. Rename rent_action_items to demands and align
ALTER TABLE rent_action_items RENAME TO demands;
ALTER TABLE demands RENAME COLUMN rent_action_id TO reservation_id;
ALTER TABLE demands ADD COLUMN event_id BIGINT REFERENCES events(id);
ALTER TABLE demands ADD COLUMN business_function VARCHAR(191);
ALTER TABLE demands ADD COLUMN eligible_duration VARCHAR(64);
ALTER TABLE demands ADD COLUMN place_id BIGINT REFERENCES places(id);
ALTER TABLE demands ADD COLUMN created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;
ALTER TABLE demands ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP;


-- 4. Port event_asset_needs to demands
-- Map item_type_id to item_id and set kind to 'item_type'
INSERT INTO demands (event_id, item_kind, item_id, requested_quantity, place_id, metadata, created_at, updated_at)
SELECT event_id, 'item_type', item_type_id, quantity, place_id, metadata, created_at, updated_at
FROM event_asset_needs;

-- Cleanup event_asset_needs
DROP TABLE event_asset_needs;

-- 5. Create check_out_actions
CREATE TABLE check_out_actions (
    id BIGSERIAL PRIMARY KEY,
    reservation_id BIGINT REFERENCES rental_reservations(id),
    asset_id BIGINT NOT NULL REFERENCES assets(id),
    agent_id BIGINT REFERENCES people(id),
    recipient_id BIGINT REFERENCES people(id),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    from_location_id BIGINT REFERENCES places(id),
    to_location_id BIGINT REFERENCES places(id),
    action_status VARCHAR(64), -- schema.org/ActionStatusType
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 6. Create return_actions
CREATE TABLE return_actions (
    id BIGSERIAL PRIMARY KEY,
    reservation_id BIGINT REFERENCES rental_reservations(id),
    asset_id BIGINT NOT NULL REFERENCES assets(id),
    agent_id BIGINT REFERENCES people(id),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    from_location_id BIGINT REFERENCES places(id),
    to_location_id BIGINT REFERENCES places(id),
    action_status VARCHAR(64),
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indices
CREATE INDEX idx_reservations_status ON rental_reservations(reservation_status);
CREATE INDEX idx_demands_reservation ON demands(reservation_id);
CREATE INDEX idx_demands_event ON demands(event_id);
CREATE INDEX idx_checkout_reservation ON check_out_actions(reservation_id);
CREATE INDEX idx_return_reservation ON return_actions(reservation_id);
