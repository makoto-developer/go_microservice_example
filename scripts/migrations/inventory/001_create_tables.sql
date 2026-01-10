-- Inventory Service Database Schema
-- 在庫管理サービスのデータベーススキーマ

-- =========================================
-- Enums
-- =========================================

CREATE TYPE reservation_status AS ENUM (
  'RESERVED',
  'CONFIRMED',
  'RELEASED',
  'EXPIRED'
);

CREATE TYPE change_type AS ENUM (
  'INITIAL',
  'RESTOCK',
  'RETURN',
  'RESERVATION',
  'RELEASE',
  'CONFIRMATION',
  'DAMAGE',
  'STOCK_TAKING',
  'MANUAL_ADJUSTMENT'
);

-- =========================================
-- Table: inventories
-- 商品在庫情報
-- =========================================

CREATE TABLE inventories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id UUID NOT NULL,
  variation_id UUID,  -- オプショナル（バリエーション商品の場合）
  shop_id UUID NOT NULL,
  quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
  reserved_quantity INTEGER NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
  available_quantity INTEGER NOT NULL DEFAULT 0 CHECK (available_quantity >= 0),
  alert_threshold INTEGER NOT NULL DEFAULT 5 CHECK (alert_threshold >= 0),
  last_alerted_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

  -- 制約
  CONSTRAINT unique_inventory_per_product UNIQUE (product_id, variation_id, shop_id),
  CONSTRAINT check_available_quantity CHECK (available_quantity = quantity - reserved_quantity)
);

-- インデックス
CREATE INDEX idx_inventories_product_id ON inventories(product_id);
CREATE INDEX idx_inventories_shop_id ON inventories(shop_id);
CREATE INDEX idx_inventories_variation_id ON inventories(variation_id) WHERE variation_id IS NOT NULL;
CREATE INDEX idx_inventories_available_quantity ON inventories(available_quantity);
CREATE INDEX idx_inventories_alert_threshold ON inventories(alert_threshold);

-- =========================================
-- Table: reservations
-- 在庫引き当て情報
-- =========================================

CREATE TABLE reservations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  inventory_id UUID NOT NULL REFERENCES inventories(id) ON DELETE CASCADE,
  order_id UUID NOT NULL,
  quantity INTEGER NOT NULL CHECK (quantity > 0),
  status reservation_status NOT NULL DEFAULT 'RESERVED',
  expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
  confirmed_at TIMESTAMP WITH TIME ZONE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

  -- 制約
  CONSTRAINT unique_reservation_per_order UNIQUE (order_id, inventory_id)
);

-- インデックス
CREATE INDEX idx_reservations_inventory_id ON reservations(inventory_id);
CREATE INDEX idx_reservations_order_id ON reservations(order_id);
CREATE INDEX idx_reservations_status ON reservations(status);
CREATE INDEX idx_reservations_expires_at ON reservations(expires_at);

-- 有効期限切れの引き当てを検索するためのインデックス
CREATE INDEX idx_reservations_expired ON reservations(expires_at, status)
  WHERE status = 'RESERVED' AND expires_at < NOW();

-- =========================================
-- Table: inventory_history
-- 在庫変動履歴
-- =========================================

CREATE TABLE inventory_history (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  inventory_id UUID NOT NULL REFERENCES inventories(id) ON DELETE CASCADE,
  change_type change_type NOT NULL,
  change_quantity INTEGER NOT NULL,  -- 正数で加算、負数で減算
  quantity_before INTEGER NOT NULL,
  quantity_after INTEGER NOT NULL,
  reason TEXT NOT NULL,
  operator VARCHAR(255) NOT NULL,  -- 操作者（ユーザーID or "system"）
  order_id UUID,  -- 注文関連の場合
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- インデックス
CREATE INDEX idx_inventory_history_inventory_id ON inventory_history(inventory_id);
CREATE INDEX idx_inventory_history_created_at ON inventory_history(created_at);
CREATE INDEX idx_inventory_history_change_type ON inventory_history(change_type);
CREATE INDEX idx_inventory_history_order_id ON inventory_history(order_id) WHERE order_id IS NOT NULL;

-- パーティショニング用（将来的に）
CREATE INDEX idx_inventory_history_created_at_month ON inventory_history(DATE_TRUNC('month', created_at));

-- =========================================
-- Table: stock_takings
-- 在庫棚卸し記録
-- =========================================

CREATE TABLE stock_takings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  inventory_id UUID NOT NULL REFERENCES inventories(id) ON DELETE CASCADE,
  shop_id UUID NOT NULL,
  system_quantity INTEGER NOT NULL,  -- システム上の在庫数
  actual_quantity INTEGER NOT NULL,  -- 実在庫数
  difference INTEGER NOT NULL,  -- 差異（actual - system）
  difference_reason TEXT,
  operator VARCHAR(255) NOT NULL,  -- 棚卸し実施者
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- インデックス
CREATE INDEX idx_stock_takings_inventory_id ON stock_takings(inventory_id);
CREATE INDEX idx_stock_takings_shop_id ON stock_takings(shop_id);
CREATE INDEX idx_stock_takings_created_at ON stock_takings(created_at);

-- =========================================
-- Functions & Triggers
-- =========================================

-- 在庫更新時のタイムスタンプ自動更新
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_inventories_updated_at
  BEFORE UPDATE ON inventories
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reservations_updated_at
  BEFORE UPDATE ON reservations
  FOR EACH ROW
  EXECUTE FUNCTION update_updated_at_column();

-- 在庫数変更時に available_quantity を自動計算
CREATE OR REPLACE FUNCTION calculate_available_quantity()
RETURNS TRIGGER AS $$
BEGIN
  NEW.available_quantity = NEW.quantity - NEW.reserved_quantity;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER calculate_inventories_available_quantity
  BEFORE INSERT OR UPDATE ON inventories
  FOR EACH ROW
  EXECUTE FUNCTION calculate_available_quantity();

-- 在庫変更履歴の自動記録
CREATE OR REPLACE FUNCTION record_inventory_history()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    INSERT INTO inventory_history (
      inventory_id,
      change_type,
      change_quantity,
      quantity_before,
      quantity_after,
      reason,
      operator
    ) VALUES (
      NEW.id,
      'INITIAL',
      NEW.quantity,
      0,
      NEW.quantity,
      'Initial inventory registration',
      'system'
    );
  ELSIF TG_OP = 'UPDATE' AND OLD.quantity != NEW.quantity THEN
    INSERT INTO inventory_history (
      inventory_id,
      change_type,
      change_quantity,
      quantity_before,
      quantity_after,
      reason,
      operator
    ) VALUES (
      NEW.id,
      'MANUAL_ADJUSTMENT',  -- デフォルト、実際の変更タイプはアプリケーション側で設定
      NEW.quantity - OLD.quantity,
      OLD.quantity,
      NEW.quantity,
      'Quantity updated',
      'system'
    );
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER record_inventories_history
  AFTER INSERT OR UPDATE ON inventories
  FOR EACH ROW
  EXECUTE FUNCTION record_inventory_history();

-- =========================================
-- Initial Data (Optional)
-- =========================================

-- 開発環境用のサンプルデータ（必要に応じて削除）
-- INSERT INTO inventories (product_id, shop_id, quantity, alert_threshold)
-- VALUES
--   ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 100, 10),
--   ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', 50, 5);

-- =========================================
-- Comments
-- =========================================

COMMENT ON TABLE inventories IS '商品在庫情報';
COMMENT ON COLUMN inventories.product_id IS '商品ID（Shop Serviceから参照）';
COMMENT ON COLUMN inventories.variation_id IS 'バリエーションID（オプショナル）';
COMMENT ON COLUMN inventories.shop_id IS 'ショップID';
COMMENT ON COLUMN inventories.quantity IS '総在庫数';
COMMENT ON COLUMN inventories.reserved_quantity IS '引き当て済み在庫数';
COMMENT ON COLUMN inventories.available_quantity IS '利用可能在庫数（quantity - reserved_quantity）';
COMMENT ON COLUMN inventories.alert_threshold IS '在庫アラート閾値';
COMMENT ON COLUMN inventories.last_alerted_at IS '最後にアラートを送信した日時';

COMMENT ON TABLE reservations IS '在庫引き当て情報';
COMMENT ON COLUMN reservations.inventory_id IS '在庫ID';
COMMENT ON COLUMN reservations.order_id IS '注文ID（Order Serviceから参照）';
COMMENT ON COLUMN reservations.quantity IS '引き当て数量';
COMMENT ON COLUMN reservations.status IS '引き当てステータス';
COMMENT ON COLUMN reservations.expires_at IS '引き当て有効期限（30分）';
COMMENT ON COLUMN reservations.confirmed_at IS '確定日時（決済完了時）';

COMMENT ON TABLE inventory_history IS '在庫変動履歴';
COMMENT ON COLUMN inventory_history.change_type IS '変動タイプ';
COMMENT ON COLUMN inventory_history.change_quantity IS '変動数量（正数で加算、負数で減算）';
COMMENT ON COLUMN inventory_history.quantity_before IS '変動前在庫数';
COMMENT ON COLUMN inventory_history.quantity_after IS '変動後在庫数';
COMMENT ON COLUMN inventory_history.reason IS '変動理由';
COMMENT ON COLUMN inventory_history.operator IS '操作者（ユーザーID or "system"）';

COMMENT ON TABLE stock_takings IS '在庫棚卸し記録';
COMMENT ON COLUMN stock_takings.system_quantity IS 'システム上の在庫数';
COMMENT ON COLUMN stock_takings.actual_quantity IS '実在庫数';
COMMENT ON COLUMN stock_takings.difference IS '差異（actual - system）';
COMMENT ON COLUMN stock_takings.difference_reason IS '差異理由';
