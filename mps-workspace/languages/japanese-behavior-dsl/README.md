# 日本語振る舞いDSL

**お客さんも読める仕様書から実装コードまでを自動生成するDSL**

---

## 📖 概要

このDSLは、日本語で記述された振る舞い仕様から以下を自動生成します:

- ✅ Go実装コード（usecase, domain, handler）
- ✅ Protobuf定義（gRPC API）
- ✅ TypeScript/Reactコード（フロントエンド）
- ✅ フロー図（Mermaid）
- ✅ データ設計（ER図、SQL）
- ✅ お客さん向け仕様書（Markdown）

---

## 🚀 クイックスタート

### 前提条件

- JetBrains MPS 2023.3以降
- Java 11以降
- Go 1.21以降（生成コード確認用）

### 開始手順

1. **ドキュメントを読む**
   ```bash
   cd docs/dsl-design/

   # まずクイックスタートで基本を学ぶ（30分）
   cat quickstart.md

   # 次にステップバイステップで実装（1-2日）
   cat mps-step-by-step-guide.md
   ```

2. **MPSプロジェクトを開く**
   ```
   JetBrains MPSを起動
   → Open Project
   → このディレクトリを選択
   ```

3. **サンプルを確認**
   ```
   Solutions > sandbox > examples
   → user-registration.usecase を開く
   ```

4. **生成実行**
   ```
   Build > Make Project

   生成先: ../../../generated/
   ```

---

## 📁 ディレクトリ構造

```
japanese-behavior-dsl/
├── README.md                    # このファイル
├── models/
│   └── (MPS models will be here after project creation)
├── solutions/
│   └── sandbox/
│       └── examples/
│           └── user-registration.usecase
└── generator/
    └── template/
        ├── go-generator.mps
        ├── proto-generator.mps
        └── ts-generator.mps
```

---

## 🎯 DSL構文例

### 日本語仕様（.usecase）

```kotlin
ユースケース: ユーザー登録 {
  ID: "UC_AUTH_001"
  概要: "新規ユーザーを登録し、メール認証リンクを送信する"

  入力 {
    メールアドレス: 文字列 "ユーザーのメールアドレス"
    パスワード: 文字列 "ログイン用パスワード（8文字以上）"
  }

  出力 {
    ユーザーID: UUID "作成されたユーザーの一意識別子"
    認証トークン: 文字列 "JWT認証トークン"
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

  エラー定義 {
    メールアドレス形式不正 {
      コード: "AUTH_E001"
      HTTPステータス: 400
      メッセージ: "メールアドレスの形式が不正です"
    }
  }
}
```

---

## 🔧 生成されるコード

### Go実装

```go
// generated/auth/usecase/user_registration.go
package usecase

type UserRegistrationUsecase interface {
    Execute(ctx context.Context, input UserRegistrationInput) (UserRegistrationOutput, error)
}

func (u *userRegistrationUsecaseImpl) Execute(...) (..., error) {
    // ステップ1: 入力値の検証
    if !isValidEmail(input.Email) {
        return output, ErrInvalidEmailFormat
    }

    // ステップ2: 重複チェック
    existingUser, err := u.repo.FindByEmail(ctx, input.Email)
    if err == nil && existingUser != nil {
        return output, ErrEmailDuplicate
    }

    // ... 以下続く
}
```

### Protobuf

```protobuf
// generated/proto/auth/v1/auth.proto
service AuthService {
  rpc Register(RegisterRequest) returns (RegisterResponse);
}

message RegisterRequest {
  string email = 1;
  string password = 2;
}
```

### TypeScript

```typescript
// generated/frontend/src/api/auth.ts
export async function register(email: string, password: string): Promise<RegisterResponse> {
  const response = await fetch('/api/v1/auth/register', {
    method: 'POST',
    body: JSON.stringify({ email, password }),
  });
  return response.json();
}
```

---

## 📚 ドキュメント

| ドキュメント | 説明 | 所要時間 |
|------------|------|---------|
| `docs/dsl-design/quickstart.md` | 30分で動作プロトタイプを作成 | 30分 |
| `docs/dsl-design/mps-step-by-step-guide.md` | 詳細な実装手順 | 1-2日 |
| `docs/dsl-design/japanese-behavior-dsl-concept.md` | 設計思想・コンセプト | 30分 |
| `docs/dsl-design/generation-examples.md` | 生成コードの例 | 15分 |

---

## 🛠️ 開発ワークフロー

### 1. 仕様をDSLで記述

```kotlin
ユースケース: 新しい機能 {
  // 日本語で仕様を記述
}
```

### 2. 生成実行

```
Build > Make Project
```

### 3. 生成コードを確認

```bash
# Goコードのコンパイル確認
cd generated/auth
go build ./...

# TypeScript型チェック
cd generated/frontend
npm run type-check
```

### 4. カスタムロジック追加（必要な場合のみ）

```go
// manual/auth/custom_logic.go
// 複雑なバリデーションや外部API連携のみ
```

---

## 🧪 テスト

### 生成コードのテスト

```bash
# Go単体テスト
cd generated/auth
go test ./...

# TypeScriptテスト
cd generated/frontend
npm test
```

### DSL自体のテスト

```
Solutions > sandbox > tests
→ テストケースを実行
```

---

## 📊 期待される効果

| 指標 | 手動実装 | DSL駆動 | 削減率 |
|------|---------|---------|--------|
| コード量 | 1,000行 | 250行（DSL） | 75%削減 |
| 実装時間 | 8時間 | 2時間 | 75%短縮 |
| バグ発生率 | 中 | 低 | 60%削減 |
| 仕様変更対応 | 3時間 | 30分 | 83%短縮 |

---

## 🆘 トラブルシューティング

### Generatorが実行されない

```
Build > Clean
Build > Rebuild Language japanese-behavior-dsl
Build > Make Project
```

### 生成コードがコンパイルエラー

```
Build > Show Generation Plan
→ ログでエラー箇所を確認
→ Generatorテンプレートを修正
```

### 日本語が文字化け

```
File > Settings > Editor > File Encodings
→ すべてをUTF-8に設定
```

---

## 🤝 コントリビューション

改善案やバグ報告は以下へ:

1. GitHub Issues
2. Pull Request
3. ドキュメント修正

---

## 📝 ライセンス

このDSLはMITライセンスです。

---

## 🎓 学習リソース

- [JetBrains MPS公式ドキュメント](https://www.jetbrains.com/help/mps/)
- [MPS入門チュートリアル](https://www.jetbrains.com/help/mps/mps-tutorial.html)
- [Language Engineering with MPS](https://languageengineering.io/)

---

**次のステップ**: `docs/dsl-design/quickstart.md` から始めましょう！
