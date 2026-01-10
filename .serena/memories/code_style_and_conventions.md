# コードスタイルと規約

## Go言語規約

### 命名規則
- **Package名**: 小文字、単語区切りなし（`domain`, `usecase`, `handler`）
- **Type名**: PascalCase（`User`, `Order`, `PaymentService`）
- **Function名**: PascalCase（公開）、camelCase（非公開）
- **Variable名**: camelCase（`userID`, `orderTotal`）
- **Constant名**: PascalCase または UPPER_SNAKE_CASE

### ファイル構成
```go
// パッケージ宣言
package domain

// import（標準ライブラリ → サードパーティ → プロジェクト内）
import (
    "context"
    "time"

    "github.com/google/uuid"

    "github.com/makoto-developer/go_microservice_example/pkg/errors"
)

// 型定義
type User struct {
    ID    uuid.UUID
    Email string
}

// メソッド
func (u *User) Validate() error {
    // ...
}
```

### エラーハンドリング
```go
// エラーは明示的に処理
result, err := someFunction()
if err != nil {
    return nil, fmt.Errorf("failed to do something: %w", err)
}

// エラーは wrap して返す
return nil, fmt.Errorf("context: %w", err)
```

## DSL定義規約

### 命名規則
- **microservice名**: PascalCase（`AuthService`, `OrderService`）
- **entity名**: PascalCase（`User`, `Order`）
- **field名**: snake_case（`user_id`, `created_at`）
- **usecase名**: PascalCase（`UserRegistration`, `CreateOrder`）
- **enum名**: PascalCase（`Role`, `OrderStatus`）
- **enum値**: UPPER_CASE（`CUSTOMER`, `ADMIN`）

### DSL構造
```kotlin
microservice ServiceName {
  version: "v1"

  // 1. エンティティ定義
  entity EntityName {
    id: UUID primary_key
    field: type constraint
  }

  // 2. Enum定義
  enum EnumName {
    VALUE1, VALUE2
  }

  // 3. ユースケース定義
  usecase UsecaseName {
    input: { field: type }
    output: { field: type }
    errors: { ErrorType }
  }

  // 4. gRPCサービス定義
  grpc_service {
    rpc MethodName(Request) returns (Response)
  }
}
```

### DSL定義の制約
- 1サービス = 100-300行以内（推奨）
- 400行超過 = 警告
- 明確な命名（曖昧な名前は避ける）
- コメントは最小限（DSL自体が自己文書化）

## ディレクトリ構造規約

### 生成コードの構造
```
generated/<service>/
├── domain/          # ドメインエンティティ、Enum
├── usecase/         # ユースケース実装
├── handler/         # gRPCハンドラー
├── infrastructure/  # Repository実装、RabbitMQ Publisher
├── tests/           # テストコード
├── main.go          # エントリーポイント
└── go.mod           # Go モジュール定義
```

### カスタムロジックの配置
```
manual/<service>/
└── custom_logic.go  # 複雑なビジネスロジック、外部API連携
```

## Git コミットメッセージ

### フォーマット
```
<type>: <subject>

<body>

🤖 Generated with Claude Code
Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

### Type
- `feat`: 新機能
- `fix`: バグ修正
- `refactor`: リファクタリング
- `docs`: ドキュメント
- `test`: テスト
- `chore`: その他（ビルド、設定）

### 例
```
feat: Auth ServiceのDSL定義を追加

- ユーザー登録、ログイン、JWT発行のユースケースを定義
- DSLから2,000行のGoコードを生成

🤖 Generated with Claude Code
Co-Authored-By: Claude Sonnet 4.5 <noreply@anthropic.com>
```

## デザインパターン

### Repository パターン
```go
type UserRepository interface {
    Create(ctx context.Context, user *domain.User) error
    FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
    Update(ctx context.Context, user *domain.User) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### Dependency Injection
```go
// main.go
userRepo := infrastructure.NewPostgresUserRepository(db)
userUsecase := usecase.NewUserRegistration(userRepo)
handler := handler.NewAuthServiceHandler(userUsecase)
```

### Event-Driven Architecture
```go
// イベント発行
eventPublisher.PublishOrderCreated(ctx, OrderCreatedEvent{
    OrderID: order.ID,
    UserID:  order.UserID,
})

// イベント購読（別サービス）
subscriber.SubscribeOrderCreated(func(event OrderCreatedEvent) {
    // 処理
})
```
