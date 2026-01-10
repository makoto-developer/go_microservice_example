# MPS Workspace

このディレクトリには、JetBrains MPSによるDSL定義とコード生成に関連するファイルが含まれています。

## ディレクトリ構成

```
mps-workspace/
├── languages/           # DSL言語定義
│   ├── microservice-dsl/  # マイクロサービスDSL
│   ├── grpc-dsl/          # gRPC サービスDSL
│   └── event-dsl/         # イベント駆動DSL
└── solutions/           # DSLを使ったサービス定義
    ├── auth-service/    # 認証認可サービス定義
    │   └── service.model
    ├── shop-service/    # ショップサービス定義
    │   └── service.model
    └── ...
```

## DSL定義ファイル (.model)

各サービスの `service.model` ファイルには、以下の要素が含まれます：

### 1. エンティティ (Entities)
データベースのテーブルに対応するドメインモデル

```kotlin
entity User {
  id: UUID primary_key
  email: string unique not_null
  password_hash: string not_null
  role: Role not_null
}
```

### 2. 列挙型 (Enums)
定数値の定義

```kotlin
enum Role {
  CUSTOMER,
  SHOP_OWNER,
  ADMIN
}
```

### 3. ユースケース (Use Cases)
ビジネスロジックの定義

```kotlin
usecase UserRegistration {
  input: { email: string, password: string }
  output: { user_id: UUID, token: string }
  errors: { EmailAlreadyExists, WeakPassword }
}
```

### 4. gRPCサービス
API エンドポイントの定義

```kotlin
grpc_service {
  rpc Register(RegisterRequest) returns (RegisterResponse)
  rpc Login(LoginRequest) returns (LoginResponse)
}
```

### 5. 依存関係
外部サービスの宣言

```kotlin
dependencies {
  database: PostgreSQL
  cache: Redis
  messaging: RabbitMQ
}
```

### 6. イベント（オプション）
イベント駆動アーキテクチャで使用

```kotlin
events {
  publish UserRegistered { user_id: UUID, email: string }
  subscribe OrderCreated { order_id: UUID }
}
```

## コード生成

### 生成コマンド

```bash
# 特定のサービス
./scripts/mps-generate.sh auth-service

# すべてのサービス
./scripts/mps-generate.sh --all
```

### 生成されるコード構成

```
generated/
└── auth-service/
    ├── domain/
    │   ├── user.go           # エンティティ
    │   ├── role.go           # Enum
    │   └── repository.go     # Repositoryインターフェース
    ├── usecase/
    │   ├── user_registration.go  # ユースケース実装
    │   └── interfaces.go         # 依存インターフェース
    ├── handler/
    │   └── grpc_handler.go   # gRPCハンドラー
    └── infrastructure/
        ├── postgres_repo.go  # Repository実装
        └── redis_cache.go    # Redisキャッシュ
```

## DSL定義のベストプラクティス

### ✅ 良い例

1. **明確な命名**
   ```kotlin
   usecase UserRegistration { ... }
   ```

2. **適切な型指定**
   ```kotlin
   entity User {
     id: UUID primary_key
     email: string unique not_null
   }
   ```

3. **簡潔な定義**
   - 1サービス = 100-300行
   - 必要な要素のみ定義

### ❌ 悪い例

1. **曖昧な命名**
   ```kotlin
   usecase DoStuff { ... }
   ```

2. **型が不明確**
   ```kotlin
   entity Thing {
     id: string
     data: string
   }
   ```

3. **冗長なコメント**
   ```kotlin
   // これはユーザーです
   // データベースのusersテーブルに対応します
   // ... (不要なコメントが多い)
   entity User { ... }
   ```

## トークン最適化

### DSL駆動開発のメリット

| アプローチ | Claudeが読む量 | トークン消費 | 削減率 |
|----------|--------------|------------|--------|
| 手動実装 | 2,000-3,000行のGoコード | ~15,000 | - |
| MPS DSL | 100-300行のDSL定義 | ~1,500 | **90%削減** |

### 12サービス全体のトークン消費見積もり

| Phase | サービス数 | トークン消費 |
|-------|----------|------------|
| Phase 1 | 2 | 6,300 |
| Phase 2 | 4 | 12,600 |
| Phase 3-4 | 6 | 18,900 |
| **合計** | **12** | **37,800** |

**1セッション（200,000トークン）で全サービス開発可能！**

## 開発フロー

```
1. 要件定義確認
   ↓
2. DSL定義作成（このディレクトリ）
   ↓
3. コード生成
   ↓
4. 生成コード確認
   ↓
5. カスタムロジック実装（必要な場合のみ）
```

## 注意事項

### ⚠️ 生成コードは編集禁止

`generated/` ディレクトリ配下のファイルは**絶対に編集しない**

- 理由: 再生成で上書きされる
- 変更が必要な場合: DSL定義を修正 → 再生成

### ⚠️ カスタムロジックの実装場所

生成コードで表現できない複雑なロジックのみ `manual/` ディレクトリに実装

## 参考ドキュメント

- [プロジェクトルート CLAUDE.md](../CLAUDE.md) - 開発ガイド
- [.claude/rules/mps-workflow.md](../.claude/rules/mps-workflow.md) - MPS開発フロー
- [.claude/rules/code-generation.md](../.claude/rules/code-generation.md) - コード生成ルール
- [.claude/rules/token-optimization.md](../.claude/rules/token-optimization.md) - トークン最適化
