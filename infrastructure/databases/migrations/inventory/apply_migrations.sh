#!/bin/bash

set -e

DB_HOST="${INVENTORY_DB_HOST:-localhost}"
DB_PORT="${INVENTORY_DB_PORT:-5435}"
DB_USER="${INVENTORY_DB_USER:-inventory_user}"
DB_PASSWORD="${INVENTORY_DB_PASSWORD:-inventory_password}"
DB_NAME="${INVENTORY_DB_NAME:-inventory_db}"

export PGPASSWORD=$DB_PASSWORD

echo "Applying Inventory Service migrations..."

for migration in $(ls -1 scripts/migrations/inventory/*.sql | sort); do
    echo "Applying migration: $migration"
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $migration
done

echo "Migrations applied successfully!"
