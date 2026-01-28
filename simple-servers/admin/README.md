# Admin Service

管理サービス - システム管理者用の機能を提供

## 概要

- 管理者アカウント管理
- 監査ログ記録
- システム全体の管理機能

## データベース

- **PostgreSQL**: localhost:22021
- **Database**: admin_service
- **User**: postgres
- **Password**: postgres_password

### テーブル

#### admin_users
- 管理者ユーザー情報
- username, email, password_hash, role
- is_active, last_login_at

#### audit_logs
- 監査ログ
- admin_user_id, action, entity_type, entity_id
- changes (JSONB), ip_address, user_agent

## セットアップ

### 1. データベース起動

```bash
cd ../../infrastructure/docker
docker compose up -d postgres_admin
```

### 2. スキーマ作成

```bash
./setup_db.sh
```

### 3. サービス起動

```bash
go run main.go
```

## 環境変数

- `DATABASE_URL`: データベース接続URL (デフォルト: postgresql://postgres:postgres_password@localhost:22021/admin_service?sslmode=disable)
- `SERVICE_PORT`: gRPCサーバーポート (デフォルト: 22111)

## アーキテクチャ

Database per Service パターンを採用:
- Admin Service専用のPostgreSQLインスタンス (port 22021)
- 他サービスとのデータベース分離
- マイクロサービスアーキテクチャの原則に準拠
