-- ============================================
-- Shop Service - Database Schema
-- ============================================
-- Description: ショップサービスのデータベーススキーマ
-- Version: 1.0.0
-- Created: 2026-01-10

-- ============================================
-- ENUM Types
-- ============================================

-- ショップステータス
CREATE TYPE shop_status AS ENUM (
  'PENDING_APPROVAL',  -- 承認待ち
  'APPROVED',          -- 承認済み
  'REJECTED',          -- 却下
  'SUSPENDED'          -- 一時停止
);

-- 注文ステータス
CREATE TYPE order_status AS ENUM (
  'RECEIVED',   -- 受注
  'PREPARING',  -- 準備中
  'SHIPPED',    -- 発送済み
  'DELIVERED',  -- 配達完了
  'CANCELLED'   -- キャンセル
);

-- 配送業者
CREATE TYPE carrier AS ENUM (
  'YAMATO',      -- ヤマト運輸
  'SAGAWA',      -- 佐川急便
  'JAPAN_POST'   -- 日本郵便
);

-- ============================================
-- Tables
-- ============================================

-- ショップテーブル
CREATE TABLE shops (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_id UUID NOT NULL,  -- Auth Serviceのuser_idを参照
  name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  logo_url VARCHAR(500),
  owner_name VARCHAR(255) NOT NULL,
  phone_number VARCHAR(20) NOT NULL,
  business_hours VARCHAR(500) NOT NULL,
  return_policy TEXT NOT NULL,
  status shop_status NOT NULL DEFAULT 'PENDING_APPROVAL',
  published BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  
  CONSTRAINT check_phone_format CHECK (phone_number ~ '^[0-9\-+() ]+$'),
  CONSTRAINT check_published_requires_approved CHECK (
    published = FALSE OR status = 'APPROVED'
  )
);

-- ショップカテゴリテーブル
CREATE TABLE shop_categories (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  shop_id UUID NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
  category_name VARCHAR(100) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  
  CONSTRAINT unique_shop_category UNIQUE(shop_id, category_name)
);

-- 商品テーブル
CREATE TABLE products (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  shop_id UUID NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  price DECIMAL(10, 2) NOT NULL,
  category VARCHAR(100) NOT NULL,
  stock_quantity INTEGER NOT NULL DEFAULT 0,
  weight DECIMAL(10, 3),  -- kg単位
  size VARCHAR(100),
  jan_code VARCHAR(13),
  published BOOLEAN NOT NULL DEFAULT FALSE,
  deleted BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  
  CONSTRAINT check_price_positive CHECK (price >= 0),
  CONSTRAINT check_stock_non_negative CHECK (stock_quantity >= 0),
  CONSTRAINT check_weight_positive CHECK (weight IS NULL OR weight > 0),
  CONSTRAINT check_jan_code_format CHECK (
    jan_code IS NULL OR jan_code ~ '^[0-9]{8}$|^[0-9]{13}$'
  )
);

-- 商品画像テーブル
CREATE TABLE product_images (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  url VARCHAR(500) NOT NULL,
  display_order INTEGER NOT NULL,
  thumbnail_200_url VARCHAR(500) NOT NULL,
  thumbnail_400_url VARCHAR(500) NOT NULL,
  thumbnail_800_url VARCHAR(500) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  
  CONSTRAINT unique_product_display_order UNIQUE(product_id, display_order),
  CONSTRAINT check_display_order_positive CHECK (display_order > 0)
);

-- 商品タグテーブル
CREATE TABLE product_tags (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  tag_name VARCHAR(50) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  
  CONSTRAINT unique_product_tag UNIQUE(product_id, tag_name)
);

-- 商品バリエーションテーブル
CREATE TABLE product_variations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  sku VARCHAR(100) UNIQUE NOT NULL,
  attribute_name VARCHAR(50) NOT NULL,  -- e.g., "Size", "Color"
  attribute_value VARCHAR(100) NOT NULL,  -- e.g., "M", "Red"
  price DECIMAL(10, 2) NOT NULL,
  stock_quantity INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  
  CONSTRAINT check_variation_price_positive CHECK (price >= 0),
  CONSTRAINT check_variation_stock_non_negative CHECK (stock_quantity >= 0)
);

-- 注文テーブル
CREATE TABLE orders (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  shop_id UUID NOT NULL REFERENCES shops(id) ON DELETE RESTRICT,
  customer_id UUID NOT NULL,  -- Customer Serviceのcustomer_idを参照
  order_number VARCHAR(20) UNIQUE NOT NULL,
  status order_status NOT NULL DEFAULT 'RECEIVED',
  total_amount DECIMAL(10, 2) NOT NULL,
  shipping_address TEXT NOT NULL,
  payment_method VARCHAR(50) NOT NULL,
  tracking_number VARCHAR(100),
  carrier carrier,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  
  CONSTRAINT check_total_amount_positive CHECK (total_amount > 0),
  CONSTRAINT check_tracking_requires_carrier CHECK (
    tracking_number IS NULL OR carrier IS NOT NULL
  )
);

-- 注文明細テーブル
CREATE TABLE order_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  product_id UUID NOT NULL,  -- 商品削除後も履歴保持のためREFERENCESなし
  product_name VARCHAR(255) NOT NULL,  -- スナップショット
  quantity INTEGER NOT NULL,
  unit_price DECIMAL(10, 2) NOT NULL,
  subtotal DECIMAL(10, 2) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  
  CONSTRAINT check_quantity_positive CHECK (quantity > 0),
  CONSTRAINT check_unit_price_positive CHECK (unit_price >= 0),
  CONSTRAINT check_subtotal_positive CHECK (subtotal >= 0),
  CONSTRAINT check_subtotal_calculation CHECK (
    subtotal = quantity * unit_price
  )
);

-- 売上レポートテーブル
CREATE TABLE sales_reports (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  shop_id UUID NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
  date DATE NOT NULL,
  total_sales DECIMAL(12, 2) NOT NULL,
  order_count INTEGER NOT NULL,
  average_order_value DECIMAL(10, 2) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  
  CONSTRAINT unique_shop_date UNIQUE(shop_id, date),
  CONSTRAINT check_total_sales_non_negative CHECK (total_sales >= 0),
  CONSTRAINT check_order_count_non_negative CHECK (order_count >= 0),
  CONSTRAINT check_average_order_value_non_negative CHECK (average_order_value >= 0)
);

-- ============================================
-- Indexes
-- ============================================

-- Shops
CREATE INDEX idx_shops_owner_id ON shops(owner_id);
CREATE INDEX idx_shops_status ON shops(status);
CREATE INDEX idx_shops_published ON shops(published) WHERE published = TRUE;

-- Shop Categories
CREATE INDEX idx_shop_categories_shop_id ON shop_categories(shop_id);

-- Products
CREATE INDEX idx_products_shop_id ON products(shop_id);
CREATE INDEX idx_products_category ON products(category);
CREATE INDEX idx_products_published ON products(published) WHERE published = TRUE AND deleted = FALSE;
CREATE INDEX idx_products_jan_code ON products(jan_code) WHERE jan_code IS NOT NULL;

-- Product Images
CREATE INDEX idx_product_images_product_id ON product_images(product_id);

-- Product Tags
CREATE INDEX idx_product_tags_product_id ON product_tags(product_id);
CREATE INDEX idx_product_tags_tag_name ON product_tags(tag_name);

-- Product Variations
CREATE INDEX idx_product_variations_product_id ON product_variations(product_id);
CREATE INDEX idx_product_variations_sku ON product_variations(sku);

-- Orders
CREATE INDEX idx_orders_shop_id ON orders(shop_id);
CREATE INDEX idx_orders_customer_id ON orders(customer_id);
CREATE INDEX idx_orders_order_number ON orders(order_number);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at);

-- Order Items
CREATE INDEX idx_order_items_order_id ON order_items(order_id);
CREATE INDEX idx_order_items_product_id ON order_items(product_id);

-- Sales Reports
CREATE INDEX idx_sales_reports_shop_id_date ON sales_reports(shop_id, date);

-- ============================================
-- Triggers
-- ============================================

-- 自動更新: updated_at カラム
CREATE OR REPLACE FUNCTION update_shop_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_shops_updated_at
BEFORE UPDATE ON shops
FOR EACH ROW
EXECUTE FUNCTION update_shop_updated_at();

CREATE TRIGGER trigger_products_updated_at
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION update_shop_updated_at();

CREATE TRIGGER trigger_product_variations_updated_at
BEFORE UPDATE ON product_variations
FOR EACH ROW
EXECUTE FUNCTION update_shop_updated_at();

CREATE TRIGGER trigger_orders_updated_at
BEFORE UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION update_shop_updated_at();

-- 注文番号の自動生成
CREATE OR REPLACE FUNCTION generate_order_number()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.order_number IS NULL OR NEW.order_number = '' THEN
    NEW.order_number := 'ORD-' || TO_CHAR(NOW(), 'YYYYMMDD') || '-' || 
                       LPAD(nextval('order_number_seq')::TEXT, 6, '0');
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE SEQUENCE order_number_seq START 1;

CREATE TRIGGER trigger_generate_order_number
BEFORE INSERT ON orders
FOR EACH ROW
EXECUTE FUNCTION generate_order_number();

-- 売上レポートの自動計算
CREATE OR REPLACE FUNCTION update_sales_report()
RETURNS TRIGGER AS $$
BEGIN
  INSERT INTO sales_reports (shop_id, date, total_sales, order_count, average_order_value)
  SELECT 
    shop_id,
    DATE(created_at),
    SUM(total_amount),
    COUNT(*),
    AVG(total_amount)
  FROM orders
  WHERE shop_id = NEW.shop_id
    AND DATE(created_at) = DATE(NEW.created_at)
    AND status != 'CANCELLED'
  GROUP BY shop_id, DATE(created_at)
  ON CONFLICT (shop_id, date) DO UPDATE
  SET 
    total_sales = EXCLUDED.total_sales,
    order_count = EXCLUDED.order_count,
    average_order_value = EXCLUDED.average_order_value;
  
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_sales_report_on_order
AFTER INSERT OR UPDATE OF status ON orders
FOR EACH ROW
WHEN (NEW.status = 'DELIVERED')
EXECUTE FUNCTION update_sales_report();

-- ============================================
-- Test Data (開発用)
-- ============================================

-- テストショップ（owner_id は Auth Serviceのユーザーを想定）

-- ============================================
-- Comments
-- ============================================

COMMENT ON TABLE shops IS 'ショップ情報';
COMMENT ON TABLE shop_categories IS 'ショップカテゴリ';
COMMENT ON TABLE products IS '商品情報';
COMMENT ON TABLE product_images IS '商品画像';
COMMENT ON TABLE product_tags IS '商品タグ';
COMMENT ON TABLE product_variations IS '商品バリエーション（サイズ・色違い等）';
COMMENT ON TABLE orders IS '注文情報';
COMMENT ON TABLE order_items IS '注文明細';
COMMENT ON TABLE sales_reports IS '売上レポート（日次集計）';

COMMENT ON COLUMN shops.status IS 'ショップステータス（承認待ち、承認済み、却下、一時停止）';
COMMENT ON COLUMN shops.published IS '公開フラグ（承認済みのみ公開可能）';
COMMENT ON COLUMN products.deleted IS '論理削除フラグ（注文履歴保持のため物理削除しない）';
COMMENT ON COLUMN product_images.display_order IS '表示順序（1から開始）';
COMMENT ON COLUMN product_variations.sku IS '在庫管理単位コード（一意）';
COMMENT ON COLUMN orders.order_number IS '注文番号（自動生成: ORD-YYYYMMDD-XXXXXX）';
COMMENT ON COLUMN order_items.product_name IS '商品名スナップショット（商品削除後も履歴表示可能）';
COMMENT ON COLUMN sales_reports.date IS '集計日（日次）';
