CREATE TABLE IF NOT EXISTS inventories (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL,
    shop_id UUID NOT NULL,
    quantity INTEGER NOT NULL CHECK (quantity >= 0),
    reserved_quantity INTEGER NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT check_reserved_lte_quantity CHECK (reserved_quantity <= quantity)
);

CREATE UNIQUE INDEX idx_inventories_product_shop ON inventories(product_id, shop_id);
CREATE INDEX idx_inventories_shop_id ON inventories(shop_id);
CREATE INDEX idx_inventories_product_id ON inventories(product_id);
