# Go Microservice Example - プロジェクト概要

## プロジェクト目的
オンラインショップ（モール型）のマイクロサービス実践例
- 会員登録、カート、ショップで商品を選んで購入、配送まで
- ショップ登録、商品一覧、注文管理
- オンラインショップの管理者機能

## 開発手法
**JetBrains MPS DSL駆動開発**
- DSL定義（100-200行）→ MPS Generator → Goコード（2,000-3,000行）自動生成
- トークン消費90%削減（Claudeは生成コードを読まず、DSL定義のみ読む）

## マイクロサービス構成（12サービス）

### Phase 1: 基盤サービス（完了）
1. **Auth Service** - 認証・認可
2. **Shop Service** - ショップ管理

### Phase 2: コアサービス
3. **Customer Service** - 顧客管理
4. **Inventory Service** - 在庫管理（完全実装済み）
5. **Order Service** - 注文管理
6. **Payment Service** - 決済処理

### Phase 3: 拡張サービス1
7. **Shipping Service** - 配送管理
8. **Notification Service** - 通知送信
9. **Review Service** - レビュー管理

### Phase 4: 拡張サービス2
10. **Chat Service** - チャット機能
11. **Search Service** - 検索機能
12. **Admin Service** - 管理機能

## ディレクトリ構成

```
mps-workspace/          # MPS DSL定義（編集可能）
├── languages/          # DSL言語定義
└── solutions/          # サービス定義（12サービス）

generated/              # 生成コード（読み取り専用、触らない）
├── auth-service/
├── shop-service/
└── ... (12サービス)

manual/                 # カスタムロジック（編集可能）
└── custom/

proto/                  # Protocol Buffers定義（DSL生成）
scripts/                # ビルド・生成スクリプト
├── mps-generate.sh     # DSLからGoコード生成
└── migrations/         # データベースマイグレーション

mock/                   # モックサービス実装
tools/                  # ビルドツール
docs/                   # ドキュメント
```

## プロジェクトルール

### 最重要原則
1. **DSL First**: すべてはDSL定義から始める
2. **Read Only Generated**: 生成コードは読まない・触らない
3. **Token Efficiency**: 常にトークン消費を意識

### ワークフロー
1. 要件定義確認（`docs/requirements/`）
2. DSL定義作成（`mps-workspace/solutions/`）
3. コード生成（`./scripts/mps-generate.sh <service>`）
4. 生成コード確認（コンパイルのみ）
5. カスタムロジック実装（`manual/`、必要な場合のみ）
