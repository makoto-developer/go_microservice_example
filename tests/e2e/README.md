# E2E Tests - Complete Purchase Flow

完全な購入フローのEnd-to-End統合テストスイート

## 概要

このディレクトリには、マイクロサービスアーキテクチャ全体を検証するE2Eテストが含まれています。

### テスト範囲

1. **Complete Purchase Flow** (`complete_purchase_flow_test.go`)
   - ユーザー登録から商品購入、レビュー投稿までの完全フロー
   - 10ステップの統合シナリオ

2. **Error Scenarios** (`error_scenarios_test.go`)
   - 在庫不足エラー
   - 決済失敗時のロールバック
   - 重複注文防止
   - 認証エラー
   - タイムアウト処理

3. **Performance Tests** (`performance_test.go`)
   - 並行ユーザー登録（100件）
   - 並行注文作成（50件）
   - 検索パフォーマンス（100並行）
   - データベースコネクションプール

## 前提条件

### 必須サービス

すべてのマイクロサービスが起動している必要があります：

```bash
# プロジェクトルートから
docker-compose up -d
```

対象サービス：
- Auth Service (8081)
- Shop Service (8082)
- Customer Service (8083)
- Inventory Service (8084)
- Order Service (8085)
- Payment Service (8086)
- Shipping Service (8089)
- Notification Service (8088)
- Review Service (8090)
- Search Service (8092)
- Chat Service (8091)
- Admin Service (8093)

## 使用方法

### 1. ヘルスチェック

すべてのサービスが起動しているか確認：

```bash
./all_services_health_check.sh
```

出力例：
```
🏥 E2E Health Check: Verifying all microservices...

Checking auth service... ✅ OK
Checking shop service... ✅ OK
Checking customer service... ✅ OK
...

✅ All services are healthy!
```

### 2. 全テスト実行

```bash
./test_runner.sh
```

このスクリプトは以下を実行します：
1. ヘルスチェック
2. 完全購入フローテスト
3. エラーシナリオテスト
4. パフォーマンステスト（オプション）

### 3. 個別テスト実行

特定のテストのみ実行する場合：

```bash
# 完全購入フロー
go test -v -run TestCompletePurchaseFlow -timeout 5m

# エラーシナリオ
go test -v -run TestErrorScenarios -timeout 5m

# パフォーマンステスト
go test -v -run TestPerformanceScenarios -timeout 10m

# 特定のサブテスト
go test -v -run TestCompletePurchaseFlow/Step1_UserRegistration
```

### 4. ベンチマーク実行

```bash
go test -bench=. -benchtime=10s -timeout 30m
```

## テストシナリオ詳細

### Complete Purchase Flow (10ステップ)

```
1. ユーザー登録 (Auth Service)
   └─ POST /api/v1/register
   └─ user_id, JWT token 取得

2. 顧客情報登録 (Customer Service)
   └─ POST /api/v1/customers
   └─ 配送先住所登録

3. 商品検索 (Search Service)
   └─ GET /api/v1/search?q=商品名
   └─ 商品一覧取得

4. 商品詳細取得 (Shop Service)
   └─ GET /api/v1/products/:id
   └─ 価格・在庫確認

5. 在庫確認 (Inventory Service)
   └─ GET /api/v1/inventory/:product_id
   └─ 在庫数確認

6. 注文作成 (Order Service)
   └─ POST /api/v1/orders
   └─ order_id 取得
   └─ 在庫引き当て

7. 決済処理 (Payment Service)
   └─ POST /api/v1/payments
   └─ クレジットカード決済
   └─ 決済成功確認

8. 配送手配 (Shipping Service)
   └─ GET /api/v1/shipments?order_id=XXX
   └─ tracking_number 取得

9. 配送通知確認 (Notification Service)
   └─ GET /api/v1/notifications?user_id=XXX
   └─ 自動通知送信確認

10. レビュー投稿 (Review Service)
    └─ POST /api/v1/reviews
    └─ 評価・コメント投稿
```

### Error Scenarios

#### 1. 在庫不足エラー
```go
注文数量: 999999（在庫超過）
期待結果: HTTP 400 Bad Request
エラーメッセージ: "insufficient stock"
```

#### 2. 決済失敗とロールバック
```go
カード番号: 0000000000000000（無効）
期待結果:
  - 決済失敗 (HTTP 400)
  - 在庫解放（ロールバック）
  - 注文ステータス: CANCELLED
```

#### 3. 重複注文防止
```go
同じidempotency_keyで2回注文
期待結果:
  - 1回目: HTTP 201 Created
  - 2回目: HTTP 200 OK（同じorder_id返却）
```

#### 4. 無効な認証
```go
JWT: "invalid-token"
期待結果: HTTP 401 Unauthorized
```

#### 5. タイムアウト処理
```go
タイムアウト: 100ms
期待結果: タイムアウトエラーまたは正常レスポンス
```

### Performance Tests

#### 並行ユーザー登録
```
並行数: 100ユーザー
目標: 90%以上成功
タイムアウト: 30秒
```

#### 並行注文作成
```
並行数: 50注文
目標: 80%以上成功
タイムアウト: 60秒
```

#### 検索パフォーマンス
```
並行検索: 100リクエスト
平均レスポンス: < 500ms
最大レスポンス: < 2秒
```

#### データベースコネクションプール
```
並行操作: 200リクエスト
成功率: 95%以上
タイムアウト: 30秒
```

## テスト環境

### Docker Compose (オプション)

テスト専用環境を起動する場合：

```bash
docker-compose -f docker-compose.test.yml up -d
```

テスト専用のポート：
- PostgreSQL: 5433
- Redis: 6380
- RabbitMQ: 5673, 15673

### 環境変数

```bash
# デフォルト設定（変更可能）
export AUTH_SERVICE_URL=http://localhost:8081
export CUSTOMER_SERVICE_URL=http://localhost:8083
export SHOP_SERVICE_URL=http://localhost:8082
# ... etc
```

## トラブルシューティング

### Services Not Ready

```bash
❌ Services are not ready

解決策:
1. すべてのサービスが起動しているか確認
   docker-compose ps

2. ヘルスチェックエンドポイント確認
   curl http://localhost:8081/health

3. ログ確認
   docker-compose logs [service-name]
```

### Test Timeout

```bash
❌ Test timed out after 5 minutes

解決策:
1. タイムアウトを延長
   go test -timeout 10m

2. サービスのパフォーマンス確認
3. データベース接続確認
```

### Database Connection Errors

```bash
❌ Error: database connection refused

解決策:
1. PostgreSQLが起動しているか確認
   docker-compose ps postgres

2. 接続情報確認（.env）
3. ネットワーク確認
   docker network ls
```

### Payment Test Failures

```bash
❌ Payment processing failed

解決策:
1. テストカード番号確認（4242424242424242）
2. Payment Serviceログ確認
3. Stripe設定確認（本番環境の場合）
```

## CI/CD統合

### GitHub Actions

```yaml
name: E2E Tests

on: [push, pull_request]

jobs:
  e2e-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Start services
        run: docker-compose up -d

      - name: Wait for services
        run: |
          cd tests/e2e
          ./all_services_health_check.sh

      - name: Run E2E tests
        run: |
          cd tests/e2e
          go test -v -timeout 10m
```

## レポート

### テストレポート生成

```bash
# JSON形式
go test -json > test-report.json

# カバレッジレポート
go test -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### パフォーマンスレポート

```bash
go test -bench=. -benchmem > bench-report.txt
```

## 開発ガイド

### 新しいテストケース追加

1. 既存のテストファイルに追加、または新しいファイル作成
2. `E2ETestContext` 構造体を活用
3. ヘルパー関数を再利用
4. `README.md` を更新

### テストデータ管理

```go
// テストユーザーの作成
testCtx := setupTestUser(t)

// テスト商品の取得
productID := getTestProductID(t, testCtx.Token)

// テスト注文の作成
orderID := createTestOrder(t, testCtx, productID)
```

## まとめ

このE2Eテストスイートは、マイクロサービスアーキテクチャ全体の動作を検証します。

**重要な検証項目：**
- ✅ サービス間通信
- ✅ データ整合性
- ✅ エラーハンドリング
- ✅ トランザクション管理
- ✅ パフォーマンス
- ✅ 並行処理

**次のステップ：**
1. すべてのサービスを起動
2. `./all_services_health_check.sh` 実行
3. `./test_runner.sh` で全テスト実行
4. レポート確認
5. 問題があれば修正

---

**作成日:** 2026-01-29
**バージョン:** 1.0.0
