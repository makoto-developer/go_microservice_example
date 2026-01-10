# .claude ディレクトリ

このディレクトリには、このプロジェクト専用のClaude動作ルールが定義されています。

## ファイル構成

```
.claude/
├── CLAUDE.md                       # プロジェクト固有設定（メイン）
├── README.md                       # このファイル
└── rules/                          # 詳細ルール
    ├── mps-workflow.md             # MPS開発フロー
    ├── code-generation.md          # コード生成ルール
    ├── token-optimization.md       # トークン最適化
    ├── phase-execution-plan.md     # Phase別実行計画
    └── autonomous-execution.md     # 自律実行ルール
```

## ファイル説明

### CLAUDE.md（319行）
プロジェクト固有のClaude動作ルールのメインファイル。

**内容**:
- プロジェクト概要（MPS DSL駆動開発）
- 3つの最重要ルール（DSL優先、生成コード不可侵、トークン最適化）
- ディレクトリ構成
- 標準的な作業フロー
- 禁止事項・推奨事項
- Claudeの自律実行権限

### rules/mps-workflow.md（306行）
MPS開発フローの詳細定義。

**内容**:
- 基本フロー（要件定義 → DSL設計 → コード生成 → 確認 → カスタムロジック → テスト）
- 6段階のPhase定義
- DSL定義のベストプラクティス
- トラブルシューティング

### rules/code-generation.md（328行）
コード生成ルールとマッピング定義。

**内容**:
- 生成コードは読み取り専用
- DSL → Go 生成マッピング
  - Entity → Go Struct
  - Enum → Go Type
  - Usecase → Interface + Implementation
  - gRPC Service → Handler
- カスタムロジックの統合方法

### rules/token-optimization.md（330行）
トークン消費最適化の戦略と測定方法。

**内容**:
- トークン消費の現実（従来 vs MPS DSL）
- ベストプラクティス（生成コードを読まない、DSL定義を簡潔に、並行開発）
- トークン無駄遣いのパターン
- プロジェクト全体のトークン消費見積もり（37,800トークン）
- 90%削減の実現方法

### rules/phase-execution-plan.md（新規作成）
Phase 2-4の詳細実行計画。

**内容**:
- 全体構成（12サービスを4つのPhaseに分割）
- Phase 2: Customer, Inventory, Order, Payment（12,600トークン）
- Phase 3: Shipping, Notification, Review（7,500トークン）
- Phase 4: Chat, Search, Admin（10,500トークン）
- サービス依存関係と並行実行戦略
- 各サービスのDSL定義内容見積もり
- 品質チェックリスト

### rules/autonomous-execution.md（新規作成）
Claudeの自律実行ルール。

**内容**:
- 自律実行可能な作業（要件読み込み、DSL作成、コード生成、品質チェック、進捗記録）
- ユーザー承認が必要な作業（Git操作、カスタムロジック実装、Phase移行）
- 並行実行の自動制御ロジック
- エラーハンドリング
- 品質保証の自動チェック
- トークン最適化の自動制御
- Phase実行コマンド

## 使い方

### プロジェクト開始時

1. まず `CLAUDE.md` を読む
2. 必要に応じて `rules/` 配下の詳細ルールを参照

### Phase実行時

1. `rules/phase-execution-plan.md` で実行計画を確認
2. `rules/autonomous-execution.md` で自律実行ルールを確認
3. Phaseを開始（Claudeが自律実行）

### トラブル発生時

1. `rules/mps-workflow.md` のトラブルシューティングを確認
2. `rules/token-optimization.md` でトークン消費を見直し

## トークン消費の目安

| ファイル | 行数 | 読み込み頻度 | トークン |
|---------|------|------------|---------|
| CLAUDE.md | 319 | 初回のみ | ~2,000 |
| mps-workflow.md | 306 | 必要時 | ~1,800 |
| code-generation.md | 328 | 必要時 | ~2,000 |
| token-optimization.md | 330 | 初回のみ | ~2,000 |
| phase-execution-plan.md | ~500 | Phase開始時 | ~3,000 |
| autonomous-execution.md | ~600 | Phase開始時 | ~3,600 |

**合計**: 約2,400行、初回読み込み時のトークン消費は約7,000-10,000トークン

## 更新履歴

### 2026-01-10
- 初回作成
- CLAUDE.md作成（319行）
- rules/mps-workflow.md作成（306行）
- rules/code-generation.md作成（328行）
- rules/token-optimization.md作成（330行）
- rules/phase-execution-plan.md作成（約500行）
- rules/autonomous-execution.md作成（約600行）

## まとめ

このディレクトリのルールに従うことで、Claudeは：

1. **DSL優先の開発**を徹底
2. **生成コードを読まず**トークンを削減
3. **Phase 2-4を自律的に実行**
4. **並行実行を自動制御**して効率化

結果として、**12サービスを37,800トークン（90%削減）で開発可能**になります。
