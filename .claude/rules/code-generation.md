# コード生成ルール

このドキュメントでは、MPS DSLからGoコードを生成する際のルールを定義します。

---

## 基本原則

### 1. 生成コードは読み取り専用

```
generated/
├── auth/
│   ├── domain/
│   ├── usecase/
│   └── handler/
└── ...

↑ このディレクトリ配下のファイルは**絶対に編集しない**
```

**理由**: 再生成で上書きされるため

---

### 2. 生成コードとカスタムコードの分離

| ディレクトリ | 用途 | 編集可否 |
|------------|------|---------|
| `generated/` | MPS生成コード | ❌ 不可 |
| `manual/` | カスタムロジック | ✅ 可能 |
| `mps-workspace/` | DSL定義 | ✅ 可能 |

---

## DSL → Go 生成マッピング

### Entity → Go Struct

**DSL定義**:
```kotlin
entity User {
  id: UUID primary_key
  email: string unique not_null
  password_hash: string not_null
  role: Role not_null
  created_at: timestamp
}
```

**生成されるGoコード**:
```go
// generated/auth/domain/user.go
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

---

### Enum → Go Type

**DSL定義**:
```kotlin
enum Role {
  CUSTOMER, SHOP_OWNER, ADMIN
}
```

**生成されるGoコード**:
```go
// generated/auth/domain/role.go
package domain

type Role string

const (
    RoleCustomer   Role = "CUSTOMER"
    RoleShopOwner  Role = "SHOP_OWNER"
    RoleAdmin      Role = "ADMIN"
)
```

---

### Usecase → Go Interface + Implementation

**DSL定義**:
```kotlin
usecase UserRegistration {
  input: { email: string, password: string, role: Role }
  output: { user_id: UUID, token: string }
  errors: { EmailAlreadyExists, WeakPassword }
}
```

**生成されるGoコード**:
```go
// generated/auth/usecase/user_registration.go
package usecase

import (
    "context"
    "github.com/google/uuid"
)

// Input
type UserRegistrationInput struct {
    Email    string
    Password string
    Role     Role
}

// Output
type UserRegistrationOutput struct {
    UserID uuid.UUID
    Token  string
}

// Errors
type EmailAlreadyExistsError struct{}
func (e EmailAlreadyExistsError) Error() string {
    return "email already exists"
}

type WeakPasswordError struct{}
func (e WeakPasswordError) Error() string {
    return "password is too weak"
}

// Interface
type UserRegistrationUsecase interface {
    Execute(ctx context.Context, input UserRegistrationInput) (UserRegistrationOutput, error)
}

// Implementation
type userRegistrationUsecaseImpl struct {
    userRepo UserRepository
    // ... dependencies
}

func NewUserRegistrationUsecase(userRepo UserRepository) UserRegistrationUsecase {
    return &userRegistrationUsecaseImpl{
        userRepo: userRepo,
    }
}

func (u *userRegistrationUsecaseImpl) Execute(
    ctx context.Context,
    input UserRegistrationInput,
) (UserRegistrationOutput, error) {
    // 生成されたビジネスロジック
    // - バリデーション
    // - Repository呼び出し
    // - エラーハンドリング
}
```

---

### gRPC Service → Handler

**DSL定義**:
```kotlin
grpc_service {
  rpc Register(RegisterRequest) returns (RegisterResponse)
  rpc Login(LoginRequest) returns (LoginResponse)
}
```

**生成されるGoコード**:
```go
// generated/auth/handler/grpc_handler.go
package handler

import (
    "context"
    pb "github.com/.../proto/auth/v1"
)

type AuthServiceHandler struct {
    pb.UnimplementedAuthServiceServer
    registerUsecase usecase.UserRegistrationUsecase
    loginUsecase    usecase.UserLoginUsecase
}

func NewAuthServiceHandler(
    registerUsecase usecase.UserRegistrationUsecase,
    loginUsecase usecase.UserLoginUsecase,
) *AuthServiceHandler {
    return &AuthServiceHandler{
        registerUsecase: registerUsecase,
        loginUsecase:    loginUsecase,
    }
}

func (h *AuthServiceHandler) Register(
    ctx context.Context,
    req *pb.RegisterRequest,
) (*pb.RegisterResponse, error) {
    // 生成されたハンドラーロジック
    // - Requestからusecaseへの変換
    // - Usecase実行
    // - Responseへの変換
}
```

---

## 生成コードの制約

### 1. 編集禁止

```bash
# ❌ 絶対にやってはいけない
vim generated/auth/domain/user.go
git add generated/auth/domain/user.go

# ✅ 正しいアプローチ
vim mps-workspace/solutions/auth-service/service.model
./scripts/mps-generate.sh auth-service
```

---

### 2. Git管理

`.gitignore` に以下を追加推奨：

```gitignore
# MPS Generated Code (オプション: コミットしない場合)
generated/
```

**コミットするかしないか？**

| 方針 | メリット | デメリット |
|------|---------|----------|
| **コミットする** | CI/CDで再生成不要 | Diffが大きい |
| **コミットしない** | Diffが小さい | CI/CDで再生成必要 |

**推奨**: プロジェクトの方針による（ユーザーに確認）

---

### 3. レビュー対象

**コードレビューで確認すべきもの**:
- ✅ DSL定義（`mps-workspace/solutions/`）
- ✅ カスタムロジック（`manual/`）
- ❌ 生成コード（`generated/`）← レビュー不要

**理由**: DSL定義が正しければ、生成コードも正しい

---

## カスタムロジックの統合

### 生成コードからカスタムコードを呼び出す

**パターン1: インターフェース経由**

```go
// generated/auth/usecase/user_registration.go
func (u *userRegistrationUsecaseImpl) Execute(...) {
    // 生成されたコード

    // カスタムバリデーション呼び出し
    if err := u.customValidator.Validate(input); err != nil {
        return output, err
    }

    // 残りの生成コード
}
```

```go
// manual/auth/validation.go
package custom

type PasswordValidator struct{}

func (v *PasswordValidator) Validate(password string) error {
    // カスタムバリデーションロジック
}
```

---

## トラブルシューティング

### 生成コードが期待と異なる

**対処**:
1. DSL定義を確認
2. Generator実装を確認
3. DSL定義を修正 → 再生成

### 生成コードがコンパイルできない

**対処**:
1. DSL定義の型を確認
2. Generator実装のバグを修正
3. 再生成

---

## まとめ

### コード生成の鉄則

1. **生成コードは読まない・触らない**
2. **変更はDSL定義から**
3. **カスタムロジックはmanual/に分離**
4. **生成コードのレビューは不要**

この原則を守ることで、安全かつ効率的に開発できます。
