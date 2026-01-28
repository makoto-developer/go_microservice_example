ALTER TABLE shops
ADD COLUMN approved_at TIMESTAMP,
ADD COLUMN approved_by UUID;

CREATE INDEX idx_shops_approved_at ON shops(approved_at);
