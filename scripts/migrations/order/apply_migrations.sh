#!/bin/bash

set -e

DB_HOST="${ORDER_DB_HOST:-localhost}"
DB_PORT="${ORDER_DB_PORT:-5436}"
DB_USER="${ORDER_DB_USER:-order_user}"
DB_PASSWORD="${ORDER_DB_PASSWORD:-order_password}"
DB_NAME="${ORDER_DB_NAME:-order_db}"

export PGPASSWORD=$DB_PASSWORD

echo "Applying Order Service migrations..."

for migration in $(ls -1 scripts/migrations/order/*.sql | sort); do
    echo "Applying migration: $migration"
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $migration
done

echo "Migrations applied successfully!"
