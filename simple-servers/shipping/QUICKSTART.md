# Shipping Service クイックスタート

## 起動コマンド

```bash
# 1. データベース起動
cd ../../infrastructure/docker
docker-compose up -d postgres_shipping

# 2. サービス起動
cd ../../simple-servers/shipping
./shipping
```

## 検証コマンド

```bash
# 包括的検証
./verify.sh

# データベース接続確認
cd ../../infrastructure/docker
docker-compose exec postgres_shipping psql -U postgres -d shipping_service -c "SELECT 1;"

# テーブル一覧
docker-compose exec postgres_shipping psql -U postgres -d shipping_service -c "\dt"

# 配送情報確認（サンプルデータがある場合）
docker-compose exec postgres_shipping psql -U postgres -d shipping_service -c "SELECT * FROM shipments LIMIT 5;"
```

## 環境変数

```bash
# デフォルト値（変更不要）
export SHIPPING_DATABASE_URL="postgresql://postgres:postgres_password@localhost:22016/shipping_service?sslmode=disable"
export SHIPPING_SERVICE_PORT="22108"
```

## トラブルシューティング

### データベース接続エラー

```bash
# コンテナ状態確認
docker-compose ps postgres_shipping

# 再起動
docker-compose restart postgres_shipping

# ログ確認
docker-compose logs postgres_shipping
```

### ポート競合

デフォルトポートが使用中の場合:

```bash
# 空きポート確認
lsof -i :22016  # DB
lsof -i :22108  # Service

# 環境変数で変更
export SHIPPING_DATABASE_URL="postgresql://postgres:postgres_password@localhost:XXXXX/shipping_service?sslmode=disable"
export SHIPPING_SERVICE_PORT="YYYYY"
```

## データベーススキーマ再作成

```bash
cd ../../infrastructure/docker

# データベース削除
docker-compose exec postgres_shipping psql -U postgres -c "DROP DATABASE shipping_service;"

# データベース作成
docker-compose exec postgres_shipping psql -U postgres -c "CREATE DATABASE shipping_service;"

# スキーマ適用
docker-compose exec -T postgres_shipping psql -U postgres -d shipping_service < ../../simple-servers/shipping/schema.sql
```

## サービス再ビルド

```bash
cd simple-servers/shipping

# クリーンビルド
rm -f shipping
go clean -cache
go build -o shipping

# 起動確認
./shipping
```

## 主要ポート

| コンポーネント | ポート | 用途 |
|--------------|--------|------|
| PostgreSQL | 22016 | データベース接続 |
| gRPC Service | 22108 | サービスAPI |

## 関連ドキュメント

- [README.md](./README.md) - 詳細ドキュメント
- [IMPLEMENTATION_SUMMARY.md](./IMPLEMENTATION_SUMMARY.md) - 実装サマリー
- [schema.sql](./schema.sql) - データベーススキーマ
