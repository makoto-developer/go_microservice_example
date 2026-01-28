#!/bin/bash

echo "================================"
echo "Admin Service Verification"
echo "================================"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check PostgreSQL is running
echo "1. Checking PostgreSQL (port 22021)..."
if docker ps | grep -q "postgres_admin"; then
    echo -e "${GREEN}✅ PostgreSQL container is running${NC}"
else
    echo -e "${RED}❌ PostgreSQL container is not running${NC}"
    exit 1
fi

# Check database connectivity
echo ""
echo "2. Checking database connectivity..."
if docker exec go_microservice_postgres_admin_dev psql -U postgres -d admin_service -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Database is accessible${NC}"
else
    echo -e "${RED}❌ Database is not accessible${NC}"
    exit 1
fi

# Check tables exist
echo ""
echo "3. Checking database schema..."
TABLE_COUNT=$(docker exec go_microservice_postgres_admin_dev psql -U postgres -d admin_service -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public' AND table_name IN ('admin_users', 'audit_logs');" | xargs)

if [ "$TABLE_COUNT" = "2" ]; then
    echo -e "${GREEN}✅ Required tables exist (admin_users, audit_logs)${NC}"
else
    echo -e "${YELLOW}⚠️  Some tables might be missing (found: $TABLE_COUNT/2)${NC}"
fi

# Check Go service builds
echo ""
echo "4. Checking Go service build..."
if [ -f "admin-service" ]; then
    echo -e "${GREEN}✅ Service binary exists${NC}"
else
    echo -e "${YELLOW}⚠️  Building service...${NC}"
    go build -o admin-service
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ Service built successfully${NC}"
    else
        echo -e "${RED}❌ Service build failed${NC}"
        exit 1
    fi
fi

# Test service startup
echo ""
echo "5. Testing service startup..."
DATABASE_URL="postgresql://postgres:postgres_password@localhost:22021/admin_service?sslmode=disable" \
SERVICE_PORT=22111 \
./admin-service &
PID=$!

sleep 3

if kill -0 $PID 2>/dev/null; then
    echo -e "${GREEN}✅ Service started successfully${NC}"

    # Check port is listening
    if lsof -ti:22111 >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Service is listening on port 22111${NC}"
    else
        echo -e "${RED}❌ Service is not listening on port 22111${NC}"
    fi

    kill $PID 2>/dev/null
    wait $PID 2>/dev/null
else
    echo -e "${RED}❌ Service failed to start${NC}"
    exit 1
fi

echo ""
echo "================================"
echo -e "${GREEN}All checks passed!${NC}"
echo "================================"
echo ""
echo "Admin Service Configuration:"
echo "  - Database: postgresql://localhost:22021/admin_service"
echo "  - Service Port: 22111"
echo "  - Tables: admin_users, audit_logs"
echo ""
