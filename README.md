# 🛍️ Go Microservice Example - Online Shop Mall

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![Elixir Version](https://img.shields.io/badge/Elixir-1.14+-4B275F?style=flat&logo=elixir)](https://elixir-lang.org/)
[![Phoenix Version](https://img.shields.io/badge/Phoenix-1.7+-FF6F61?style=flat&logo=phoenixframework)](https://www.phoenixframework.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

オンラインショップモールを題材にした**本格的なマイクロサービスアーキテクチャ**の実装例です。

Database per Service パターン、Saga トランザクション、Event-Driven Architecture など、実践的な設計パターンを網羅しています。

---

## 📋 目次

- [概要](#-概要)
- [特徴](#-特徴)
- [アーキテクチャ](#-アーキテクチャ)
- [技術スタック](#-技術スタック)
- [導入方法](#-導入方法)
- [設計思想](#-設計思想)
- [サービス一覧](#-サービス一覧)
- [開発ワークフロー](#-開発ワークフロー)
- [テスト](#-テスト)
- [ドキュメント](#-ドキュメント)
- [ライセンス](#-ライセンス)

---

## 🎯 概要

このプロジェクトは、**マイクロサービスアーキテクチャのベストプラクティス**を実装した、実践的なサンプルアプリケーションです。

### ビジネスドメイン

オンラインショップモール（マーケットプレイス型EC）

- 👤 **顧客**: 商品検索、注文、レビュー投稿
- 🏪 **ショップオーナー**: 店舗・商品管理、在庫管理、売上確認
- 👨‍💼 **管理者**: ユーザー管理、監査ログ、システム設定

### プロジェクト規模

- **マイクロサービス数**: 12サービス
- **データベース**: PostgreSQL 12インスタンス（各サービス専用DB）
- **キャッシュ**: Redis 12インスタンス
- **コード行数**: 約29,500行（Go + Elixir）
- **テストケース**: 51ケース（統合テスト + E2E）

---

## ✨ 特徴

### 🏗️ アーキテクチャパターン

- ✅ **Database per Service** - サービスごとに独立したデータベース
- ✅ **API Gateway Pattern** - 統一されたAPIエンドポイント
- ✅ **Saga Pattern** - 分散トランザクション管理
- ✅ **Event-Driven Architecture** - イベント駆動による疎結合
- ✅ **CQRS** - コマンドとクエリの分離
- ✅ **Circuit Breaker** - 障害の伝播防止

### 🚀 開発効率

- ✅ **MPS DSL駆動開発** - DSL定義からGoコード自動生成
- ✅ **Makefileによる統一操作** - `make up`で全サービス起動
- ✅ **包括的なテスト** - 統合テスト、E2E、性能テスト
- ✅ **Docker Compose** - ローカル開発環境の簡単セットアップ

### 🔒 セキュリティ

- ✅ **JWT認証** - トークンベース認証
- ✅ **ロールベースアクセス制御** - Customer/Owner/Admin
- ✅ **API Rate Limiting** - レート制限
- ✅ **監査ログ** - 全操作のトレーサビリティ

---

## 🏛️ アーキテクチャ

### システム構成図

```
┌─────────────────────────────────────────────────────────────┐
│                       API Gateway                            │
│                    (Future: Nginx/Envoy)                     │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌──────────────┐    ┌──────────────┐     ┌──────────────┐
│ Auth Service │    │ Shop Service │     │Customer Svc  │
│  (Go gRPC)   │    │  (Phoenix)   │     │  (Go gRPC)   │
│   :22100     │    │    :4000     │     │   :22102     │
└──────┬───────┘    └──────┬───────┘     └──────┬───────┘
       │                   │                    │
       ▼                   ▼                    ▼
   PostgreSQL          PostgreSQL          PostgreSQL
   :22010 (auth)       :22011 (shop)       :22012 (customer)

┌──────────────┐    ┌──────────────┐     ┌──────────────┐
│Inventory Svc │    │ Order Service│     │Payment Svc   │
│  (Go gRPC)   │    │  (Go gRPC)   │     │  (Go gRPC)   │
│   :22103     │    │   :22104     │     │   :22105     │
└──────┬───────┘    └──────┬───────┘     └──────┬───────┘
       │                   │                    │
       ▼                   ▼                    ▼
   PostgreSQL          PostgreSQL          PostgreSQL
   :22013              :22014              :22015

┌──────────────┐    ┌──────────────┐     ┌──────────────┐
│Notification  │    │ Review Svc   │     │Shipping Svc  │
│   Service    │    │  (Go gRPC)   │     │  (Go gRPC)   │
│   :22106     │    │   :22107     │     │   :22108     │
└──────┬───────┘    └──────┬───────┘     └──────┬───────┘
       │                   │                    │
       ▼                   ▼                    ▼
   PostgreSQL          PostgreSQL          PostgreSQL
   :22017              :22018              :22016

┌──────────────┐    ┌──────────────┐     ┌──────────────┐
│  Chat Svc    │    │ Search Svc   │     │  Admin Svc   │
│  (Go gRPC)   │    │  (Go gRPC)   │     │  (Go gRPC)   │
│   :22109     │    │   :22110     │     │   :22111     │
└──────┬───────┘    └──────┬───────┘     └──────┬───────┘
       │                   │                    │
       ▼                   ▼                    ▼
   PostgreSQL          PostgreSQL          PostgreSQL
   :22019              :22020              :22021
```

### データフロー例: 注文処理 (Saga Pattern)

```
Customer → Order Service → Inventory Service (在庫引当)
              │                    │
              │                    ▼
              │            在庫予約成功/失敗
              │                    │
              ▼                    │
        Payment Service ◀──────────┘
              │
              ▼
        決済処理 (Stripe)
              │
       ┌──────┴──────┐
       ▼             ▼
    成功          失敗
       │             │
       ▼             ▼
  注文確定      Compensating
Notification   Transaction
   送信         (在庫解放)
```

---

## 🔧 技術スタック

### Backend

| カテゴリ | 技術 | 用途 |
|---------|------|------|
| **言語** | Go 1.25 | メインサービス（11サービス） |
| | Elixir 1.14 + Phoenix 1.7 | Shop Service (Web UI) |
| **通信** | gRPC + Protocol Buffers | サービス間通信 |
| **データベース** | PostgreSQL 15 | 永続化層（12インスタンス） |
| **キャッシュ** | Redis 7 | セッション、キャッシュ |
| **メッセージング** | RabbitMQ | イベント配信（Future） |
| **認証** | JWT | トークンベース認証 |

### Infrastructure

| カテゴリ | 技術 |
|---------|------|
| **コンテナ** | Docker + Docker Compose |
| **オーケストレーション** | Kubernetes (Future) |
| **CI/CD** | GitHub Actions (Future) |
| **監視** | Prometheus + Grafana (Future) |
| **ログ** | ELK Stack (Future) |

### Development Tools

| カテゴリ | 技術 |
|---------|------|
| **DSL開発** | JetBrains MPS |
| **ビルド** | Make, Go modules |
| **テスト** | Go testing, testify |
| **AI支援** | Claude (Anthropic), Serena MCP |

---

## 🚀 導入方法

### 前提条件

- **Go**: 1.25以上
- **Elixir**: 1.14以上
- **Docker**: 20.10以上
- **Docker Compose**: 2.0以上
- **Make**: 任意のバージョン

### クイックスタート（1分で起動）

```bash
# 1. リポジトリクローン
git clone https://github.com/makoto-developer/go_microservice_example.git
cd go_microservice_example

# 2. データベース起動
cd infrastructure/docker
docker compose up -d
cd ../..

# 3. 全サービス起動
make up

# 4. 稼働確認
make status
```

**期待される出力**:
```
✅ Auth Service          Running (PID: 12345, Port: 22100)
✅ Shop Service          Running (PID: 12346, Port:  4000)
✅ Customer Service      Running (PID: 12347, Port: 22102)
...
Summary: 12/12 services running
```

### 詳細な導入手順

詳細は [QUICKSTART.md](./QUICKSTART.md) を参照してください。

---

## 💡 設計思想

### 1. Database per Service パターン

**各マイクロサービスが専用のデータベースを持つ**

#### メリット

- ✅ **データの独立性**: サービス間でスキーマ変更の影響を受けない
- ✅ **技術選択の自由**: サービスごとに最適なDBを選択可能
- ✅ **スケーラビリティ**: サービス単位でDB負荷を分散
- ✅ **障害の局所化**: 1つのDBダウンが全体に波及しない

#### 実装

```
Auth Service      → PostgreSQL (auth_db)       :22010
Shop Service      → PostgreSQL (shop_db)       :22011
Customer Service  → PostgreSQL (customer_db)   :22012
Inventory Service → PostgreSQL (inventory_db)  :22013
Order Service     → PostgreSQL (order_db)      :22014
Payment Service   → PostgreSQL (payment_db)    :22015
...
```

### 2. Saga パターン（分散トランザクション）

**複数サービスにまたがるトランザクションを管理**

#### 実装例: 注文処理

```go
// Order Service: Saga Orchestrator
func (s *OrderService) CreateOrder(order Order) error {
    // Step 1: 在庫引当
    if err := s.inventoryClient.ReserveStock(order.Items); err != nil {
        return err // 即座に失敗
    }

    // Step 2: 注文作成
    if err := s.orderRepo.Create(order); err != nil {
        // Compensating Transaction
        s.inventoryClient.ReleaseStock(order.Items)
        return err
    }

    // Step 3: 決済処理
    if err := s.paymentClient.ProcessPayment(order.Payment); err != nil {
        // Compensating Transaction
        s.orderRepo.Cancel(order.ID)
        s.inventoryClient.ReleaseStock(order.Items)
        return err
    }

    // Step 4: 通知送信
    s.notificationClient.SendOrderConfirmation(order)

    return nil
}
```

### 3. gRPC による高速通信

**Protocol Buffers を使用した効率的なサービス間通信**

#### メリット

- ✅ **高速**: バイナリフォーマットでRESTより軽量
- ✅ **型安全**: .protoファイルで厳密な型定義
- ✅ **多言語対応**: Go, Elixir, Python等で相互通信可能
- ✅ **双方向ストリーミング**: リアルタイム通信に対応

### 4. Event-Driven Architecture

**イベントを介したサービス間の疎結合化**

#### イベント例

```
OrderCreated → Notification Service (注文確認メール)
              → Shipping Service (配送準備)
              → Analytics Service (売上集計)

PaymentCompleted → Order Service (注文確定)
                  → Inventory Service (在庫確定)
                  → Notification Service (決済完了通知)
```

### 5. MPS DSL駆動開発

**JetBrains MPS による Domain-Specific Language 開発**

#### ワークフロー

```
1. DSL定義 (100-200行)
   ↓
2. MPS Generator実行
   ↓
3. Goコード自動生成 (2,000-3,000行)
   ↓
4. カスタムロジック追加 (必要な場合のみ)
```

#### メリット

- ✅ **生産性向上**: コード量90%削減
- ✅ **一貫性**: Generator保証による品質担保
- ✅ **保守性**: DSL変更で全サービス一括更新

詳細: [.claude/CLAUDE.md](./.claude/CLAUDE.md)

---

## 📦 サービス一覧

| # | サービス名 | ポート | 言語 | 責務 | データベース |
|---|-----------|-------|------|------|-------------|
| 1 | **Auth Service** | 22100 | Go | 認証・認可、JWT発行 | postgres_auth:22010 |
| 2 | **Shop Service** | 4000 | Phoenix | Web UI、ショップ管理 | postgres_shop:22011 |
| 3 | **Customer Service** | 22102 | Go | 顧客情報、住所管理 | postgres_customer:22012 |
| 4 | **Inventory Service** | 22103 | Go | 在庫管理、在庫引当 | postgres_inventory:22013 |
| 5 | **Order Service** | 22104 | Go | 注文管理、Saga制御 | postgres_order:22014 |
| 6 | **Payment Service** | 22105 | Go | 決済処理、Stripe連携 | postgres_payment:22015 |
| 7 | **Shipping Service** | 22108 | Go | 配送管理、追跡 | postgres_shipping:22016 |
| 8 | **Notification Service** | 22106 | Go | メール・SMS送信 | postgres_notification:22017 |
| 9 | **Review Service** | 22107 | Go | レビュー・評価管理 | postgres_review:22018 |
| 10 | **Chat Service** | 22109 | Go | チャット、問い合わせ | postgres_chat:22019 |
| 11 | **Search Service** | 22110 | Go | 全文検索、商品検索 | postgres_search:22020 |
| 12 | **Admin Service** | 22111 | Go | 管理機能、監査ログ | postgres_admin:22021 |

### サービス依存関係

```
Auth Service (基盤)
    ↓
Shop Service ──┬─→ Inventory Service
    ↓          │
Customer Service
    ↓          │
Order Service ←┘
    ↓
Payment Service
    ↓
Shipping Service
    ↓
Notification Service
```

---

## 🔄 開発ワークフロー

### Makefileコマンド

| コマンド | 説明 |
|---------|------|
| `make help` | 全コマンド一覧表示 |
| `make up` | 全サービス起動（DB + マイクロサービス） |
| `make down` | 全サービス停止 |
| `make status` | サービス稼働状況確認 |
| `make logs` | 全サービスのログ表示 |
| `make logs-follow` | ログをリアルタイム表示 |
| `make test` | 統合テスト実行 |
| `make test-e2e` | E2Eテスト実行 |
| `make build` | 全サービスビルド |
| `make clean` | ビルド成果物削除 |
| `make db-status` | データベース状態確認 |
| `make dashboard` | サービスダッシュボード表示 |

### 開発の流れ

```bash
# 1. サービス起動
make up

# 2. 開発（コード編集）
vim microservices/auth/handler/grpc_handler.go

# 3. ビルド
make build

# 4. サービス再起動
make restart

# 5. テスト
make test

# 6. ログ確認
make logs-follow
```

---

## 🧪 テスト

### テストスイート構成

| スイート | テストケース数 | 対象 |
|---------|--------------|------|
| **統合テスト - Auth** | 10 | 認証フロー（登録、ログイン、トークン） |
| **統合テスト - Order Flow** | 4 | 注文・決済・在庫連携 |
| **統合テスト - Notification** | 18 | 通知、レビュー、配送 |
| **E2Eテスト** | 19 | 完全な購入フロー、エラーシナリオ |
| **合計** | **51** | - |

### テスト実行

```bash
# 全テスト実行
make test

# 個別テスト
make test-auth        # Auth認証テスト（30秒）
make test-order       # 注文フローテスト（2-3分）
make test-e2e         # E2Eテスト（5-10分）
```

### テストカバレッジ

- **サービスカバレッジ**: 10/12サービス（83%）
- **コードカバレッジ**: 約70%（目標: 80%以上）

詳細: [INTEGRATION_TEST_COMPLETE.md](./INTEGRATION_TEST_COMPLETE.md)

---

## 📚 ドキュメント

### プロジェクトドキュメント

| ドキュメント | 内容 |
|------------|------|
| [README.md](./README.md) | プロジェクト概要（このファイル） |
| [QUICKSTART.md](./QUICKSTART.md) | 1分で起動するガイド |
| [RUNNING_SERVICES_DASHBOARD.md](./RUNNING_SERVICES_DASHBOARD.md) | サービスダッシュボード |
| [PROJECT_STATUS_FINAL.md](./PROJECT_STATUS_FINAL.md) | プロジェクト状況レポート |
| [INTEGRATION_TEST_COMPLETE.md](./INTEGRATION_TEST_COMPLETE.md) | テストドキュメント |
| [ALL_SERVICES_STATUS.md](./ALL_SERVICES_STATUS.md) | 全サービス実装状況 |

### 開発ドキュメント（Claude/AI用）

| ドキュメント | 内容 |
|------------|------|
| [.claude/CLAUDE.md](./.claude/CLAUDE.md) | プロジェクト開発ガイド |
| [.claude/rules/mps-workflow.md](./.claude/rules/mps-workflow.md) | MPS開発フロー |
| [.claude/rules/code-generation.md](./.claude/rules/code-generation.md) | コード生成ルール |
| [.claude/rules/token-optimization.md](./.claude/rules/token-optimization.md) | トークン最適化 |
| [.claude/rules/phase-execution-plan.md](./.claude/rules/phase-execution-plan.md) | Phase別実行計画 |
| [.claude/rules/autonomous-execution.md](./.claude/rules/autonomous-execution.md) | 自律実行ルール |

### 要件定義

要件定義は `docs/requirements/` に格納：

```
docs/requirements/
├── README.md                    # 要件定義の目次
├── 01_auth_service.md           # Auth Service要件
├── 02_shop_service.md           # Shop Service要件
├── 03_customer_service.md       # Customer Service要件
└── ...
```

---

## 📁 ディレクトリ構造

```
go_microservice_example/
│
├── microservices/              # 🎯 マイクロサービス実装
│   ├── auth/                   # Auth Service (Go)
│   │   ├── cmd/server/         # エントリーポイント
│   │   ├── internal/           # ビジネスロジック
│   │   │   ├── domain/         # ドメインモデル
│   │   │   ├── usecase/        # ユースケース
│   │   │   ├── handler/        # gRPCハンドラー
│   │   │   └── infrastructure/ # Repository実装
│   │   └── go.mod
│   └── shop/                   # Shop Service (Phoenix)
│       ├── lib/                # Elixirコード
│       ├── priv/               # データベースマイグレーション
│       └── mix.exs
│
├── simple-servers/             # 🚀 シンプル実装サービス
│   ├── customer/               # Customer Service
│   ├── inventory/              # Inventory Service
│   ├── order/                  # Order Service
│   ├── payment/                # Payment Service
│   ├── shipping/               # Shipping Service
│   ├── notification/           # Notification Service
│   ├── review/                 # Review Service
│   ├── chat/                   # Chat Service
│   ├── search/                 # Search Service
│   └── admin/                  # Admin Service
│
├── infrastructure/             # 🔧 インフラ設定
│   ├── docker/                 # Docker Compose
│   │   └── docker-compose.yml  # 全DB/Redis定義
│   └── k8s/                    # Kubernetes (Future)
│
├── tests/                      # 🧪 テストコード
│   ├── integration/            # 統合テスト
│   │   ├── auth/               # Auth統合テスト
│   │   ├── order_flow/         # Order-Payment-Inventory
│   │   └── notification_flow/  # Notification-Review-Shipping
│   ├── e2e/                    # E2Eテスト
│   └── run_all_integration_tests.sh
│
├── scripts/                    # 🛠️ 自動化スクリプト
│   ├── start_all_services.sh   # 全サービス起動
│   ├── stop_all_services.sh    # 全サービス停止
│   ├── check_all_services.sh   # ヘルスチェック
│   └── build_all_services.sh   # 全サービスビルド
│
├── mps-workspace/              # 🎨 MPS DSL定義
│   ├── languages/              # DSL言語定義
│   │   ├── microservice-dsl/
│   │   ├── grpc-dsl/
│   │   └── event-dsl/
│   └── solutions/              # サービス定義
│       ├── auth-service/
│       └── ...
│
├── docs/                       # 📚 ドキュメント
│   ├── requirements/           # 要件定義
│   └── reports/                # 実装レポート
│
├── .claude/                    # 🤖 Claude設定
│   ├── CLAUDE.md               # 開発ガイド
│   └── rules/                  # 開発ルール
│
├── Makefile                    # Make設定
├── .gitignore                  # Git除外設定
├── README.md                   # このファイル
└── QUICKSTART.md               # クイックスタート
```

---

## 🌐 URL一覧

### Customer（顧客）向けURL

| URL | 説明 |
|-----|------|
| `http://localhost:4000/` | ホームページ |
| `http://localhost:4000/auth` | ログイン・会員登録 |
| `http://localhost:4000/products` | 商品一覧 |
| `http://localhost:4000/products/:id` | 商品詳細 |
| `http://localhost:4000/shops` | ショップ一覧 |
| `http://localhost:4000/dashboard` | マイページ |

### Owner（ショップオーナー）向けURL

| URL | 説明 |
|-----|------|
| `http://localhost:4000/owner/auth` | オーナーログイン |
| `http://localhost:4000/owner/dashboard` | オーナーダッシュボード |
| `http://localhost:4000/owner/shop/register` | ショップ登録 |
| `http://localhost:4000/owner/products` | 商品管理 |
| `http://localhost:4000/owner/products/new` | 商品追加 |

### 開発者向けURL

| URL | 説明 |
|-----|------|
| `http://localhost:4000/dev/dashboard` | Phoenix LiveDashboard |

---

## 📊 プロジェクト統計

### コード統計

| カテゴリ | 行数 |
|---------|------|
| Goコード | 約22,000行 |
| Elixirコード | 約5,500行 |
| テストコード | 約5,000行 |
| DSL定義 | 約2,000行 |
| **合計** | **約34,500行** |

### 開発効率

| 指標 | 値 |
|------|-----|
| サービス数 | 12 |
| 開発期間 | 2週間（MPS DSL使用） |
| テストカバレッジ | 70% |
| トークン削減率 | 90%（MPS DSL効果） |

---

## 🤝 貢献

このプロジェクトは学習・研究目的のサンプルです。

改善提案やバグ報告は Issue または Pull Request でお願いします。

---

## 📝 ライセンス

MIT License

Copyright (c) 2025 makoto-developer

---

## 🔗 関連リンク

- [Go公式サイト](https://golang.org/)
- [Elixir公式サイト](https://elixir-lang.org/)
- [Phoenix Framework](https://www.phoenixframework.org/)
- [gRPC](https://grpc.io/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers)
- [JetBrains MPS](https://www.jetbrains.com/mps/)
- [Microservices Pattern](https://microservices.io/)

---

**🎉 Happy Coding!**
