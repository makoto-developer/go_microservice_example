# Docker インフラ構成

このドキュメントは、新しいインフラ構成（PostgreSQL×12、Redis×12、Elasticsearch、MinIO）の使用方法を説明します。

---

## 新しい構成の概要

### インフラサービス
- **Elasticsearch** (kuromoji plugin) - 実際のElasticsearch（ポート: 9200）
- **RabbitMQ** - メッセージキュー（ポート: 5672, 管理UI: 15672）
- **MinIO** - オブジェクトストレージ（ポート: 9000, コンソール: 9001）
- **MailHog** - SMTPテスト（SMTP: 1025, UI: 8025）

### PostgreSQL（サービスごと）
| サービス | ポート | データベース名 |
|---------|-------|--------------|
| Auth | 5432 | auth_service |
| Shop | 5433 | shop_service |
| Customer | 5434 | customer_service |
| Inventory | 5435 | inventory_service |
| Order | 5436 | order_service |
| Payment | 5437 | payment_service |
| Shipping | 5438 | shipping_service |
| Notification | 5439 | notification_service |
| Review | 5440 | review_service |
| Chat | 5441 | chat_service |
| Search | 5442 | search_service |
| Admin | 5443 | admin_service |

### Redis（サービスごと）
| サービス | ポート |
|---------|-------|
| Auth | 6379 |
| Shop | 6380 |
| Customer | 6381 |
| Inventory | 6382 |
| Order | 6383 |
| Payment | 6384 |
| Shipping | 6385 |
| Notification | 6386 |
| Review | 6387 |
| Chat | 6388 |
| Search | 6389 |
| Admin | 6390 |

---

## セットアップ手順

### 1. 環境ファイルの作成

```bash
# 新しい.envファイルをコピー
cp .env.new .env

# 必要に応じて編集
vim .env
```

### 2. 新しいdocker-compose.ymlの適用

```bash
# 既存の構成をバックアップ（すでに完了）
# cp docker-compose.yml docker-compose.yml.backup

# 新しい構成を有効化
mv docker-compose.new.yml docker-compose.yml
mv .env.new .env
```

### 3. インフラサービスの起動

```bash
# すべて起動
docker-compose up -d

# ログ確認
docker-compose logs -f

# 特定のサービスのみ起動
docker-compose up -d elasticsearch rabbitmq minio

# PostgreSQL/Redisのみ起動
docker-compose up -d postgres_auth postgres_shop redis_auth redis_shop
```

---

## 接続情報

### Elasticsearch

```bash
# ヘルスチェック
curl http://localhost:9200/_cluster/health

# kuromoji pluginの確認
curl http://localhost:9200/_cat/plugins

# インデックス作成（kuromoji analyzer使用）
curl -X PUT "localhost:9200/products" -H 'Content-Type: application/json' -d'
{
  "settings": {
    "analysis": {
      "analyzer": {
        "kuromoji_analyzer": {
          "type": "custom",
          "tokenizer": "kuromoji_tokenizer"
        }
      }
    }
  }
}
'
```

### PostgreSQL

```bash
# Auth Serviceデータベースに接続
psql -h localhost -p 5432 -U postgres -d auth_service

# Shop Serviceデータベースに接続
psql -h localhost -p 5433 -U postgres -d shop_service
```

### Redis

```bash
# Auth Service Redisに接続
redis-cli -p 6379 -a redis_password

# Shop Service Redisに接続
redis-cli -p 6380 -a redis_password
```

### MinIO

```bash
# MinIOコンソール
open http://localhost:9001

# 認証情報
# ユーザー: minioadmin
# パスワード: minioadmin
```

### RabbitMQ

```bash
# 管理UI
open http://localhost:15672

# 認証情報
# ユーザー: admin
# パスワード: change_me_in_production
```

### MailHog

```bash
# メールUI
open http://localhost:8025
```

---

## トラブルシューティング

### Elasticsearchが起動しない

```bash
# vm.max_map_countを増やす
sudo sysctl -w vm.max_map_count=262144

# macOSの場合（Docker Desktop設定で）
# Resources > Advanced > Memory を 4GB 以上に設定
```

### kuromoji pluginがインストールされない

```bash
# elasticsearch-initコンテナのログを確認
docker-compose logs elasticsearch-init

# 手動でインストール
docker-compose exec elasticsearch elasticsearch-plugin install analysis-kuromoji
docker-compose restart elasticsearch
```

### PostgreSQL接続エラー

```bash
# コンテナの状態確認
docker-compose ps | grep postgres

# ログ確認
docker-compose logs postgres_auth

# 再起動
docker-compose restart postgres_auth
```

### Redisパスワードエラー

```bash
# パスワード付きで接続
redis-cli -p 6379 -a redis_password

# .envファイルのREDIS_PASSWORDを確認
grep REDIS_PASSWORD .env
```

---

## データ永続化

データは以下のボリュームに保存されます：

```bash
# ボリューム一覧
docker volume ls | grep go_microservice

# 特定のボリュームを削除（注意: データが消えます）
docker volume rm go_microservice_postgres_auth_data_dev

# すべてのボリュームを削除（注意: すべてのデータが消えます）
docker-compose down -v
```

---

## リソース使用量

### 推奨スペック

- **CPU**: 4コア以上
- **メモリ**: 8GB以上
- **ディスク**: 20GB以上の空き容量

### 現在の使用量確認

```bash
# コンテナのリソース使用量
docker stats

# ディスク使用量
docker system df
```

---

## 旧構成との違い

### 旧構成（docker-compose.yml.backup）
- 単一PostgreSQL（すべてのサービスが共有）
- 単一Redis（すべてのサービスが共有）
- mock-elasticsearch（モック）
- MinIOなし

### 新構成（docker-compose.yml）
- PostgreSQL × 12（サービスごとに分離）
- Redis × 12（サービスごとに分離）
- 実際のElasticsearch（kuromoji plugin）
- MinIO（オブジェクトストレージ）

### 移行方法

旧構成に戻す場合：

```bash
# 新構成を停止
docker-compose down

# 旧構成を復元
mv docker-compose.yml docker-compose.new.yml
mv docker-compose.yml.backup docker-compose.yml

# 起動
docker-compose up -d
```

---

## まとめ

新しいインフラ構成により：

✅ **サービス分離**: 各サービスが独立したPostgreSQL・Redisを持つ
✅ **実際のElasticsearch**: kuromoji日本語解析を使用可能
✅ **MinIO**: 画像・ファイルストレージを提供
✅ **スケーラビリティ**: サービスごとに個別にスケール可能

この構成で本格的なマイクロサービス開発が可能になりました。
