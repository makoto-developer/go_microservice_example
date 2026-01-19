CREATE TYPE reservation_status AS ENUM ('pending', 'confirmed', 'released', 'expired');

CREATE TABLE IF NOT EXISTS reservations (
    id UUID PRIMARY KEY,
    inventory_id UUID NOT NULL REFERENCES inventories(id) ON DELETE CASCADE,
    order_id UUID NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    status reservation_status NOT NULL DEFAULT 'pending',
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_reservations_inventory_id ON reservations(inventory_id);
CREATE INDEX idx_reservations_order_id ON reservations(order_id);
CREATE INDEX idx_reservations_status ON reservations(status);
CREATE INDEX idx_reservations_expires_at ON reservations(expires_at) WHERE status = 'pending';
