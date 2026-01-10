# プロジェクトローカル Claude 設定

このファイルは、このプロジェクト専用のClaude動作ルールを定義します。

## プロジェクト概要

- **プロジェクト名**: Go MicroService 実践例（オンラインショップモール）
- **開発手法**: JetBrains MPS DSL駆動開発
- **言語**: Go 1.25
- **アーキテクチャ**: マイクロサービス（12サービス）
- **最重要課題**: トークン消費の最適化

---

## 🎯 このプロジェクトの最重要ルール

### 1. **MPS DSL優先の原則**

**すべての開発はDSL定義から始める**

```
❌ NG: いきなりGoコードを書く
✅ OK: DSL定義 → MPS生成 → カスタムロジック実装
```

**ワークフロー**:
1. `docs/requirements/` で要件を確認
2. `mps-workspace/solutions/` でDSL定義
3. `./scripts/mps-generate.sh` でコード生成
4. `manual/` でカスタムロジック実装（必要な場合のみ）

---

### 2. **生成コード不可侵の原則**

**`generated/` ディレクトリは絶対に触らない**

```
❌ 禁止: generated/ 配下のファイルを編集
✅ 許可: DSL定義を修正 → 再生成
✅ 許可: manual/ にカスタムロジック追加
```

**理由**: 再生成で上書きされるため

---

### 3. **トークン最適化の原則**

**Claudeは生成コードを読まない、DSL定義のみ読む**

| 読むもの | トークン消費 |
|---------|------------|
| ❌ 生成コード（2,000-3,000行） | ~15,000 |
| ✅ DSL定義（100-200行） | ~1,500 |

**削減率**: **90%削減**

---

## 📂 ディレクトリ構成の理解

```
mps-workspace/
├── languages/           # DSL言語定義（Claudeが編集可能）
│   ├── microservice-dsl/
│   ├── grpc-dsl/
│   └── event-dsl/
└── solutions/           # サービス定義（Claudeが編集可能）
    ├── auth-service/
    └── ...

generated/               # 生成コード（読み取り専用、編集禁止）
├── auth/
└── ...

manual/                  # カスタムロジック（Claudeが編集可能）
└── custom/

docs/requirements/       # 要件定義（必ず参照）
├── 01_auth_service.md
└── ...
```

---

## 🔄 標準的な作業フロー

### 新しいサービスを追加する場合

1. **要件確認**
   ```bash
   # 該当する要件定義を読む
   cat docs/requirements/XX_service_name.md
   ```

2. **DSL定義作成**
   ```
   mps-workspace/solutions/service-name/
   └── service.model  # ここにDSL定義
   ```

3. **コード生成**
   ```bash
   ./scripts/mps-generate.sh service-name
   ```

4. **生成結果確認**
   ```bash
   ls -la generated/service-name/
   ```

5. **カスタムロジック実装（必要な場合のみ）**
   ```
   manual/service-name/
   └── custom_logic.go
   ```

---

### 既存サービスに機能追加する場合

1. **既存DSL定義を読む**
   ```
   mps-workspace/solutions/service-name/service.model
   ```

2. **DSL定義を更新**
   - entity追加
   - usecase追加
   - grpc_service更新

3. **再生成**
   ```bash
   ./scripts/mps-generate.sh service-name
   ```

4. **カスタムロジック追加**（複雑な処理のみ）
   ```
   manual/service-name/
   ```

---

## 📖 Claudeが参照すべきドキュメント

### 必ず参照
- `docs/requirements/README.md` - 要件定義の入り口
- `docs/requirements/XX_service.md` - 各サービスの要件
- `CLAUDE.md`（プロジェクトルート） - 開発ガイド全体

### 必要に応じて参照
- `README.md` - プロジェクト概要・技術スタック
- `SETUP.md` - 環境構築・MPS使用方法
- `.claude/rules/*.md` - 詳細ルール

---

## 🚫 禁止事項

### 1. 生成コードの編集
```bash
❌ vim generated/auth/domain/user.go  # 絶対NG
✅ vim mps-workspace/solutions/auth-service/service.model  # OK
```

### 2. DSLを経由しない実装
```
❌ いきなりmanual/にコード書く
✅ まずDSL定義 → 足りない部分をmanual/に実装
```

### 3. トークン無駄遣い
```
❌ generated/ 配下のファイルを読む
✅ DSL定義のみ読む
```

---

## ✅ 推奨事項

### 1. DSL定義の簡潔さ
- 1サービス = 100-200行以内
- 明確な命名（UserRegistration, UserLogin等）
- 必要最小限の定義

### 2. カスタムロジックの限定
- 複雑なバリデーション
- 外部API連携（Stripe等）
- ビジネスルール実装

これらのみ `manual/` に実装

### 3. 要件定義の活用
- 実装前に必ず `docs/requirements/` を読む
- 要件に書かれていることのみ実装
- 不明点はユーザーに確認

---

## 📝 詳細ルール

詳細は `.claude/rules/` を参照：

- [mps-workflow.md](./.claude/rules/mps-workflow.md) - MPS開発フロー
- [code-generation.md](./.claude/rules/code-generation.md) - コード生成ルール
- [token-optimization.md](./.claude/rules/token-optimization.md) - トークン最適化

---

## 🎓 まとめ

### このプロジェクトでClaude が守るべき3原則

1. **DSL First**: すべてはDSL定義から始める
2. **Read Only Generated**: 生成コードは読まない・触らない
3. **Token Efficiency**: 常にトークン消費を意識

これらを守ることで、12個のマイクロサービスを効率的に開発できます。
