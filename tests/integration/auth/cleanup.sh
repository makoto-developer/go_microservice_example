#!/bin/bash

# Cleanup script for Auth Service integration tests
# This script removes all test data from the database

set -e

# Database configuration
DB_HOST=${AUTH_DB_HOST:-localhost}
DB_PORT=${AUTH_DB_PORT:-22010}
DB_USER=${AUTH_DB_USER:-postgres}
DB_PASSWORD=${AUTH_DB_PASSWORD:-postgres}
DB_NAME=${AUTH_DB_NAME:-auth_service}

echo "🧹 Cleaning up test data from Auth Service database..."

# Set PGPASSWORD for non-interactive authentication
export PGPASSWORD=$DB_PASSWORD

# Delete test data from customer_users table
echo "  - Cleaning customer_users table..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c \
  "DELETE FROM customer_users WHERE email LIKE 'test_%';" 2>/dev/null || echo "    ⚠️  Warning: Could not clean customer_users table"

# Delete test data from owner_users table
echo "  - Cleaning owner_users table..."
psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c \
  "DELETE FROM owner_users WHERE email LIKE 'test_%';" 2>/dev/null || echo "    ⚠️  Warning: Could not clean owner_users table"

# Count remaining test users
CUSTOMER_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c \
  "SELECT COUNT(*) FROM customer_users WHERE email LIKE 'test_%';" 2>/dev/null | tr -d ' ' || echo "0")

OWNER_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c \
  "SELECT COUNT(*) FROM owner_users WHERE email LIKE 'test_%';" 2>/dev/null | tr -d ' ' || echo "0")

echo ""
echo "✅ Cleanup complete!"
echo "   Customer test users remaining: $CUSTOMER_COUNT"
echo "   Owner test users remaining: $OWNER_COUNT"

# Unset password
unset PGPASSWORD
