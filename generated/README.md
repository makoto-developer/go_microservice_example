# Generated Code Directory

このディレクトリには、MPS DSL定義から自動生成されたGoコードが配置されます。

## ⚠️ 重要な注意事項

### このディレクトリ配下のファイルは**絶対に編集しない**

- **理由**: MPS Generatorによる再生成で上書きされます
- **変更が必要な場合**: DSL定義を修正 → 再生成

```bash
# ❌ 禁止
vim generated/auth-service/domain/user.go

# ✅ 正しい方法
vim mps-workspace/solutions/auth-service/service.model
./scripts/mps-generate.sh auth-service
```

## ディレクトリ構成

各サービスは以下の構成で生成されます：

```
generated/
└── service-name/
    ├── domain/              # ドメイン層
    │   ├── entity.go        # エンティティ
    │   ├── enum.go          # Enum定義
    │   ├── repository.go    # Repositoryインターフェース
    │   └── value_object.go  # 値オブジェクト
    ├── usecase/             # ユースケース層
    │   ├── usecase_impl.go  # ユースケース実装
    │   └── interfaces.go    # 依存インターフェース
    ├── handler/             # ハンドラー層
    │   └── grpc_handler.go  # gRPCハンドラー
    ├── infrastructure/      # インフラ層
    │   ├── postgres_repo.go # Repository実装
    │   ├── redis_cache.go   # キャッシュ実装
    │   └── rabbitmq_pub.go  # メッセージング実装
    ├── tests/               # テストコード
    │   └── usecase_test.go  # ユニットテスト
    └── go.mod               # Go module定義
```

## 生成されるコードの例

### Entity (domain/user.go)

```go
package domain

import (
    "time"
    "github.com/google/uuid"
)

type User struct {
    ID           uuid.UUID `db:"id" json:"id"`
    Email        string    `db:"email" json:"email"`
    PasswordHash string    `db:"password_hash" json:"-"`
    Role         Role      `db:"role" json:"role"`
    CreatedAt    time.Time `db:"created_at" json:"created_at"`
}
```

### Enum (domain/role.go)

```go
package domain

type Role string

const (
    RoleCustomer   Role = "CUSTOMER"
    RoleShopOwner  Role = "SHOP_OWNER"
    RoleAdmin      Role = "ADMIN"
)
```

### Usecase (usecase/user_registration.go)

```go
package usecase

import (
    "context"
    "github.com/google/uuid"
)

type UserRegistrationInput struct {
    Email    string
    Password string
    Role     Role
}

type UserRegistrationOutput struct {
    UserID  uuid.UUID
    Message string
}

type UserRegistrationUsecase interface {
    Execute(ctx context.Context, input UserRegistrationInput) (UserRegistrationOutput, error)
}

type userRegistrationUsecaseImpl struct {
    userRepo UserRepository
    // ... dependencies
}

func (u *userRegistrationUsecaseImpl) Execute(
    ctx context.Context,
    input UserRegistrationInput,
) (UserRegistrationOutput, error) {
    // Generated business logic
    // - Validation
    // - Repository calls
    // - Error handling
}
```

### gRPC Handler (handler/grpc_handler.go)

```go
package handler

import (
    "context"
    pb "github.com/.../proto/auth/v1"
)

type AuthServiceHandler struct {
    pb.UnimplementedAuthServiceServer
    registerUsecase usecase.UserRegistrationUsecase
}

func (h *AuthServiceHandler) Register(
    ctx context.Context,
    req *pb.RegisterRequest,
) (*pb.RegisterResponse, error) {
    // Generated handler logic
    // - Request to usecase input conversion
    // - Usecase execution
    // - Response conversion
}
```

## コード生成の実行

```bash
# 特定のサービスのみ生成
./scripts/mps-generate.sh auth-service

# すべてのサービスを生成
./scripts/mps-generate.sh --all
```

## カスタムロジックの実装

生成コードで表現できない複雑なロジックは `manual/` ディレクトリに実装します。

```
manual/
└── auth/
    ├── password_validator.go  # カスタムバリデーション
    └── email_sender.go         # 外部API連携
```

## Git管理について

### オプション1: コミットしない（推奨）

`.gitignore` に追加:
```gitignore
generated/
```

**メリット**:
- Diffが小さい
- コンフリクトが起きにくい

**デメリット**:
- CI/CDで再生成が必要

### オプション2: コミットする

**メリット**:
- CI/CDで再生成不要
- コード確認が容易

**デメリット**:
- Diffが大きい
- コンフリクトが起きやすい

## トークン最適化

Claudeは生成コードを読む必要はありません。DSL定義のみで十分です。

| 読むファイル | 行数 | トークン消費 |
|------------|------|------------|
| ❌ generated/auth-service/**/*.go | 2,000-3,000 | ~15,000 |
| ✅ mps-workspace/solutions/auth-service/service.model | 100-300 | ~1,500 |

**削減率: 90%**

## 参考ドキュメント

- [mps-workspace/README.md](../mps-workspace/README.md) - DSL定義ガイド
- [.claude/rules/code-generation.md](../.claude/rules/code-generation.md) - コード生成ルール
- [.claude/rules/mps-workflow.md](../.claude/rules/mps-workflow.md) - MPS開発フロー
