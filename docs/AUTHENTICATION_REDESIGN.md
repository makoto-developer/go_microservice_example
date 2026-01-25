# 認証設計の見直し - Amazon型アカウント分離

**問題**: 現在の設計では、カスタマーとオーナーが同じUserテーブルでRoleで区別されている。
これはAmazonのような実際のECサイトと異なる。

**Amazonの例**:
- 一般顧客: amazon.com でショッピング
- 出品者: seller.amazon.com (Amazon Seller Central) で販売管理
- **完全に別のアカウント体系**
- 同じメールアドレスでも両方登録可能

---

## 🎯 新しい設計方針

### 基本原則

1. **カスタマーとオーナーは完全に別のアカウント**
2. **同じメールアドレスで両方登録可能**
3. **認証エンドポイントも分離**
4. **JWT issuerで区別**

---

## 📐 設計パターンの比較

### パターンA: Auth Service内で完全分離（推奨）

```
Auth Service
├── customer_users テーブル
│   ├── id (UUID)
│   ├── email
│   ├── password_hash
│   └── email_verified
│
└── owner_users テーブル
    ├── id (UUID)
    ├── email (customer_usersと重複可)
    ├── password_hash
    └── email_verified

エンドポイント:
- POST /api/v1/auth/customer/register
- POST /api/v1/auth/customer/login
- POST /api/v1/auth/owner/register
- POST /api/v1/auth/owner/login

JWT Claims:
{
  "user_id": "...",
  "email": "...",
  "type": "customer" | "owner",  // ← 新規追加
  "iss": "auth-service"
}
```

**メリット**:
- ✅ 認証ロジックが統一（JWT発行等）
- ✅ パスワード管理が一元化
- ✅ メール送信機能が共通
- ✅ 実装がシンプル

**デメリット**:
- ❌ Auth Serviceが肥大化

**採用理由**: このプロジェクトの規模に最適

---

### パターンB: 各サービスに認証を持たせる

```
Customer Service
└── customer_accounts
    ├── id
    ├── email
    └── password_hash

Shop Service
└── owner_accounts
    ├── id
    ├── email
    └── password_hash

Auth Service
└── JWT発行・検証のみ（共通ライブラリ的な役割）
```

**メリット**:
- ✅ 完全な独立性
- ✅ 各サービスが自律的

**デメリット**:
- ❌ 認証ロジックの重複
- ❌ セキュリティ更新が煩雑

---

### パターンC: 別々のAuth Service

```
Customer Auth Service
└── customer_users

Owner Auth Service
└── owner_users
```

**メリット**:
- ✅ 完全分離
- ✅ スケーリングが独立

**デメリット**:
- ❌ サービス数が増える
- ❌ コード重複

---

## ✅ 採用する設計: パターンA

**Auth Service内で完全分離**

---

## 📋 新しいDB設計

### customer_users テーブル

```sql
CREATE TABLE customer_users (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email                           VARCHAR(255) NOT NULL UNIQUE,
    password_hash                   VARCHAR(255) NOT NULL,
    email_verified                  BOOLEAN NOT NULL DEFAULT FALSE,
    email_verification_token        VARCHAR(255),
    email_verification_expires_at   TIMESTAMP,
    password_reset_token            VARCHAR(255),
    password_reset_expires_at       TIMESTAMP,
    created_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_customer_users_email ON customer_users(email);
CREATE INDEX idx_customer_users_verification_token ON customer_users(email_verification_token);
CREATE INDEX idx_customer_users_reset_token ON customer_users(password_reset_token);
```

### owner_users テーブル

```sql
CREATE TABLE owner_users (
    id                              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email                           VARCHAR(255) NOT NULL UNIQUE,
    password_hash                   VARCHAR(255) NOT NULL,
    email_verified                  BOOLEAN NOT NULL DEFAULT FALSE,
    email_verification_token        VARCHAR(255),
    email_verification_expires_at   TIMESTAMP,
    password_reset_token            VARCHAR(255),
    password_reset_expires_at       TIMESTAMP,
    business_verified               BOOLEAN NOT NULL DEFAULT FALSE,  -- ← オーナー固有
    business_verification_status    VARCHAR(50) DEFAULT 'pending',   -- ← オーナー固有
    created_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at                      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_owner_users_email ON owner_users(email);
CREATE INDEX idx_owner_users_verification_token ON owner_users(email_verification_token);
CREATE INDEX idx_owner_users_reset_token ON owner_users(password_reset_token);
```

**重要**: 同じメールアドレスが`customer_users`と`owner_users`の両方に存在可能。

---

## 🔧 新しいドメインモデル

### auth/internal/domain/customer_user.go

```go
package domain

import (
    "time"
    "github.com/google/uuid"
)

type CustomerUser struct {
    ID                        uuid.UUID
    Email                     string
    PasswordHash              string
    EmailVerified             bool
    EmailVerificationToken    string
    EmailVerificationExpiresAt *time.Time
    PasswordResetToken        string
    PasswordResetExpiresAt    *time.Time
    CreatedAt                 time.Time
    UpdatedAt                 time.Time
}

func NewCustomerUser(email, passwordHash string) *CustomerUser {
    now := time.Now()
    return &CustomerUser{
        ID:            uuid.New(),
        Email:         email,
        PasswordHash:  passwordHash,
        EmailVerified: false,
        CreatedAt:     now,
        UpdatedAt:     now,
    }
}
```

### auth/internal/domain/owner_user.go

```go
package domain

import (
    "time"
    "github.com/google/uuid"
)

type OwnerUser struct {
    ID                        uuid.UUID
    Email                     string
    PasswordHash              string
    EmailVerified             bool
    EmailVerificationToken    string
    EmailVerificationExpiresAt *time.Time
    PasswordResetToken        string
    PasswordResetExpiresAt    *time.Time
    BusinessVerified          bool    // ← オーナー固有
    BusinessVerificationStatus string // ← オーナー固有
    CreatedAt                 time.Time
    UpdatedAt                 time.Time
}

func NewOwnerUser(email, passwordHash string) *OwnerUser {
    now := time.Now()
    return &OwnerUser{
        ID:                        uuid.New(),
        Email:                     email,
        PasswordHash:              passwordHash,
        EmailVerified:             false,
        BusinessVerified:          false,
        BusinessVerificationStatus: "pending",
        CreatedAt:                 now,
        UpdatedAt:                 now,
    }
}
```

---

## 🔐 JWT設計

### JWTペイロード

```go
type Claims struct {
    UserID    string `json:"user_id"`
    Email     string `json:"email"`
    UserType  string `json:"user_type"`  // "customer" or "owner"
    jwt.StandardClaims
}
```

### JWT発行

```go
// カスタマー用
func GenerateCustomerAccessToken(userID, email string) (string, error) {
    claims := &Claims{
        UserID:   userID,
        Email:    email,
        UserType: "customer",
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
            Issuer:    "auth-service",
            Audience:  "customer",  // ← 重要
        },
    }
    // ...
}

// オーナー用
func GenerateOwnerAccessToken(userID, email string) (string, error) {
    claims := &Claims{
        UserID:   userID,
        Email:    email,
        UserType: "owner",
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
            Issuer:    "auth-service",
            Audience:  "owner",  // ← 重要
        },
    }
    // ...
}
```

---

## 🌐 API設計

### Customer用エンドポイント

```
POST /api/v1/auth/customer/register
POST /api/v1/auth/customer/login
POST /api/v1/auth/customer/logout
POST /api/v1/auth/customer/verify-email
POST /api/v1/auth/customer/request-password-reset
POST /api/v1/auth/customer/reset-password
POST /api/v1/auth/customer/refresh-token
```

### Owner用エンドポイント

```
POST /api/v1/auth/owner/register
POST /api/v1/auth/owner/login
POST /api/v1/auth/owner/logout
POST /api/v1/auth/owner/verify-email
POST /api/v1/auth/owner/verify-business  // ← オーナー固有
POST /api/v1/auth/owner/request-password-reset
POST /api/v1/auth/owner/reset-password
POST /api/v1/auth/owner/refresh-token
```

---

## 📊 サービス間の関連

### Customer Service

```go
type Customer struct {
    ID              uuid.UUID
    AuthUserID      uuid.UUID  // customer_users.id を参照
    FirstName       string
    LastName        string
    // ...
}
```

### Shop Service

```go
type Shop struct {
    ID          uuid.UUID
    OwnerAuthID uuid.UUID  // owner_users.id を参照
    Name        string
    // ...
}
```

---

## 🔄 マイグレーション計画

### Step 1: 新テーブル作成

```sql
-- customer_users テーブル作成
CREATE TABLE customer_users (...);

-- owner_users テーブル作成
CREATE TABLE owner_users (...);
```

### Step 2: 既存データ移行

```sql
-- 既存のusersテーブルからcustomer_usersへ
INSERT INTO customer_users (id, email, password_hash, email_verified, created_at, updated_at)
SELECT id, email, password_hash, email_verified, created_at, updated_at
FROM users
WHERE role = 'CUSTOMER';

-- 既存のusersテーブルからowner_usersへ
INSERT INTO owner_users (id, email, password_hash, email_verified, created_at, updated_at)
SELECT id, email, password_hash, email_verified, created_at, updated_at
FROM users
WHERE role = 'SHOP_OWNER';
```

### Step 3: Customer/Shopテーブル更新

```sql
-- Customerテーブルの外部キーは維持（user_id → auth_user_id）
ALTER TABLE customers RENAME COLUMN user_id TO auth_user_id;

-- Shopテーブルの外部キーは維持（owner_id → owner_auth_id）
ALTER TABLE shops RENAME COLUMN owner_id TO owner_auth_id;
```

### Step 4: 旧テーブル削除

```sql
-- usersテーブルとuser_rolesテーブルを削除（バックアップ後）
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS users;
```

---

## 🎯 実装タスク

### Task 1: DB Migration

- [ ] customer_users テーブル作成
- [ ] owner_users テーブル作成
- [ ] 既存データ移行スクリプト
- [ ] 外部キー更新

### Task 2: Domain層

- [ ] CustomerUser エンティティ
- [ ] OwnerUser エンティティ
- [ ] 旧User削除

### Task 3: Repository層

- [ ] CustomerUserRepository
- [ ] OwnerUserRepository
- [ ] 旧UserRepository削除

### Task 4: Usecase層

- [ ] CustomerRegistrationUsecase
- [ ] CustomerLoginUsecase
- [ ] OwnerRegistrationUsecase
- [ ] OwnerLoginUsecase
- [ ] 旧usecaseリファクタ

### Task 5: Handler層

- [ ] CustomerAuthHandler (gRPC)
- [ ] OwnerAuthHandler (gRPC)
- [ ] 旧AuthHandler削除

### Task 6: JWT Service

- [ ] GenerateCustomerAccessToken
- [ ] GenerateOwnerAccessToken
- [ ] ValidateToken (type判定追加)

### Task 7: Proto定義

- [ ] customer_auth.proto
- [ ] owner_auth.proto
- [ ] 旧auth.proto削除

### Task 8: 他サービス更新

- [ ] Customer Service: auth_user_id参照に変更
- [ ] Shop Service: owner_auth_id参照に変更

---

## 📝 実装例

### CustomerRegistrationUsecase

```go
package usecase

import (
    "context"
    "fmt"
    "golang.org/x/crypto/bcrypt"
    "github.com/makoto-developer/go_microservice_example/microservices/auth/internal/domain"
    "github.com/makoto-developer/go_microservice_example/microservices/auth/internal/repository"
)

type CustomerRegistrationUsecase struct {
    customerUserRepo repository.CustomerUserRepository
    jwtService       *JWTService
    emailService     *EmailService
}

func NewCustomerRegistrationUsecase(
    repo repository.CustomerUserRepository,
    jwt *JWTService,
    email *EmailService,
) *CustomerRegistrationUsecase {
    return &CustomerRegistrationUsecase{
        customerUserRepo: repo,
        jwtService:       jwt,
        emailService:     email,
    }
}

func (u *CustomerRegistrationUsecase) Execute(
    ctx context.Context,
    email, password string,
) (userID, accessToken, refreshToken string, err error) {
    // メールアドレス重複チェック（customer_usersテーブル内のみ）
    existing, _ := u.customerUserRepo.FindByEmail(ctx, email)
    if existing != nil {
        return "", "", "", fmt.Errorf("email already registered as customer")
    }

    // パスワードハッシュ化
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", "", "", fmt.Errorf("failed to hash password: %w", err)
    }

    // CustomerUser作成
    user := domain.NewCustomerUser(email, string(hashedPassword))

    // メール認証トークン生成
    token, _ := generateRandomToken(32)
    user.EmailVerificationToken = token
    expiresAt := time.Now().Add(24 * time.Hour)
    user.EmailVerificationExpiresAt = &expiresAt

    // DB保存
    if err := u.customerUserRepo.Create(ctx, user); err != nil {
        return "", "", "", fmt.Errorf("failed to create user: %w", err)
    }

    // JWT発行（customer用）
    accessToken, err = u.jwtService.GenerateCustomerAccessToken(user.ID.String(), user.Email)
    if err != nil {
        return "", "", "", err
    }

    refreshToken, err = u.jwtService.GenerateCustomerRefreshToken(user.ID.String())
    if err != nil {
        return "", "", "", err
    }

    // メール送信
    u.emailService.SendVerificationEmail(email, token)

    return user.ID.String(), accessToken, refreshToken, nil
}
```

---

## 🔒 セキュリティ考慮事項

### 1. メールアドレスの重複

**許可**: 同じメールアドレスでcustomerとownerに登録可能

**実装**:
```go
// customer_users.email にUNIQUE制約
// owner_users.email にUNIQUE制約
// ただし、テーブル間での重複はOK
```

### 2. JWT検証

各サービスでJWT検証時に`user_type`と`audience`をチェック：

```go
func ValidateCustomerToken(tokenString string) (*Claims, error) {
    claims, err := parseToken(tokenString)
    if err != nil {
        return nil, err
    }

    if claims.UserType != "customer" {
        return nil, errors.New("invalid user type")
    }

    if claims.Audience != "customer" {
        return nil, errors.New("invalid audience")
    }

    return claims, nil
}
```

### 3. パスワードポリシー

Customer/Owner共通:
- 8文字以上
- 英大文字、小文字、数字、記号を含む

---

## 📖 まとめ

### 変更点

| 項目 | Before | After |
|------|--------|-------|
| アカウント | 1つのUserテーブル | customer_users と owner_users に分離 |
| Role | RoleフィールドでCUSTOMER/OWNER | テーブル自体で分離 |
| メール重複 | 不可 | テーブル間では可能 |
| JWT | roleクレーム | user_typeとaudienceクレーム |
| エンドポイント | /auth/register | /auth/customer/register, /auth/owner/register |

### 期待される効果

- ✅ Amazonと同様の明確な分離
- ✅ 同じメールで顧客と出品者の両方になれる
- ✅ セキュリティ向上（JWT audienceチェック）
- ✅ 将来的な拡張が容易（business_verification等）

---

## 次のステップ

1. この設計レビュー
2. 承認後、マイグレーションスクリプト作成
3. 実装（Domain → Repository → Usecase → Handler）
4. テスト
5. デプロイ
