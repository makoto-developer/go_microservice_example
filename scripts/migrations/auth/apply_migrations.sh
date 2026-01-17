#!/bin/bash

set -e

DOCKER_CONTAINER="go_microservice_postgres_auth_dev"
DB_USER="postgres"
DB_NAME="auth_db"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

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

echo "Running Auth Service migrations..."
run_sql "$SCRIPT_DIR/001_create_auth_tables.sql" "Auth Schema Migration (001)"

echo "========================================="
echo "Auth Service Migration Completed!"
echo "========================================="
