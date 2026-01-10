# MPS DSL駆動開発ワークフロー

## 基本フロー

```
要件定義 → DSL定義 → コード生成 → 確認 → カスタムロジック → テスト
```

## Phase 1: 要件定義の確認

### 必須アクション
1. **該当サービスの要件を読む**
   ```bash
   cat docs/requirements/XX_service_name.md
   
   # または Serena で読む
   mcp__serena__read_file(
     relative_path="docs/requirements/XX_service_name.md",
     max_answer_chars=10000
   )
   ```

2. **要件から抽出すべき情報**
   - サービス概要
   - 機能要件
   - エンティティ（データ構造）
   - ユースケース（ビジネスロジック）
   - gRPC API（外部インターフェース）

3. **要件の確認**
   - 不明点があればユーザーに質問
   - 実装範囲の確認
   - 優先順位の確認

## Phase 2: DSL定義

### DSL定義の構造
```kotlin
microservice ServiceName {
  version: "v1"

  // 1. エンティティ定義
  entity EntityName {
    id: UUID primary_key
    field1: string not_null
    field2: int
    created_at: timestamp
  }

  // 2. Enum定義
  enum StatusType {
    ACTIVE, INACTIVE
  }

  // 3. ユースケース定義
  usecase UsecaseName {
    input: { field1: string, field2: int }
    output: { result_id: UUID, status: StatusType }
    errors: { ErrorType1, ErrorType2 }
  }

  // 4. gRPCサービス定義
  grpc_service {
    rpc MethodName(Request) returns (Response)
  }

  // 5. 依存関係
  dependencies {
    database: PostgreSQL
    cache: Redis
    messaging: RabbitMQ
  }

  // 6. イベント定義（オプション）
  events {
    publish EventName { field1: string }
    subscribe OtherEvent { field2: int }
  }
}
```

### DSL定義のベストプラクティス

#### ✅ 良い例
```kotlin
// 明確な命名
usecase UserRegistration {
  input: { email: string, password: string }
  output: { user_id: UUID, token: string }
  errors: { EmailAlreadyExists, WeakPassword }
}

// 適切な型指定
entity User {
  id: UUID primary_key
  email: string unique not_null
  created_at: timestamp not_null
}
```

#### ❌ 悪い例
```kotlin
// 曖昧な命名
usecase DoStuff {
  input: { data: string }
  output: { result: string }
}

// 型が不明確
entity Thing {
  id: string
  data: string
}
```

### Serenaを使ったDSL定義確認
```python
# DSL定義の概要取得
mcp__serena__get_symbols_overview(
  relative_path="mps-workspace/solutions/auth-service/service.model"
)

# 特定シンボルの詳細確認
mcp__serena__find_symbol(
  name_path_pattern="User",
  relative_path="mps-workspace/solutions/auth-service",
  include_body=True
)
```

## Phase 3: コード生成

### 生成コマンド
```bash
# 特定のサービス
./scripts/mps-generate.sh service-name

# すべてのサービス
./scripts/mps-generate.sh --all
```

### 生成されるファイル構成
```
generated/service-name/
├── domain/
│   ├── entity.go           # エンティティ
│   ├── enum.go             # Enum
│   └── repository.go       # Repositoryインターフェース
├── usecase/
│   ├── usecase_impl.go     # ユースケース実装
│   └── interfaces.go       # 依存インターフェース
├── handler/
│   └── grpc_handler.go     # gRPCハンドラー
├── infrastructure/
│   ├── postgres_repo.go    # Repository実装
│   ├── redis_cache.go      # Redisキャッシュ
│   └── rabbitmq_pub.go     # イベント発行
└── tests/
    └── usecase_test.go     # テストコード
```

## Phase 4: 生成コードの確認

### 確認項目

1. **ファイル構成の確認（Serena使用）**
   ```python
   mcp__serena__list_dir(
     relative_path="generated/service-name",
     recursive=True,
     skip_ignored_files=True
   )
   ```

2. **コンパイルチェック**
   ```bash
   cd generated/service-name
   go build ./...
   ```

3. **生成内容の確認（Serena使用）**
   ```python
   # シンボル概要のみ確認（トークン節約）
   mcp__serena__get_symbols_overview(
     relative_path="generated/service-name/domain/user.go"
   )
   ```

### ⚠️ 注意
**生成コードを全部読む必要はありません**
- DSL定義が正しければ、生成コードも正しい
- コンパイルが通れば基本的にOK
- 詳細な確認はテスト実行で行う

## Phase 5: カスタムロジックの実装

### 実装が必要な場合（manual/に追加）

1. **複雑なバリデーション**
   ```go
   // manual/auth/validation.go
   func ValidatePassword(password string) error {
       // 複雑なパスワードルール
   }
   ```

2. **外部API連携**
   ```go
   // manual/payment/stripe.go
   func ProcessStripePayment(amount int) error {
       // Stripe API呼び出し
   }
   ```

3. **ビジネスロジック**
   ```go
   // manual/order/saga.go
   func ExecuteOrderSaga(order Order) error {
       // 分散トランザクション処理
   }
   ```

### Serenaを使ったコード編集
```python
# ファイル作成
mcp__serena__create_text_file(
  relative_path="manual/auth/validation.go",
  content="package custom\n\n..."
)

# 既存ファイル編集（正規表現）
mcp__serena__replace_content(
  relative_path="manual/auth/validation.go",
  needle="old_pattern",
  repl="new_pattern",
  mode="regex"
)
```

## Phase 6: テスト

### テスト実行
```bash
# ユニットテスト
go test ./generated/service-name/...

# カスタムロジックのテスト
go test ./manual/service-name/...

# 統合テスト
go test -tags=integration ./tests/integration/...
```

## 修正が必要な場合

### DSL定義の修正（Serena使用）
```bash
1. DSL定義を編集
2. 再生成
3. 再度テスト
```

```python
# Serenaでの編集例
mcp__serena__replace_content(
  relative_path="mps-workspace/solutions/service-name/service.model",
  needle="old_definition",
  repl="new_definition",
  mode="regex"
)
```

### ❌ やってはいけないこと
```bash
# 生成コードを直接編集
vim generated/service-name/domain/user.go  # NG!

# 理由: 再生成で上書きされる
```

## トークン最適化のポイント

### ✅ トークン節約
1. **DSL定義のみ読む**（生成コードは読まない）
2. **Serenaのシンボル概要を活用**（`get_symbols_overview`）
3. **必要な部分のみ読む**（`find_symbol`）

### ❌ トークン無駄遣い
1. **生成コードを全部読む**（2,000-3,000行）
2. **同じファイルを複数回読む**
3. **不要なドキュメントを読む**

## まとめ

### MPS開発フローの鉄則
1. **要件 → DSL → 生成 → 確認 の順**
2. **生成コードは読まない・触らない**
3. **カスタムロジックは最小限**
4. **問題があればDSL定義を修正 → 再生成**
5. **Serenaを活用してトークン節約**
