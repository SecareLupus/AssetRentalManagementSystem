CREATE TABLE scheduled_deliveries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    target_date DATETIME NOT NULL,
    notes TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE scheduled_delivery_items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scheduled_delivery_id INTEGER NOT NULL,
    item_kind TEXT NOT NULL,
    item_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    FOREIGN KEY(scheduled_delivery_id) REFERENCES scheduled_deliveries(id) ON DELETE CASCADE
);

CREATE TABLE shipments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    scheduled_delivery_id INTEGER,
    provider_id INTEGER NOT NULL,
    ship_date DATETIME NOT NULL,
    carrier TEXT,
    tracking_number TEXT,
    status TEXT NOT NULL DEFAULT 'Preparing',
    notes TEXT,
    direction TEXT NOT NULL CHECK(direction IN ('outbound', 'inbound')),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(scheduled_delivery_id) REFERENCES scheduled_deliveries(id) ON DELETE SET NULL
);

ALTER TABLE checkout_actions ADD COLUMN shipment_id INTEGER REFERENCES shipments(id) ON DELETE SET NULL;
ALTER TABLE checkout_actions ADD COLUMN scheduled_delivery_id INTEGER REFERENCES scheduled_deliveries(id) ON DELETE SET NULL;

ALTER TABLE return_actions ADD COLUMN shipment_id INTEGER REFERENCES shipments(id) ON DELETE SET NULL;
