-- ============================================
-- Shop Service - Test Data Schema
-- ============================================
-- Description: テストデータ専用スキーマとデータ投入
-- Version: 1.0.0
-- Created: 2026-01-11
-- Note: 開発用データと混在しないように、test_dataスキーマで管理

-- ============================================
-- Create Test Data Schema
-- ============================================

CREATE SCHEMA IF NOT EXISTS test_data;

-- ============================================
-- Grant Privileges
-- ============================================

GRANT ALL PRIVILEGES ON SCHEMA test_data TO admin;

-- ============================================
-- Test Data Insertion
-- ============================================

-- テストショップ（test_dataスキーマ内ではなく、publicスキーマに挿入）
-- ※ ショップはpublicスキーマのテーブルを使用

INSERT INTO shops (id, owner_id, name, description, logo_url, owner_name, phone_number, business_hours, return_policy, status, published) VALUES
(
  '11111111-1111-1111-1111-111111111111',
  '00000000-0000-0000-0000-000000000001',  -- テスト用owner UUID
  'テクノショップ',
  '最新テクノロジー商品を取り扱うショップです',
  '',
  '山田太郎',
  '03-1234-5678',
  '10:00-20:00',
  '30日間返品可能',
  'APPROVED',
  TRUE
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO shops (id, owner_id, name, description, logo_url, owner_name, phone_number, business_hours, return_policy, status, published) VALUES
(
  '22222222-2222-2222-2222-222222222222',
  '00000000-0000-0000-0000-000000000002',  -- テスト用owner UUID
  'ファッションストア',
  'おしゃれな衣類を揃えています',
  '',
  '佐藤花子',
  '03-9876-5432',
  '11:00-21:00',
  '14日間返品可能',
  'APPROVED',
  TRUE
)
ON CONFLICT (id) DO NOTHING;

-- ショップカテゴリ
INSERT INTO shop_categories (shop_id, category_name) VALUES
('11111111-1111-1111-1111-111111111111', 'electronics'),
('11111111-1111-1111-1111-111111111111', 'gadgets'),
('22222222-2222-2222-2222-222222222222', 'fashion'),
('22222222-2222-2222-2222-222222222222', 'casual')
ON CONFLICT (shop_id, category_name) DO NOTHING;

-- テスト商品
-- ※ 商品もpublicスキーマのテーブルを使用（データベース設計上の理由）

INSERT INTO products (id, shop_id, name, description, price, category, stock_quantity, weight, size, jan_code, published, deleted) VALUES
(
  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
  '11111111-1111-1111-1111-111111111111',
  'ワイヤレスイヤホン Pro',
  '高音質Bluetooth 5.0対応のワイヤレスイヤホン。ノイズキャンセリング機能付き、防水IPX5対応で運動時も安心。最大24時間の連続再生が可能です。',
  29800.00,
  'electronics',
  50,
  0.050,
  'コンパクト',
  '4560123456789',
  TRUE,
  FALSE
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO products (id, shop_id, name, description, price, category, stock_quantity, weight, size, jan_code, published, deleted) VALUES
(
  'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
  '11111111-1111-1111-1111-111111111111',
  'スマートウォッチ X1',
  '健康管理機能付きスマートウォッチ。心拍数・歩数計測、睡眠モニタリング対応。スマートフォン通知機能も搭載し、日常生活をサポートします。',
  45000.00,
  'electronics',
  30,
  0.045,
  'フリーサイズ',
  '4560123456790',
  TRUE,
  FALSE
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO products (id, shop_id, name, description, price, category, stock_quantity, weight, size, jan_code, published, deleted) VALUES
(
  'cccccccc-cccc-cccc-cccc-cccccccccccc',
  '22222222-2222-2222-2222-222222222222',
  'カジュアルTシャツ',
  '100%コットンの快適な着心地。カラーバリエーション豊富で、普段着に最適。洗濯機で洗えてお手入れ簡単です。',
  3980.00,
  'fashion',
  100,
  0.200,
  'M, L, XL',
  '4560123456791',
  TRUE,
  FALSE
)
ON CONFLICT (id) DO NOTHING;

-- 商品画像
INSERT INTO product_images (product_id, url, display_order, thumbnail_200_url, thumbnail_400_url, thumbnail_800_url) VALUES
(
  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
  'https://example.com/products/wireless-earphone-pro-1.jpg',
  1,
  'https://example.com/products/thumbnails/wireless-earphone-pro-1-200.jpg',
  'https://example.com/products/thumbnails/wireless-earphone-pro-1-400.jpg',
  'https://example.com/products/thumbnails/wireless-earphone-pro-1-800.jpg'
)
ON CONFLICT (product_id, display_order) DO NOTHING;

INSERT INTO product_images (product_id, url, display_order, thumbnail_200_url, thumbnail_400_url, thumbnail_800_url) VALUES
(
  'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
  'https://example.com/products/wireless-earphone-pro-2.jpg',
  2,
  'https://example.com/products/thumbnails/wireless-earphone-pro-2-200.jpg',
  'https://example.com/products/thumbnails/wireless-earphone-pro-2-400.jpg',
  'https://example.com/products/thumbnails/wireless-earphone-pro-2-800.jpg'
)
ON CONFLICT (product_id, display_order) DO NOTHING;

-- 商品タグ
INSERT INTO product_tags (product_id, tag_name) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'ワイヤレス'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'ノイズキャンセリング'),
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '防水'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'スマートウォッチ'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '健康管理'),
('cccccccc-cccc-cccc-cccc-cccccccccccc', 'カジュアル'),
('cccccccc-cccc-cccc-cccc-cccccccccccc', 'コットン')
ON CONFLICT (product_id, tag_name) DO NOTHING;

-- 商品バリエーション
INSERT INTO product_variations (product_id, sku, attribute_name, attribute_value, price, stock_quantity) VALUES
('cccccccc-cccc-cccc-cccc-cccccccccccc', 'TSHIRT-M-BLK', 'サイズ-色', 'M-ブラック', 3980.00, 30),
('cccccccc-cccc-cccc-cccc-cccccccccccc', 'TSHIRT-M-WHT', 'サイズ-色', 'M-ホワイト', 3980.00, 25),
('cccccccc-cccc-cccc-cccc-cccccccccccc', 'TSHIRT-L-BLK', 'サイズ-色', 'L-ブラック', 3980.00, 20),
('cccccccc-cccc-cccc-cccc-cccccccccccc', 'TSHIRT-L-WHT', 'サイズ-色', 'L-ホワイト', 3980.00, 15),
('cccccccc-cccc-cccc-cccc-cccccccccccc', 'TSHIRT-XL-BLK', 'サイズ-色', 'XL-ブラック', 3980.00, 10)
ON CONFLICT (sku) DO NOTHING;

-- 注文（テスト用）
INSERT INTO orders (id, shop_id, customer_id, order_number, status, total_amount, shipping_address, payment_method, tracking_number, carrier) VALUES
(
  'dddddddd-dddd-dddd-dddd-dddddddddddd',
  '11111111-1111-1111-1111-111111111111',
  '10000000-0000-0000-0000-000000000001',
  'ORD-20260110-000001',
  'SHIPPED',
  74800.00,
  '東京都渋谷区1-2-3',
  'クレジットカード',
  '1234567890',
  'YAMATO'
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO orders (id, shop_id, customer_id, order_number, status, total_amount, shipping_address, payment_method, tracking_number, carrier) VALUES
(
  'eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee',
  '22222222-2222-2222-2222-222222222222',
  '10000000-0000-0000-0000-000000000002',
  'ORD-20260110-000002',
  'PREPARING',
  7960.00,
  '大阪府大阪市北区4-5-6',
  '代金引換',
  NULL,
  NULL
)
ON CONFLICT (id) DO NOTHING;

-- 注文明細
INSERT INTO order_items (order_id, product_id, product_name, quantity, unit_price, subtotal) VALUES
('dddddddd-dddd-dddd-dddd-dddddddddddd', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'ワイヤレスイヤホン Pro', 1, 29800.00, 29800.00),
('dddddddd-dddd-dddd-dddd-dddddddddddd', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'スマートウォッチ X1', 1, 45000.00, 45000.00),
('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'cccccccc-cccc-cccc-cccc-cccccccccccc', 'カジュアルTシャツ (M-ブラック)', 2, 3980.00, 7960.00)
ON CONFLICT DO NOTHING;

-- 売上レポート（テスト用）
INSERT INTO sales_reports (shop_id, date, total_sales, order_count, average_order_value) VALUES
('11111111-1111-1111-1111-111111111111', '2026-01-10', 74800.00, 1, 74800.00),
('22222222-2222-2222-2222-222222222222', '2026-01-10', 0.00, 0, 0.00)
ON CONFLICT (shop_id, date) DO NOTHING;

-- ============================================
-- Test Data Summary
-- ============================================

\echo ''
\echo '========================================='
\echo 'Test Data Schema Initialized'
\echo '========================================='
\echo ''
\echo 'Created schema: test_data'
\echo ''
\echo 'Inserted test data:'
\echo '  - 2 shops (テクノショップ, ファッションストア)'
\echo '  - 3 products (ワイヤレスイヤホン Pro, スマートウォッチ X1, カジュアルTシャツ)'
\echo '  - 2 orders (発送済み、準備中)'
\echo ''
\echo 'Note: テストデータは ON CONFLICT DO NOTHING により'
\echo '      既存データがあれば上書きせずスキップします'
\echo '========================================='
