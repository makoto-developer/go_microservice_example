#!/bin/bash

# Test execution script for Auth Service integration tests

set -e

# Configuration
AUTH_SERVICE_ADDR=${AUTH_SERVICE_ADDR:-localhost:22100}
DB_HOST=${AUTH_DB_HOST:-localhost}
DB_PORT=${AUTH_DB_PORT:-22010}

echo "🧪 Auth Service Integration Tests"
echo "=================================="
echo ""
echo "Configuration:"
echo "  - Auth Service: $AUTH_SERVICE_ADDR"
echo "  - Database: $DB_HOST:$DB_PORT"
echo ""

# Check if Auth Service is running
echo "🔍 Checking Auth Service availability..."
if ! nc -z ${AUTH_SERVICE_ADDR/:/ } 2>/dev/null; then
    echo "❌ Error: Auth Service is not running at $AUTH_SERVICE_ADDR"
    echo "   Please start the Auth Service first:"
    echo "   cd microservices/auth && ./auth-server"
    exit 1
fi
echo "✅ Auth Service is running"
echo ""

# Check if database is running
echo "🔍 Checking database availability..."
if ! nc -z $DB_HOST ${DB_PORT} 2>/dev/null; then
    echo "❌ Error: Database is not running at $DB_HOST:$DB_PORT"
    echo "   Please start the database first"
    exit 1
fi
echo "✅ Database is running"
echo ""

# Install dependencies
echo "📦 Installing dependencies..."
go mod download
echo "✅ Dependencies installed"
echo ""

# Run cleanup before tests
echo "🧹 Cleaning up previous test data..."
./cleanup.sh
echo ""

# Run tests
echo "🚀 Running integration tests..."
echo ""
go test -v -timeout 30s ./...

TEST_EXIT_CODE=$?

echo ""
if [ $TEST_EXIT_CODE -eq 0 ]; then
    echo "✅ All tests passed!"
else
    echo "❌ Some tests failed (exit code: $TEST_EXIT_CODE)"
fi

# Cleanup after tests
echo ""
echo "🧹 Cleaning up test data..."
./cleanup.sh

exit $TEST_EXIT_CODE
