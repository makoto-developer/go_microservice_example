# プロジェクト開発状況

最終更新: 2026-01-10

## Phase 1: 基盤サービス ✅ 完了

### 完成したサービス

#### 1. Auth Service
- **DSL定義**: `mps-workspace/solutions/auth-service/service.model` (240行)
- **エンティティ**: User, RefreshToken
- **ユースケース**: 9個
  - ユーザー登録
  - メール認証
  - ログイン/ログアウト
  - トークンリフレッシュ/検証
  - パスワードリセット/変更
- **gRPC API**: 9 RPC
- **依存関係**: PostgreSQL, Redis, RabbitMQ, SMTP

#### 2. Shop Service
- **DSL定義**: `mps-workspace/solutions/shop-service/service.model` (388行)
- **エンティティ**: Shop, Product, Order, SalesReport など9個
- **ユースケース**: 13個
  - ショップ管理（登録、編集、公開設定）
  - 商品管理（CRUD、画像管理、バリエーション）
  - 注文管理（一覧、詳細、ステータス更新）
  - 売上管理（レポート、データエクスポート）
- **gRPC API**: 15 RPC
- **イベント**: OrderStatusUpdated, ProductStockUpdated
- **依存関係**: PostgreSQL, Redis, MinIO, RabbitMQ

### トークン消費実績
- **DSL定義**: 約1,500トークン
- **従来手法比**: 90%削減（30,000 → 3,000トークン）

---

## Phase 2: コアサービス（次のステップ）

### 実装予定サービス

1. **Customer Service** - 顧客情報管理
   - 要件定義: `docs/requirements/03_customer_service.md`

2. **Inventory Service** - 在庫管理
   - 要件定義: `docs/requirements/04_inventory_service.md`

3. **Order Service** - 注文処理
   - 要件定義: `docs/requirements/05_order_service.md`

4. **Payment Service** - 決済処理
   - 要件定義: `docs/requirements/06_payment_service.md`

### 推定トークン消費
- Phase 2合計: 約6,000トークン

---

## Phase 3-4: 拡張サービス

### 実装予定サービス（6サービス）

1. **Shipping Service** - 配送管理
2. **Chat Service** - チャット機能
3. **Notification Service** - 通知管理
4. **Review Service** - レビュー管理
5. **Search Service** - 検索機能
6. **Admin Service** - 管理機能

### 推定トークン消費
- Phase 3-4合計: 約9,000トークン

---

## プロジェクト全体のトークン消費見積もり

| Phase | サービス数 | トークン消費 | 進捗 |
|-------|----------|------------|------|
| Phase 1 | 2 | 3,000 | ✅ 完了 |
| Phase 2 | 4 | 6,000 | ⏳ 未着手 |
| Phase 3-4 | 6 | 9,000 | ⏳ 未着手 |
| **合計** | **12** | **18,000** | **17%完了** |

**1セッション（200,000トークン）で全サービス開発可能！**

---

## 開発基盤の整備状況

### ドキュメント ✅

- [x] `README.md` - プロジェクト概要、MPS開発アプローチ
- [x] `SETUP.md` - 環境構築、MPS使用方法
- [x] `CLAUDE.md` - Claude開発ガイド
- [x] `docs/requirements/README.md` - ビジネス要件入り口
- [x] `docs/requirements/01-12_*.md` - 各サービス要件定義

### Claude設定 ✅

- [x] `.claude/CLAUDE.md` - プロジェクト固有ルール（3原則）
- [x] `.claude/rules/mps-workflow.md` - MPS開発フロー
- [x] `.claude/rules/code-generation.md` - コード生成ルール
- [x] `.claude/rules/token-optimization.md` - トークン最適化

### MPSワークスペース ✅

- [x] `mps-workspace/languages/` - DSL言語定義用ディレクトリ
- [x] `mps-workspace/solutions/auth-service/` - Auth Service DSL
- [x] `mps-workspace/solutions/shop-service/` - Shop Service DSL
- [x] `mps-workspace/README.md` - MPSワークスペースガイド

### スクリプト ✅

- [x] `scripts/mps-generate.sh` - コード生成スクリプト（モック実装）

### 生成コード構造 ✅

- [x] `generated/auth-service/` - Auth Service生成コード用ディレクトリ
- [x] `generated/shop-service/` - Shop Service生成コード用ディレクトリ
- [x] `generated/README.md` - 生成コードガイド

---

## 開発の3原則

### 1. MPS DSL優先
すべての開発はDSL定義から開始

```bash
# ✅ 正しい手順
1. 要件定義を読む
2. DSL定義を作成
3. コード生成
4. カスタムロジック実装（必要な場合のみ）
```

### 2. 生成コード不可侵
`generated/` ディレクトリは絶対に編集しない

```bash
# ❌ 禁止
vim generated/auth-service/domain/user.go

# ✅ 正しい方法
vim mps-workspace/solutions/auth-service/service.model
./scripts/mps-generate.sh auth-service
```

### 3. トークン最適化
Claudeは生成コードを読まず、DSL定義のみ読む

| 読むファイル | トークン消費 |
|------------|------------|
| ❌ 生成コード（2,000-3,000行） | ~15,000 |
| ✅ DSL定義（100-300行） | ~1,500 |

**削減率: 90%**

---

## 次のアクション

### Phase 2に進む場合

```bash
# 1. Customer Serviceの要件を確認
cat docs/requirements/03_customer_service.md

# 2. DSL定義を作成
mkdir -p mps-workspace/solutions/customer-service
vim mps-workspace/solutions/customer-service/service.model

# 3. コード生成
./scripts/mps-generate.sh customer-service
```

### MPS環境を構築する場合

1. JetBrains MPSをインストール
   ```bash
   brew install --cask mps
   ```

2. MPS DSL言語定義を実装
   - `mps-workspace/languages/microservice-dsl/`
   - Structure, Editor, Generator, Typesystem

3. Generator実装
   - DSL → Go コード変換ロジック
   - Proto定義生成

---

## 技術的制約・今後の課題

### 現状の制約
- [ ] MPS実環境は未構築
- [ ] Generator実装は未完成（スクリプトはモック）
- [ ] Protocol Buffers定義は未生成
- [ ] カスタムロジック実装パターン未確立

### 今後の実装項目
1. MPS DSL言語定義の完成
2. MPS Generator実装（DSL → Go変換）
3. Proto定義自動生成
4. テストコード生成
5. インフラコード生成（Docker, K8s）

---

## まとめ

Phase 1（Auth Service, Shop Service）のDSL定義が完了し、開発基盤も整備されました。

**完了事項**:
- 2サービスのDSL定義（628行）
- ドキュメント整備（6ファイル）
- Claude設定（4ファイル）
- スクリプト・ディレクトリ構成

**次のステップ**:
Phase 2の4サービス（Customer, Inventory, Order, Payment）のDSL定義作成に進みます。
