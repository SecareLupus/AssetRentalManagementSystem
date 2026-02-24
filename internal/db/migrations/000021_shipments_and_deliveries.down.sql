ALTER TABLE checkout_actions DROP COLUMN scheduled_delivery_id;
ALTER TABLE checkout_actions DROP COLUMN shipment_id;

ALTER TABLE return_actions DROP COLUMN shipment_id;

DROP TABLE IF EXISTS shipments;
DROP TABLE IF EXISTS scheduled_delivery_items;
DROP TABLE IF EXISTS scheduled_deliveries;
