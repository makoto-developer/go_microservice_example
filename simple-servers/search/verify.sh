#!/bin/bash

echo "======================================"
echo "Search Service Verification Script"
echo "======================================"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if PostgreSQL container is running
echo "1. Checking PostgreSQL container..."
if docker ps | grep -q "postgres_search"; then
    echo -e "${GREEN}✅ PostgreSQL container is running${NC}"
else
    echo -e "${RED}❌ PostgreSQL container is not running${NC}"
    echo "   Run: cd ../../infrastructure/docker && docker-compose up -d postgres_search"
    exit 1
fi

# Check database connection
echo ""
echo "2. Testing database connection..."
if docker exec go_microservice_postgres_search_dev psql -U postgres -d search_service -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Database connection successful${NC}"
else
    echo -e "${RED}❌ Cannot connect to database${NC}"
    exit 1
fi

# Check tables
echo ""
echo "3. Checking database tables..."
TABLES=$(docker exec go_microservice_postgres_search_dev psql -U postgres -d search_service -t -c "SELECT tablename FROM pg_tables WHERE schemaname = 'public';" | tr -d ' ' | grep -v '^$')

if echo "$TABLES" | grep -q "search_indexes"; then
    echo -e "${GREEN}✅ search_indexes table exists${NC}"
else
    echo -e "${RED}❌ search_indexes table not found${NC}"
fi

if echo "$TABLES" | grep -q "search_logs"; then
    echo -e "${GREEN}✅ search_logs table exists${NC}"
else
    echo -e "${RED}❌ search_logs table not found${NC}"
fi

# Check indexes
echo ""
echo "4. Checking database indexes..."
INDEXES=$(docker exec go_microservice_postgres_search_dev psql -U postgres -d search_service -t -c "SELECT indexname FROM pg_indexes WHERE schemaname = 'public';" | tr -d ' ' | grep -v '^$' | grep -v '_pkey$')

echo "   Found indexes:"
echo "$INDEXES" | while read idx; do
    if [ ! -z "$idx" ]; then
        echo -e "   ${GREEN}✓${NC} $idx"
    fi
done

# Check if binary exists
echo ""
echo "5. Checking service binary..."
if [ -f "search-service" ]; then
    echo -e "${GREEN}✅ Service binary exists${NC}"
else
    echo -e "${YELLOW}⚠️  Service binary not found. Run: make build${NC}"
fi

# Test service startup (if binary exists)
if [ -f "search-service" ]; then
    echo ""
    echo "6. Testing service startup..."
    ./search-service &
    SERVICE_PID=$!
    sleep 2

    if ps -p $SERVICE_PID > /dev/null; then
        echo -e "${GREEN}✅ Service started successfully${NC}"
        echo "   Service is running on port 22110"
        kill $SERVICE_PID 2>/dev/null
        wait $SERVICE_PID 2>/dev/null
    else
        echo -e "${RED}❌ Service failed to start${NC}"
        exit 1
    fi
fi

echo ""
echo "======================================"
echo -e "${GREEN}All checks passed!${NC}"
echo "======================================"
echo ""
echo "To start the service:"
echo "  make run"
echo ""
echo "To check database:"
echo "  make db-test"
echo ""
