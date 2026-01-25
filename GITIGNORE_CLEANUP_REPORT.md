# .gitignore クリーンアップレポート

## 概要

.gitignoreファイルの更新と、コミット済みだが除外すべきファイルのチェックを実施しました。

---

## 🚨 重大な問題

### 1. **Node.jsの依存関係がコミット済み**

**場所**: `web/shop_mall_web/node_modules/` (13MB)

**問題**:
- node_modulesは通常gitignoreすべき
- リポジトリサイズが肥大化
- セキュリティリスク（依存パッケージのバージョン固定）

**対処**:
```bash
git rm -r --cached web/shop_mall_web/node_modules
git rm --cached web/shop_mall_web/package-lock.json
git commit -m "Remove node_modules from git tracking"
```

---

### 2. **機密情報の可能性があるファイル**

**ファイル**: `.env.new`

**問題**:
- 環境変数ファイルがコミット済み
- 機密情報（API KEY等）が含まれる可能性

**対処**:
```bash
# まず内容を確認
cat .env.new

# 機密情報が含まれる場合は履歴から完全削除が必要
# BFG Repo-Cleaner または git filter-branch を使用
git filter-branch --force --index-filter \
  'git rm --cached --ignore-unmatch .env.new' \
  --prune-empty --tag-name-filter cat -- --all
```

**⚠️ 警告**: 機密情報が含まれている場合、GitHubにpush済みなら以下も必要:
1. GitHub上のリポジトリ設定でシークレットをローテーション
2. API KEYを無効化・再発行

---

### 3. **AIツール設定ファイルがコミット済み**

**ディレクトリ**:
- `.claude/` (Claude設定)
- `.serena/` (Serena設定)

**問題**:
- .gitignoreで除外されているが、過去にコミット済み
- プロジェクト固有の設定が含まれる

**対処**:
```bash
git rm -r --cached .claude .serena
git commit -m "Remove AI tool config from git tracking"
```

**注意**: これらのディレクトリは実際には有用なプロジェクト設定を含むため、.gitignoreから削除してコミットし続けることも検討してください。

---

## ⚠️ その他の問題

### 4. **ログファイルがコミット済み**

**ファイル**:
- `terminal1.log`
- `microservices/auth/auth-server.log`
- `microservices/shop/shop-server.log`
- `web/shop_mall_web/phoenix.log`
- `logs/*.log`

**対処**:
```bash
git rm --cached terminal1.log
git rm --cached microservices/auth/auth-server.log
git rm --cached microservices/shop/shop-server.log
git rm --cached web/shop_mall_web/phoenix.log
git rm --cached logs/*.log
git commit -m "Remove log files from git tracking"
```

---

### 5. **.DS_Storeファイルがコミット済み**

**場所**: 複数のディレクトリ

**対処**:
```bash
find . -name .DS_Store -print0 | xargs -0 git rm --cached
git commit -m "Remove .DS_Store files from git tracking"
```

---

### 6. **バックアップファイルがコミット済み**

**ファイル**:
- `docker-compose.yml.backup`
- `infrastructure/databases/migrations/shop/001_create_tables.sql.backup`

**対処**:
```bash
git rm --cached docker-compose.yml.backup
git rm --cached infrastructure/databases/migrations/shop/001_create_tables.sql.backup
git commit -m "Remove backup files from git tracking"
```

---

### 7. **バイナリファイルがコミット済み**

**ファイル**:
- `microservices/inventory/server`
- `microservices/shop/server`

**対処**:
```bash
git rm --cached microservices/inventory/server
git rm --cached microservices/shop/server
git commit -m "Remove binary files from git tracking"
```

---

### 8. **テスト成果物がコミット済み**

**ファイル**:
- `web/shop_mall_web/test-results/.last-run.json`
- `web/shop_mall_web/playwright-report/`

**対処**:
```bash
git rm --cached web/shop_mall_web/test-results/.last-run.json
git rm -r --cached web/shop_mall_web/playwright-report
git commit -m "Remove test artifacts from git tracking"
```

---

## ✅ 実施済み対応

### .gitignoreファイルの更新

以下を追加しました:

```gitignore
# Elixir/Phoenix
_build/
deps/
*.beam
*.ez
/priv/static/assets/
/priv/static/cache_manifest.json
erl_crash.dump

# Node.js
node_modules/
npm-debug.log
package-lock.json

# Playwright/テスト
test-results/
playwright-report/
.last-run.json

# 環境変数
.env.new
.env.*.local
```

---

## 📋 推奨される実行手順

### ステップ1: 機密情報の確認（最優先）

```bash
# .env.newの内容を確認
cat .env.new

# 機密情報が含まれている場合は別途対応が必要
```

### ステップ2: 不要ファイルの削除

```bash
# node_modules（最大の問題）
git rm -r --cached web/shop_mall_web/node_modules
git rm --cached web/shop_mall_web/package-lock.json

# ログファイル
git rm --cached terminal1.log
git rm --cached microservices/auth/auth-server.log
git rm --cached microservices/shop/shop-server.log
git rm --cached web/shop_mall_web/phoenix.log
git rm -r --cached logs/ 2>/dev/null || true

# .DS_Store
find . -name .DS_Store -print0 | xargs -0 git rm --cached 2>/dev/null || true

# バイナリ
git rm --cached microservices/inventory/server 2>/dev/null || true
git rm --cached microservices/shop/server 2>/dev/null || true

# バックアップ
git rm --cached docker-compose.yml.backup
git rm --cached infrastructure/databases/migrations/shop/001_create_tables.sql.backup 2>/dev/null || true

# テスト成果物
git rm --cached web/shop_mall_web/test-results/.last-run.json 2>/dev/null || true
git rm -r --cached web/shop_mall_web/playwright-report 2>/dev/null || true

# .env.new（機密情報確認後）
# git rm --cached .env.new

# AI設定（プロジェクトに必要な場合はスキップ）
# git rm -r --cached .claude .serena
```

### ステップ3: コミット

```bash
git add .gitignore
git commit -m "chore: update .gitignore and remove tracked files that should be ignored

- Remove node_modules from tracking (13MB)
- Remove log files
- Remove .DS_Store files
- Remove binary files (server executables)
- Remove backup files
- Remove test artifacts
- Update .gitignore for Elixir, Node.js, and test artifacts"
```

### ステップ4: プッシュ前の確認

```bash
# リポジトリサイズの確認
du -sh .git

# 削除されたファイルの確認
git status
```

---

## 🔒 セキュリティ上の注意

### .env.newファイルについて

**重要**: このファイルに機密情報が含まれている場合:

1. **GitHubにpush済みの場合**:
   - すべてのAPI KEY、パスワード、トークンを無効化・再発行
   - GitHubのSettings → Secrets でシークレットを更新

2. **履歴から完全削除**:
   ```bash
   # BFG Repo-Cleaner（推奨）
   bfg --delete-files .env.new

   # または git filter-branch
   git filter-branch --force --index-filter \
     'git rm --cached --ignore-unmatch .env.new' \
     --prune-empty --tag-name-filter cat -- --all

   # 強制プッシュ（チーム全員に通知が必要）
   git push origin --force --all
   ```

---

## 📊 改善効果の見積もり

### リポジトリサイズ削減

- **node_modules削除**: 約13MB削減
- **ログファイル削除**: 約1-2MB削減
- **.DS_Store削除**: 約100KB削減

**合計**: 約15MB削減見込み

### その他の効果

- クリーンアップ速度向上
- git操作の高速化
- セキュリティリスク低減

---

## まとめ

### 必須対応

1. ✅ .gitignore更新（完了）
2. ⚠️ node_modules削除（最優先）
3. ⚠️ .env.newの機密情報確認（最優先）
4. ⚠️ ログファイル削除

### 推奨対応

5. .DS_Store削除
6. バイナリファイル削除
7. バックアップファイル削除
8. テスト成果物削除

### 検討事項

- `.claude/`と`.serena/`はプロジェクト固有設定として残すか判断
