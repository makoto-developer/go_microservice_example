#!/bin/bash

# Service availability check script

set -e

AUTH_SERVICE_ADDR=${AUTH_SERVICE_ADDR:-localhost:22100}
DB_HOST=${AUTH_DB_HOST:-localhost}
DB_PORT=${AUTH_DB_PORT:-22010}

echo "🔍 Checking service availability..."
echo ""

# Check Auth Service
echo -n "Auth Service ($AUTH_SERVICE_ADDR): "
if nc -z ${AUTH_SERVICE_ADDR/:/ } 2>/dev/null; then
    echo "✅ Running"
else
    echo "❌ Not running"
    echo ""
    echo "To start Auth Service:"
    echo "  cd microservices/auth && ./auth-server"
    exit 1
fi

# Check Database
echo -n "PostgreSQL ($DB_HOST:$DB_PORT): "
if nc -z $DB_HOST $DB_PORT 2>/dev/null; then
    echo "✅ Running"
else
    echo "❌ Not running"
    echo ""
    echo "To start database:"
    echo "  docker-compose up -d postgres_auth_dev"
    exit 1
fi

# Check database connectivity
echo -n "Database connectivity: "
export PGPASSWORD=${AUTH_DB_PASSWORD:-postgres}
if psql -h $DB_HOST -p $DB_PORT -U ${AUTH_DB_USER:-postgres} -d ${AUTH_DB_NAME:-auth_service} -c "SELECT 1" >/dev/null 2>&1; then
    echo "✅ Connected"
else
    echo "❌ Cannot connect"
    echo ""
    echo "Please check database credentials"
    exit 1
fi
unset PGPASSWORD

echo ""
echo "✅ All services are ready!"
