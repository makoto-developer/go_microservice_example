-- Create cart_items table
CREATE TABLE IF NOT EXISTS cart_items (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    variation_id UUID,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_cart_items_customer_id ON cart_items(customer_id);
CREATE INDEX idx_cart_items_expires_at ON cart_items(expires_at);
CREATE UNIQUE INDEX idx_cart_items_customer_product ON cart_items(customer_id, product_id, COALESCE(variation_id, '00000000-0000-0000-0000-000000000000'::UUID));

-- Create guest_cart_items table
CREATE TABLE IF NOT EXISTS guest_cart_items (
    id UUID PRIMARY KEY,
    session_id VARCHAR(255) NOT NULL,
    product_id UUID NOT NULL,
    variation_id UUID,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_guest_cart_items_session_id ON guest_cart_items(session_id);
CREATE INDEX idx_guest_cart_items_expires_at ON guest_cart_items(expires_at);
CREATE UNIQUE INDEX idx_guest_cart_items_session_product ON guest_cart_items(session_id, product_id, COALESCE(variation_id, '00000000-0000-0000-0000-000000000000'::UUID));
