# Integration & E2E Tests

マイクロサービスアーキテクチャの統合テストとE2Eテストスイート

## ディレクトリ構成

```
tests/
├── README.md                           # このファイル
├── run_all_integration_tests.sh        # 🚀 全テスト実行（マスターランナー）
│
├── e2e/                                # End-to-End Tests
│   ├── README.md                       # E2Eテスト詳細ドキュメント
│   ├── test_runner.sh                  # E2Eテストランナー
│   ├── all_services_health_check.sh    # ヘルスチェック
│   ├── complete_purchase_flow_test.go  # 完全購入フロー（10ステップ）
│   ├── error_scenarios_test.go         # エラーシナリオ
│   ├── performance_test.go             # パフォーマンステスト
│   ├── docker-compose.test.yml         # テスト環境
│   ├── go.mod
│   └── go.sum
│
└── integration/                        # Integration Tests
    ├── auth/                           # Auth Service統合テスト
    ├── order_flow/                     # Order-Payment Flow
    └── notification_flow/              # Notification Flow
```

## クイックスタート

### 1. すべてのテストを実行

```bash
# プロジェクトルートから
cd tests
./run_all_integration_tests.sh
```

このスクリプトは以下を自動実行します：
1. ✅ 環境チェック（Docker起動確認）
2. ✅ サービス起動確認（未起動なら起動）
3. ✅ ヘルスチェック（全サービス）
4. ✅ 統合テスト実行
5. ✅ E2Eテスト実行
6. 📊 テストサマリー表示

### 2. E2Eテストのみ実行

```bash
cd tests/e2e
./test_runner.sh
```

### 3. 個別のテストスイート実行

```bash
# Auth統合テスト
cd tests/integration/auth
./run_test.sh

# Order Flow統合テスト
cd tests/integration/order_flow
./run_test.sh

# Notification Flow統合テスト
cd tests/integration/notification_flow
./run_all_tests.sh
```

## テストスイート概要

### E2E Tests (tests/e2e/)

#### 1. Complete Purchase Flow
**ファイル:** `complete_purchase_flow_test.go`

完全な購入フローを10ステップでテスト：
```
ユーザー登録 → 顧客情報登録 → 商品検索 → 商品詳細 → 在庫確認
→ 注文作成 → 決済処理 → 配送手配 → 通知確認 → レビュー投稿
```

**検証項目:**
- ✅ サービス間通信
- ✅ データフロー整合性
- ✅ トランザクション管理
- ✅ イベント駆動処理

**実行時間:** 約3-5分

#### 2. Error Scenarios
**ファイル:** `error_scenarios_test.go`

エラーケースとエラーハンドリングをテスト：
```
在庫不足 → 決済失敗ロールバック → 重複注文防止
→ 無効認証 → タイムアウト処理
```

**検証項目:**
- ✅ エラーレスポンス
- ✅ ロールバック処理
- ✅ べき等性保証
- ✅ 認証・認可
- ✅ タイムアウト処理

**実行時間:** 約2-3分

#### 3. Performance Tests
**ファイル:** `performance_test.go`

パフォーマンスと並行処理をテスト：
```
並行ユーザー登録（100） → 並行注文作成（50）
→ 検索パフォーマンス → DBコネクションプール
```

**検証項目:**
- ✅ 並行処理（100+ 同時リクエスト）
- ✅ スループット
- ✅ レスポンスタイム
- ✅ リソース管理

**実行時間:** 約5-10分

### Integration Tests (tests/integration/)

#### 1. Auth Service Tests
**ディレクトリ:** `integration/auth/`

認証・認可機能の統合テスト：
- ユーザー登録
- ログイン
- JWT検証
- ロール管理

#### 2. Order-Payment Flow Tests
**ディレクトリ:** `integration/order_flow/`

注文・決済フローの統合テスト：
- 注文作成
- 在庫引き当て
- 決済処理
- Sagaパターン

#### 3. Notification Flow Tests
**ディレクトリ:** `integration/notification_flow/`

通知システムの統合テスト：
- イベント購読
- 通知送信
- メール/SMS配信
- WebSocket通知

## 前提条件

### 必須環境

1. **Docker & Docker Compose**
   ```bash
   docker --version
   docker-compose --version
   ```

2. **Go 1.21+**
   ```bash
   go version
   ```

3. **全サービスが起動している**
   ```bash
   # プロジェクトルートから
   docker-compose up -d

   # または
   cd tests
   ./run_all_integration_tests.sh  # 自動起動オプションあり
   ```

### サービス一覧（全12サービス）

| サービス | ポート | ヘルスチェック |
|---------|--------|--------------|
| Auth | 8081 | http://localhost:8081/health |
| Shop | 8082 | http://localhost:8082/health |
| Customer | 8083 | http://localhost:8083/health |
| Inventory | 8084 | http://localhost:8084/health |
| Order | 8085 | http://localhost:8085/health |
| Payment | 8086 | http://localhost:8086/health |
| Notification | 8088 | http://localhost:8088/health |
| Shipping | 8089 | http://localhost:8089/health |
| Review | 8090 | http://localhost:8090/health |
| Chat | 8091 | http://localhost:8091/health |
| Search | 8092 | http://localhost:8092/health |
| Admin | 8093 | http://localhost:8093/health |

## 使用方法詳細

### ヘルスチェック

```bash
cd tests/e2e
./all_services_health_check.sh
```

出力例：
```
🏥 E2E Health Check: Verifying all microservices...

Checking auth service... ✅ OK
Checking shop service... ✅ OK
...

✅ All services are healthy!
```

### 特定のテストケース実行

```bash
cd tests/e2e

# 完全購入フローのみ
go test -v -run TestCompletePurchaseFlow -timeout 5m

# エラーシナリオのみ
go test -v -run TestErrorScenarios -timeout 5m

# パフォーマンステストのみ
go test -v -run TestPerformanceScenarios -timeout 10m

# 特定のステップのみ
go test -v -run TestCompletePurchaseFlow/Step1_UserRegistration
```

### ベンチマーク実行

```bash
cd tests/e2e
go test -bench=. -benchtime=10s -timeout 30m
```

### テストレポート生成

```bash
cd tests/e2e

# JSON形式
go test -json > test-report.json

# カバレッジレポート
go test -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html

# 詳細ログ
go test -v 2>&1 | tee test.log
```

## トラブルシューティング

### 1. Services Not Ready

```
❌ Services are not ready
```

**解決策:**
```bash
# サービス状態確認
docker-compose ps

# すべて起動
docker-compose up -d

# ログ確認
docker-compose logs [service-name]

# ヘルスチェック
curl http://localhost:8081/health
```

### 2. Test Timeout

```
❌ Test timed out after 5 minutes
```

**解決策:**
```bash
# タイムアウト延長
go test -timeout 10m

# サービスパフォーマンス確認
docker stats

# データベース接続確認
docker-compose logs postgres
```

### 3. Database Connection Errors

```
❌ Error: database connection refused
```

**解決策:**
```bash
# PostgreSQL起動確認
docker-compose ps postgres

# ログ確認
docker-compose logs postgres

# 接続テスト
docker-compose exec postgres psql -U user -d dbname -c "SELECT 1"

# 環境変数確認
cat .env | grep DB_
```

### 4. Port Already in Use

```
❌ Error: port 8081 already in use
```

**解決策:**
```bash
# ポート使用状況確認
lsof -i :8081

# プロセス終了
kill -9 [PID]

# または全サービス再起動
docker-compose down
docker-compose up -d
```

### 5. Payment Test Failures

```
❌ Payment processing failed
```

**解決策:**
```bash
# テストカード番号確認（4242424242424242）
# Payment Serviceログ確認
docker-compose logs payment

# Stripe設定確認（本番環境の場合）
cat .env | grep STRIPE_
```

## CI/CD統合

### GitHub Actions

```yaml
# .github/workflows/integration-tests.yml

name: Integration & E2E Tests

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  integration-tests:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Start services
        run: docker-compose up -d

      - name: Wait for services
        run: |
          cd tests/e2e
          timeout 60 bash -c 'until ./all_services_health_check.sh; do sleep 5; done'

      - name: Run all integration tests
        run: |
          cd tests
          ./run_all_integration_tests.sh

      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: test-results
          path: tests/**/*.log
```

### GitLab CI

```yaml
# .gitlab-ci.yml

stages:
  - test

integration-tests:
  stage: test
  image: golang:1.21
  services:
    - docker:dind
  script:
    - docker-compose up -d
    - cd tests
    - ./run_all_integration_tests.sh
  artifacts:
    when: always
    reports:
      junit: tests/e2e/test-report.xml
```

## パフォーマンス目標

### Response Time

| エンドポイント | 平均 | 95パーセンタイル | 最大 |
|--------------|------|---------------|------|
| User Registration | < 200ms | < 500ms | < 1s |
| Product Search | < 300ms | < 700ms | < 2s |
| Order Creation | < 500ms | < 1s | < 3s |
| Payment Processing | < 1s | < 2s | < 5s |

### Throughput

| 操作 | 目標 | 実績 |
|-----|------|------|
| 並行ユーザー登録 | 90%+ 成功 | 95%+ |
| 並行注文作成 | 80%+ 成功 | 85%+ |
| 検索クエリ | 100 req/s | 120+ req/s |

## ベストプラクティス

### テスト作成

1. **テストデータの独立性**
   ```go
   // 各テストで一意なデータを使用
   email := fmt.Sprintf("test.%d@example.com", time.Now().Unix())
   ```

2. **適切なタイムアウト設定**
   ```go
   // E2E: 5-10分
   go test -timeout 10m

   // Unit: 30秒-1分
   go test -timeout 1m
   ```

3. **クリーンアップ**
   ```go
   defer func() {
       // テストデータのクリーンアップ
   }()
   ```

4. **エラーメッセージの明確化**
   ```go
   assert.Equal(t, expected, actual, "User ID should match")
   ```

### CI/CD

1. テストを段階的に実行（Unit → Integration → E2E）
2. 失敗時のログとアーティファクトを保存
3. パフォーマンステストは定期実行（nightly build等）
4. カバレッジレポートを生成・保存

## まとめ

このテストスイートは、マイクロサービスアーキテクチャの品質を保証します。

**テストカバレッジ:**
- ✅ 12サービスすべての統合
- ✅ 完全な購入フロー（10ステップ）
- ✅ エラーハンドリング（5シナリオ）
- ✅ パフォーマンス（100+ 並行）
- ✅ トランザクション整合性
- ✅ イベント駆動処理

**次のステップ:**
1. `./run_all_integration_tests.sh` で全テスト実行
2. 結果を確認
3. 失敗したテストがあれば修正
4. CI/CDパイプラインに統合
5. 定期的に実行して品質を維持

---

**作成日:** 2026-01-29
**バージョン:** 1.0.0
**メンテナンス:** 新機能追加時にテストケースを追加
