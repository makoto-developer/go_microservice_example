# 日本語振る舞いDSL コンセプト設計書

**作成日**: 2026-01-31
**対象**: 開発者、お客さん
**目的**: お客さんも読める仕様書から実装コードまでを自動生成

---

## 1. 設計思想

### 1.1 基本コンセプト

```
お客さんが読める仕様
    ↓ (自動変換)
開発者が理解できる技術仕様
    ↓ (自動生成)
動作する実装コード
```

### 1.2 ターゲット読者層

| 読者 | 読むもの | 理解すべきこと |
|------|---------|--------------|
| **お客さん** | `.usecase` ファイル | ビジネスロジック、条件分岐、エラー処理 |
| **開発者** | `.model` ファイル | エンティティ、API、データ構造 |
| **システム** | 生成コード | 実装の詳細 |

---

## 2. DSL構文設計

### 2.1 キーワード体系

#### 日本語キーワード（お客さん向け）

```
構造系:
  - ユースケース
  - 入力
  - 出力
  - 処理フロー
  - ステップN
  - エラー定義

制御系:
  - もし〜場合
  - そうでなければ
  - 繰り返し
  - それぞれについて

データ操作系:
  - データベースから検索
  - データベースに保存
  - 作成する
  - 更新する
  - 削除する

検証系:
  - 検証項目
  - 前提条件
  - 期待結果
```

#### 技術キーワード（開発者向け）

```
構造系:
  - microservice
  - entity
  - usecase
  - grpc_service

型系:
  - string, int, boolean, UUID
  - timestamp, decimal
  - list, map

制約系:
  - primary_key, foreign_key
  - unique, not_null
  - min_length, max_length
```

### 2.2 構文例

#### 日本語仕様（.usecase）

```kotlin
ユースケース: ユーザー登録 {
  概要: "新規ユーザーを登録し、メール認証リンクを送信する"

  入力 {
    メールアドレス: 文字列
    パスワード: 文字列
  }

  出力 {
    ユーザーID: UUID
    認証トークン: 文字列
  }

  処理フロー {
    ステップ1: "入力値の検証" {
      もし メールアドレスが正しい形式でない場合:
        → エラー「メールアドレスの形式が不正です」を返す
    }

    ステップ2: "重複チェック" {
      既存ユーザー = データベースから メールアドレスで検索する

      もし 既存ユーザーが見つかった場合:
        → エラー「このメールアドレスは既に登録されています」を返す
    }
  }
}
```

#### 技術仕様（.model）- 既存

```kotlin
microservice AuthService {
  entity User {
    id: UUID primary_key
    email: string unique not_null
  }

  usecase UserRegistration {
    input: { email: string, password: string }
    output: { user_id: UUID, token: string }
  }
}
```

---

## 3. 生成物の種類

### 3.1 コード生成

| 生成物 | 形式 | 生成元 |
|--------|------|--------|
| Go実装 | `.go` | `.usecase` + `.model` |
| Protobuf | `.proto` | `.usecase` + `.model` |
| TypeScript | `.ts`, `.tsx` | `.usecase` |
| テストコード | `_test.go` | `.usecase` のテストケース |

### 3.2 ドキュメント生成

| ドキュメント | 形式 | 対象読者 |
|------------|------|---------|
| フロー図 | Mermaid | 開発者、お客さん |
| ER図 | Mermaid | 開発者、DBA |
| 仕様書 | Markdown | お客さん |
| API仕様 | OpenAPI | 開発者 |

---

## 4. 型システム

### 4.1 基本型の対応表

| 日本語表記 | 技術仕様 | Go型 | Protobuf型 | TypeScript型 |
|-----------|---------|------|-----------|-------------|
| 文字列 | string | string | string | string |
| 整数 | int | int64 | int64 | number |
| 真偽値 | boolean | bool | bool | boolean |
| UUID | UUID | uuid.UUID | string | string |
| 日時 | timestamp | time.Time | google.protobuf.Timestamp | Date |
| 小数 | decimal | decimal.Decimal | string | number |

### 4.2 複合型

| 日本語表記 | 技術仕様 | Go型 | TypeScript型 |
|-----------|---------|------|-------------|
| リスト | list<T> | []T | T[] |
| マップ | map<K,V> | map[K]V | Record<K,V> |
| 列挙型 | enum | type T string | enum |

---

## 5. 制御構造

### 5.1 条件分岐

#### 日本語DSL

```kotlin
もし ユーザーが管理者である場合:
  管理者権限を付与する
そうでなければ もし ユーザーがショップオーナーである場合:
  ショップ管理権限を付与する
そうでなければ:
  顧客権限を付与する
```

#### 生成されるGoコード

```go
if user.Role == RoleAdmin {
    grantAdminPermissions()
} else if user.Role == RoleShopOwner {
    grantShopOwnerPermissions()
} else {
    grantCustomerPermissions()
}
```

### 5.2 繰り返し

#### 日本語DSL

```kotlin
注文の各商品について繰り返し {
  在庫を引き当てる(商品ID, 数量)
}
```

#### 生成されるGoコード

```go
for _, item := range order.Items {
    reserveStock(item.ProductID, item.Quantity)
}
```

### 5.3 トランザクション

#### 日本語DSL

```kotlin
トランザクション開始 {
  注文を保存する
  在庫を減らす
  決済を記録する
} トランザクション確定

エラー時の処理 {
  トランザクションをロールバックする
  エラー「注文処理に失敗しました」を返す
}
```

#### 生成されるGoコード

```go
tx, err := db.BeginTx(ctx)
if err != nil {
    return err
}
defer tx.Rollback()

if err := saveOrder(tx, order); err != nil {
    return fmt.Errorf("注文処理に失敗しました: %w", err)
}

if err := decreaseStock(tx, productID, quantity); err != nil {
    return fmt.Errorf("注文処理に失敗しました: %w", err)
}

if err := recordPayment(tx, payment); err != nil {
    return fmt.Errorf("注文処理に失敗しました: %w", err)
}

if err := tx.Commit(); err != nil {
    return fmt.Errorf("注文処理に失敗しました: %w", err)
}
```

---

## 6. データベース操作の抽象化

### 6.1 検索操作

| 日本語表記 | SQL | Go（Repository） |
|-----------|-----|-----------------|
| データベースから IDで検索 | SELECT * FROM ... WHERE id = ? | repo.FindByID(ctx, id) |
| データベースから メールアドレスで検索 | SELECT * FROM ... WHERE email = ? | repo.FindByEmail(ctx, email) |
| データベースから すべて取得 | SELECT * FROM ... | repo.FindAll(ctx) |

### 6.2 更新操作

| 日本語表記 | SQL | Go（Repository） |
|-----------|-----|-----------------|
| データベースに保存する | INSERT INTO ... | repo.Save(ctx, entity) |
| データベースを更新する | UPDATE ... | repo.Update(ctx, entity) |
| データベースから削除する | DELETE FROM ... | repo.Delete(ctx, id) |

---

## 7. エラーハンドリング

### 7.1 エラー定義

#### 日本語DSL

```kotlin
エラー定義 {
  メールアドレス重複 {
    コード: "AUTH_E003"
    HTTPステータス: 409
    メッセージ: "このメールアドレスは既に登録されています"
  }
}
```

#### 生成されるGoコード

```go
type EmailDuplicateError struct {
    Code       string
    Message    string
    HTTPStatus int
}

func (e *EmailDuplicateError) Error() string {
    return e.Message
}

var ErrEmailDuplicate = &EmailDuplicateError{
    Code:       "AUTH_E003",
    Message:    "このメールアドレスは既に登録されています",
    HTTPStatus: 409,
}
```

### 7.2 エラー処理フロー

```kotlin
もし エラーが発生した場合:
  エラーログに記録する
  エラー「処理に失敗しました」を返す
  処理を終了する
```

---

## 8. イベント駆動

### 8.1 イベント発行

#### 日本語DSL

```kotlin
非同期で実行する {
  ユーザー登録完了イベントを発行する {
    ユーザーID: 新規ユーザーのID
    メールアドレス: 新規ユーザーのメールアドレス
  }
}
```

#### 生成されるGoコード

```go
go func() {
    eventBus.Publish(UserRegisteredEvent{
        UserID: newUser.ID,
        Email:  newUser.Email,
    })
}()
```

---

## 9. テストケースの生成

### 9.1 テストシナリオ定義

#### 日本語DSL

```kotlin
テストケース {
  正常系: "新規ユーザーの登録成功" {
    前提 {
      メールアドレス「test@example.com」が未登録である
    }

    実行 {
      入力値 {
        メールアドレス: "test@example.com"
        パスワード: "SecurePass123!"
      }
    }

    期待結果 {
      ユーザーIDが返される
      認証トークンが返される
      データベースにユーザーが保存されている
    }
  }
}
```

#### 生成されるGoテストコード

```go
func TestUserRegistration_Success(t *testing.T) {
    // 前提
    ctx := context.Background()
    repo := setupTestRepository(t)
    usecase := NewUserRegistrationUsecase(repo)

    // 実行
    output, err := usecase.Execute(ctx, UserRegistrationInput{
        Email:    "test@example.com",
        Password: "SecurePass123!",
    })

    // 期待結果
    require.NoError(t, err)
    assert.NotEmpty(t, output.UserID)
    assert.NotEmpty(t, output.Token)

    // データベース確認
    savedUser, err := repo.FindByID(ctx, output.UserID)
    require.NoError(t, err)
    assert.Equal(t, "test@example.com", savedUser.Email)
}
```

---

## 10. 実装優先順位

### Phase 1: 基本構文（必須）

- [ ] ユースケース定義
- [ ] 入力・出力定義
- [ ] 処理フロー
- [ ] 条件分岐（もし〜場合）
- [ ] エラー定義

### Phase 2: データベース操作

- [ ] データベース検索
- [ ] データベース保存
- [ ] トランザクション

### Phase 3: 高度な機能

- [ ] 繰り返し処理
- [ ] イベント発行
- [ ] 非同期処理

### Phase 4: ドキュメント生成

- [ ] フロー図生成（Mermaid）
- [ ] ER図生成
- [ ] 仕様書生成（Markdown）

### Phase 5: フロントエンド生成

- [ ] TypeScript型定義
- [ ] Reactコンポーネント
- [ ] API呼び出し関数

---

## 11. 次のステップ

1. **MPS言語プロジェクトの作成**
   - JetBrains MPSで新規言語プロジェクト作成
   - 言語名: `japanese-behavior-dsl`

2. **構文構造の定義（Structure）**
   - コンセプト定義
   - プロパティ定義
   - 参照定義

3. **エディタの定義（Editor）**
   - 視覚的な表現
   - 入力支援

4. **Generatorの実装**
   - Go生成ロジック
   - Protobuf生成ロジック
   - TypeScript生成ロジック

5. **サンプル実装・検証**
   - ユーザー登録ユースケースで検証
   - 生成コードのテスト

---

この設計に基づいて、次のPhaseでMPS言語実装を開始します。
