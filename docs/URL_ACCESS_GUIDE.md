# URL アクセスガイド

このドキュメントは、役割別（顧客、オーナー、管理者）のアクセス可能なURLを整理したものです。

---

## 🛍️ 顧客（Customer）向けURL

### フロントエンド（予定）

| サービス | URL | ポート | 説明 | 実装状況 |
|---------|-----|--------|------|---------|
| メインWebアプリ | http://localhost:22200 | 22200 | 商品閲覧・購入・レビュー | 🔜 未実装 |

### マイクロサービス（予定）

| サービス | URL | ポート | 説明 | 実装状況 |
|---------|-----|--------|------|---------|
| Auth Service | grpc://localhost:22100 | 22100 | 顧客登録・ログイン | 🔜 未実装 |
| Shop Service | grpc://localhost:22101 | 22101 | ショップ一覧・商品閲覧 | 🔜 未実装 |
| Customer Service | grpc://localhost:22102 | 22102 | 顧客情報管理・住所管理 | 🔜 未実装 |
| Order Service | grpc://localhost:22104 | 22104 | 注文作成・注文履歴 | 🔜 未実装 |
| Payment Service | grpc://localhost:22105 | 22105 | 決済処理 | 🔜 未実装 |
| Review Service | grpc://localhost:22108 | 22108 | レビュー投稿・閲覧 | 🔜 未実装 |
| Chat Service (gRPC) | grpc://localhost:22109 | 22109 | ショップとのチャット | 🔜 未実装 |
| Chat Service (WebSocket) | ws://localhost:22110 | 22110 | リアルタイムチャット | 🔜 未実装 |
| Search Service | grpc://localhost:22111 | 22111 | 商品検索 | 🔜 未実装 |

### 顧客が利用できる機能

- ✅ アカウント登録・ログイン（Customer role）
- ✅ ショップ・商品の閲覧
- ✅ 商品の検索
- ✅ カートへの追加・注文作成
- ✅ 決済処理
- ✅ 注文履歴の確認
- ✅ 配送状況の確認
- ✅ レビューの投稿・閲覧
- ✅ ショップとのチャット
- ✅ 通知の受信

---

## 🏪 オーナー（Shop Owner）向けURL

### フロントエンド（予定）

| サービス | URL | ポート | 説明 | 実装状況 |
|---------|-----|--------|------|---------|
| オーナーダッシュボード | http://localhost:22200/owner | 22200 | ショップ管理画面 | 🔜 未実装 |

### マイクロサービス（予定）

| サービス | URL | ポート | 説明 | 実装状況 |
|---------|-----|--------|------|---------|
| Auth Service | grpc://localhost:22100 | 22100 | オーナー登録・ログイン | 🔜 未実装 |
| Shop Service | grpc://localhost:22101 | 22101 | ショップ作成・商品管理 | 🔜 未実装 |
| Inventory Service | grpc://localhost:22103 | 22103 | 在庫管理 | 🔜 未実装 |
| Order Service | grpc://localhost:22104 | 22104 | 注文管理・発送処理 | 🔜 未実装 |
| Payment Service | grpc://localhost:22105 | 22105 | 売上確認 | 🔜 未実装 |
| Shipping Service | grpc://localhost:22106 | 22106 | 配送管理 | 🔜 未実装 |
| Notification Service | grpc://localhost:22107 | 22107 | 通知管理 | 🔜 未実装 |
| Review Service | grpc://localhost:22108 | 22108 | レビュー管理・返信 | 🔜 未実装 |
| Chat Service (gRPC) | grpc://localhost:22109 | 22109 | 顧客とのチャット | 🔜 未実装 |
| Chat Service (WebSocket) | ws://localhost:22110 | 22110 | リアルタイムチャット | 🔜 未実装 |

### オーナーが利用できる機能

- ✅ アカウント登録・ログイン（Owner role）
- ✅ ショップの作成・編集
- ✅ 商品の登録・編集・削除
- ✅ 在庫管理（追加・更新）
- ✅ 注文の確認・処理
- ✅ 配送状況の更新
- ✅ 売上の確認
- ✅ レビューの閲覧・返信
- ✅ 顧客とのチャット
- ✅ 通知の送信・受信

---

## 👨‍💼 管理者（Admin）向けURL

### フロントエンド（予定）

| サービス | URL | ポート | 説明 | 実装状況 |
|---------|-----|--------|------|---------|
| 管理者ダッシュボード | http://localhost:22201 | 22201 | システム全体の管理画面 | 🔜 未実装 |

### マイクロサービス（予定）

| サービス | URL | ポート | 説明 | 実装状況 |
|---------|-----|--------|------|---------|
| Auth Service | grpc://localhost:22100 | 22100 | 管理者ログイン | 🔜 未実装 |
| Admin Service | grpc://localhost:22112 | 22112 | システム全体の管理 | 🔜 未実装 |
| **全サービスへのアクセス** | grpc://localhost:22100-22111 | 22100-22111 | 全マイクロサービス | 🔜 未実装 |

### 管理者が利用できる機能

- ✅ アカウント登録・ログイン（Admin role）
- ✅ ユーザー管理（顧客・オーナーの管理）
- ✅ ショップの承認・停止
- ✅ 商品の審査・削除
- ✅ 注文の管理・キャンセル
- ✅ 決済の管理・返金処理
- ✅ レビューの削除
- ✅ チャット内容の監視
- ✅ システムログの確認
- ✅ 統計・レポートの閲覧
- ✅ 全サービスへのアクセス権限

---

## 🔧 開発・管理ツールURL（全員）

### 現在稼働中のインフラサービス

| サービス | URL | ポート | 説明 | 実装状況 |
|---------|-----|--------|------|---------|
| **Elasticsearch** | http://localhost:22000 | 22000 | 検索エンジンAPI | ✅ 稼働中 |
| **RabbitMQ Management** | http://localhost:22003 | 22003 | メッセージキュー管理画面 | ✅ 稼働中 |
| **MinIO Console** | http://localhost:22005 | 22005 | オブジェクトストレージ管理画面 | ✅ 稼働中 |
| **MailHog UI** | http://localhost:22007 | 22007 | メールテスト確認画面 | ✅ 稼働中 |

#### RabbitMQ Management
```
URL: http://localhost:22003
ユーザー名: admin
パスワード: rabbitmq_password
```

#### MinIO Console
```
URL: http://localhost:22005
ユーザー名: admin
パスワード: minio_password
```

#### MailHog UI
```
URL: http://localhost:22007
（認証不要）
```

#### Elasticsearch
```
URL: http://localhost:22000
（認証無効化済み）

# Health Check
curl http://localhost:22000/_cluster/health

# Index一覧
curl http://localhost:22000/_cat/indices?v
```

---

## 📊 データベース接続情報

### PostgreSQL（役割別）

#### 顧客関連データベース

| サービス | ポート | データベース名 | 用途 |
|---------|--------|--------------|------|
| Auth Service | 22010 | auth_service | 顧客アカウント情報 |
| Customer Service | 22012 | customer_service | 顧客プロフィール・住所 |
| Order Service | 22014 | order_service | 注文履歴 |
| Payment Service | 22015 | payment_service | 決済履歴 |
| Review Service | 22018 | review_service | レビュー |
| Chat Service | 22019 | chat_service | チャット履歴 |

接続例:
```bash
psql -h localhost -p 22012 -U postgres -d customer_service
```

#### オーナー関連データベース

| サービス | ポート | データベース名 | 用途 |
|---------|--------|--------------|------|
| Auth Service | 22010 | auth_service | オーナーアカウント情報 |
| Shop Service | 22011 | shop_service | ショップ・商品情報 |
| Inventory Service | 22013 | inventory_service | 在庫管理 |
| Shipping Service | 22016 | shipping_service | 配送管理 |
| Notification Service | 22017 | notification_service | 通知管理 |

接続例:
```bash
psql -h localhost -p 22011 -U postgres -d shop_service
```

#### 管理者関連データベース

| サービス | ポート | データベース名 | 用途 |
|---------|--------|--------------|------|
| Admin Service | 22021 | admin_service | システム管理 |
| Search Service | 22020 | search_service | 検索インデックス管理 |

接続例:
```bash
psql -h localhost -p 22021 -U postgres -d admin_service
```

### Redis（役割別）

#### 顧客関連キャッシュ

| サービス | ポート | 用途 |
|---------|--------|------|
| Auth Service | 22030 | セッション・JWT blacklist |
| Customer Service | 22032 | 顧客データキャッシュ |
| Order Service | 22034 | 注文キャッシュ |
| Payment Service | 22035 | 決済キャッシュ |
| Review Service | 22038 | レビューキャッシュ |
| Chat Service | 22039 | チャットメッセージキャッシュ |

接続例:
```bash
redis-cli -h localhost -p 22032 -a redis_password
```

#### オーナー関連キャッシュ

| サービス | ポート | 用途 |
|---------|--------|------|
| Shop Service | 22031 | ショップデータキャッシュ |
| Inventory Service | 22033 | 在庫キャッシュ |
| Shipping Service | 22036 | 配送キャッシュ |
| Notification Service | 22037 | 通知キュー |

接続例:
```bash
redis-cli -h localhost -p 22031 -a redis_password
```

#### 管理者関連キャッシュ

| サービス | ポート | 用途 |
|---------|--------|------|
| Search Service | 22040 | 検索キャッシュ |
| Admin Service | 22041 | 管理データキャッシュ |

接続例:
```bash
redis-cli -h localhost -p 22041 -a redis_password
```

---

## 🔐 認証・認可の仕組み

### 役割（Role）の種類

```go
enum Role {
  CUSTOMER,    // 顧客
  SHOP_OWNER,  // オーナー
  ADMIN        // 管理者
}
```

### アクセス制御

#### Customer（顧客）
- ✅ 自分の情報のみ閲覧・編集可能
- ✅ 全ショップ・商品の閲覧可能
- ✅ 自分の注文のみ閲覧・管理可能
- ✅ レビューの投稿・編集・削除（自分のもののみ）
- ❌ 他ユーザーの情報は閲覧不可
- ❌ ショップ管理機能は使用不可

#### Shop Owner（オーナー）
- ✅ 自分のショップ情報の閲覧・編集可能
- ✅ 自分のショップの商品管理可能
- ✅ 自分のショップへの注文管理可能
- ✅ 自分のショップのレビュー閲覧・返信可能
- ❌ 他のショップの情報は閲覧不可
- ❌ 顧客の個人情報（住所等）は閲覧不可
- ❌ システム管理機能は使用不可

#### Admin（管理者）
- ✅ 全ユーザー・ショップ・商品の管理可能
- ✅ 全注文・決済の管理可能
- ✅ システムログの閲覧可能
- ✅ 統計・レポートの閲覧可能
- ✅ 全サービスへのアクセス可能

---

## 📝 実装状況

### Phase 1: 基盤サービス ✅ 完了
- ✅ Auth Service (DSL定義完了)
- ✅ Shop Service (DSL定義完了)
- ✅ インフラストラクチャ（Docker Compose）

### Phase 2: コアサービス 🔜 未実装
- 🔜 Customer Service
- 🔜 Inventory Service
- 🔜 Order Service
- 🔜 Payment Service

### Phase 3: 拡張サービス1 🔜 未実装
- 🔜 Shipping Service
- 🔜 Notification Service
- 🔜 Review Service

### Phase 4: 拡張サービス2 🔜 未実装
- 🔜 Chat Service
- 🔜 Search Service
- 🔜 Admin Service

### フロントエンド 🔜 未実装
- 🔜 Phoenix Web (顧客・オーナー用)
- 🔜 Admin Dashboard (管理者用)

---

## 🚀 サービス起動方法

### インフラストラクチャの起動
```bash
# 全サービス起動
make up

# または
cd infrastructure/docker
docker compose up -d

# 起動確認
docker compose ps
```

### マイクロサービスの起動（未実装）
```bash
# Auth Service
cd microservices/auth
go run cmd/main.go

# Shop Service
cd microservices/shop
go run cmd/main.go

# ... 他のサービスも同様
```

---

## 📖 関連ドキュメント

- [PORT_ASSIGNMENT.md](./PORT_ASSIGNMENT.md) - ポート割り当て詳細
- [SETUP.md](./SETUP.md) - 環境構築手順
- [CLAUDE.md](./CLAUDE.md) - 開発ガイド
- [requirements/](./requirements/) - 各サービスの要件定義

---

**最終更新**: 2026-01-25
**管理者**: Claude
