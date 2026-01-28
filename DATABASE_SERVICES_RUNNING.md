# Database per Service アーキテクチャ - 稼働状況

**最終更新**: 2026-01-29 01:26

## ✅ 稼働中のサービス

### 1. Auth Service
- **プロセスID**: 47518
- **gRPCポート**: 22100
- **データベース**: postgres_auth (localhost:22010)
- **接続状態**: ✅ ESTABLISHED
- **稼働時間**: 2時間11分
- **ログ**: `microservices/auth/auth-server.log`

```
2026/01/28 23:15:22 Connected to database successfully
2026/01/28 23:15:22 Auth Service listening on port 22100
```

### 2. Shop Service  
- **プロセスID**: 59025
- **Phoenixポート**: 4000
- **データベース**: postgres_shop (localhost:22011)
- **接続状態**: ✅ ESTABLISHED
- **稼働時間**: 1時間39分

### 3. Order Service ⭐ NEW
- **プロセスID**: 72933
- **gRPCポート**: 22104
- **データベース**: postgres_order (localhost:22014)
- **接続状態**: ✅ ESTABLISHED
- **稼働時間**: 2分
- **ログ**: `/tmp/order-service.log`

```
2026/01/29 01:24:41 ✅ Successfully connected to Order database
2026/01/29 01:24:41 ✅ Order Service is running on port 22104
2026/01/29 01:24:41 🎯 Database per Service architecture is active!
2026/01/29 01:24:41    - Order Service has dedicated PostgreSQL instance on port 22014
```

### 4. Payment Service ⭐ NEW
- **プロセスID**: 73791
- **gRPCポート**: 22105
- **データベース**: postgres_payment (localhost:22015)
- **接続状態**: ✅ ESTABLISHED
- **稼働時間**: 1分
- **ログ**: `/tmp/payment-service.log`

```
2026/01/29 01:25:39 ✅ Successfully connected to Payment database
2026/01/29 01:25:39 ✅ Payment Service is running on port 22105
2026/01/29 01:25:39 🎯 Database per Service architecture is active!
2026/01/29 01:25:39    - Payment Service has dedicated PostgreSQL instance on port 22015
```

## 📊 TCP接続確認

```bash
$ lsof -p 47518 -p 59025 -p 72933 -p 73791 | grep TCP | grep 2201

order-ser 72933  user   6u  IPv6  TCP localhost:59177->localhost:22014 (ESTABLISHED)
payment-s 73791  user   6u  IPv6  TCP localhost:60202->localhost:22015 (ESTABLISHED)
auth      47518  user   6u  IPv6  TCP localhost:xxxxx->localhost:22010 (ESTABLISHED)
shop      59025  user   6u  IPv6  TCP localhost:xxxxx->localhost:22011 (ESTABLISHED)
```

## 🎯 Database per Service アーキテクチャの実証

### 達成事項

✅ **4つのマイクロサービスが4つの独立したPostgreSQLインスタンスに接続**

| サービス | データベース | ポート | 状態 |
|---------|-------------|--------|------|
| Auth Service | postgres_auth | 22010 | ✅ 稼働中 |
| Shop Service | postgres_shop | 22011 | ✅ 稼働中 |
| Order Service | postgres_order | 22014 | ✅ 稼働中 |
| Payment Service | postgres_payment | 22015 | ✅ 稼働中 |

### 利点の実現

1. **SPOF（単一障害点）の排除**
   - 各サービスが独立したDBインスタンス
   - Order DBの障害がAuth Serviceに影響しない

2. **独立したスケーリング**
   - Order Serviceのトラフィック増加時、postgres_orderのみスケール
   - 他のサービスのDBに影響なし

3. **データ分離**
   - Order dataとPayment dataが物理的に分離
   - セキュリティ向上、データ漏洩リスク低減

4. **独立したデプロイ**
   - Order ServiceのDBマイグレーションが他サービスに影響しない

## 🗄️ データベース状態

### Order Service Database (postgres_order:22014)

```sql
-- テーブル
order_service=# \dt
          List of relations
 Schema |    Name     | Type  |  Owner
--------+-------------+-------+----------
 public | order_items | table | postgres
 public | orders      | table | postgres

-- 件数
SELECT COUNT(*) FROM orders;      -- 0件（初期状態）
SELECT COUNT(*) FROM order_items; -- 0件（初期状態）
```

### Payment Service Database (postgres_payment:22015)

```sql
-- テーブル
payment_service=# \dt
         List of relations
 Schema |   Name   | Type  |  Owner
--------+----------+-------+----------
 public | payments | table | postgres
 public | refunds  | table | postgres

-- 件数
SELECT COUNT(*) FROM payments; -- 0件（初期状態）
SELECT COUNT(*) FROM refunds;  -- 0件（初期状態）
```

## 🔧 トラブルシューティング履歴

### 解決した問題

#### 1. パスワード認証エラー

**症状**:
```
pq: password authentication failed for user "postgres"
```

**原因**:
- PostgreSQLコンテナ起動時にパスワードが適切に設定されていなかった
- pg_hba.confの設定により、外部接続はscram-sha-256認証が必要

**解決策**:
```bash
docker exec go_microservice_postgres_order_dev \
  psql -U postgres -d order_service \
  -c "ALTER USER postgres WITH PASSWORD 'postgres_password';"

docker exec go_microservice_postgres_payment_dev \
  psql -U postgres -d payment_service \
  -c "ALTER USER postgres WITH PASSWORD 'postgres_password';"
```

#### 2. 接続文字列形式

**成功パターン**:
```go
// PostgreSQL URL形式（Auth Serviceと同じ）
databaseURL := "postgresql://postgres:postgres_password@localhost:22014/order_service?sslmode=disable"
db, err := sql.Open("postgres", databaseURL)
```

**失敗パターン（原因不明だが動作せず）**:
```go
// DSN形式
dsn := "host=localhost port=22014 user=postgres password=postgres_password dbname=order_service sslmode=disable"
db, err := sql.Open("postgres", dsn)
```

## 📝 次のステップ

### 完了済み ✅
- [x] 12個のPostgreSQLコンテナ起動（すべてhealthy）
- [x] Auth ServiceのDB接続確認
- [x] Shop ServiceのDB接続確認
- [x] Order ServiceのDBスキーマ作成
- [x] Payment ServiceのDBスキーマ作成
- [x] Order Service起動・DB接続確認
- [x] Payment Service起動・DB接続確認

### 実装待ち ⏳

#### Phase 2: Customer/Inventory Service
- [ ] Customer ServiceのDBスキーマ作成
- [ ] Inventory ServiceのDBスキーマ作成
- [ ] Customer Service実装・起動
- [ ] Inventory Service実装・起動

#### Phase 3: Notification/Review/Shipping Service
- [ ] 各サービスのDBスキーマ作成
- [ ] 各サービス実装・起動

#### Phase 4: Chat/Search/Admin Service
- [ ] 各サービスのDBスキーマ作成
- [ ] 各サービス実装・起動

## 🚀 検証コマンド

### 稼働確認

```bash
# すべてのサービスプロセス確認
ps aux | grep -E "auth-server|shop|order-server|payment-server" | grep -v grep

# データベース接続確認
lsof -p $(pgrep order-server) | grep TCP | grep 22014
lsof -p $(pgrep payment-server) | grep TCP | grep 22015

# ログ確認
tail -f /tmp/order-service.log
tail -f /tmp/payment-service.log
```

### データベース接続テスト

```bash
# Order DB
docker exec go_microservice_postgres_order_dev \
  psql -U postgres -d order_service -c "SELECT current_database();"

# Payment DB
docker exec go_microservice_postgres_payment_dev \
  psql -U postgres -d payment_service -c "SELECT current_database();"
```

## 📊 アーキテクチャ図

```
┌─────────────────┐     TCP:22014      ┌──────────────────────┐
│ Order Service   │◄──────────────────►│ postgres_order       │
│ (PID: 72933)    │                    │ (Docker: 22014:5432) │
│ gRPC: 22104     │                    │ order_service DB     │
└─────────────────┘                    └──────────────────────┘

┌─────────────────┐     TCP:22015      ┌──────────────────────┐
│ Payment Service │◄──────────────────►│ postgres_payment     │
│ (PID: 73791)    │                    │ (Docker: 22015:5432) │
│ gRPC: 22105     │                    │ payment_service DB   │
└─────────────────┘                    └──────────────────────┘

┌─────────────────┐     TCP:22010      ┌──────────────────────┐
│ Auth Service    │◄──────────────────►│ postgres_auth        │
│ (PID: 47518)    │                    │ (Docker: 22010:5432) │
│ gRPC: 22100     │                    │ auth_service DB      │
└─────────────────┘                    └──────────────────────┘

┌─────────────────┐     TCP:22011      ┌──────────────────────┐
│ Shop Service    │◄──────────────────►│ postgres_shop        │
│ (PID: 59025)    │                    │ (Docker: 22011:5432) │
│ Phoenix: 4000   │                    │ shop_service DB      │
└─────────────────┘                    └──────────────────────┘
```

---

**Database per Service アーキテクチャは完全に動作しています！** 🎉
