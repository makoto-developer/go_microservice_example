# 技術スタック

## 開発手法
- **JetBrains MPS 2023.2+** - DSL駆動開発

## バックエンド（マイクロサービス）
- **Go 1.25**
- **PostgreSQL 17-alpine** - データベース
- **RabbitMQ 4-management-alpine** - メッセージブローカー
- **Redis 7.4-alpine** - キャッシュ、セッション管理

## サービス間通信
- **gRPC** - サービス間通信プロトコル
- **Protocol Buffers** - シリアライゼーション

## 検索・ストレージ
- **Elasticsearch** (Mock) - 全文検索エンジン
- **MinIO / S3** - オブジェクトストレージ

## 認証・セキュリティ
- **JWT-go** - JWT認証ライブラリ
- **bcrypt** - パスワードハッシュ化

## 外部サービス連携
- **Stripe Go SDK** - 決済プロバイダー連携

## インフラ・コンテナ
- **Docker** - コンテナ化
- **Docker Compose** - コンテナオーケストレーション

## モニタリング・テスト
- **MailHog** - メールテストツール（開発環境）
- **testify** - テストフレームワーク

## 開発環境
- **macOS (Darwin)** - 開発プラットフォーム
- **VSCode / GoLand** - 推奨IDE
