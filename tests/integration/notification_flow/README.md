# Notification/Review/Shipping Flow Integration Tests

このディレクトリには、Notification、Review、Shipping サービスの結合テストが含まれています。

## 📁 ディレクトリ構成

```
notification_flow/
├── clients/                    # サービスクライアント
│   ├── notification_client.go  # Notification Service クライアント
│   ├── review_client.go        # Review Service クライアント
│   └── shipping_client.go      # Shipping Service クライアント
├── notification_test.go        # 通知フローテスト
├── review_test.go             # レビューフローテスト
├── shipping_test.go           # 配送フローテスト
├── go.mod                     # Go モジュール定義
├── run_all_tests.sh           # テスト実行スクリプト
└── README.md                  # このファイル
```

## 🎯 テスト対象

### Notification Service テスト

1. **注文確認通知** (`TestOrderConfirmationNotification`)
   - `order_confirmation` テンプレート取得
   - 通知レコード作成（status: pending）
   - ステータス更新（pending → sent）

2. **決済完了通知** (`TestPaymentSuccessNotification`)
   - `payment_success` テンプレート取得
   - 決済完了通知送信

3. **配送開始通知** (`TestShippingUpdateNotification`)
   - `shipping_update` テンプレート取得
   - 配送開始通知送信

4. **通知失敗ハンドリング** (`TestNotificationFailureHandling`)
   - 通知送信失敗シナリオ
   - エラーメッセージ記録

5. **受信者別通知取得** (`TestMultipleNotificationsForRecipient`)
   - 受信者ごとの通知一覧取得

### Review Service テスト

1. **レビュー作成** (`TestCreateReview`)
   - レビュー投稿
   - 評価（rating 1-5）
   - コメント

2. **評価バリデーション** (`TestReviewRatingValidation`)
   - 無効な評価（0, 6）の拒否
   - 有効な評価（1-5）の許可

3. **商品別レビュー取得** (`TestGetReviewsByProduct`)
   - 商品IDでレビュー一覧取得
   - 作成日時順ソート

4. **平均評価計算** (`TestAverageRating`)
   - 商品の平均評価算出

5. **レビュー更新** (`TestUpdateReview`)
   - レビュー内容の変更
   - 編集可能期間の検証

6. **編集可能期限** (`TestReviewEditableUntil`)
   - `editable_until` タイムスタンプ検証

7. **重複レビュー防止** (`TestDuplicateReviewPrevention`)
   - 同一注文・商品の重複投稿拒否

### Shipping Service テスト

1. **配送情報作成** (`TestCreateShipment`)
   - shipment レコード作成
   - 初期ステータス: preparing

2. **配送ステータスフロー** (`TestShipmentStatusFlow`)
   - ステータス遷移: preparing → shipped → in_transit → delivered
   - tracking_number 発行
   - 配送業者設定

3. **追跡イベント** (`TestTrackingEvents`)
   - tracking_events にイベント記録
   - 時系列順ソート

4. **配送と通知の連携** (`TestShippingWithNotification`)
   - Shipping Service で配送開始
   - Notification Service で通知送信

5. **追跡イベントの時系列検証** (`TestMultipleTrackingEventsChronology`)
   - イベントの時間順整合性

6. **注文IDの一意性** (`TestShipmentUniqueOrderID`)
   - 1注文につき1配送情報の制約

## 🚀 実行方法

### 前提条件

各サービスのデータベースが起動していること:

- Notification Service: PostgreSQL on port `22017`
- Review Service: PostgreSQL on port `22018`
- Shipping Service: PostgreSQL on port `22016`

### 全テスト実行

```bash
cd tests/integration/notification_flow
./run_all_tests.sh
```

### 個別テスト実行

```bash
# Notification テスト
go test -v -run TestOrderConfirmationNotification

# Review テスト
go test -v -run TestCreateReview

# Shipping テスト
go test -v -run TestShipmentStatusFlow
```

### カバレッジ確認

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📊 テストシナリオ詳細

### シナリオ 1: 注文から配送までの通知フロー

```
1. 注文作成 (Order Service)
   ↓
2. 注文確認通知 (Notification Service)
   - template: order_confirmation
   - status: pending → sent
   ↓
3. 決済完了 (Payment Service)
   ↓
4. 決済完了通知 (Notification Service)
   - template: payment_success
   ↓
5. 配送情報作成 (Shipping Service)
   - status: preparing
   - tracking_number 発行
   ↓
6. 配送開始通知 (Notification Service)
   - template: shipping_update
   ↓
7. 追跡イベント記録 (Shipping Service)
   - preparing → shipped → in_transit → delivered
```

### シナリオ 2: レビュー投稿フロー

```
1. 注文完了・商品受取
   ↓
2. レビュー投稿 (Review Service)
   - rating: 1-5
   - review_text
   - editable_until: 30日後
   ↓
3. レビュー取得
   - 商品別レビュー一覧
   - 平均評価計算
   ↓
4. レビュー更新（編集可能期間内）
   - rating/review_text 変更
```

## 🗄️ データベーススキーマ

### Notification Service

```sql
-- notification_templates
id, name, channel, subject, body_template

-- notifications
id, template_id, recipient, channel, subject, body, status, sent_at, error_message
```

### Review Service

```sql
-- reviews
id, customer_id, product_id, order_id, rating, review_text, editable_until
```

### Shipping Service

```sql
-- shipments
id, order_id, tracking_number, carrier, status, shipped_at, delivered_at,
recipient_name, recipient_phone, shipping_address

-- tracking_events
id, shipment_id, status, location, description, event_time
```

## 🔧 トラブルシューティング

### データベース接続エラー

```bash
# データベース起動確認
docker ps | grep postgres

# ポート確認
netstat -an | grep 22017
netstat -an | grep 22018
netstat -an | grep 22016
```

### テスト失敗時

```bash
# 詳細ログ出力
go test -v -run TestName

# 特定のテストのみ実行
go test -v -run "TestNotification"
```

### データクリーンアップ

テストデータが残っている場合:

```sql
-- Notification Service
DELETE FROM notifications WHERE recipient LIKE '%test%' OR recipient LIKE '%example.com';

-- Review Service
DELETE FROM reviews WHERE review_text LIKE 'Test%';

-- Shipping Service
DELETE FROM tracking_events WHERE shipment_id IN (
    SELECT id FROM shipments WHERE recipient_name LIKE 'Test%'
);
DELETE FROM shipments WHERE recipient_name LIKE 'Test%';
```

## 📝 テスト追加ガイド

新しいテストを追加する場合:

1. **クライアント拡張** (`clients/`)
   - 新しいメソッドを追加
   - エラーハンドリング実装

2. **テストケース作成**
   - `*_test.go` ファイルに追加
   - テスト名: `Test<機能名>`
   - AAA パターン: Arrange, Act, Assert

3. **テストスクリプト更新**
   - `run_all_tests.sh` に新しいテスト追加

## 🎓 ベストプラクティス

1. **独立性**: 各テストは他のテストに依存しない
2. **クリーンアップ**: テスト後はデータを削除
3. **べき等性**: 同じテストを何度実行しても同じ結果
4. **明確な命名**: テスト名から目的がわかる
5. **適切なアサーション**: `require` vs `assert` を使い分け

## 📈 今後の拡張

- [ ] gRPC エンドポイントのテスト
- [ ] イベント駆動テスト（RabbitMQ連携）
- [ ] パフォーマンステスト
- [ ] エンドツーエンドテスト（全サービス連携）
- [ ] カオスエンジニアリング（障害注入）

## 📚 参考資料

- [Notification Service README](../../../simple-servers/notification/README.md)
- [Review Service README](../../../simple-servers/review/README.md)
- [Shipping Service README](../../../simple-servers/shipping/README.md)
