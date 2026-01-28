#!/bin/bash

set -e

DB_HOST="${CUSTOMER_DB_HOST:-localhost}"
DB_PORT="${CUSTOMER_DB_PORT:-5434}"
DB_USER="${CUSTOMER_DB_USER:-customer_user}"
DB_PASSWORD="${CUSTOMER_DB_PASSWORD:-customer_password}"
DB_NAME="${CUSTOMER_DB_NAME:-customer_db}"

export PGPASSWORD=$DB_PASSWORD

echo "Applying Customer Service migrations..."

for migration in $(ls -1 scripts/migrations/customer/*.sql | sort); do
    echo "Applying migration: $migration"
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $migration
done

echo "Migrations applied successfully!"
