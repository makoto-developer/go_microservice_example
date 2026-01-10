# 推奨コマンド

## MPS DSL関連

### コード生成
```bash
# 全サービス生成
./scripts/mps-generate.sh --all

# 特定サービス生成
./scripts/mps-generate.sh auth-service
./scripts/mps-generate.sh shop-service
```

### MPS IDE起動
```bash
# MPS プロジェクトを開く
open mps-workspace/

# または
/Applications/MPS.app/Contents/MacOS/mps mps-workspace/
```

## Docker関連

### 起動・停止
```bash
make up              # 全サービス起動
make up-infra        # インフラのみ起動
make down            # 全サービス停止
make restart         # 全サービス再起動
make dev             # 開発環境起動（インフラ+モック）
```

### ログ確認
```bash
make logs            # 全ログ表示
make logs-infra      # インフラログのみ
make ps              # 稼働状況確認
make health          # ヘルスチェック
```

### データベース
```bash
make db-init         # データベース初期化
make db-connect      # PostgreSQL接続
make db-list         # データベース一覧
make db-reset        # データベースリセット（要注意）
make backup-db       # バックアップ作成
```

### ビルド
```bash
make build           # 全サービスビルド
make build-infra     # インフラビルド
make build-mocks     # モックサービスビルド
make build-services  # マイクロサービスビルド
```

### ユーティリティ
```bash
make open-mailhog    # MailHog UI をブラウザで開く
make open-rabbitmq   # RabbitMQ UI をブラウザで開く
make shell-postgres  # PostgreSQLシェル
make shell-redis     # Redisシェル
make go-tidy         # 全サービスで go mod tidy 実行
```

## Go関連

### ビルド・実行
```bash
# 特定サービス実行
cd generated/auth-service && go run .
cd generated/shop-service && go run .

# テスト実行
go test ./...
go test -cover ./...

# 依存関係整理
go mod tidy
go mod download
```

## Git関連

### ブランチ操作
```bash
git status
git add .
git commit -m "message"
git push origin <branch>

# ブランチ作成
git checkout -b feature/new-feature

# ブランチ切り替え
git checkout main
```

## macOS固有コマンド

### ファイル操作
```bash
# ファイル検索
find . -name "*.go" -type f

# ファイル内検索
grep -r "pattern" .

# ディレクトリリスト
ls -la

# ポート確認
lsof -i :20000  # PostgreSQL
lsof -i :20100  # Auth Service
```

## 初回セットアップ

```bash
# 1. 環境変数ファイル作成
cp .env.example .env

# 2. DSLからコード生成
./scripts/mps-generate.sh --all

# 3. Docker環境起動
make up

# 4. データベース初期化
make db-init

# 5. サービス起動確認
make ps
make logs
```
