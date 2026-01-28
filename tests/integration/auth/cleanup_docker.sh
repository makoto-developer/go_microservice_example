#!/bin/bash

# Cleanup script for Auth Service integration tests (using Docker)
# This script removes all test data from the database

set -e

CONTAINER_NAME="go_microservice_postgres_auth_dev"
DB_USER="postgres"
DB_NAME="auth_service"

echo "🧹 Cleaning up test data from Auth Service database..."

# Delete test data from customer_users table
echo "  - Cleaning customer_users table..."
docker exec $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -c \
  "DELETE FROM customer_users WHERE email LIKE 'test_%';" 2>/dev/null || echo "    ⚠️  Warning: Could not clean customer_users table"

# Delete test data from owner_users table
echo "  - Cleaning owner_users table..."
docker exec $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -c \
  "DELETE FROM owner_users WHERE email LIKE 'test_%';" 2>/dev/null || echo "    ⚠️  Warning: Could not clean owner_users table"

# Count remaining test users
CUSTOMER_COUNT=$(docker exec $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -t -c \
  "SELECT COUNT(*) FROM customer_users WHERE email LIKE 'test_%';" 2>/dev/null | tr -d ' ' || echo "0")

OWNER_COUNT=$(docker exec $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME -t -c \
  "SELECT COUNT(*) FROM owner_users WHERE email LIKE 'test_%';" 2>/dev/null | tr -d ' ' || echo "0")

echo ""
echo "✅ Cleanup complete!"
echo "   Customer test users remaining: $CUSTOMER_COUNT"
echo "   Owner test users remaining: $OWNER_COUNT"
