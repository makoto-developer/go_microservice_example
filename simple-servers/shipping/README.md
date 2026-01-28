# Shipping Service

配送管理サービス - 注文の配送追跡、配送業者管理、配送ステータス管理を提供します。

## 概要

Shipping Serviceは、オンラインショップの配送業務を管理するマイクロサービスです。

### 主な機能
- 配送情報の登録・管理
- 配送追跡番号の生成・管理
- 配送ステータスの更新
- 配送履歴の追跡（トラッキングイベント）
- 配送業者（キャリア）の管理

## アーキテクチャ

### Database per Service Pattern
- **専用PostgreSQLインスタンス**: `postgres_shipping` (port 22016)
- **データベース**: `shipping_service`
- **gRPCサービスポート**: 22108

### データベーススキーマ

#### shipments テーブル
```sql
- id (UUID, Primary Key)
- order_id (UUID, Unique, 注文ID)
- tracking_number (VARCHAR(100), Unique, 追跡番号)
- carrier (VARCHAR(100), 配送業者名)
- status (VARCHAR(50), 配送ステータス: preparing/shipped/in_transit/delivered/failed)
- shipped_at (TIMESTAMP, 発送日時)
- delivered_at (TIMESTAMP, 配達完了日時)
- recipient_name (VARCHAR(255), 受取人名)
- recipient_phone (VARCHAR(50), 受取人電話番号)
- shipping_address (TEXT, 配送先住所)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

#### tracking_events テーブル
```sql
- id (UUID, Primary Key)
- shipment_id (UUID, Foreign Key → shipments.id ON DELETE CASCADE)
- status (VARCHAR(50), イベントステータス)
- location (VARCHAR(255), 現在地)
- description (TEXT, イベント詳細)
- event_time (TIMESTAMP, イベント発生時刻)
- created_at (TIMESTAMP)
```

### インデックス
- `idx_shipments_order_id` - 注文IDでの検索を高速化
- `idx_shipments_tracking_number` - 追跡番号での検索を高速化
- `idx_shipments_status` - ステータスでの検索を高速化（配送中の荷物一覧取得など）
- `idx_tracking_events_shipment_id` - 配送履歴の取得を高速化

## セットアップ

### 前提条件
- Go 1.25.0+
- Docker & Docker Compose
- PostgreSQL 16 (Dockerコンテナで提供)

### データベースセットアップ

1. PostgreSQLコンテナ起動（infrastructure/docker）
```bash
cd ../../infrastructure/docker
docker-compose up -d postgres_shipping
```

2. パスワード設定
```bash
docker-compose exec postgres_shipping psql -U postgres -c "ALTER USER postgres WITH PASSWORD 'postgres_password';"
```

3. スキーマ作成
```bash
docker-compose exec -T postgres_shipping psql -U postgres -d shipping_service < ../../simple-servers/shipping/schema.sql
```

### ビルド・起動

```bash
cd simple-servers/shipping

# ビルド
go build -o shipping

# 起動
./shipping
```

### 環境変数

```bash
SHIPPING_DATABASE_URL=postgresql://postgres:postgres_password@localhost:22016/shipping_service?sslmode=disable
SHIPPING_SERVICE_PORT=22108
```

## 検証

```bash
# 包括的な検証スクリプト実行
./verify.sh
```

検証内容:
- データベース接続確認
- テーブル存在確認（shipments, tracking_events）
- インデックス確認（4つのインデックス）
- 外部キー制約確認
- サービスビルド確認
- サービス起動確認

## API設計（将来実装）

### gRPC Services（予定）

```protobuf
service ShippingService {
  rpc CreateShipment(CreateShipmentRequest) returns (CreateShipmentResponse);
  rpc GetShipment(GetShipmentRequest) returns (GetShipmentResponse);
  rpc UpdateShipmentStatus(UpdateShipmentStatusRequest) returns (UpdateShipmentStatusResponse);
  rpc GetTrackingEvents(GetTrackingEventsRequest) returns (GetTrackingEventsResponse);
  rpc AddTrackingEvent(AddTrackingEventRequest) returns (AddTrackingEventResponse);
}
```

## 他サービスとの連携

### 依存関係
- **Order Service**: 注文完了時に配送情報を受け取る
- **Notification Service**: 配送ステータス変更時に通知を送信

### イベント駆動（将来実装）
- `OrderCompleted` イベントを購読 → 配送情報作成
- `ShipmentStatusChanged` イベントを発行 → 通知サービスへ
- `ShipmentDelivered` イベントを発行 → 注文サービスへ

## 配送ステータスフロー

```
preparing (準備中)
    ↓
shipped (発送済み)
    ↓
in_transit (輸送中)
    ↓
delivered (配達完了)

※ failed (配達失敗) への遷移も可能
```

## 開発状況

### ✅ 完了
- [x] データベーススキーマ設計
- [x] PostgreSQL専用インスタンス設定
- [x] テーブル・インデックス作成
- [x] 外部キー制約設定
- [x] Goサービススケルトン作成
- [x] データベース接続確認
- [x] ビルド・起動確認
- [x] 検証スクリプト作成

### 🚧 今後の実装予定
- [ ] gRPC API実装
- [ ] Repository層実装
- [ ] Usecase層実装
- [ ] イベント駆動実装（Order Service連携）
- [ ] 追跡番号生成ロジック
- [ ] 配送業者API連携（ヤマト運輸、佐川急便など）
- [ ] ユニットテスト
- [ ] 統合テスト

## トラブルシューティング

### データベース接続エラー
```bash
# コンテナ状態確認
cd ../../infrastructure/docker
docker-compose ps postgres_shipping

# ログ確認
docker-compose logs postgres_shipping

# 接続テスト
docker-compose exec postgres_shipping psql -U postgres -d shipping_service -c "SELECT 1;"
```

### ポート競合
デフォルトポート:
- Database: 22016
- gRPC Service: 22108

他のサービスと競合する場合は環境変数で変更してください。

## 参考資料
- [Database per Service Pattern](../../docs/architecture/database-per-service.md)
- [配送管理要件定義](../../docs/requirements/shipping_service.md)
