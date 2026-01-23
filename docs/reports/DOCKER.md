# Docker Compose 構成

このプロジェクトは Docker Compose を使用して、すべてのマイクロサービスと PostgreSQL データベースを一度に起動できます。

## 構成

### サービス一覧

| サービス | ポート | データベース | DB ポート |
|---------|--------|------------|----------|
| Auth Service | 50051 | postgres-auth | 5432 |
| Shop Service | 50052 | postgres-shop | 5433 |
| Customer Service | 50053 | postgres-customer | 5434 |
| Inventory Service | 50054 | postgres-inventory | 5435 |

### データベース

各サービスは独立した PostgreSQL 16 データベースを使用します：

- `postgres-auth`: auth_db（ポート 5432）
- `postgres-shop`: shop_db（ポート 5433）
- `postgres-customer`: customer_db（ポート 5434）
- `postgres-inventory`: inventory_db（ポート 5435）

## 使用方法

### すべてのサービスを起動

```bash
docker-compose up -d
```

### ログを確認

```bash
# すべてのサービス
docker-compose logs -f

# 特定のサービス
docker-compose logs -f auth-service
docker-compose logs -f shop-service
docker-compose logs -f customer-service
docker-compose logs -f inventory-service
```

### サービスの状態を確認

```bash
docker-compose ps
```

### 特定のサービスのみ起動

```bash
# Auth Service とそのデータベースのみ
docker-compose up -d postgres-auth auth-service

# Customer Service とそのデータベースのみ
docker-compose up -d postgres-customer customer-service
```

### サービスを停止

```bash
# すべて停止
docker-compose down

# データも削除する場合
docker-compose down -v
```

### サービスを再ビルド

```bash
# すべて再ビルド
docker-compose build

# 特定のサービスのみ
docker-compose build auth-service
docker-compose build customer-service
```

### サービスを再起動

```bash
# すべて再起動
docker-compose restart

# 特定のサービスのみ
docker-compose restart auth-service
```

## データベースマイグレーション

各サービスのデータベースマイグレーションは、サービス起動後に手動で実行する必要があります。

### Customer Service のマイグレーション例

```bash
# マイグレーションツールをインストール
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# マイグレーション実行
cd generated/customer
migrate -path migrations -database "postgresql://postgres:postgres@localhost:5434/customer_db?sslmode=disable" up
```

### 他のサービスも同様に

```bash
# Auth Service (ポート 5432)
migrate -path generated/auth/migrations -database "postgresql://postgres:postgres@localhost:5432/auth_db?sslmode=disable" up

# Shop Service (ポート 5433)
migrate -path generated/shop/migrations -database "postgresql://postgres:postgres@localhost:5433/shop_db?sslmode=disable" up

# Inventory Service (ポート 5435)
migrate -path generated/inventory/migrations -database "postgresql://postgres:postgres@localhost:5435/inventory_db?sslmode=disable" up
```

## トラブルシューティング

### ポート競合

すでに使用中のポートがある場合は、`docker-compose.yml` のポート設定を変更してください。

### データベース接続エラー

データベースのヘルスチェックが完了するまで待ってから、サービスが起動します。
エラーが発生する場合は、以下を確認：

```bash
# データベースのヘルスチェック状態を確認
docker-compose ps

# データベースのログを確認
docker-compose logs postgres-auth
docker-compose logs postgres-customer
```

### イメージの再ビルド

コードを変更した場合は、イメージを再ビルドしてください：

```bash
docker-compose build auth-service
docker-compose up -d auth-service
```

## 開発時の便利なコマンド

### すべてをクリーンに再起動

```bash
docker-compose down -v
docker-compose build
docker-compose up -d
```

### 特定のサービスのシェルに入る

```bash
docker-compose exec auth-service sh
docker-compose exec customer-service sh
```

### データベースに直接接続

```bash
# Auth DB
docker-compose exec postgres-auth psql -U postgres -d auth_db

# Customer DB
docker-compose exec postgres-customer psql -U postgres -d customer_db
```

## ネットワーク

すべてのサービスは `microservices` ネットワーク上で通信します。

サービス間の通信には、コンテナ名を使用します：
- `auth-service:50051`
- `shop-service:50052`
- `customer-service:50053`
- `inventory-service:50054`

## ボリューム

データベースのデータは永続化されます：
- `postgres-auth-data`
- `postgres-shop-data`
- `postgres-customer-data`
- `postgres-inventory-data`

データを完全に削除する場合は：

```bash
docker-compose down -v
```
