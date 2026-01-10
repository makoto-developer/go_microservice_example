# 自律実行ルール

このドキュメントは、Claudeが各Phaseを自律的に実行するためのルールを定義します。

---

## 基本原則

### Claude の自律実行権限

Claudeは以下の作業を**ユーザー承認なし**で自律実行できます：

#### ✅ 自動実行可能な作業

1. **ファイル読み込み**
   - `docs/requirements/*.md` - 要件定義
   - `mps-workspace/solutions/*/service.model` - 既存DSL定義
   - `docs/PROJECT_STATUS.md` - プロジェクト状況

2. **DSL定義作成**
   - `mps-workspace/solutions/<service>/service.model` 作成
   - エンティティ、ユースケース、gRPC定義

3. **コード生成実行**
   - `./scripts/mps-generate.sh <service>` 実行
   - 生成コードのディレクトリ構造作成

4. **ドキュメント更新**
   - `docs/PROJECT_STATUS.md` 更新
   - トークン消費記録

5. **進捗確認**
   - チェックリスト確認
   - 品質検証

#### ⚠️ ユーザー承認が必要な作業

1. **Git操作**
   - `git commit`
   - `git push`

2. **カスタムロジック実装**
   - `manual/` ディレクトリへのコード追加
   - 外部API連携実装

3. **Phase移行**
   - Phase 2 → Phase 3
   - Phase 3 → Phase 4

---

## Phase自律実行フロー

### Phase開始時の自動アクション

```
1. Phase情報表示
   ├─ 対象サービス一覧
   ├─ トークン見積
   └─ 依存関係確認

2. 並行実行可能サービスの特定
   ├─ 依存関係解析
   └─ 並行実行グループ作成

3. 実行開始確認
   └─ ユーザーに "Phase X を開始します。よろしいですか？" と確認
```

### サービス実装の自動アクション

各サービスごとに以下を自動実行：

```
1. 要件定義読み込み
   └─ cat docs/requirements/XX_service.md

2. DSL定義作成
   ├─ エンティティ抽出
   ├─ ユースケース定義
   ├─ gRPC定義
   └─ mps-workspace/solutions/<service>/service.model 作成

3. コード生成
   └─ ./scripts/mps-generate.sh <service>

4. 品質チェック
   ├─ DSL定義の行数確認（100-300行目標）
   ├─ コード生成成功確認
   └─ ディレクトリ構造確認

5. 進捗記録
   └─ docs/PROJECT_STATUS.md 更新
```

### Phase完了時の自動アクション

```
1. 完了サマリー表示
   ├─ 作成したサービス数
   ├─ トークン消費実績
   └─ 削減率計算

2. 次Phase提案
   └─ "Phase X が完了しました。Phase Y に進みますか？"
```

---

## 並行実行の自動制御

### 並行実行の判定ロジック

```python
def can_parallel_execute(service_a, service_b):
    """
    2つのサービスが並行実行可能か判定
    """
    dependencies_a = get_dependencies(service_a)
    dependencies_b = get_dependencies(service_b)

    # 相互依存チェック
    if service_b in dependencies_a:
        return False
    if service_a in dependencies_b:
        return False

    # 共通依存は問題なし
    return True
```

### 並行実行の実装例

#### Phase 2の並行実行

```
[Customer Service] と [Inventory Service] は並行実行可能

理由:
- Customer → Inventory 依存なし
- Inventory → Customer 依存なし
- 共通依存: Auth Service（完成済み）

実行:
1. Customer Service要件読み込み
2. Inventory Service要件読み込み
3. Customer ServiceのDSL定義作成
4. Inventory ServiceのDSL定義作成
5. 両方のコード生成を実行
```

---

## エラーハンドリング

### エラー発生時の自動対処

#### 1. DSL定義エラー

**検知**:
- DSL定義が400行を超過
- 必須セクション欠落（entity, usecase, grpc_service）

**自動対処**:
1. 警告メッセージ表示
2. 問題箇所の特定
3. 修正案の提示
4. ユーザー確認後に修正

#### 2. コード生成エラー

**検知**:
- `./scripts/mps-generate.sh` がエラー終了

**自動対処**:
1. エラーログ確認
2. DSL定義の見直し
3. 再実行

#### 3. トークン超過警告

**検知**:
- 消費トークンが見積の120%超過

**自動対処**:
1. 警告表示
2. トークン削減策の提案
   - 生成コードを読んでいないか確認
   - 同じファイルを複数回読んでいないか確認
3. 継続/中断の判断をユーザーに確認

---

## 品質保証の自動チェック

### DSL定義の自動検証

```
チェック項目:
1. ファイル存在確認
   └─ mps-workspace/solutions/<service>/service.model

2. 必須セクション確認
   ├─ microservice <ServiceName> { ... }
   ├─ entity定義が1つ以上
   ├─ usecase定義が1つ以上
   └─ grpc_service定義

3. 行数確認
   ├─ 最小: 50行
   ├─ 推奨: 100-300行
   └─ 最大: 400行（警告）

4. 命名規則確認
   ├─ エンティティ: PascalCase
   ├─ ユースケース: PascalCase
   └─ フィールド: snake_case
```

### コード生成の自動検証

```
チェック項目:
1. ディレクトリ構造
   ├─ generated/<service>/domain/
   ├─ generated/<service>/usecase/
   ├─ generated/<service>/handler/
   ├─ generated/<service>/infrastructure/
   └─ generated/<service>/tests/

2. go.modファイル
   └─ generated/<service>/go.mod

3. コンパイル確認（モック段階では省略可）
   └─ cd generated/<service> && go build ./...
```

---

## トークン最適化の自動制御

### 読み込みファイルの自動管理

```
ルール:
1. 同じファイルを2回読まない
   └─ 読み込み済みリストを保持

2. 生成コードは読まない
   ├─ generated/**/*.go は読まない
   └─ DSL定義のみ読む

3. 要件定義は機能要件のみ
   ├─ grep -A 50 "機能要件" を使用
   └─ 全文読み込みを避ける
```

### トークン消費の自動記録

```
各サービス完成時に記録:
- サービス名
- DSL定義行数
- 推定トークン消費
- 実績トークン消費

Phase完了時に集計:
- Phase合計トークン
- 見積との差分
- 削減率
```

---

## 進捗可視化

### Phase進捗の自動表示

各サービス完成時に表示：

```
## Phase 2進捗

| サービス | ステータス | DSL行数 | トークン |
|---------|----------|---------|---------|
| Customer | ✅ 完了 | 180 | 1,500 |
| Inventory | ✅ 完了 | 210 | 1,800 |
| Order | ⏳ 作業中 | - | - |
| Payment | ⏳ 未着手 | - | - |

**進捗**: 50% (2/4サービス)
**トークン**: 3,300 / 12,600 (26%)
```

---

## 自律実行の開始条件

### Phase 2開始の条件

```
前提条件:
1. Phase 1完了
   ├─ Auth Service DSL定義完了
   └─ Shop Service DSL定義完了

2. ユーザー承認
   └─ "Phase 2を開始します。よろしいですか？" → "はい"

3. トークン残量確認
   └─ 残りトークン >= 15,000
```

### Phase 3開始の条件

```
前提条件:
1. Phase 2完了
   ├─ Customer Service完了
   ├─ Inventory Service完了
   ├─ Order Service完了
   └─ Payment Service完了

2. ユーザー承認
   └─ "Phase 3を開始します。よろしいですか？" → "はい"

3. トークン残量確認
   └─ 残りトークン >= 10,000
```

---

## 自律実行のコマンド

### ユーザーが使用できるコマンド

#### Phase実行コマンド

```
"Phase 2を開始して"
→ Phase 2を自律的に実行開始

"Phase 2を並行実行して"
→ 並行可能なサービスを並行実行

"Phase 2の進捗を確認"
→ 現在の進捗を表示
```

#### 制御コマンド

```
"一時停止"
→ 現在の作業を完了後に停止

"次のサービスをスキップ"
→ 現在のサービスをスキップして次へ

"トークン消費を確認"
→ 現在までのトークン消費を表示
```

---

## 自律実行の制約

### 実行しない作業

1. **コードの手動実装**
   - `manual/` へのコード追加はユーザー承認必要

2. **DSL言語定義の変更**
   - `mps-workspace/languages/` の変更は対象外

3. **Generator実装**
   - MPS Generatorのコード変更は対象外

4. **インフラ構築**
   - Docker, Kubernetes設定は対象外

---

## まとめ

### 自律実行の範囲

Claudeは以下を自律的に実行：

1. **要件定義の読み込み**
2. **DSL定義の作成**
3. **コード生成の実行**
4. **品質チェック**
5. **進捗記録**
6. **並行実行の制御**

### ユーザーの役割

ユーザーは以下を担当：

1. **Phase開始の承認**
2. **カスタムロジックの実装判断**
3. **Git操作の承認**
4. **最終的な品質確認**

この役割分担により、効率的かつ安全に開発を進めることができます。
