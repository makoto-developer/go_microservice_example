# 認証設計レビュー - カスタマーとオーナーのアカウント分離

**タスク**: カスタマーとオーナーは認証を分けていますか？(アカウント)を分けているか？設計をチェック

**レビュー日**: 2026-01-25

---

## 📋 現在の設計

### 1. Auth Service（認証サービス）

#### User Entity (auth/internal/domain/user.go)

```go
type Role string

const (
    RoleCustomer  Role = "CUSTOMER"
    RoleShopOwner Role = "SHOP_OWNER"
    RoleAdmin     Role = "ADMIN"
)

type User struct {
    ID           uuid.UUID
    Email        string
    PasswordHash string
    Role         Role        // ← ロールで区別
    // ...
}
```

**特徴**:
- ✅ **1つのUserテーブル**でカスタマーとオーナーを管理
- ✅ **Roleフィールド**で区別（RBAC: Role-Based Access Control）
- ✅ 同一のメールアドレスで登録（ロールごとに1アカウント）

---

### 2. Customer Service（顧客サービス）

#### Customer Entity (customer/internal/domain/customer.go)

```go
type Customer struct {
    ID              uuid.UUID
    UserID          uuid.UUID  // ← Auth ServiceのUser IDを参照
    FirstName       string
    LastName        string
    PhoneNumber     string
    // ... 顧客固有情報
}
```

**特徴**:
- ✅ **UserID**でAuth ServiceのUserを参照
- ✅ 顧客固有の情報（名前、電話番号、住所等）を管理
- ✅ Role=CUSTOMERのUserに対応

---

### 3. Shop Service（ショップサービス）

#### Shop Entity (shop/internal/domain/shop.go)

```go
type Shop struct {
    ID          uuid.UUID
    OwnerID     uuid.UUID  // ← Auth ServiceのUser IDを参照
    Name        string
    Description string
    // ... ショップ固有情報
}
```

**特徴**:
- ✅ **OwnerID**でAuth ServiceのUserを参照
- ✅ ショップ固有の情報（名前、説明、営業時間等）を管理
- ✅ Role=SHOP_OWNERのUserに対応

---

## ✅ 設計評価

### 採用されているパターン

**パターン**: **ロールベースの単一認証 + サービス分離**

```
┌─────────────────────────────────────────┐
│       Auth Service (認証統合)            │
│  ┌───────────────────────────────────┐  │
│  │ User (email, password, role)      │  │
│  │  - Role: CUSTOMER                 │  │
│  │  - Role: SHOP_OWNER               │  │
│  │  - Role: ADMIN                    │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
           ↓                    ↓
    ┌──────────┐         ┌──────────┐
    │ Customer │         │   Shop   │
    │ Service  │         │ Service  │
    │(user_id) │         │(owner_id)│
    └──────────┘         └──────────┘
```

---

## 🔍 設計の妥当性チェック

### ✅ 良い点（推奨される設計）

#### 1. **認証の統一**
- ✅ 1つのAuth Serviceで全ユーザーを管理
- ✅ JWT発行ロジックが統一
- ✅ パスワードリセット、メール認証等の機能が共通化

#### 2. **関心の分離**
- ✅ Auth Service: 認証・認可のみ
- ✅ Customer Service: 顧客情報管理
- ✅ Shop Service: ショップ情報管理

#### 3. **柔軟な権限管理**
- ✅ JWTトークンにroleを含める
- ✅ 各サービスでroleベースのアクセス制御が可能
- ✅ 将来的なロール追加が容易（例: MODERATOR, SUPPORT）

#### 4. **マイクロサービスアーキテクチャに適合**
- ✅ 各サービスが独立したデータベースを持つ
- ✅ サービス間の依存は最小限（User IDの参照のみ）
- ✅ スケーラビリティが高い

---

### ⚠️ 確認すべき点

#### 1. **1人が複数ロールを持てるか？**

**現状**: 1つのUserは1つのRoleのみ

**問題**:
- ❌ ショップオーナーが顧客として商品を購入できない
- ❌ 1人で複数のショップを運営する場合に対応できない

**解決策A**: User-Role多対多関係に変更

```go
// 変更案
type User struct {
    ID           uuid.UUID
    Email        string
    PasswordHash string
    Roles        []Role  // ← 複数ロール対応
}

// または
type UserRole struct {
    UserID uuid.UUID
    Role   Role
}
```

**解決策B**: 現状維持（1ロール1アカウント）

- ショップオーナー用とカスタマー用で別々にアカウント登録
- メールアドレスが重複する場合はエラー

**推奨**: **解決策A（複数ロール対応）**
- ユーザビリティが向上
- 実運用に即している

---

#### 2. **CustomerとShopの関係性**

**現状の関係**:

```
User (Role=CUSTOMER)
  ↓
Customer (user_id)

User (Role=SHOP_OWNER)
  ↓
Shop (owner_id)
```

**問題点**:
- ❌ Role=CUSTOMERのUserはShopを持てない
- ❌ Role=SHOP_OWNERのUserはCustomer情報を持てない

**実際の要件**:
- ✅ ショップオーナーも商品を購入したい
- ✅ 1人のユーザーがCustomerとShopOwner両方の役割を持つ

**解決策**:

```
User (Roles=[CUSTOMER, SHOP_OWNER])
  ↓              ↓
Customer      Shop
(user_id)   (owner_id)
```

両方のテーブルにレコードを作成する。

---

#### 3. **ショップ登録フロー**

**現状の想定フロー**:

```
1. ユーザー登録 (Role=SHOP_OWNER)
   ↓
2. Auth Serviceでアカウント作成
   ↓
3. ショップ登録
   ↓
4. Shop ServiceでShop作成
```

**問題**:
- ❌ 既存のカスタマーがショップを開設できない
- ❌ Role変更の仕組みがない

**推奨フロー**:

```
1. ユーザー登録 (Role=CUSTOMER)
   ↓
2. 顧客として利用
   ↓
3. ショップ開設申請
   ↓
4. User.RolesにSHOP_OWNERを追加
   ↓
5. Shop作成
```

---

#### 4. **顧客情報とショップオーナー情報の重複**

**問題**:
- Customer.FirstName, LastName
- Shop.OwnerName, OwnerPhone

これらは同じUserの情報を異なるサービスで管理している。

**解決策**:
- Shop Serviceは最小限の情報のみ（ショップ情報）
- オーナーの個人情報はCustomer Serviceまたは専用のProfile Serviceで管理

---

## 📊 推奨される設計改善

### 改善1: 複数ロール対応

#### Before（現在）

```go
type User struct {
    Role Role  // 1つのみ
}
```

#### After（推奨）

```go
type User struct {
    // Rolesフィールドを追加（後方互換性のためRoleも残す）
}

type UserRole struct {
    ID        uuid.UUID
    UserID    uuid.UUID
    Role      Role
    CreatedAt time.Time
}
```

**マイグレーション**:
```sql
CREATE TABLE user_roles (
    id         UUID PRIMARY KEY,
    user_id    UUID NOT NULL REFERENCES users(id),
    role       VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    UNIQUE(user_id, role)
);

-- 既存データの移行
INSERT INTO user_roles (id, user_id, role, created_at)
SELECT gen_random_uuid(), id, role, created_at FROM users;
```

---

### 改善2: ショップ開設フロー

#### 新しいユースケース

**Shop Service**:
```go
// RegisterShopUsecase
// 既存のユーザーがショップを開設
func (u *RegisterShopUsecase) Execute(ctx context.Context, userID uuid.UUID, shopData ShopData) error {
    // 1. Auth Serviceにユーザー確認 & ロール追加を依頼
    err := u.authClient.AddUserRole(ctx, userID, "SHOP_OWNER")

    // 2. ショップ作成
    shop := &Shop{
        OwnerID: userID,
        Name:    shopData.Name,
        // ...
    }

    return u.shopRepo.Create(ctx, shop)
}
```

**Auth Service**:
```go
// AddUserRoleUsecase
// ユーザーにロールを追加
func (u *AddUserRoleUsecase) Execute(ctx context.Context, userID uuid.UUID, role Role) error {
    // user_rolesテーブルに追加
    userRole := &UserRole{
        UserID: userID,
        Role:   role,
    }
    return u.userRoleRepo.Create(ctx, userRole)
}
```

---

### 改善3: JWT Claimsの更新

#### Before（現在）

```go
type Claims struct {
    UserID string `json:"user_id"`
    Role   string `json:"role"`  // 1つのみ
    jwt.StandardClaims
}
```

#### After（推奨）

```go
type Claims struct {
    UserID string   `json:"user_id"`
    Roles  []string `json:"roles"`  // 複数ロール
    jwt.StandardClaims
}
```

**チェック例**:
```go
// 顧客機能へのアクセス
if contains(claims.Roles, "CUSTOMER") {
    // OK
}

// ショップ管理機能へのアクセス
if contains(claims.Roles, "SHOP_OWNER") {
    // OK
}
```

---

## 🎯 結論と推奨アクション

### 現在の設計の評価

**✅ 基本設計は適切**
- ロールベースアクセス制御（RBAC）
- マイクロサービスアーキテクチャに適合
- 認証の統一

**⚠️ 改善が必要な点**
- 1ユーザー1ロールの制約
- ショップ開設フローの不明確さ
- 情報の重複管理

---

### 推奨される実装順序

#### Phase 1: 最低限の対応（すぐに実装）

1. **ドキュメント更新**
   - ユーザー登録時のロール選択の意味を明確化
   - 「ショップオーナーとして登録 = ショップ管理のみ」を明記
   - 「カスタマーとして登録 = 購入のみ」を明記

2. **UI/UX改善**
   - 登録画面で「どちらで登録しますか？」を明確に
   - 後からロール変更できることを案内

#### Phase 2: 複数ロール対応（推奨）

1. **DB設計変更**
   - `user_roles`テーブル追加
   - マイグレーション実行

2. **Auth Service更新**
   - `AddUserRole` usecase追加
   - `RemoveUserRole` usecase追加
   - JWT Claims更新（rolesを配列化）

3. **各サービス更新**
   - JWTパース処理をrolesに対応
   - アクセス制御をrolesベースに変更

4. **Shop Service更新**
   - ショップ開設時にAuth Serviceへロール追加を依頼

#### Phase 3: 情報統合（オプション）

1. **Profile Service新設**
   - ユーザーの基本情報を統合管理
   - Customer/Shopから個人情報を分離

---

## 📝 回答: 設計チェック結果

### 質問: カスタマーとオーナーは認証を分けていますか？

**回答**: **分けていません（統一されています）**

**詳細**:
- ✅ **1つのAuth Serviceで統合管理**
- ✅ **Roleフィールドで区別**（CUSTOMER / SHOP_OWNER）
- ✅ **JWT認証は共通**

**この設計は適切か？**: **はい、基本的に適切です**

**ただし**:
- ⚠️ 1ユーザー1ロールの制約がある
- ⚠️ ショップオーナーが顧客として買い物できない
- ⚠️ 将来的に複数ロール対応が必要

**推奨**: Phase 2の複数ロール対応を実装

---

## 📚 参考: 設計パターン比較

### パターンA: 統合認証 + ロール区分（現在の設計）

**メリット**:
- ✅ 認証ロジックが統一
- ✅ 実装がシンプル
- ✅ パスワード管理が一元化

**デメリット**:
- ❌ 1ユーザー1ロールの制約（現状）
- ❌ ロール固有の情報管理が複雑

**採用例**: GitHub, Slack, Google

---

### パターンB: 認証分離（不採用）

**メリット**:
- ✅ カスタマーとオーナーで完全に分離
- ✅ 各ロール固有の認証フローが可能

**デメリット**:
- ❌ 認証ロジックの重複
- ❌ 1人が両方になる場合に複雑
- ❌ パスワード管理が煩雑

**採用例**: 一部の企業向けSaaS

---

## まとめ

**現在の設計は基本的に適切**ですが、**1ユーザー1ロールの制約**を解消するために、**複数ロール対応**の実装を推奨します。

**優先度**:
1. **High**: ドキュメント更新（すぐ）
2. **Medium**: 複数ロール対応（Phase 2）
3. **Low**: Profile Service新設（将来）
