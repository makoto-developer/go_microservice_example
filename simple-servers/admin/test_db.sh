#!/bin/bash

echo "================================"
echo "Admin Service Database Test"
echo "================================"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test 1: Check tables
echo -e "${BLUE}1. Checking database tables...${NC}"
docker exec go_microservice_postgres_admin_dev psql -U postgres -d admin_service -c "\dt"

echo ""
echo -e "${BLUE}2. Checking admin_users table structure...${NC}"
docker exec go_microservice_postgres_admin_dev psql -U postgres -d admin_service -c "\d admin_users"

echo ""
echo -e "${BLUE}3. Checking audit_logs table structure...${NC}"
docker exec go_microservice_postgres_admin_dev psql -U postgres -d admin_service -c "\d audit_logs"

echo ""
echo -e "${BLUE}4. Querying admin_users data...${NC}"
docker exec go_microservice_postgres_admin_dev psql -U postgres -d admin_service -c "
SELECT id, username, email, role, is_active, created_at
FROM admin_users
ORDER BY created_at DESC
LIMIT 5;
"

echo ""
echo -e "${BLUE}5. Querying audit_logs data...${NC}"
docker exec go_microservice_postgres_admin_dev psql -U postgres -d admin_service -c "
SELECT id, operation_type, operator_name, target_type, ip_address, created_at
FROM audit_logs
ORDER BY created_at DESC
LIMIT 5;
"

echo ""
echo -e "${GREEN}Database test complete!${NC}"
