-- Create favorites table
CREATE TABLE IF NOT EXISTS favorites (
    id UUID PRIMARY KEY,
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,
    product_id UUID NOT NULL,
    notify_on_restock BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_favorites_customer_id ON favorites(customer_id);
CREATE INDEX idx_favorites_product_id ON favorites(product_id);
CREATE UNIQUE INDEX idx_favorites_customer_product ON favorites(customer_id, product_id);
CREATE INDEX idx_favorites_notify_on_restock ON favorites(notify_on_restock) WHERE notify_on_restock = TRUE;
