# Go MicroService 実践例

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

