# Notification/Review/Shipping Flow Tests - Quick Start Guide

## 🚀 5分で始める

### Step 1: ディレクトリ移動

```bash
cd /Users/user/work/repositories/github.com/makoto-developer/go_microservice_example/tests/integration/notification_flow
```

### Step 2: サービス起動確認

```bash
./check_services.sh
```

**期待される出力**:
```
✅ Notification Service is accessible
✅ Review Service is accessible
✅ Shipping Service is accessible
✅ All services are ready for testing
```

### Step 3: テスト実行

```bash
./run_all_tests.sh
```

**期待される出力**:
```
Running Notification Tests...
PASS: TestOrderConfirmationNotification
PASS: TestPaymentSuccessNotification
PASS: TestShippingUpdateNotification
PASS: TestNotificationFailureHandling
PASS: TestMultipleNotificationsForRecipient

Running Review Tests...
PASS: TestCreateReview
PASS: TestReviewRatingValidation
PASS: TestGetReviewsByProduct
PASS: TestAverageRating
PASS: TestUpdateReview
PASS: TestReviewEditableUntil
PASS: TestDuplicateReviewPrevention

Running Shipping Tests...
PASS: TestCreateShipment
PASS: TestShipmentStatusFlow
PASS: TestTrackingEvents
PASS: TestShippingWithNotification
PASS: TestMultipleTrackingEventsChronology
PASS: TestShipmentUniqueOrderID

✅ All Tests Completed
Coverage: 85.2% of statements
```

---

## 📊 個別テスト実行

### Notification テストのみ

```bash
go test -v -run TestNotification
```

### Review テストのみ

```bash
go test -v -run TestReview
```

### Shipping テストのみ

```bash
go test -v -run TestShipping
```

### 特定のテスト1つ

```bash
go test -v -run TestCreateReview
```

---

## 🔧 トラブルシューティング

### エラー: データベース接続失敗

```bash
# サービスが起動しているか確認
docker ps | grep postgres

# 該当ポートが開いているか確認
netstat -an | grep 22017  # Notification
netstat -an | grep 22018  # Review
netstat -an | grep 22016  # Shipping
```

**解決策**:
```bash
# データベースを起動
cd /path/to/project/root
docker-compose up -d
```

### エラー: go.mod not found

```bash
# 依存関係を初期化
go mod tidy
```

### エラー: テスト失敗

```bash
# 詳細ログを表示
go test -v -run TestFailedTestName

# テストデータをクリーンアップ
psql -h localhost -p 22017 -U postgres -d notification_service -c "DELETE FROM notifications WHERE recipient LIKE '%test%';"
```

---

## 📈 カバレッジレポート表示

### HTML形式で表示

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

ブラウザで `coverage.out` がHTML形式で表示されます。

### テキスト形式で表示

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

---

## 🎯 テストの意味

### Notification Service

- **注文確認通知**: 注文が確定したことをユーザーに通知
- **決済完了通知**: 決済が成功したことを通知
- **配送開始通知**: 商品が発送されたことを通知
- **失敗ハンドリング**: 通知送信失敗時のエラー処理

### Review Service

- **レビュー作成**: 商品の評価とコメントを投稿
- **評価検証**: 1-5の範囲で評価を制限
- **平均評価計算**: 商品の総合評価を算出
- **重複防止**: 同じ注文・商品への重複レビューを防止

### Shipping Service

- **配送情報作成**: 配送先情報を登録
- **ステータス遷移**: preparing → shipped → delivered
- **追跡イベント**: 配送の進捗を時系列で記録
- **通知連携**: 配送開始時に通知サービスと連携

---

## 📚 詳細情報

詳細なドキュメントは以下を参照:

- [README.md](./README.md) - 完全なドキュメント
- [IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md) - 実装サマリー

---

## ✅ チェックリスト

テスト実行前の確認項目:

- [ ] PostgreSQL (port 22017) が起動している
- [ ] PostgreSQL (port 22018) が起動している
- [ ] PostgreSQL (port 22016) が起動している
- [ ] Go 1.21+ がインストールされている
- [ ] `go mod tidy` を実行済み

すべてチェックが完了したら:

```bash
./run_all_tests.sh
```

---

## 🎉 完了！

テストが全てPASSすれば、Notification/Review/Shipping の統合フローが正常に動作していることが確認できます。
