# Claude開発ガイド - Go MicroService 実践例

このファイルは、Claudeがこのプロジェクトの開発を支援する際の指針を定義します。

## プロジェクト固有の開発方針

### 🎯 JetBrains MPS による DSL駆動開発

このプロジェクトでは、**JetBrains MPS（Meta Programming System）を使用したDSL駆動開発**を採用しています。

#### 開発アプローチ

1. **MPS DSLでサービスを定義**
   - マイクロサービスの仕様をDSLで記述
   - ビジネスロジック、エンティティ、ユースケースを宣言的に定義
   - 100-200行のDSL定義で1サービスを表現

2. **Goコードを自動生成**
   - MPS Generatorが2000-3000行のGoコードを生成
   - gRPCハンドラー、Repository、Usecase、Domain層を自動生成
   - Protocol Buffers定義も自動生成

3. **カスタムロジックのみ手動実装**
   - 複雑なビジネスロジックのみ手動で実装
   - 生成コードは触らない（再生成で上書き）

#### トークン削減効果

- **従来**: 2000-3000行のGoコードを読む → ~15,000トークン
- **MPS**: 100-200行のDSL定義を読む → ~1,000-1,500トークン
- **削減率**: **90%削減**

---

## Claude への指示

### 1. MPS DSL定義の支援

Claudeは以下を支援します：

#### DSL設計フェーズ
- マイクロサービスDSLの言語設計
- Structure、Editor、Generator、Typesystemの実装
- DSL制約の定義

#### サービス定義フェーズ
- 要件定義からDSL定義への変換
- DSL定義のレビュー・最適化
- Generator実装の支援

**例**: Auth Service の DSL定義
```
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

### 2. カスタムロジックの実装

生成されたコード（`microservices/`ディレクトリ）は読み取り専用です。
カスタムロジックは `manual/` ディレクトリに実装してください。

```
microservices/      # マイクロサービス（MPS生成コード含む）
├── auth/
│   ├── internal/   # 生成コード（触らない）
│   │   ├── domain/
│   │   ├── usecase/
│   │   └── handler/
│   ├── cmd/
│   ├── proto/
│   └── migrations/
manual/             # 手動実装（Claudeが編集）
├── auth/
│   └── custom_logic.go
```

### 3. ワークフロー

#### Phase 1: DSL設計（1回のみ）
1. `mps-workspace/languages/microservice-dsl/` でDSL言語を設計
2. Generator実装（DSL → Go変換ロジック）
3. Typesystem、制約を定義

#### Phase 2: サービス定義（各サービス）
1. 要件定義（`docs/requirements/`）を読む
2. DSL定義を作成（`mps-workspace/solutions/auth-service/`）
3. MPSでコード生成
4. 生成コードを確認
5. 必要なカスタムロジックを実装

#### Phase 3: テスト・統合
1. テストコード実装（生成 or 手動）
2. 統合テスト実行
3. 必要に応じてDSL定義を調整 → 再生成

---

## ディレクトリ構成（MPS使用時）

```
.
├── microservices/           # 全マイクロサービス（12サービス）
│   ├── auth/                # Auth Service
│   │   ├── cmd/             # メインアプリケーション
│   │   ├── internal/        # 生成コード（触らない）
│   │   │   ├── domain/
│   │   │   ├── usecase/
│   │   │   └── handler/
│   │   ├── proto/           # Protocol Buffers
│   │   ├── migrations/      # DBマイグレーション
│   │   ├── Dockerfile
│   │   └── go.mod
│   └── ... (全12サービス)
│
├── infrastructure/          # インフラ・デプロイ設定
│   ├── docker/              # Docker Compose設定
│   ├── config/              # 共通設定
│   └── databases/           # DBマイグレーション
│
├── docs/                    # ドキュメント
│   ├── README.md
│   ├── SETUP.md
│   ├── CLAUDE.md
│   ├── requirements/        # 要件定義
│   └── reports/             # レポート
│
├── build/                   # ビルド成果物
│   ├── bin/                 # バイナリ
│   └── proto/               # Protocol Buffers生成ファイル
│
├── tools/                   # 開発ツール
│   ├── scripts/
│   │   └── mps-generate.sh  # MPS生成スクリプト
│   ├── test-client/
│   └── mock/
│
├── mps-workspace/           # MPS専用ワークスペース
│   ├── languages/           # DSL定義
│   └── solutions/           # DSLを使ったサービス定義
│
├── manual/                  # カスタムロジック
│   └── custom/
│
└── web/                     # フロントエンド
```

---

## 開発時の注意事項

### ✅ すべきこと

1. **DSL定義を優先**
   - まずDSLで表現できないか検討
   - 定型的なコードはDSL + Generatorで自動化

2. **生成コードは読み取り専用**
   - `microservices/*/internal/` 配下は絶対に手動編集しない
   - 変更が必要ならDSL定義を修正 → 再生成

3. **要件変更はDSLから**
   - ビジネスロジック変更 → DSL定義を更新
   - API変更 → DSL定義を更新
   - DB変更 → DSL定義を更新

### ❌ してはいけないこと

1. **生成コードの手動編集**
   - 再生成で上書きされます

2. **生成コードを直接Git管理**
   - `.gitignore` で除外（オプション）

3. **DSLを経由しない実装**
   - 定型的なコードを手動で書くのはNG

---

## MPS関連コマンド

### コード生成
```bash
# すべてのサービスを生成
./tools/scripts/mps-generate.sh --all

# 特定のサービスのみ生成
./tools/scripts/mps-generate.sh auth-service
./tools/scripts/mps-generate.sh shop-service

# 生成コードの確認
ls -la generated/auth/
```

### MPS起動
```bash
# MPS IDE起動
open mps-workspace/

# または
/Applications/MPS.app/Contents/MacOS/mps mps-workspace/
```

---

## Claude との協働例

### 例1: 新しいサービスを追加

**ユーザー**: 「Payment Serviceを実装して」

**Claudeの作業**:
1. `docs/requirements/06_payment_service.md` を読む
2. DSL定義を作成（`mps-workspace/solutions/payment-service/`）
3. DSL定義をユーザーに提示 → 承認待ち
4. MPS生成スクリプトを実行
5. 生成コードを確認
6. 複雑なビジネスロジックがあれば `manual/payment/` に実装

### 例2: 既存サービスの機能追加

**ユーザー**: 「Auth Serviceにソーシャルログイン機能を追加」

**Claudeの作業**:
1. 既存DSL定義を読む（`mps-workspace/solutions/auth-service/`）
2. ソーシャルログイン用のusecase、entityを追加
3. DSL定義を更新
4. 再生成
5. OAuth連携などのカスタムロジックを `manual/auth/` に実装

### 例3: API仕様の変更

**ユーザー**: 「Register APIのレスポンスにroleを追加」

**Claudeの作業**:
1. DSL定義で `output` に `role: Role` を追加
2. 再生成
3. 完了（gRPC定義、ハンドラー、すべて自動更新）

---

## トラブルシューティング

### MPS生成エラー

**エラー**: Generator実行エラー

**対処**:
1. DSL定義の構文エラーを確認
2. MPS IDEでエラーメッセージを確認
3. Generator実装を見直し

### 生成コードのコンパイルエラー

**対処**:
1. DSL Generatorのテンプレートを修正
2. 再生成
3. （一時的に）手動で修正 → Generator修正後に削除

---

## まとめ

- **このプロジェクトはMPS DSL駆動開発です**
- **Claudeは主にDSL定義とカスタムロジック実装を支援します**
- **生成コードは触らない = トークン削減 + 開発効率化**

詳細な開発手順は [SETUP.md](./SETUP.md) を参照してください。
