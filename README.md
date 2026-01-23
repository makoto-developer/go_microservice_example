# Go MicroService 実践例

オンラインショップ（モール型）を題材にした Go マイクロサービスの実装例です。

## 📖 ドキュメント

- **[プロジェクト概要・セットアップ](./docs/README.md)** - 詳細な説明と始め方
- **[開発ガイド (Claude用)](./docs/CLAUDE.md)** - MPS DSL駆動開発の手順
- **[セットアップガイド](./docs/SETUP.md)** - 環境構築の詳細
- **[プロジェクト状況](./docs/PROJECT_STATUS.md)** - 実装進捗

## 📁 ディレクトリ構造

```
go_microservice_example/
│
├── microservices/         # 🎯 全マイクロサービス（12サービス）
│   ├── auth/              # 認証・認可
│   ├── shop/              # ショップ管理
│   ├── customer/          # 顧客管理
│   ├── inventory/         # 在庫管理
│   ├── order/             # 注文管理
│   ├── payment/           # 決済処理
│   ├── shipping/          # 配送管理
│   ├── notification/      # 通知送信
│   ├── review/            # レビュー管理
│   ├── chat/              # チャット機能
│   ├── search/            # 検索機能
│   └── admin/             # 管理機能
│
├── infrastructure/        # 🔧 インフラ・デプロイ設定
│   ├── docker/            # Docker Compose設定
│   ├── config/            # 共通設定
│   └── databases/         # DBマイグレーション
│
├── docs/                  # 📚 ドキュメント
│   ├── requirements/      # 要件定義
│   └── reports/           # 実装レポート
│
├── build/                 # 🏗️ ビルド成果物
│   ├── bin/               # バイナリ
│   └── proto/             # Protocol Buffers生成ファイル
│
├── tools/                 # 🛠️ 開発ツール
│   ├── scripts/           # ビルドスクリプト
│   ├── test-client/       # テストクライアント
│   └── mock/              # モックサービス
│
├── mps-workspace/         # 🎨 MPS DSL定義
│   ├── languages/         # DSL言語定義
│   └── solutions/         # サービス定義
│
└── .claude/               # 🤖 Claude設定
    └── rules/             # 開発ルール
```

## 🚀 Quick Start

```bash
# 1. Docker インフラ起動
cd infrastructure/docker
docker-compose up -d

# 2. Auth Service 起動
cd ../../microservices/auth
go mod tidy
go build -o auth-server ./cmd/server
./auth-server
```

詳細は [docs/README.md](./docs/README.md) を参照してください。

## 🎯 開発手法

このプロジェクトは **JetBrains MPS による DSL駆動開発**を採用しています。

- DSL定義（100-200行）→ Goコード（2,000-3,000行）自動生成
- トークン消費90%削減
- 詳細: [docs/CLAUDE.md](./docs/CLAUDE.md)

## 🏗️ アーキテクチャ

- **パターン**: マイクロサービスアーキテクチャ
- **通信**: gRPC + Protocol Buffers
- **データベース**: PostgreSQL (各サービス独立DB)
- **キャッシュ**: Redis
- **メッセージング**: RabbitMQ
- **検索**: Elasticsearch

## 📝 ライセンス

MIT License
