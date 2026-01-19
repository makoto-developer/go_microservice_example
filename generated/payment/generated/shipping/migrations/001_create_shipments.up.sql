CREATE TYPE shipment_status AS ENUM ('pending', 'preparing', 'shipped', 'in_transit', 'delivered', 'failed');

CREATE TABLE IF NOT EXISTS shipments (
    id UUID PRIMARY KEY,
    order_id UUID NOT NULL UNIQUE,
    customer_id UUID NOT NULL,
    status shipment_status NOT NULL DEFAULT 'pending',
    tracking_number VARCHAR(100),
    carrier VARCHAR(100) NOT NULL,
    shipping_address TEXT NOT NULL,
    estimated_delivery TIMESTAMP NOT NULL,
    actual_delivery TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_shipments_order_id ON shipments(order_id);
CREATE INDEX idx_shipments_customer_id ON shipments(customer_id);
CREATE INDEX idx_shipments_status ON shipments(status);
CREATE INDEX idx_shipments_tracking_number ON shipments(tracking_number);
