# トークン最適化ルール

このドキュメントでは、Claude のトークン消費を最小限に抑えるためのルールを定義します。

---

## 🎯 目標

**12個のマイクロサービスを効率的に開発するため、トークン消費を90%削減**

---

## 📊 トークン消費の現実

### 従来の開発（手動実装）

| サービス | コード行数 | トークン消費 |
|---------|----------|------------|
| 1サービス | 2,000-3,000行 | ~15,000 |
| 12サービス | 24,000-36,000行 | ~180,000 |

**問題**: 1セッションで2-3サービスしか開発できない

---

### MPS DSL駆動開発

| サービス | DSL定義 | 生成コード | Claudeが読む | トークン消費 |
|---------|---------|-----------|------------|------------|
| 1サービス | 100-200行 | 2,000-3,000行 | DSLのみ | ~1,500 |
| 12サービス | 1,200-2,400行 | 24,000-36,000行 | DSLのみ | ~18,000 |

**効果**: **90%削減**、1セッションで全サービス開発可能

---

## ✅ トークン最適化のベストプラクティス

### 1. 生成コードを読まない

```bash
# ❌ トークン無駄遣い（~15,000トークン）
cat generated/auth/domain/user.go
cat generated/auth/usecase/user_registration.go
cat generated/auth/handler/grpc_handler.go
# ... 合計2,000-3,000行

# ✅ トークン節約（~1,500トークン）
cat mps-workspace/solutions/auth-service/service.model
# 100-200行のみ
```

**削減率**: **90%削減**

---

### 2. DSL定義を簡潔に保つ

#### ✅ 良い例（簡潔）

```kotlin
microservice AuthService {
  version: "v1"

  entity User {
    id: UUID primary_key
    email: string unique not_null
    role: Role not_null
  }

  enum Role { CUSTOMER, SHOP_OWNER, ADMIN }

  usecase UserRegistration {
    input: { email: string, password: string }
    output: { user_id: UUID, token: string }
    errors: { EmailAlreadyExists }
  }
}
```

**行数**: ~50行
**トークン**: ~300

---

#### ❌ 悪い例（冗長）

```kotlin
microservice AuthService {
  // これは認証サービスです
  // ユーザー登録、ログイン、JWT発行を行います
  version: "v1"

  // ユーザーエンティティ
  // データベースのusersテーブルに対応します
  entity User {
    // 主キー
    id: UUID primary_key
    // メールアドレス、ユニーク制約あり
    email: string unique not_null
    // ... 冗長なコメント
  }
  // ... 以下同様
}
```

**行数**: ~200行
**トークン**: ~1,200

**問題**: コメントが冗長、トークン4倍消費

---

### 3. 要件定義を効率的に読む

#### ✅ 良い例（必要な部分のみ）

```bash
# 機能要件のみ確認
grep -A 20 "## 機能要件" docs/requirements/01_auth_service.md
```

**トークン**: ~500-700

---

#### ❌ 悪い例（全文読む）

```bash
# 全文読む（非機能要件、API設計、DB設計等も含む）
cat docs/requirements/01_auth_service.md
```

**トークン**: ~2,000-3,000

**問題**: 不要な情報も読んでしまう

**注**: このプロジェクトでは要件定義を簡潔化済み（非機能要件等を削除）

---

### 4. 並行開発の活用

#### ✅ 良い例（並行開発）

```bash
# 独立したサービスを並行開発
- Task 1: Auth Service定義
- Task 2: Shop Service定義
- Task 3: Inventory Service定義

# トークン消費: 1,500 × 3 = 4,500（並行実行）
```

---

#### ❌ 悪い例（逐次開発）

```bash
# 1つずつ開発
1. Auth Service完成 → 次へ
2. Shop Service完成 → 次へ
3. Inventory Service完成

# トークン消費: 同じだが時間がかかる
```

---

## 📏 トークン消費の測定

### ファイルサイズとトークンの目安

| ファイル | 行数 | トークン（目安） |
|---------|------|--------------|
| DSL定義 | 100行 | ~600 |
| DSL定義 | 200行 | ~1,200 |
| Go生成コード | 1,000行 | ~6,000 |
| Go生成コード | 3,000行 | ~18,000 |
| 要件定義 | 150行 | ~900 |

**計算式**: 行数 × 6 ≈ トークン数（目安）

---

## 🚫 トークン無駄遣いのパターン

### 1. 生成コードを読む

```bash
❌ cat generated/auth/domain/user.go
❌ cat generated/auth/usecase/user_registration.go
❌ cat generated/auth/handler/grpc_handler.go

トークン消費: ~15,000
```

**対策**: DSL定義のみ読む

---

### 2. 同じファイルを複数回読む

```bash
❌ 1回目: cat mps-workspace/solutions/auth-service/service.model
❌ 2回目: cat mps-workspace/solutions/auth-service/service.model
❌ 3回目: cat mps-workspace/solutions/auth-service/service.model

トークン消費: ~4,500（3回分）
```

**対策**: 最初に読んだ内容を記憶する

---

### 3. 不要なドキュメントを読む

```bash
❌ cat README.md（プロジェクト概要は既知）
❌ cat SETUP.md（環境構築は完了済み）
❌ cat docs/architecture.md（今回不要）

トークン消費: ~5,000
```

**対策**: 必要なドキュメントのみ読む

---

## ✅ トークン節約のテクニック

### 1. ファイルの部分読み込み

```bash
# ✅ 必要な部分のみ
head -n 50 docs/requirements/01_auth_service.md
grep -A 10 "機能要件" docs/requirements/01_auth_service.md

# ❌ 全文読み込み
cat docs/requirements/01_auth_service.md
```

---

### 2. 要約の活用

```bash
# ✅ 要約を作成・保存（最初の1回のみ）
# 要件定義を読む → 要約作成 → Serena Memoryに保存

# 2回目以降は要約のみ参照
serena.read_memory("auth-service-summary.md")
```

**トークン削減**: 70-80%

---

### 3. 差分のみ確認

```bash
# ✅ 変更箇所のみ確認
git diff mps-workspace/solutions/auth-service/service.model

# ❌ 全ファイル再読み込み
cat mps-workspace/solutions/auth-service/service.model
```

---

## 📊 プロジェクト全体のトークン消費見積もり

### Phase 1: 基盤サービス（Auth, Shop）

| タスク | トークン |
|--------|---------|
| 要件定義読み込み | 1,800 |
| DSL定義作成 | 3,000 |
| 生成コード確認（コンパイルのみ） | 500 |
| カスタムロジック実装 | 1,000 |
| **合計** | **6,300** |

---

### Phase 2: コアサービス（Customer, Inventory, Order, Payment）

| タスク | トークン |
|--------|---------|
| 要件定義読み込み | 3,600 |
| DSL定義作成 | 6,000 |
| 生成コード確認 | 1,000 |
| カスタムロジック実装 | 2,000 |
| **合計** | **12,600** |

---

### Phase 3-4: 残りのサービス（6サービス）

| タスク | トークン |
|--------|---------|
| 要件定義読み込み | 5,400 |
| DSL定義作成 | 9,000 |
| 生成コード確認 | 1,500 |
| カスタムロジック実装 | 3,000 |
| **合計** | **18,900** |

---

### プロジェクト全体（12サービス）

| Phase | トークン |
|-------|---------|
| Phase 1 | 6,300 |
| Phase 2 | 12,600 |
| Phase 3-4 | 18,900 |
| **合計** | **37,800** |

**1セッション（200,000トークン）で全サービス開発可能！**

---

## まとめ

### トークン最適化の3原則

1. **生成コードは読まない** → DSL定義のみ
2. **要件定義は必要な部分のみ** → 機能要件に絞る
3. **並行開発を活用** → 独立サービスは同時に

これらを守ることで、**90%のトークン削減**を実現できます。
