# セットアップガイド

このドキュメントでは、ローカル開発環境のセットアップ方法と開発の始め方を説明します。

## 目次

- [前提条件](#前提条件)
- [依存関係のインストール](#依存関係のインストール)
- [Protocol Buffers のコンパイル](#protocol-buffers-のコンパイル)
- [開発環境の起動](#開発環境の起動)
- [各サービスの起動](#各サービスの起動)
- [テスト実行](#テスト実行)
- [ビルド](#ビルド)
- [トラブルシューティング](#トラブルシューティング)
- [開発Tips](#開発tips)

---

## 前提条件

以下のツールがインストールされている必要があります：

### 必須
- **Go**: 1.25 以上
  ```bash
  go version
  # go version go1.25.x darwin/amd64
  ```

- **Docker**: 最新版
  ```bash
  docker --version
  # Docker version 24.x.x
  ```

- **Docker Compose**: 最新版
  ```bash
  docker compose version
  # Docker Compose version v2.x.x
  ```

- **Protocol Buffers コンパイラ**: 3.20 以上
  ```bash
  protoc --version
  # libprotoc 3.20.x
  ```

### 推奨
- **Make**: ビルド自動化
- **Git**: バージョン管理
- **VSCode** または **GoLand**: 推奨IDE

### MPS（DSL駆動開発用）
- **JetBrains MPS**: 2023.2 以上
  ```bash
  # MPS のダウンロード
  # https://www.jetbrains.com/mps/download/

  # インストール確認
  /Applications/MPS.app/Contents/MacOS/mps --version
  ```

**注意**: このプロジェクトはMPS DSL駆動開発を採用しています。
詳細は [CLAUDE.md](./CLAUDE.md) を参照してください。

---

## 依存関係のインストール

### 1. Go モジュールのダウンロード

```bash
# プロジェクトルートで実行
go mod download
```

### 2. Protocol Buffers の Go プラグインをインストール

```bash
# protoc-gen-go のインストール
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# protoc-gen-go-grpc のインストール
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# PATH に追加（~/.zshrc または ~/.bashrc に追記）
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 3. 開発ツールのインストール

```bash
# Makefile を使う場合（推奨）
make install

# または手動で
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/cosmtrek/air@latest  # ホットリロード
```

---

## Protocol Buffers のコンパイル

Protocol Buffers 定義から Go コードを生成します。

```bash
# すべての .proto ファイルをコンパイル
make proto

# または手動で
./scripts/proto-gen.sh
```

生成されたファイルは `gen/go/` ディレクトリに配置されます。

---

## MPS によるコード生成（DSL駆動開発）

このプロジェクトでは、JetBrains MPSを使用してDSLからGoコードを生成します。

### MPS ワークスペースの初期化

```bash
# MPSプロジェクトを開く
open mps-workspace/

# または
/Applications/MPS.app/Contents/MacOS/mps mps-workspace/
```

### DSL定義の作成

1. MPS IDEで `mps-workspace/solutions/` を開く
2. 新しいSolutionを作成（例: `auth-service`）
3. DSLでサービスを定義

**例**: Auth Service の DSL定義
```kotlin
microservice AuthService {
  version: "v1"

  entity User {
    id: UUID primary_key
    email: string unique not_null
    password_hash: string not_null
    role: Role not_null
  }

  enum Role { CUSTOMER, SHOP_OWNER, ADMIN }

  usecase UserRegistration {
    input: { email: string, password: string, role: Role }
    output: { user_id: UUID, token: string }
    errors: { EmailAlreadyExists, WeakPassword }
  }

  grpc_service {
    rpc Register(RegisterRequest) returns (RegisterResponse)
    rpc Login(LoginRequest) returns (LoginResponse)
  }
}
```

### Goコードの生成

```bash
# 特定のサービスを生成
./scripts/mps-generate.sh auth-service

# すべてのサービスを生成
./scripts/mps-generate.sh --all

# 生成結果の確認
ls -la generated/auth/
```

**生成されるもの**:
- ドメイン層（`domain/`）
- ユースケース層（`usecase/`）
- インフラ層（`infrastructure/`）
- gRPCハンドラー（`handler/`）
- Protocol Buffers定義（`proto/`）
- テストコード（`tests/`）

### カスタムロジックの実装

生成されたコード（`generated/`）は触らず、カスタムロジックは `manual/` に実装します。

```bash
# カスタムロジック用のディレクトリ作成
mkdir -p manual/auth

# カスタムロジックを実装
vim manual/auth/custom_validation.go
```

**詳細**: [CLAUDE.md](./CLAUDE.md) の「Claude への指示」セクションを参照

---

## 開発環境の起動

### Docker Compose で一括起動（推奨）

```bash
# すべてのインフラサービスを起動
make dev-up

# または
docker compose -f deployments/docker/docker-compose.dev.yml up -d
```

起動されるサービス：
- PostgreSQL (複数DB)
- Redis
- RabbitMQ
- MinIO
- Elasticsearch
- Prometheus
- Grafana

### 起動確認

```bash
# コンテナの状態確認
docker compose ps

# ログ確認
docker compose logs -f
```

### 停止

```bash
make dev-down

# または
docker compose -f deployments/docker/docker-compose.dev.yml down
```

---

## 各サービスの起動

### 方法1: Makefile を使用（推奨）

```bash
# Auth Service を起動
make run-auth

# Shop Service を起動
make run-shop

# すべてのサービスを起動
make run-all
```

### 方法2: 直接起動

```bash
# Auth Service
go run cmd/auth/main.go

# Shop Service
go run cmd/shop/main.go

# 環境変数を指定して起動
ENV=development go run cmd/auth/main.go
```

### 方法3: ホットリロード（開発時）

```bash
# air を使用
cd cmd/auth
air

# または各サービスで
make watch-auth
```

---

## テスト実行

### すべてのテスト

```bash
# すべてのテストを実行
make test

# または
go test ./...
```

### カバレッジ付き

```bash
# カバレッジレポート生成
make test-coverage

# または
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 特定のサービスのテスト

```bash
# Auth Service のみ
go test ./internal/auth/...

# Shop Service のみ
go test ./internal/shop/...
```

### E2E テスト

```bash
# E2E テスト実行（Docker 環境が必要）
make test-e2e

# または
go test ./tests/e2e/...
```

### 統合テスト

```bash
# 統合テスト実行
make test-integration

# または
go test -tags=integration ./tests/integration/...
```

---

## ビルド

### すべてのサービスをビルド

```bash
make build
```

バイナリは `bin/` ディレクトリに生成されます。

### 特定のサービスをビルド

```bash
# Auth Service
make build-auth

# Shop Service
make build-shop
```

### Docker イメージのビルド

```bash
# すべてのサービス
make docker-build

# 特定のサービス
docker build -t auth-service -f deployments/docker/Dockerfile.auth .
```

---

## トラブルシューティング

### Protocol Buffers のコンパイルエラー

**エラー**: `protoc-gen-go: program not found or is not executable`

**解決策**:
```bash
# PATH を確認
echo $PATH | grep $(go env GOPATH)/bin

# PATH に追加されていない場合
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Docker コンテナが起動しない

**エラー**: `port is already allocated`

**解決策**:
```bash
# ポートを使用しているプロセスを確認
lsof -i :5432  # PostgreSQL
lsof -i :6379  # Redis

# 既存のコンテナを削除
docker compose down -v
```

### データベース接続エラー

**エラー**: `connection refused`

**解決策**:
```bash
# PostgreSQL が起動しているか確認
docker compose ps postgres

# ログを確認
docker compose logs postgres

# 再起動
docker compose restart postgres
```

### Go モジュールのエラー

**エラー**: `cannot find package`

**解決策**:
```bash
# モジュールを整理
go mod tidy

# vendor ディレクトリを使用する場合
go mod vendor
```

---

## 開発Tips

### 1. ホットリロード

開発時は `air` を使用すると、ファイル変更時に自動でリロードされます。

```bash
# air のインストール
go install github.com/cosmtrek/air@latest

# サービスディレクトリで実行
cd cmd/auth
air
```

### 2. Linter の使用

```bash
# すべてのコードを検証
make lint

# または
golangci-lint run ./...
```

### 3. デバッグ

VSCode の場合、`.vscode/launch.json` を作成：

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Auth Service",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/auth",
      "env": {
        "ENV": "development"
      }
    }
  ]
}
```

### 4. データベースマイグレーション

```bash
# マイグレーションの実行
make migrate-up

# ロールバック
make migrate-down

# マイグレーションの作成
make migrate-create NAME=add_users_table
```

### 5. ログレベルの変更

環境変数でログレベルを変更できます：

```bash
# デバッグレベル
LOG_LEVEL=debug go run cmd/auth/main.go

# 本番環境（info レベル）
LOG_LEVEL=info go run cmd/auth/main.go
```

### 6. gRPC クライアントのテスト

```bash
# grpcurl のインストール
brew install grpcurl

# サービスの確認
grpcurl -plaintext localhost:50051 list

# メソッドの実行
grpcurl -plaintext -d '{"email":"test@example.com","password":"password"}' \
  localhost:50051 auth.v1.AuthService/Login
```

---

## 次のステップ

- [要件定義を確認](./docs/requirements/README.md)
- [API仕様を確認](./proto/README.md)
- [アーキテクチャガイドを読む](./docs/architecture/README.md)
- [コントリビューションガイドを確認](./CONTRIBUTING.md)
