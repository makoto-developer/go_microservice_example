#!/bin/bash
set -e

# MPS Code Generation Script
# このスクリプトは、MPSのDSL定義からGoコードを生成します

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
MPS_WORKSPACE="$PROJECT_ROOT/mps-workspace"
GENERATED_DIR="$PROJECT_ROOT/generated"

# カラー定義
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ヘルプメッセージ
show_help() {
  cat << EOF
Usage: $0 [OPTIONS] [SERVICE_NAME]

MPS DSL定義からGoコードを生成します。

OPTIONS:
  --all           すべてのサービスを生成
  -h, --help      このヘルプを表示

EXAMPLES:
  $0 auth-service         # Auth Serviceのみ生成
  $0 --all                # すべてのサービスを生成

NOTES:
  - このスクリプトはJetBrains MPSがインストールされている必要があります
  - 生成されたコードは generated/ ディレクトリに配置されます
  - 生成コードは編集しないでください（再生成で上書きされます）
EOF
}

# DSLジェネレーター確認
check_generator() {
  local generator="$PROJECT_ROOT/bin/dsl-generator"
  if [ ! -f "$generator" ]; then
    echo -e "${YELLOW}DSL generator が見つかりません。ビルドします...${NC}"
    cd "$PROJECT_ROOT/tools/dsl-generator" && go build -o ../../bin/dsl-generator ./cmd/main.go
    if [ $? -ne 0 ]; then
      echo -e "${RED}Error: ジェネレーターのビルドに失敗しました${NC}"
      exit 1
    fi
    echo -e "${GREEN}✓ ジェネレーターをビルドしました${NC}"
    echo ""
  fi
}

# サービス一覧を取得
get_services() {
  if [ -d "$MPS_WORKSPACE/solutions" ]; then
    find "$MPS_WORKSPACE/solutions" -mindepth 1 -maxdepth 1 -type d -exec basename {} \;
  fi
}

# 単一サービスを生成
generate_service() {
  local service_name="$1"
  local service_dir="$MPS_WORKSPACE/solutions/$service_name"
  local output_dir="$GENERATED_DIR/$service_name"

  if [ ! -d "$service_dir" ]; then
    echo -e "${RED}Error: サービス '$service_name' が見つかりません${NC}"
    echo "利用可能なサービス:"
    get_services | sed 's/^/  - /'
    exit 1
  fi

  echo -e "${BLUE}Generating code for $service_name...${NC}"

  # 出力ディレクトリを作成
  mkdir -p "$output_dir"

  # ========================================
  # DSL Generator 実行
  # ========================================
  local dsl_file="$service_dir/service.model"
  local generator="$PROJECT_ROOT/bin/dsl-generator"

  if [ ! -f "$generator" ]; then
    echo -e "${RED}Error: DSL generator が見つかりません${NC}"
    echo "ジェネレーターをビルドしてください:"
    echo "  cd tools/dsl-generator && go build -o ../../bin/dsl-generator ./cmd/main.go"
    exit 1
  fi

  if [ ! -f "$dsl_file" ]; then
    echo -e "${RED}Error: DSL定義ファイルが見つかりません: $dsl_file${NC}"
    exit 1
  fi

  echo "DSL定義: $dsl_file"
  echo "生成先: $output_dir"
  echo ""

  # DSLジェネレーター実行
  "$generator" \
    -input "$dsl_file" \
    -output "$GENERATED_DIR" \
    -service "$service_name"

  if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓ $service_name のコード生成が完了しました${NC}"
  else
    echo ""
    echo -e "${RED}✗ $service_name のコード生成に失敗しました${NC}"
    exit 1
  fi

  echo ""
}

# すべてのサービスを生成
generate_all() {
  echo -e "${BLUE}すべてのサービスのコードを生成します...${NC}"
  echo ""

  local services=($(get_services))

  if [ ${#services[@]} -eq 0 ]; then
    echo -e "${RED}Error: サービスが見つかりません${NC}"
    exit 1
  fi

  for service in "${services[@]}"; do
    generate_service "$service"
  done

  echo -e "${GREEN}✓ すべてのサービスを生成しました${NC}"
}

# メイン処理
main() {
  # 引数チェック
  if [ $# -eq 0 ]; then
    show_help
    exit 1
  fi

  case "$1" in
    -h|--help)
      show_help
      exit 0
      ;;
    --all)
      check_generator
      generate_all
      ;;
    *)
      check_generator
      generate_service "$1"
      ;;
  esac
}

main "$@"
