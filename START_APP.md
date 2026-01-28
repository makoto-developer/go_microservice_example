# アプリケーション起動ガイド

このドキュメントでは、Shop Mall Webアプリケーションを起動する方法を説明します。

---

## 📋 前提条件

以下がインストールされていることを確認してください：

- **Elixir**: 1.18.4+
- **Erlang**: 28.3.1+
- **Node.js**: 18+
- **Go**: 1.25+ （バックエンドサービス用）
- **PostgreSQL**: 17+ （任意: Docker使用も可）
- **Redis**: 7+ （任意: Docker使用も可）

---

## 🚀 起動方法

### 1. Phoenixアプリケーション（フロントエンド）のみ起動

最もシンプルな方法です。バックエンドサービスなしでUIを確認できます。

```bash
# プロジェクトルートから
cd web/shop_mall_web

# 依存関係のインストール（初回のみ）
mix setup

# Phoenixサーバー起動
mix phx.server

# または、IExで起動（デバッグ時に便利）
iex -S mix phx.server
```

**アクセス**: http://localhost:4000

**ポート変更**（ポート22200で起動する場合）:
```bash
PORT=22200 mix phx.server
```

**アクセス**: http://localhost:22200

---

### 2. Phoenixアプリ + バックエンドサービス起動

完全な機能を利用するには、バックエンドのGoマイクロサービスも起動します。

#### 手順A: 個別にサービスを起動

**ターミナル1: Phoenixアプリ**
```bash
cd web/shop_mall_web
PORT=22200 mix phx.server
```

**ターミナル2: Auth Service**
```bash
cd microservices/auth
go run cmd/main.go
```

**ターミナル3: Shop Service**
```bash
cd microservices/shop
go run cmd/main.go
```

**ターミナル4: その他のサービス（必要に応じて）**
```bash
# Customer Service
cd microservices/customer && go run cmd/main.go &

# Inventory Service
cd microservices/inventory && go run cmd/main.go &

# Order Service
cd microservices/order && go run cmd/main.go &

# ... 他のサービス
```

---

#### 手順B: Docker Compose使用（推奨）

**注意**: Docker Composeファイルが未設定の場合は、手順Aを使用してください。

```bash
# インフラサービスのみ起動（PostgreSQL、Redis等）
docker-compose up -d

# Phoenixアプリ起動
cd web/shop_mall_web
PORT=22200 mix phx.server
```

---

## 🔧 E2Eテスト実行

E2Eテストを実行する場合：

```bash
# Phoenixサーバーが起動していることを確認
# 別ターミナルで:

cd web/shop_mall_web

# 全テスト実行（Headless）
npx playwright test

# 特定のテストを実行
npx playwright test e2e/owner_flow.spec.js
npx playwright test e2e/customer_flow.spec.js

# Headed モード（ブラウザを表示して実行）
npx playwright test --headed

# UIモード（インタラクティブ）
npx playwright test --ui

# レポート表示
npx playwright show-report
```

---

## 📊 各サービスのポート

| サービス | ポート | 用途 |
|---------|-------|------|
| Phoenix LiveView | 4000 または 22200 | Webアプリ |
| Auth Service | 50051 | 認証gRPC |
| Shop Service | 50052 | ショップgRPC |
| Customer Service | 50053 | 顧客gRPC |
| Inventory Service | 50054 | 在庫gRPC |
| Order Service | 50055 | 注文gRPC |
| Payment Service | 50056 | 決済gRPC |
| PostgreSQL | 5432 | データベース |
| Redis | 6379 | キャッシュ |

**注**: 実際のポート番号は各サービスの設定ファイルで確認してください。

---

## 🐛 トラブルシューティング

### 1. ポートが既に使用されている

```bash
# プロセスを確認
lsof -i :4000
lsof -i :22200

# プロセスを終了
kill -9 <PID>
```

### 2. 依存関係のエラー

```bash
# Elixir依存関係を再インストール
cd web/shop_mall_web
mix deps.clean --all
mix setup

# Node.js依存関係を再インストール
npm install
```

### 3. コンパイルエラー

```bash
# クリーンビルド
mix clean
mix compile
```

### 4. データベース接続エラー

```bash
# データベースを作成
mix ecto.create

# マイグレーション実行
mix ecto.migrate
```

### 5. gRPCサービス接続エラー

**症状**: 商品一覧等が表示されない

**原因**: バックエンドサービスが起動していない

**対処**: Auth Service、Shop Service等を起動してください（上記「手順A」参照）

---

## 📝 開発時のワークフロー

### 通常の開発

```bash
# Phoenixアプリのみ起動して開発
cd web/shop_mall_web
PORT=22200 mix phx.server

# ファイル変更時に自動リロード
# → LiveReloadが有効なので自動的に反映されます
```

### バックエンド含む開発

```bash
# ターミナルを分割して各サービスを起動
# または tmux/screen を使用

# Phoenix
PORT=22200 mix phx.server

# Auth Service
cd microservices/auth && go run cmd/main.go

# Shop Service
cd microservices/shop && go run cmd/main.go
```

---

## 🎯 現在の実装状況

### 実装済み
- ✅ 認証画面（オーナー/顧客）
- ✅ ダッシュボード（オーナー/顧客）
- ✅ ショップ登録・管理
- ✅ 商品登録・管理
- ✅ 商品一覧・詳細閲覧
- ✅ ショップ一覧・詳細
- ✅ パスワードリセット
- ✅ レスポンシブデザイン

### 未実装（バックエンドgRPCサービス必要）
- ⚠️ カート機能
- ⚠️ 購入フロー
- ⚠️ 注文履歴
- ⚠️ レビュー投稿
- ⚠️ 管理者機能

---

## 🔗 関連ドキュメント

- [E2Eテストレポート](web/shop_mall_web/e2e/TEST_REPORT.md)
- [オーナーシナリオ](web/shop_mall_web/e2e/scenarios/OWNER_SCENARIO.md)
- [顧客シナリオ](web/shop_mall_web/e2e/scenarios/CUSTOMER_SCENARIO.md)
- [Phoenix LiveView公式ドキュメント](https://hexdocs.pm/phoenix_live_view/)

---

## ✅ 起動確認

アプリが正常に起動したら、以下のURLにアクセスして確認してください：

- **トップページ**: http://localhost:22200/
- **顧客認証**: http://localhost:22200/auth
- **顧客ダッシュボード**: http://localhost:22200/dashboard
- **商品一覧**: http://localhost:22200/products
- **ショップ一覧**: http://localhost:22200/shops
- **オーナー認証**: http://localhost:22200/owner/auth
- **オーナーダッシュボード**: http://localhost:22200/owner/dashboard
- **商品管理**: http://localhost:22200/owner/products

---

**最終更新**: 2026-01-28
