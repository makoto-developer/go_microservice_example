#!/bin/bash

# ============================================
# Shop Service - Migration Script
# ============================================
# Description: Shop Serviceのマイグレーションを実行
# Usage:
#   ./apply_migrations.sh              # スキーマ+テストデータ
#   ./apply_migrations.sh --schema-only # スキーマのみ
#   ./apply_migrations.sh --test-only   # テストデータのみ

set -e

# ============================================
# Configuration
# ============================================

DOCKER_CONTAINER="go_microservice_postgres_dev"
DB_USER="admin"
DB_NAME="shop_db"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# ============================================
# Functions
# ============================================

run_sql() {
    local sql_file=$1
    local description=$2

    echo "========================================="
    echo "$description"
    echo "========================================="

    docker exec -i "$DOCKER_CONTAINER" psql -U "$DB_USER" -d "$DB_NAME" < "$sql_file"

    if [ $? -eq 0 ]; then
        echo "✓ $description completed successfully"
    else
        echo "✗ $description failed"
        exit 1
    fi
    echo ""
}

# ============================================
# Main
# ============================================

MODE="${1:---all}"

case "$MODE" in
    --schema-only)
        echo "Running schema migration only..."
        run_sql "$SCRIPT_DIR/001_create_tables.sql" "Schema Migration (001)"
        ;;

    --test-only)
        echo "Running test data insertion only..."
        run_sql "$SCRIPT_DIR/002_create_test_schema.sql" "Test Data Migration (002)"
        ;;

    --all|*)
        echo "Running all migrations..."
        run_sql "$SCRIPT_DIR/001_create_tables.sql" "Schema Migration (001)"
        run_sql "$SCRIPT_DIR/002_create_test_schema.sql" "Test Data Migration (002)"
        ;;
esac

echo "========================================="
echo "Migration Completed!"
echo "========================================="
echo ""
echo "Next steps:"
echo "  1. Verify data: docker exec -i $DOCKER_CONTAINER psql -U $DB_USER -d $DB_NAME -c 'SELECT id, name FROM products LIMIT 5;'"
echo "  2. Restart Shop Service: docker restart go_microservice_shop_service_dev"
echo "  3. Test gRPC: grpcurl -plaintext -d '{\"shop_id\": \"11111111-1111-1111-1111-111111111111\", \"published_only\": true}' localhost:20101 shop_service.v1.ShopService/ListProducts"
echo ""
