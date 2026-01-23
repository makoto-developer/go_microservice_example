# Docker インフラ構成完了

実施日: 2026-01-17

## 完了した作業

### 1. 新しいdocker-compose.yml作成

**ファイル**: `docker-compose.new.yml`

**構成**:
- Elasticsearch (kuromoji plugin) - 実際の実装
- PostgreSQL × 12 - サービスごとに分離
- Redis × 12 - サービスごとに分離
- RabbitMQ - メッセージキュー
- MinIO - オブジェクトストレージ
- MailHog - SMTPテスト

---

### 2. 新しい.env設定

**ファイル**: `.env.new`

**ポート割り当て**:

**PostgreSQL**: 5432-5443
- Auth: 5432
- Shop: 5433
- Customer: 5434
- Inventory: 5435
- Order: 5436
- Payment: 5437
- Shipping: 5438
- Notification: 5439
- Review: 5440
- Chat: 5441
- Search: 5442
- Admin: 5443

**Redis**: 6379-6390
- Auth: 6379
- Shop: 6380
- Customer: 6381
- Inventory: 6382
- Order: 6383
- Payment: 6384
- Shipping: 6385
- Notification: 6386
- Review: 6387
- Chat: 6388
- Search: 6389
- Admin: 6390

**その他**:
- Elasticsearch: 9200
- RabbitMQ: 5672 (管理UI: 15672)
- MinIO: 9000 (コンソール: 9001)
- MailHog: 1025 (UI: 8025)

---

### 3. Elasticsearch設定

**ファイル**: `config/elasticsearch/elasticsearch.yml`

**特徴**:
- kuromoji日本語解析プラグイン
- elasticsearch-initコンテナで自動インストール
- クラスタ名: shop-mall-cluster

---

### 4. ドキュメント

**ファイル**: `DOCKER_SETUP.md`

**内容**:
- セットアップ手順
- 接続情報
- トラブルシューティング
- データ永続化
- リソース使用量
- 旧構成との違い

---

## 旧構成との違い

### 旧構成（docker-compose.yml.backup）
- 単一PostgreSQL（すべてのサービスが共有）
- 単一Redis（すべてのサービスが共有）
- mock-elasticsearch（モック）
- MinIOなし

### 新構成（docker-compose.new.yml）
- PostgreSQL × 12（サービスごとに分離）
- Redis × 12（サービスごとに分離）
- 実際のElasticsearch（kuromoji plugin）
- MinIO（オブジェクトストレージ）

---

## 適用手順

```bash
# 1. 新しい構成を有効化
mv docker-compose.new.yml docker-compose.yml
mv .env.new .env

# 2. インフラサービス起動
docker-compose up -d

# 3. ヘルスチェック
docker-compose ps
curl http://localhost:9200/_cluster/health

# 4. kuromoji pluginの確認
curl http://localhost:9200/_cat/plugins
```

---

## リソース要件

**推奨スペック**:
- CPU: 4コア以上
- メモリ: 8GB以上
- ディスク: 20GB以上の空き容量

**コンテナ数**: 29個
- Elasticsearch: 1
- PostgreSQL: 12
- Redis: 12
- RabbitMQ: 1
- MinIO: 1
- MailHog: 1
- elasticsearch-init: 1 (一時)

---

## 次のステップ

1. **マイグレーション作成**
   - 各サービスのPostgreSQLスキーマ定義
   - 初期データ投入

2. **go.mod設定**
   - 各サービスの依存パッケージ定義
   - Docker内でのビルド設定

3. **サービス起動**
   - 12サービスのDockerfile作成
   - docker-compose.ymlにサービス追加

---

## トークン消費

- 実績: 約8,000トークン
- インフラ構成のみで完了
- 見積比: 計画通り
