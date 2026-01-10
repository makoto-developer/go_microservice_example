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

# MPSインストール確認
check_mps_installation() {
  if [ ! -d "/Applications/MPS.app" ] && [ ! -f "/usr/local/bin/mps" ]; then
    echo -e "${RED}Error: JetBrains MPS が見つかりません${NC}"
    echo ""
    echo "MPSをインストールしてください:"
    echo "  https://www.jetbrains.com/mps/download/"
    echo ""
    echo "または、Homebrewでインストール:"
    echo "  brew install --cask mps"
    exit 1
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
  # MPS Generator 実行
  # ========================================
  # 注意: 以下は仮実装です
  # 実際のMPS Generatorの実行方法はMPSプロジェクト設定に依存します

  # TODO: 実際のMPS Generator呼び出しに置き換える
  # 例:
  # mps --generate "$MPS_WORKSPACE" --model "$service_name" --output "$output_dir"

  echo -e "${YELLOW}Warning: MPS Generator の実装が必要です${NC}"
  echo "このスクリプトは現在モックです。"
  echo ""
  echo "生成先: $output_dir"
  echo "DSL定義: $service_dir/service.model"
  echo ""

  # モック: ディレクトリ構成のみ作成
  mkdir -p "$output_dir/domain"
  mkdir -p "$output_dir/usecase"
  mkdir -p "$output_dir/handler"
  mkdir -p "$output_dir/infrastructure"
  mkdir -p "$output_dir/tests"

  # モック: go.modファイル作成
  if [ ! -f "$output_dir/go.mod" ]; then
    cat > "$output_dir/go.mod" << EOL
module github.com/makoto-developer/go_microservice_example/generated/$service_name

go 1.25

require (
	github.com/google/uuid v1.6.0
	google.golang.org/grpc v1.67.1
	google.golang.org/protobuf v1.35.1
)
EOL
  fi

  echo -e "${GREEN}✓ $service_name の構造を作成しました（モック）${NC}"
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
      check_mps_installation
      generate_all
      ;;
    *)
      check_mps_installation
      generate_service "$1"
      ;;
  esac
}

main "$@"
