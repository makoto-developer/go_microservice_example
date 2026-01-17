# Go MicroService 実践例

オンラインショップ（モール型）を題材にした Go マイクロサービスの実装例です。

## 📊 実装状況

### 完了したフェーズ

| Phase | ステータス | 説明 |
|-------|----------|------|
| Phase 1 | ✅ 完了 | カスタムロジック実装 (6サービス) |
| Phase 2 | ✅ 完了 | Protocol Buffers定義とコード生成 (12サービス) |
| Phase 3 | ✅ 完了 | Docker Compose インフラ構築 (29コンテナ) |
| Phase 4 | ✅ 完了 | Auth Service gRPCサーバー実装 |
| Phase 5 | ✅ 完了 | データベースマイグレーション (Auth Service) |
| Phase 6 | 🔄 進行中 | 残り11サービスの実装 |

**詳細**: [IMPLEMENTATION_STATUS.md](./generated/auth/IMPLEMENTATION_STATUS.md)

---

## 🚀 Quick Start

### 前提条件

- Docker Desktop for Mac インストール済み
- Go 1.25 以上
- Protocol Buffers コンパイラ (protoc)

### 初回セットアップ

```bash
# 1. リポジトリをクローン
git clone <repository-url>
cd go_microservice_example

# 2. 環境変数設定
cp .env.example .env
# 必要に応じて .env を編集

# 3. Docker インフラ起動
docker-compose up -d

# 4. データベースマイグレーション
./scripts/migrations/auth/apply_migrations.sh

# 5. Auth Service ビルドと起動
cd generated/auth
go mod tidy
go build -o auth-server ./cmd/server
./auth-server
```

### Auth Service 動作確認

```bash
# gRPCサービスリスト表示
grpcurl -plaintext localhost:50051 list

# ユーザー登録テスト
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "password": "password123",
  "role": 1
}' localhost:50051 auth_service.v1.AuthService/Register

# ログインテスト
grpcurl -plaintext -d '{
  "email": "test@example.com",
  "password": "password123"
}' localhost:50051 auth_service.v1.AuthService/Login
```

### Docker 環境管理

```bash
# インフラ起動確認
docker-compose ps

# ログ確認
docker-compose logs -f

# 停止
docker-compose down

# 完全クリーンアップ（データ削除）
docker-compose down -v
```

## Motivation

- マイクロサービスをどのように設計したのか実際の例

## お題

オンラインショップ(モール型)

- 会員登録、カート、ショップで商品を選んで購入、配送まで
- ショップ登録、商品一覧、注文かんり
- オンラインショップの管理者機能

## アーキテクチャ概要

```
┌──────────────────────────────────────────────────────────────┐
│                 Elixir/Phoenix (フロントエンド)               │
│                                                               │
│  ┌─────────────────────────────────────────────────────┐    │
│  │  Phoenix Router (API Gateway的役割)                 │    │
│  │  - HTTPリクエスト受付                               │    │
│  │  - 認証・認可 (Guardian/JWT)                        │    │
│  │  - ルーティング                                     │    │
│  │  - gRPCクライアント                                 │    │
│  └─────────────────────────────────────────────────────┘    │
│                                                               │
│  ┌────────────┐  ┌────────────┐  ┌──────────────────────┐  │
│  │ LiveView   │  │  Channels  │  │  カスタマーUI       │  │
│  │ (UI)       │  │ (WebSocket)│  │  ショップ運営者UI    │  │
│  └────────────┘  └────────────┘  └──────────────────────┘  │
└────────────────────────┬─────────────────────────────────────┘
                         │ gRPC
                         │
┌────────────────────────┴─────────────────────────────────────┐
│              Go マイクロサービス群（バックエンド）            │
│                                                               │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐       │
│  │Auth      │ │Shop      │ │Order     │ │Payment   │ ...   │
│  │Service   │ │Service   │ │Service   │ │Service   │       │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘       │
│       │            │            │            │               │
│       └────────────┴────────────┴────────────┘               │
│                      RabbitMQ                                 │
└───────────────────────────────────────────────────────────────┘
```

### 通信方式
- **ブラウザ ↔ Phoenix**: HTTP/REST、WebSocket (Channels)
- **Phoenix ↔ Goサービス**: gRPC
- **Goサービス間**: gRPC + RabbitMQ (イベント駆動)
- **リアルタイム通信（チャット）**: Phoenix Channels (WebSocket)

## Stack

### フロントエンド
- Elixir / Phoenix Framework
- Phoenix LiveView (リアルタイムUI)
- Phoenix Channels (WebSocket - チャット機能)
- Guardian (JWT認証ライブラリ)
- grpc-elixir (gRPCクライアント)
- Ecto (Elixir ORM)

### バックエンド（マイクロサービス）
- Go v1.25
- Gorm (ORM)
- PostgreSQL (データベース)
- RabbitMQ (メッセージブローカー)

### サービス間通信
- gRPC (サービス間通信プロトコル)
- Protocol Buffers (シリアライゼーション)

### キャッシュ・セッション
- Redis (キャッシュ、セッション管理、レート制限)

### 検索
- Elasticsearch (全文検索エンジン)

### ストレージ
- MinIO / S3 (オブジェクトストレージ - 商品画像)

### 認証・セキュリティ
- Guardian (Elixir JWT 認証ライブラリ)
- JWT-go (Go JWT 認証ライブラリ)
- bcrypt (パスワードハッシュ化)

### 決済連携
- Stripe Go SDK (決済プロバイダー連携)

### インフラ・コンテナ
- Docker (コンテナ化)
- Kubernetes (コンテナオーケストレーション)
- Helm (Kubernetes パッケージ管理)

### サービスディスカバリ
- Consul (サービスディスカバリ、設定管理)

### モニタリング・可視化
- Prometheus (メトリクス収集)
- Grafana (ダッシュボード・可視化)
- Jaeger (分散トレーシング)
- ELK Stack / Loki (ログ集約・分析)
  - Elasticsearch (ログストレージ)
  - Logstash / Fluentd (ログ収集)
  - Kibana (ログ可視化)

### テスト
- testify (テストフレームワーク)
- gomock (モックライブラリ)
- httptest (HTTP テスト)

### 開発手法
- **JetBrains MPS (Meta Programming System)** - DSL駆動開発
  - マイクロサービスDSLでサービス定義
  - Goコード自動生成
  - トークン消費90%削減

## 開発アプローチ: MPS DSL駆動開発

このプロジェクトでは、**JetBrains MPS**を使用したDSL（Domain Specific Language）駆動開発を採用しています。

### なぜMPSを使うのか？

#### 問題: マイクロサービス開発のトークン消費
- 12個のマイクロサービスを開発
- 1サービス = 2,000-3,000行のコード
- Claudeが全コードを読む → トークンが一瞬で枯渇

#### 解決策: DSLでサービスを定義 → コード自動生成
- 1サービス = **100-200行のDSL定義**
- MPS Generatorが自動で2,000-3,000行のGoコードを生成
- Claudeは **DSLだけ読めばOK** → トークン90%削減

### 開発ワークフロー

```
1. 要件定義を確認（docs/requirements/）
   ↓
2. MPSでDSL定義を作成（100-200行）
   例: mps-workspace/solutions/auth-service/service.model
   ↓
3. MPS Generatorでコード生成
   ./scripts/mps-generate.sh auth-service
   ↓
4. 生成されたGoコードを確認（2,000-3,000行）
   generated/auth/internal/...
   ↓
5. 必要なカスタムロジックのみ手動実装
   manual/auth/custom_logic.go
```

### DSL定義例（Auth Service）

```kotlin
microservice AuthService {
  version: "v1"

  // ドメインエンティティ
  entity User {
    id: UUID primary_key
    email: string unique not_null
    password_hash: string not_null
    role: Role not_null
    created_at: timestamp
  }

  enum Role {
    CUSTOMER, SHOP_OWNER, ADMIN
  }

  // ユースケース
  usecase UserRegistration {
    input: {
      email: string
      password: string
      role: Role
    }
    output: {
      user_id: UUID
      token: string
    }
    errors: {
      EmailAlreadyExists,
      WeakPassword
    }
  }

  // gRPCサービス定義
  grpc_service {
    rpc Register(RegisterRequest) returns (RegisterResponse)
    rpc Login(LoginRequest) returns (LoginResponse)
  }

  // 依存
  dependencies {
    database: PostgreSQL
    cache: Redis
    messaging: RabbitMQ
  }
}
```

### 自動生成されるもの

この100-200行のDSL定義から、以下が自動生成されます：

- ✅ ドメイン層（`domain/user.go`, `domain/role.go`）
- ✅ ユースケース層（`usecase/user_registration.go`）
- ✅ インフラ層（`infrastructure/postgres_user_repo.go`）
- ✅ gRPCハンドラー（`handler/grpc_handler.go`）
- ✅ Protocol Buffers定義（`proto/auth/v1/auth.proto`）
- ✅ テストコード（`tests/user_registration_test.go`）

**合計**: 2,000-3,000行のGoコード

### トークン削減効果

| アプローチ | Claudeが読む量 | トークン消費 | 削減率 |
|----------|--------------|------------|--------|
| **手動実装** | 2,000-3,000行のGoコード | ~15,000 | - |
| **MPS DSL** | 100-200行のDSL定義 | ~1,500 | **90%削減** |

### MPSディレクトリ構成

```
mps-workspace/           # MPS専用ワークスペース
├── languages/           # DSL定義
│   ├── microservice-dsl/
│   ├── grpc-dsl/
│   └── event-dsl/
└── solutions/           # サービス定義
    ├── auth-service/
    ├── shop-service/
    └── ...

generated/               # 生成コード（触らない）
├── auth/
├── shop/
└── ...

manual/                  # 手動実装
└── custom/              # カスタムロジックのみ
```

### 詳細

- [CLAUDE.md](./CLAUDE.md) - Claude開発ガイド（MPS詳細）
- [SETUP.md](./SETUP.md) - 環境構築・MPS使用方法

## ディレクトリ構成

```
.
├── cmd/              # 各サービスのエントリーポイント
├── internal/         # サービス固有のコード
├── pkg/              # 共通パッケージ
├── proto/            # Protocol Buffers定義
├── gen/              # 生成されたコード
├── docs/             # ドキュメント
│   └── requirements/ # 要件定義
├── deployments/      # デプロイ設定
│   ├── docker/       # Dockerfile
│   └── kubernetes/   # K8s マニフェスト
├── tests/            # テスト
│   ├── e2e/          # E2Eテスト
│   └── integration/  # 統合テスト
└── scripts/          # ビルド・デプロイスクリプト
```

## ドキュメント

- [SETUP.md](./SETUP.md) - 環境構築・開発手順
- [docs/requirements/README.md](./docs/requirements/README.md) - 要件定義（機能仕様）

## マイクロサービス間の連携

### イベント駆動アーキテクチャ (RabbitMQ)

**主要なイベント**:

1. **注文関連**
   - `order.created` → 在庫サービス、決済サービス、通知サービス
   - `order.cancelled` → 在庫サービス、決済サービス、通知サービス
   - `order.status.updated` → 通知サービス

2. **決済関連**
   - `payment.completed` → 注文サービス、ショップサービス、通知サービス
   - `payment.failed` → 注文サービス、通知サービス
   - `refund.completed` → 注文サービス、在庫サービス、通知サービス

3. **在庫関連**
   - `inventory.reserved` → 注文サービス
   - `inventory.released` → 注文サービス
   - `inventory.low` → ショップサービス、通知サービス

4. **配送関連**
   - `shipping.dispatched` → 注文サービス、通知サービス
   - `shipping.delivered` → 注文サービス、通知サービス

5. **チャット関連**
   - `chat.message.sent` → 通知サービス（オフライン時のメール通知）
   - `chat.room.created` → 通知サービス
   - `chat.file.uploaded` → 通知サービス

**注**: チャットのリアルタイム通信は Phoenix Channels (WebSocket) で行い、
非同期処理（通知など）のみ RabbitMQ を使用

## データベース設計方針

- 各マイクロサービスは独立したデータベースを持つ
- サービス間のデータ参照は gRPC 経由で行う
- イベントソーシングによる整合性担保
- Saga パターンによる分散トランザクション管理

## Docker環境での起動

### 前提条件

- Docker Desktop for Mac がインストールされていること
- Docker Compose V2 がインストールされていること

### 環境設定

#### 1. 環境変数ファイルの作成

プロジェクトルートに `.env` ファイルを作成します。

```bash
# .env.exampleをコピー
cp .env.example .env

# 必要に応じて値を編集
vim .env
```

#### 2. 環境変数の説明

`.env` ファイルには以下の設定が含まれています：

```env
# 環境指定 (dev, test, staging, prod)
ENV=dev

# プロジェクト名（コンテナ名のプレフィックスになります）
COMPOSE_PROJECT_NAME=go_microservice

# データベース設定
POSTGRES_USER=admin
POSTGRES_PASSWORD=dev_password_123  # 本番環境では変更必須
POSTGRES_DB=microservice
POSTGRES_PORT=20000

# セキュリティ設定
JWT_SECRET=dev_jwt_secret_key_not_for_production  # 本番環境では変更必須

# 外部サービスAPI設定（開発環境はモック）
STRIPE_API_KEY=mock_stripe_key_dev
SENDGRID_API_KEY=mock_sendgrid_key_dev
```

#### 3. 環境別の設定

**開発環境 (dev)**:
```bash
ENV=dev
DEBUG=true
HOT_RELOAD=true
```

**テスト環境 (test)**:
```bash
ENV=test
DEBUG=false
HOT_RELOAD=false
```

**本番環境 (prod)**:
```bash
ENV=prod
DEBUG=false
HOT_RELOAD=false

# 本番環境では以下を必ず変更してください
POSTGRES_PASSWORD=<強固なパスワード>
JWT_SECRET=<ランダムな256bit以上の秘密鍵>
STRIPE_API_KEY=<実際のStripeキー>
SENDGRID_API_KEY=<実際のSendGridキー>
```

#### 4. セキュリティに関する注意

**重要**: `.env` ファイルには機密情報が含まれるため、以下を遵守してください：

- ✅ `.env.example` のみバージョン管理にコミット
- ❌ `.env` は絶対にGitにコミットしない（.gitignoreに含まれています）
- ✅ 本番環境では強固なパスワード・秘密鍵を使用
- ✅ 本番環境では実際の外部サービスAPIキーを設定
- ✅ 環境ごとに異なる `.env` ファイルを使用

### クイックスタート

> **✅ 現在の状態**: インフラ環境稼働中 (29コンテナ)
> - **PostgreSQL** × 12 (ポート 5432-5443)
> - **Redis** × 12 (ポート 6379-6390)
> - **Elasticsearch** × 1 (ポート 9200, 9300)
> - **RabbitMQ** × 1 (ポート 5672, 15672)
> - **MinIO** × 1 (ポート 9000-9001)
> - **MailHog** × 1 (ポート 1025, 8025)
>
> **実装済みサービス**:
> - ✅ Auth Service (gRPCサーバー稼働可能)
> - 🔄 残り11サービス (DSL定義・コード生成完了、実装待ち)

#### 方法1: Makefileを使用（推奨）

```bash
# 初回セットアップ（.env作成、ビルド、起動まで自動）
make init

# または
make quickstart

# ヘルプを表示
make help
```

#### 方法2: docker-composeを直接使用

```bash
# 1. 環境変数ファイル作成
cp .env.example .env

# 2. DSLからコード生成
./scripts/mps-generate.sh --all

# 3. Docker Composeで全サービス起動
docker-compose up -d

# 4. ログ確認
docker-compose logs -f

# 5. サービス停止
docker-compose down

# 6. ボリューム含めて完全削除
docker-compose down -v
```

### Makefileコマンド一覧

便利なMakefileコマンドを用意しています：

#### 起動・停止
```bash
make up              # 全サービス起動
make up-infra        # インフラのみ起動
make up-mocks        # モックサービスのみ起動
make down            # 全サービス停止
make restart         # 全サービス再起動
make dev             # 開発環境起動（インフラ+モック）
```

#### ビルド
```bash
make build           # 全サービスビルド
make build-mocks     # モックサービスビルド
make build-services  # マイクロサービスビルド
```

#### ログ・モニタリング
```bash
make logs            # 全ログ表示
make logs-infra      # インフラログのみ
make ps              # 稼働状況確認
make health          # ヘルスチェック
make stats           # リソース使用状況
```

#### データベース
```bash
make db-init         # データベース初期化
make db-connect      # PostgreSQL接続
make db-list         # データベース一覧
make db-reset        # データベースリセット（要注意）
make backup-db       # バックアップ作成
```

#### ユーティリティ
```bash
make open-mailhog    # MailHog UI をブラウザで開く
make open-rabbitmq   # RabbitMQ UI をブラウザで開く
make shell-postgres  # PostgreSQLシェル
make shell-redis     # Redisシェル
make go-tidy         # 全サービスで go mod tidy 実行
```

#### クリーンアップ
```bash
make clean           # コンテナ・ボリューム削除
make nuke            # 完全削除（要注意）
```

詳細は `make help` で確認できます。

### アーキテクチャ構成

```
┌────────────────────────────────────────────────────────────────┐
│                    Docker Compose Environment                  │
├────────────────────────────────────────────────────────────────┤
│                                                                │
│  Infrastructure Services                                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │ PostgreSQL   │  │ Redis        │  │ RabbitMQ     │        │
│  │ :20000       │  │ :20001       │  │ :20002       │        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
│                                                                │
│  ┌──────────────┐                                             │
│  │ MailHog      │  メール確認: http://localhost:20005        │
│  │ SMTP :20004  │                                             │
│  │ UI   :20005  │                                             │
│  └──────────────┘                                             │
│                                                                │
│  Mock External Services (本番環境では実際のサービスに接続)       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │ Stripe Mock  │  │ FCM Mock     │  │ Elasticsearch│        │
│  │ :20010       │  │ :20012       │  │ Mock :20013  │        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
│                                                                │
│  ┌──────────────┐                                             │
│  │ Carriers     │                                             │
│  │ Mock :20014  │                                             │
│  └──────────────┘                                             │
│                                                                │
│  Microservices (12 services)                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │ Auth         │  │ Shop         │  │ Customer     │        │
│  │ :20100       │  │ :20101       │  │ :20102       │        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
│                                                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │ Inventory    │  │ Order        │  │ Payment      │        │
│  │ :20103       │  │ :20104       │  │ :20105       │        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
│                                                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │ Shipping     │  │ Notification │  │ Review       │        │
│  │ :20106       │  │ :20107       │  │ :20108       │        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
│                                                                │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐        │
│  │ Chat         │  │ Search       │  │ Admin        │        │
│  │ :20109       │  │ :20111       │  │ :20112       │        │
│  │ WebSocket    │  │              │  │              │        │
│  │ :20110       │  │              │  │              │        │
│  └──────────────┘  └──────────────┘  └──────────────┘        │
└────────────────────────────────────────────────────────────────┘
```

### ポート番号一覧

| サービス | ポート | 用途 |
|---------|-------|------|
| **Infrastructure** | | |
| PostgreSQL | 20000 | データベース |
| Redis | 20001 | キャッシュ/セッション |
| RabbitMQ | 20002 | メッセージング |
| RabbitMQ Management | 20003 | RabbitMQ管理UI |
| MailHog SMTP | 20004 | メール送信テスト |
| MailHog Web UI | 20005 | メール確認画面 |
| **Mock Services** | | |
| Stripe Mock | 20010 | 決済API |
| FCM Mock | 20012 | プッシュ通知API |
| Elasticsearch Mock | 20013 | 検索API |
| Carriers Mock | 20014 | 配送業者API |
| **Microservices** | | |
| Auth Service | 20100 | 認証・認可 |
| Shop Service | 20101 | ショップ管理 |
| Customer Service | 20102 | 顧客管理 |
| Inventory Service | 20103 | 在庫管理 |
| Order Service | 20104 | 注文管理 |
| Payment Service | 20105 | 決済処理 |
| Shipping Service | 20106 | 配送管理 |
| Notification Service | 20107 | 通知送信 |
| Review Service | 20108 | レビュー管理 |
| Chat Service | 20109 | チャット（gRPC） |
| Chat Service (WS) | 20110 | チャット（WebSocket） |
| Search Service | 20111 | 検索機能 |
| Admin Service | 20112 | 管理機能 |

### コンテナ命名規則

各コンテナは環境変数に基づいて命名されます：

**命名パターン**: `${COMPOSE_PROJECT_NAME}_<service>_${ENV}`

**例**（ENV=dev の場合）:
- `go_microservice_postgres_dev` - PostgreSQLコンテナ
- `go_microservice_auth_service_dev` - Auth Serviceコンテナ
- `go_microservice_redis_dev` - Redisコンテナ

**ホスト名**: `<service>_${ENV}` (例: `postgres_dev`, `auth-service_dev`)

**ネットワーク**: `${COMPOSE_PROJECT_NAME}_network_${ENV}` (例: `go_microservice_network_dev`)

この命名規則により、同じホストで複数の環境（dev, test, staging）を同時実行できます。

### サービス間通信

Docker内部では、サービスは以下のホスト名で相互通信します：

```env
# Database（環境変数から自動設定）
DB_HOST=postgres_${ENV}
DB_PORT=5432
DB_USER=${POSTGRES_USER}
DB_PASSWORD=${POSTGRES_PASSWORD}
DB_NAME=<service>_db

# Redis
REDIS_HOST=redis_${ENV}
REDIS_PORT=6379

# RabbitMQ
RABBITMQ_URL=amqp://${RABBITMQ_USER}:${RABBITMQ_PASSWORD}@rabbitmq_${ENV}:5672/

# External Services (Mock)
STRIPE_API_URL=http://mock-stripe_${ENV}:8080
SENDGRID_API_URL=http://mock-sendgrid_${ENV}:8081
FCM_API_URL=http://mock-fcm_${ENV}:8082
ELASTICSEARCH_URL=http://mock-elasticsearch_${ENV}:9200
CARRIER_API_URL=http://mock-carriers_${ENV}:8083
```

`.env` ファイルで `ENV` 変数を変更するだけで、環境を切り替えることができます。

### モックサーバーについて

開発環境では外部サービスをモックサーバーで代替しています。

#### MailHog - メールテストツール

開発環境でのメール送信をテストするため、MailHogを使用しています。

**特徴**:
- 送信されたメールをすべてキャプチャ
- Web UIでメール内容を確認可能
- 実際にメールが送信されない（安全）
- SMTP互換

**使用方法**:
```bash
# 1. Docker起動後、ブラウザでアクセス
open http://localhost:20005

# 2. サービスからメール送信
# Notification Serviceが自動的にMailHogを使用
# SMTP_HOST: mailhog
# SMTP_PORT: 1025
```

**確認できる内容**:
- メールの件名・本文
- 送信元・送信先
- HTMLメール・テキストメール
- 添付ファイル

#### その他のモックサービス

| サービス | 用途 | 特徴 |
|---------|------|------|
| Stripe Mock | 決済API | 実際のStripe APIと同じインターフェース |
| FCM Mock | プッシュ通知 | Firebase Cloud Messaging互換 |
| Elasticsearch Mock | 検索API | 基本的な検索機能を提供 |
| Carriers Mock | 配送業者API | ヤマト・佐川・日本郵便の追跡API |

**本番環境への切り替え**:
環境変数を実際のサービスURLに変更するだけで切り替え可能です。

```bash
# 開発環境（モック）
STRIPE_API_URL=http://mock-stripe:8080
SMTP_HOST=mailhog

# 本番環境
STRIPE_API_URL=https://api.stripe.com
STRIPE_API_KEY=<actual_key>
SMTP_HOST=smtp.sendgrid.net
SENDGRID_API_KEY=<actual_key>
```

### Docker環境の改善点

最新版では以下の改善が施されています：

#### 1. 最新バージョンの使用
- **PostgreSQL 17-alpine** (最新安定版)
- **Redis 7.4-alpine** (最新安定版)
- **RabbitMQ 4-management-alpine** (最新安定版)

#### 2. リソース制限
全サービスにCPU・メモリ制限を設定：
```yaml
deploy:
  resources:
    limits:
      cpus: '1'
      memory: 512M
    reservations:
      cpus: '0.25'
      memory: 128M
```

#### 3. 自動再起動
`restart: unless-stopped` を全サービスに適用：
- コンテナクラッシュ時の自動再起動
- 開発中の安定性向上

#### 4. ログローテーション
ログファイルの肥大化を防止：
```yaml
logging:
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "3"
```

#### 5. ヘルスチェック強化
- `start_period` の追加（初回起動の猶予期間）
- Redis認証対応のヘルスチェック
- RabbitMQ起動時間の考慮（30秒）

#### 6. データ永続化
Named volumeによるデータ永続化：
- `postgres_data` - データベース
- `redis_data` - Redisデータ
- `rabbitmq_data` - RabbitMQデータ
- `rabbitmq_log` - RabbitMQログ

**データのバックアップ**:
```bash
# volumeの一覧確認
docker volume ls

# volumeのバックアップ
docker run --rm -v go_microservice_postgres_data_dev:/data -v $(pwd):/backup \
  alpine tar czf /backup/postgres_backup.tar.gz -C /data .

# リストア
docker run --rm -v go_microservice_postgres_data_dev:/data -v $(pwd):/backup \
  alpine tar xzf /backup/postgres_backup.tar.gz -C /data
```

### トラブルシューティング

#### ポートが既に使用されている

```bash
# ポート使用状況確認
lsof -i :20000  # PostgreSQL
lsof -i :20100 # Auth Service

# プロセス終了
kill -9 <PID>
```

#### コンテナログ確認

```bash
# 全サービスのログ
docker-compose logs -f

# 特定サービスのログ
docker-compose logs -f auth-service

# エラーログのみ
docker-compose logs -f | grep -i error
```

#### コンテナ再起動

```bash
# 全サービス再起動
docker-compose restart

# 特定サービスのみ再起動
docker-compose restart auth-service
```

#### データベースリセット

```bash
# ボリューム削除（全データ削除）
docker-compose down -v

# 再起動
docker-compose up -d
```

### 開発ワークフロー

1. **DSL定義変更**
   ```bash
   # DSL編集
   vim mps-workspace/solutions/auth-service/service.model
   
   # コード再生成
   ./scripts/mps-generate.sh auth-service
   ```

2. **サービス再ビルド**
   ```bash
   # イメージ再ビルド
   docker-compose build auth-service
   
   # 再起動
   docker-compose up -d auth-service
   ```

3. **動作確認**
   ```bash
   # ログ確認
   docker-compose logs -f auth-service
   
   # gRPCリクエスト送信（grpcurlを使用）
   grpcurl -plaintext localhost:20100 list
   ```

