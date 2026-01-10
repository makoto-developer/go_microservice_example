# MPS 開発ワークフロー

このドキュメントでは、JetBrains MPSを使用した開発フローを定義します。

---

## 基本フロー

```
要件定義 → DSL設計 → コード生成 → 確認 → カスタムロジック → テスト
```

---

## Phase 1: 要件定義の確認

### 必須アクション

1. **該当サービスの要件を読む**
   ```bash
   cat docs/requirements/XX_service_name.md
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

---

## Phase 2: DSL設計

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

---

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

---

## Phase 4: 生成コードの確認

### 確認項目

1. **ファイル構成の確認**
   ```bash
   ls -la generated/service-name/
   ```

2. **コンパイルチェック**
   ```bash
   cd generated/service-name
   go build ./...
   ```

3. **生成内容の確認**
   - エンティティが正しく生成されているか
   - ユースケースが期待通りか
   - gRPCハンドラーが適切か

### ⚠️ 注意

**生成コードを読む必要はありません**

- DSL定義が正しければ、生成コードも正しい
- コンパイルが通れば基本的にOK
- 詳細な確認はテスト実行で行う

---

## Phase 5: カスタムロジックの実装

### 実装が必要な場合

以下の場合のみ `manual/` にコード追加：

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

### 実装不要な場合

以下は生成コードで十分（manual/不要）：

- CRUD操作
- 単純なgRPCハンドラー
- データベースアクセス
- 標準的なエラーハンドリング

---

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

### テストカバレッジ

```bash
go test -cover -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 修正が必要な場合

### DSL定義の修正

```bash
1. mps-workspace/solutions/service-name/service.model を編集
2. ./scripts/mps-generate.sh service-name で再生成
3. 再度テスト
```

### ❌ やってはいけないこと

```bash
# 生成コードを直接編集
vim generated/service-name/domain/user.go  # NG!

# 理由: 再生成で上書きされる
```

---

## トラブルシューティング

### 生成エラー

**症状**: MPS Generator実行時にエラー

**対処**:
1. DSL定義の構文エラーを確認
2. MPS IDEでエラーメッセージを確認
3. Generator実装を見直し

### コンパイルエラー

**症状**: 生成コードがコンパイルできない

**対処**:
1. DSL定義を見直し（型の不一致等）
2. Generator実装を修正
3. 再生成

---

## まとめ

### MPS開発フローの鉄則

1. **要件 → DSL → 生成 → 確認 の順**
2. **生成コードは読まない・触らない**
3. **カスタムロジックは最小限**
4. **問題があればDSL定義を修正 → 再生成**

このフローを守ることで、効率的かつ安全に開発できます。
