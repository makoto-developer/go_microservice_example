#!/bin/bash

# Service Availability Checker
# Checks if all required databases are running and accessible

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}Checking Service Availability${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Function to check PostgreSQL connection
check_postgres() {
    local name=$1
    local port=$2
    local database=$3

    echo -e "${YELLOW}Checking $name (port $port)...${NC}"

    # Check if port is listening
    if ! nc -z localhost $port 2>/dev/null; then
        echo -e "${RED}❌ Port $port is not listening${NC}"
        return 1
    fi

    # Try to connect to database
    if PGPASSWORD=postgres_password psql -h localhost -p $port -U postgres -d $database -c "SELECT 1;" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ $name is accessible${NC}"

        # Get table count
        table_count=$(PGPASSWORD=postgres_password psql -h localhost -p $port -U postgres -d $database -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';")
        echo -e "   Tables: $table_count"
        return 0
    else
        echo -e "${RED}❌ Cannot connect to $name database${NC}"
        return 1
    fi
}

# Check all services
all_ok=true

check_postgres "Notification Service" 22017 "notification_service" || all_ok=false
echo ""

check_postgres "Review Service" 22018 "review_service" || all_ok=false
echo ""

check_postgres "Shipping Service" 22016 "shipping_service" || all_ok=false
echo ""

# Summary
echo -e "${BLUE}========================================${NC}"
if [ "$all_ok" = true ]; then
    echo -e "${GREEN}✅ All services are ready for testing${NC}"
    echo ""
    echo -e "${YELLOW}You can now run:${NC}"
    echo -e "  ./run_all_tests.sh"
    exit 0
else
    echo -e "${RED}❌ Some services are not ready${NC}"
    echo ""
    echo -e "${YELLOW}Please ensure all services are running:${NC}"
    echo -e "  docker-compose up -d"
    exit 1
fi
