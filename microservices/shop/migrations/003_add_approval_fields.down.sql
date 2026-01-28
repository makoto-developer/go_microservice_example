DROP INDEX IF EXISTS idx_shops_approved_at;

ALTER TABLE shops
DROP COLUMN IF EXISTS approved_by,
DROP COLUMN IF EXISTS approved_at;
