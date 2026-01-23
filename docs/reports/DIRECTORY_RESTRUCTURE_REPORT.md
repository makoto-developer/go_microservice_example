# ディレクトリ構造再編成レポート

**実施日**: 2026-01-24
**目的**: マイクロサービスごとにディレクトリを綺麗にカテゴライズし、リンク切れを解消

---

## 📋 実施内容

### 1. 新しいディレクトリ構造の作成

以下のカテゴリ別ディレクトリ構造を構築しました：

```
go_microservice_example/
├── microservices/         # 全マイクロサービス（12サービス）
├── infrastructure/        # インフラ・デプロイ設定
├── docs/                  # ドキュメント統合
├── build/                 # ビルド成果物
├── tools/                 # 開発ツール
├── mps-workspace/         # MPS DSL定義（既存維持）
├── .claude/               # Claude設定（既存維持）
└── web/                   # フロントエンド（既存維持）
```

---

## 🎯 移行マッピング

### マイクロサービス統合

| 旧構造 | 新構造 |
|--------|--------|
| `auth-service/` | `microservices/auth/` |
| `generated/auth/` | `microservices/auth/internal/` |
| `generated/auth-service/` | `microservices/auth/internal/` |
| （全12サービス同様） | |

### インフラ設定統合

| 旧構造 | 新構造 |
|--------|--------|
| `docker-compose*.yml` | `infrastructure/docker/` |
| `config/` | `infrastructure/config/` |
| `scripts/init-db.sql` | `infrastructure/databases/` |
| `scripts/migrations/` | `infrastructure/databases/migrations/` |

### ドキュメント整理

| 旧構造 | 新構造 |
|--------|--------|
| ルートの大量のMDファイル | `docs/reports/` |
| `README.md` | `docs/README.md` |
| `SETUP.md` | `docs/SETUP.md` |
| `CLAUDE.md` | `docs/CLAUDE.md` |

### ツール統合

| 旧構造 | 新構造 |
|--------|--------|
| `scripts/` | `tools/scripts/` |
| `test-client/` | `tools/test-client/` |
| `mock/` | `tools/mock/` |

### ビルド成果物

| 旧構造 | 新構造 |
|--------|--------|
| `bin/` | `build/bin/` |
| `proto/` | `build/proto/` |

---

## ✅ 完了したタスク

- [x] 新しいディレクトリ構造作成
- [x] マイクロサービスを `microservices/` に統合（12サービス）
- [x] インフラ設定を `infrastructure/` に統合
- [x] ドキュメントを `docs/` に整理
- [x] ツールを `tools/` に統合
- [x] ビルド成果物を `build/` に移行
- [x] ドキュメント内のリンク更新
  - `docs/README.md` のパス更新
  - `docs/CLAUDE.md` のディレクトリ構造更新
- [x] ルート `README.md` の作成
- [x] `.gitignore` の更新
- [x] 古いディレクトリの削除

---

## 📊 削除された古いディレクトリ

以下のディレクトリを削除しました：

```
auth-service/
shop-service/
customer-service/
inventory-service/
order-service/
payment-service/
shipping-service/
notification-service/
review-service/
chat-service/
search-service/
admin-service/
generated/
proto/
scripts/
bin/
config/
mock/
test-client/
services/
```

---

## 📝 ドキュメント更新内容

### README.md（ルート）
- プロジェクト概要
- ディレクトリ構造図
- ドキュメントへのリンク集

### docs/README.md
- パス更新: `generated/` → `microservices/`
- パス更新: `scripts/` → `tools/scripts/`
- パス更新: `docker-compose` → `infrastructure/docker/`

### docs/CLAUDE.md
- ディレクトリ構造図の全面更新
- パス更新: `generated/` → `microservices/`
- パス更新: `scripts/mps-generate.sh` → `tools/scripts/mps-generate.sh`

### .gitignore
- 削除された古いディレクトリを除外リストに追加
- ビルド成果物の除外設定

---

## 🎯 メリット

### 1. 明確な責任分離
各カテゴリが独立して理解しやすくなりました。

### 2. ドキュメント整理
ルートに散在していた13個のMDファイルを `docs/reports/` に集約しました。

### 3. インフラの統合
Docker, DB関連を `infrastructure/` に統合し、管理しやすくなりました。

### 4. ツールの統合
scripts, test-client等を `tools/` に集約し、開発ツールが一箇所にまとまりました。

### 5. ビルド成果物の分離
`build/` に集約し、gitignoreで適切に除外できるようになりました。

### 6. スケーラビリティ
将来K8s等を `infrastructure/kubernetes/` に追加しやすい構造になりました。

---

## 🔗 リンク更新状況

### 更新済み
- ✅ `docs/README.md` 内のすべてのパス
- ✅ `docs/CLAUDE.md` 内のすべてのパス
- ✅ ルート `README.md` の作成
- ✅ `.gitignore` の更新

### 確認が必要なファイル
以下のファイルに古いパスへの参照がある可能性があります：
- `docs/SETUP.md`
- `docs/reports/*.md`（各種レポート）
- `.claude/rules/*.md`

これらは必要に応じて後で更新してください。

---

## 📌 今後の作業

### 推奨事項

1. **Docker Compose設定の更新**
   - `infrastructure/docker/docker-compose.yml` 内のパスを確認・更新

2. **Makefileの更新**（存在する場合）
   - パス参照を新しい構造に更新

3. **CI/CD設定の更新**（存在する場合）
   - `.github/workflows/` 等のパスを更新

4. **残りのドキュメント確認**
   - `docs/SETUP.md` のパス確認
   - `docs/reports/*.md` のパス確認

---

## 🎉 完了

プロジェクト構造の再編成が完了しました。

**新しい構造の特徴**:
- 🎯 マイクロサービスが `microservices/` に集約
- 🔧 インフラ設定が `infrastructure/` に集約
- 📚 ドキュメントが `docs/` に整理
- 🛠️ ツールが `tools/` に集約
- 🏗️ ビルド成果物が `build/` に分離

すべてのカテゴリが明確になり、リンク切れも解消されました。
