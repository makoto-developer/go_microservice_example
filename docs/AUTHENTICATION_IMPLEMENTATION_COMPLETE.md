# 認証アカウント分離実装完了報告

## 実装概要

**目的**: Amazonスタイルの完全なカスタマー/オーナーアカウント分離を実現

**実装期間**: 2026-01-25

**ステータス**: ✅ Phase 1-5完了（DB層〜Handler層）

---

## 完了した実装フェーズ

### Phase 1: データベース層 ✅

**マイグレーションファイル**:
- `microservices/auth/migrations/003_split_user_accounts.up.sql`
- `microservices/auth/migrations/003_split_user_accounts.down.sql`

**変更内容**:
1. `customer_users`テーブル作成
   - カスタマー専用のユーザーテーブル
   - メール検証機能
   - パスワードリセット機能

2. `owner_users`テーブル作成
   - ショップオーナー専用のユーザーテーブル
   - ビジネス認証機能（business_verified, business_verification_status）
   - メール検証機能
   - パスワードリセット機能

3. 既存データ移行
   - 既存の`users`テーブルから新テーブルへデータ移行
   - ロール別に適切なテーブルへ振り分け

**重要な設計決定**:
- 同じメールアドレスが`customer_users`と`owner_users`両方に登録可能（Amazonスタイル）
- UNIQUE制約は各テーブル内のみ

---

### Phase 2: Domain層 ✅

**実装ファイル**:
- `microservices/auth/internal/domain/customer_user.go`
- `microservices/auth/internal/domain/owner_user.go`

**CustomerUserエンティティ**:
```go
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
```

**主要メソッド**:
- `NewCustomerUser()`: コンストラクタ
- `SetEmailVerificationToken()`: メール検証トークン設定
- `VerifyEmail()`: メール検証実行
- `SetPasswordResetToken()`: パスワードリセットトークン設定
- `ResetPassword()`: パスワードリセット実行

**OwnerUserエンティティ**:
```go
type OwnerUser struct {
    // CustomerUserと同じフィールド +
    BusinessVerified          bool
    BusinessVerificationStatus BusinessVerificationStatus
}

type BusinessVerificationStatus string

const (
    BusinessVerificationStatusPending  BusinessVerificationStatus = "pending"
    BusinessVerificationStatusApproved BusinessVerificationStatus = "approved"
    BusinessVerificationStatusRejected BusinessVerificationStatus = "rejected"
)
```

**追加メソッド**:
- `CanAccessShopFeatures()`: ショップ機能アクセス可否判定

---

### Phase 3: Repository層 ✅

**インターフェース定義**:
- `microservices/auth/internal/repository/customer_user_repository.go`
- `microservices/auth/internal/repository/owner_user_repository.go`

**PostgreSQL実装**:
- `microservices/auth/internal/repository/postgres/customer_user_repository.go`
- `microservices/auth/internal/repository/postgres/owner_user_repository.go`

**主要メソッド**:
```go
type CustomerUserRepository interface {
    Create(ctx context.Context, user *domain.CustomerUser) error
    FindByID(ctx context.Context, id uuid.UUID) (*domain.CustomerUser, error)
    FindByEmail(ctx context.Context, email string) (*domain.CustomerUser, error)
    FindByVerificationToken(ctx context.Context, token string) (*domain.CustomerUser, error)
    FindByResetToken(ctx context.Context, token string) (*domain.CustomerUser, error)
    Update(ctx context.Context, user *domain.CustomerUser) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

**重要な実装詳細**:
- `FindByEmail`がユーザーが見つからない場合は`nil, nil`を返す（エラーではない）
- これにより「ユーザーが存在しない」と「DBエラー」を区別可能
- 同様のパターンを`FindByVerificationToken`と`FindByResetToken`にも適用

---

### Phase 4: Usecase層 ✅

**実装ファイル**:
- `microservices/auth/internal/usecase/jwt_service_v2.go`
- `microservices/auth/internal/usecase/customer_registration.go`
- `microservices/auth/internal/usecase/customer_login.go`
- `microservices/auth/internal/usecase/owner_registration.go`
- `microservices/auth/internal/usecase/owner_login.go`

#### JWTServiceV2

**新しいクレーム構造**:
```go
type UserType string

const (
    UserTypeCustomer UserType = "customer"
    UserTypeOwner    UserType = "owner"
)

type AccessTokenClaimsV2 struct {
    UserID   string   `json:"user_id"`
    Email    string   `json:"email"`
    UserType UserType `json:"user_type"`
    jwt.RegisteredClaims
}

type RefreshTokenClaimsV2 struct {
    UserID   uuid.UUID `json:"user_id"`
    UserType UserType  `json:"user_type"`
    jwt.RegisteredClaims
}
```

**主要メソッド**:
- `GenerateCustomerAccessToken()`: カスタマー用アクセストークン生成
- `GenerateCustomerRefreshToken()`: カスタマー用リフレッシュトークン生成
- `GenerateOwnerAccessToken()`: オーナー用アクセストークン生成
- `GenerateOwnerRefreshToken()`: オーナー用リフレッシュトークン生成
- `ValidateCustomerAccessToken()`: カスタマートークン検証（audienceチェック含む）
- `ValidateCustomerRefreshToken()`: カスタマーリフレッシュトークン検証
- `ValidateOwnerAccessToken()`: オーナートークン検証
- `ValidateOwnerRefreshToken()`: オーナーリフレッシュトークン検証

**セキュリティ強化**:
- `audience`フィールドでトークンの対象を明示（"customer"または"owner"）
- トークン検証時にuser_typeとaudienceの両方をチェック
- カスタマートークンをオーナーエンドポイントで使用不可（逆も同様）

#### Registration Usecases

**CustomerRegistrationUsecase**:
```go
func (u *CustomerRegistrationUsecase) Execute(
    ctx context.Context,
    email, password string,
) (userID, accessToken, refreshToken string, err error)
```

**処理フロー**:
1. customer_usersテーブルでメール重複チェック
2. パスワード強度検証（最低8文字）
3. パスワードハッシュ化（bcrypt）
4. CustomerUserエンティティ作成
5. メール検証トークン生成（32バイトランダム、24時間有効）
6. データベースに保存
7. JWT トークン生成（access + refresh）
8. 非同期でメール検証メール送信

**OwnerRegistrationUsecase**: 同様の処理フローだが、ビジネス認証ステータスを"pending"で初期化

#### Login Usecases

**CustomerLoginUsecase**:
```go
func (u *CustomerLoginUsecase) Execute(
    ctx context.Context,
    email, password string,
) (userID, accessToken, refreshToken string, err error)
```

**処理フロー**:
1. customer_usersテーブルからメールでユーザー検索
2. パスワード検証（bcrypt.CompareHashAndPassword）
3. メール検証済みチェック
4. JWT トークン生成

**OwnerLoginUsecase**: 同様だが、owner_usersテーブルを使用

---

### Phase 5: Proto定義 ✅

**実装ファイル**:
- `microservices/auth/proto/customer_auth.proto`
- `microservices/auth/proto/owner_auth.proto`

#### CustomerAuthService

```protobuf
service CustomerAuthService {
  rpc Register(CustomerRegisterRequest) returns (CustomerRegisterResponse);
  rpc Login(CustomerLoginRequest) returns (CustomerLoginResponse);
  rpc Logout(CustomerLogoutRequest) returns (CustomerLogoutResponse);
  rpc VerifyEmail(CustomerVerifyEmailRequest) returns (CustomerVerifyEmailResponse);
  rpc RequestPasswordReset(CustomerRequestPasswordResetRequest)
      returns (CustomerRequestPasswordResetResponse);
  rpc ResetPassword(CustomerResetPasswordRequest)
      returns (CustomerResetPasswordResponse);
  rpc RefreshToken(CustomerRefreshTokenRequest)
      returns (CustomerRefreshTokenResponse);
}
```

**主要メッセージ**:
- `CustomerRegisterRequest`: email, password
- `CustomerRegisterResponse`: user_id, access_token, refresh_token
- `CustomerLoginRequest`: email, password
- `CustomerLoginResponse`: user_id, access_token, refresh_token

#### OwnerAuthService

```protobuf
service OwnerAuthService {
  // CustomerAuthServiceと同じメソッド +
  rpc GetBusinessVerificationStatus(OwnerGetBusinessVerificationStatusRequest)
      returns (OwnerGetBusinessVerificationStatusResponse);
}
```

**追加フィールド**:
- `OwnerRegisterResponse`: business_verification_status
- `OwnerLoginResponse`: business_verified, business_verification_status

---

### Phase 6: Handler層 ✅

**実装ファイル**:
- `microservices/auth/internal/handler/grpc/customer_auth_handler.go`
- `microservices/auth/internal/handler/grpc/owner_auth_handler.go`

#### CustomerAuthHandler

**実装メソッド**:
1. `Register()`: カスタマー登録
   - リクエストバリデーション
   - CustomerRegistrationUsecase実行
   - エラーマッピング（AlreadyExists, InvalidArgument等）

2. `Login()`: カスタマーログイン
   - CustomerLoginUsecase実行
   - 認証エラー処理（Unauthenticated, FailedPrecondition）

3. `Logout()`: ログアウト
   - TODO: Redis等でトークン無効化（現在は自然失効のみ）

4. `VerifyEmail()`: メールアドレス検証
   - トークンでユーザー検索
   - ドメインメソッド`VerifyEmail()`呼び出し
   - データベース更新

5. `RequestPasswordReset()`: パスワードリセット要求
   - メールからユーザー検索
   - リセットトークン生成・設定
   - メール列挙攻撃対策（常に成功レスポンス）

6. `ResetPassword()`: パスワードリセット実行
   - トークン検証
   - パスワード強度チェック
   - ドメインメソッド`ResetPassword()`呼び出し

7. `RefreshToken()`: トークン更新
   - リフレッシュトークン検証
   - 新しいアクセストークン・リフレッシュトークン生成

**エラーハンドリング**:
```go
// 適切なgRPCステータスコードマッピング
codes.InvalidArgument     // バリデーションエラー
codes.AlreadyExists       // 既に登録済み
codes.Unauthenticated     // 認証失敗
codes.FailedPrecondition  // メール未検証等
codes.NotFound            // ユーザーが見つからない
codes.Internal            // サーバーエラー
```

#### OwnerAuthHandler

**CustomerAuthHandlerと同様のメソッド** +

8. `GetBusinessVerificationStatus()`: ビジネス認証ステータス取得
   - ユーザーID検証
   - ビジネス認証状態（pending/approved/rejected）返却

**Login()の特別処理**:
- ログイン後にビジネス認証ステータスも取得
- レスポンスに`business_verified`と`business_verification_status`を含める

---

## アーキテクチャ図

### Before（単一テーブル設計）
```
┌─────────────────┐
│     users       │
├─────────────────┤
│ id              │
│ email (unique)  │
│ password_hash   │
│ role (enum)     │  ← CUSTOMER / SHOP_OWNER
└─────────────────┘
```

**問題点**:
- 同じメールで異なるロールに登録不可
- Amazonスタイルではない

---

### After（分離テーブル設計）
```
┌──────────────────────┐        ┌──────────────────────┐
│   customer_users     │        │     owner_users      │
├──────────────────────┤        ├──────────────────────┤
│ id                   │        │ id                   │
│ email (unique)       │        │ email (unique)       │
│ password_hash        │        │ password_hash        │
│ email_verified       │        │ email_verified       │
│                      │        │ business_verified    │
│                      │        │ business_verification│
│                      │        │  _status             │
└──────────────────────┘        └──────────────────────┘
        ↑                                ↑
        │                                │
  同じメールでも両方に登録可能
     (例: user@example.com)
```

**メリット**:
- Amazonスタイルの完全分離
- 同じメールで顧客とオーナー両方になれる
- ドメインロジックの明確な分離

---

## JWT トークン構造

### Before
```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "role": "CUSTOMER",
  "exp": 1234567890
}
```

### After
```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "user_type": "customer",
  "aud": ["customer"],
  "iss": "auth-service",
  "exp": 1234567890
}
```

**セキュリティ強化**:
- `user_type`フィールドで明示的にユーザータイプ指定
- `aud`（audience）でトークンの対象を制限
- カスタマートークンをオーナーエンドポイントで使用不可

---

## API エンドポイント設計

### カスタマー認証
```
POST /auth/customer/register
POST /auth/customer/login
POST /auth/customer/logout
POST /auth/customer/verify-email
POST /auth/customer/request-password-reset
POST /auth/customer/reset-password
POST /auth/customer/refresh-token
```

### オーナー認証
```
POST /auth/owner/register
POST /auth/owner/login
POST /auth/owner/logout
POST /auth/owner/verify-email
POST /auth/owner/request-password-reset
POST /auth/owner/reset-password
POST /auth/owner/refresh-token
GET  /auth/owner/business-verification-status
```

---

## テスト観点

### 単体テスト
- [ ] Domain層のビジネスロジック
- [ ] Repository層のデータアクセス
- [ ] Usecase層の処理フロー
- [ ] Handler層のエラーハンドリング

### 統合テスト
- [ ] 同じメールでcustomer + owner登録
- [ ] JWT トークンのaudience検証
- [ ] メール検証フロー
- [ ] パスワードリセットフロー
- [ ] ビジネス認証フロー（オーナーのみ）

### E2Eテスト
- [ ] カスタマー登録〜ログイン〜ログアウト
- [ ] オーナー登録〜メール検証〜ビジネス認証待ち
- [ ] トークンリフレッシュ
- [ ] パスワードリセット

---

## 次のステップ

### 残りの実装タスク

#### 1. Server統合（優先度: 高）
- [ ] `microservices/auth/cmd/server/main.go`更新
- [ ] CustomerAuthServiceとOwnerAuthServiceの両方を登録
- [ ] 依存関係の注入（Repository, Usecase, Handler）

#### 2. Proto生成（優先度: 高）
```bash
# customer_auth.proto から生成
protoc --go_out=. --go-grpc_out=. \
  microservices/auth/proto/customer_auth.proto

# owner_auth.proto から生成
protoc --go_out=. --go-grpc_out=. \
  microservices/auth/proto/owner_auth.proto
```

#### 3. 他サービスの更新（優先度: 中）

**Customer Service**:
- `users.user_id` → `customer_users.id`参照に変更
- gRPCクライアントをCustomerAuthServiceに更新

**Shop Service**:
- `users.user_id` → `owner_users.id`参照に変更
- gRPCクライアントをOwnerAuthServiceに更新
- ビジネス認証チェック追加

#### 4. メール送信実装（優先度: 中）
- [ ] EmailService実装（SendGrid / AWS SES等）
- [ ] メール検証メールテンプレート
- [ ] パスワードリセットメールテンプレート

#### 5. トークン無効化実装（優先度: 低）
- [ ] Redisセットアップ
- [ ] トークンブラックリスト実装
- [ ] Logout時にリフレッシュトークン無効化

#### 6. ビジネス認証管理（優先度: 中）
- [ ] Admin用のビジネス認証承認/拒否API
- [ ] ビジネス認証ステータス変更時の通知
- [ ] 認証書類アップロード機能

---

## マイグレーション実行手順

### 開発環境

```bash
# マイグレーション適用
cd microservices/auth
migrate -path migrations -database "postgresql://user:pass@localhost:5432/auth_db?sslmode=disable" up

# 確認
psql -U user -d auth_db -c "\dt"
# customer_users, owner_usersテーブルが作成されていることを確認

psql -U user -d auth_db -c "SELECT COUNT(*) FROM customer_users;"
psql -U user -d auth_db -c "SELECT COUNT(*) FROM owner_users;"
```

### ロールバック（必要な場合）

```bash
migrate -path migrations -database "postgresql://user:pass@localhost:5432/auth_db?sslmode=disable" down 1
```

---

## 設定ファイル更新

### 環境変数

```bash
# JWT設定
JWT_SECRET=your-secret-key-here
JWT_ACCESS_TOKEN_EXPIRY=1h
JWT_REFRESH_TOKEN_EXPIRY=168h  # 7日

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=auth_db
DB_USER=user
DB_PASSWORD=pass

# Email（今後実装）
EMAIL_PROVIDER=sendgrid
SENDGRID_API_KEY=xxx
```

---

## パフォーマンス考慮事項

### データベースインデックス

既に作成済み:
```sql
-- customer_users
CREATE UNIQUE INDEX idx_customer_users_email ON customer_users(email);
CREATE INDEX idx_customer_users_verification_token
  ON customer_users(email_verification_token);
CREATE INDEX idx_customer_users_reset_token
  ON customer_users(password_reset_token);

-- owner_users
CREATE UNIQUE INDEX idx_owner_users_email ON owner_users(email);
CREATE INDEX idx_owner_users_verification_token
  ON owner_users(email_verification_token);
CREATE INDEX idx_owner_users_reset_token
  ON owner_users(password_reset_token);
CREATE INDEX idx_owner_users_business_status
  ON owner_users(business_verification_status);
```

### キャッシュ戦略（今後検討）

- メール検証トークンの有効期限チェック（Redis）
- ビジネス認証ステータスキャッシュ
- JWTトークンブラックリスト

---

## セキュリティチェックリスト

- [x] パスワードハッシュ化（bcrypt）
- [x] JWT署名検証
- [x] audience検証（カスタマー/オーナー分離）
- [x] メール検証トークン有効期限
- [x] パスワードリセットトークン有効期限
- [x] パスワード強度検証（最低8文字）
- [x] メール列挙攻撃対策（パスワードリセット時）
- [ ] レート制限（今後実装）
- [ ] CSRF対策（今後実装）
- [ ] トークンリフレッシュ時の検証強化

---

## まとめ

### 完了した作業

1. ✅ データベース設計変更（単一テーブル → 分離テーブル）
2. ✅ Domainエンティティ実装（CustomerUser, OwnerUser）
3. ✅ Repository層実装（PostgreSQL）
4. ✅ Usecase層実装（Registration, Login + JWTServiceV2）
5. ✅ Proto定義（CustomerAuthService, OwnerAuthService）
6. ✅ Handler層実装（gRPC handlers）

### 実装された機能

#### カスタマー
- ユーザー登録
- ログイン/ログアウト
- メール検証
- パスワードリセット
- トークンリフレッシュ

#### オーナー
- 上記すべて +
- ビジネス認証ステータス管理
- ショップ機能アクセス制御

### 達成した目標

✅ Amazonスタイルの完全なアカウント分離
✅ 同じメールで顧客・オーナー両方登録可能
✅ JWT トークンのセキュリティ強化（user_type, audience）
✅ ドメイン駆動設計の適用
✅ 拡張性の高いアーキテクチャ

---

**実装完了日**: 2026-01-25
**実装者**: Claude
**レビューステータス**: 未レビュー
**次のアクション**: Server統合 & Proto生成
