#!/bin/bash

set -e

# Database connection parameters
DB_HOST="${SHOP_DB_HOST:-localhost}"
DB_PORT="${SHOP_DB_PORT:-5433}"
DB_USER="${SHOP_DB_USER:-shop_user}"
DB_PASSWORD="${SHOP_DB_PASSWORD:-shop_password}"
DB_NAME="${SHOP_DB_NAME:-shop_db}"

# Export password for psql
export PGPASSWORD=$DB_PASSWORD

echo "Applying Shop Service migrations..."

# Apply migrations
for migration in $(ls -1 scripts/migrations/shop/*.sql | sort); do
    echo "Applying migration: $migration"
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $migration
done

echo "Migrations applied successfully!"
