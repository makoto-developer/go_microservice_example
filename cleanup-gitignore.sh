#!/bin/bash

# .gitignore クリーンアップスクリプト
# 使用方法: bash cleanup-gitignore.sh

set -e

echo "========================================="
echo ".gitignore クリーンアップスクリプト"
echo "========================================="
echo ""

# カラー設定
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 確認関数
confirm() {
    read -p "$1 [y/N]: " response
    case "$response" in
        [yY][eE][sS]|[yY])
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

echo -e "${YELLOW}⚠️  このスクリプトは以下のファイルをGit追跡から削除します:${NC}"
echo "  - node_modules/ (13MB)"
echo "  - ログファイル (*.log)"
echo "  - .DS_Store"
echo "  - バイナリファイル (server executables)"
echo "  - バックアップファイル (*.backup)"
echo "  - テスト成果物 (test-results/, playwright-report/)"
echo ""
echo -e "${RED}注意: ファイルはワーキングディレクトリから削除されません（--cachedオプション使用）${NC}"
echo ""

if ! confirm "続行しますか？"; then
    echo "キャンセルしました"
    exit 0
fi

echo ""
echo "========================================="
echo "ステップ1: node_modules削除"
echo "========================================="

if [ -d "web/shop_mall_web/node_modules" ]; then
    echo "node_modulesを削除中..."
    git rm -r --cached web/shop_mall_web/node_modules 2>/dev/null || true
    echo -e "${GREEN}✓ node_modules削除完了${NC}"
else
    echo "node_modulesは既に追跡されていません"
fi

if [ -f "web/shop_mall_web/package-lock.json" ]; then
    git rm --cached web/shop_mall_web/package-lock.json 2>/dev/null || true
    echo -e "${GREEN}✓ package-lock.json削除完了${NC}"
fi

echo ""
echo "========================================="
echo "ステップ2: ログファイル削除"
echo "========================================="

LOG_FILES=(
    "terminal1.log"
    "microservices/auth/auth-server.log"
    "microservices/shop/shop-server.log"
    "web/shop_mall_web/phoenix.log"
)

for log_file in "${LOG_FILES[@]}"; do
    if git ls-files --error-unmatch "$log_file" > /dev/null 2>&1; then
        echo "削除中: $log_file"
        git rm --cached "$log_file" 2>/dev/null || true
        echo -e "${GREEN}✓ $log_file 削除完了${NC}"
    fi
done

# logs/ディレクトリ
if git ls-files logs/ > /dev/null 2>&1; then
    git rm -r --cached logs/ 2>/dev/null || true
    echo -e "${GREEN}✓ logs/ディレクトリ削除完了${NC}"
fi

echo ""
echo "========================================="
echo "ステップ3: .DS_Store削除"
echo "========================================="

DS_STORE_COUNT=$(find . -name .DS_Store | wc -l | tr -d ' ')
if [ "$DS_STORE_COUNT" -gt 0 ]; then
    echo ".DS_Storeファイルを削除中... ($DS_STORE_COUNT 個)"
    find . -name .DS_Store -print0 | xargs -0 git rm --cached 2>/dev/null || true
    echo -e "${GREEN}✓ .DS_Store削除完了${NC}"
else
    echo ".DS_Storeは見つかりませんでした"
fi

echo ""
echo "========================================="
echo "ステップ4: バイナリファイル削除"
echo "========================================="

BINARY_FILES=(
    "microservices/inventory/server"
    "microservices/shop/server"
)

for binary_file in "${BINARY_FILES[@]}"; do
    if git ls-files --error-unmatch "$binary_file" > /dev/null 2>&1; then
        echo "削除中: $binary_file"
        git rm --cached "$binary_file" 2>/dev/null || true
        echo -e "${GREEN}✓ $binary_file 削除完了${NC}"
    fi
done

echo ""
echo "========================================="
echo "ステップ5: バックアップファイル削除"
echo "========================================="

BACKUP_FILES=(
    "docker-compose.yml.backup"
    "infrastructure/databases/migrations/shop/001_create_tables.sql.backup"
)

for backup_file in "${BACKUP_FILES[@]}"; do
    if git ls-files --error-unmatch "$backup_file" > /dev/null 2>&1; then
        echo "削除中: $backup_file"
        git rm --cached "$backup_file" 2>/dev/null || true
        echo -e "${GREEN}✓ $backup_file 削除完了${NC}"
    fi
done

echo ""
echo "========================================="
echo "ステップ6: テスト成果物削除"
echo "========================================="

if git ls-files --error-unmatch "web/shop_mall_web/test-results/.last-run.json" > /dev/null 2>&1; then
    git rm --cached web/shop_mall_web/test-results/.last-run.json 2>/dev/null || true
    echo -e "${GREEN}✓ .last-run.json削除完了${NC}"
fi

if git ls-files web/shop_mall_web/playwright-report/ > /dev/null 2>&1; then
    git rm -r --cached web/shop_mall_web/playwright-report 2>/dev/null || true
    echo -e "${GREEN}✓ playwright-report/削除完了${NC}"
fi

echo ""
echo "========================================="
echo "ステップ7: .env.newファイルの確認"
echo "========================================="

if git ls-files --error-unmatch ".env.new" > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  .env.newファイルが見つかりました${NC}"
    echo ""
    if confirm ".env.newファイルも削除しますか？（機密情報が含まれている場合は削除を推奨）"; then
        git rm --cached .env.new 2>/dev/null || true
        echo -e "${GREEN}✓ .env.new削除完了${NC}"
        echo ""
        echo -e "${RED}⚠️  重要: .env.newに機密情報が含まれている場合:${NC}"
        echo "  1. API KEY等を無効化・再発行してください"
        echo "  2. GitHubにpush済みの場合、git履歴からも削除が必要です"
        echo "     詳細: GITIGNORE_CLEANUP_REPORT.md を参照"
    else
        echo ".env.newの削除をスキップしました"
    fi
fi

echo ""
echo "========================================="
echo "完了"
echo "========================================="
echo ""
echo -e "${GREEN}✓ クリーンアップが完了しました${NC}"
echo ""
echo "次のステップ:"
echo "  1. 変更を確認: git status"
echo "  2. .gitignoreを追加: git add .gitignore"
echo "  3. コミット: git commit -m 'chore: update .gitignore and remove ignored files from tracking'"
echo "  4. プッシュ: git push origin main"
echo ""
echo "詳細なレポート: GITIGNORE_CLEANUP_REPORT.md"
