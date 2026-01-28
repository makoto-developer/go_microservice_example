#!/bin/bash

# Integration Test Runner for Order-Payment-Inventory Flow
# This script runs the integration tests for the order payment flow

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Project root (3 levels up from this script)
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/../../.." && pwd )"

echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}   Order-Payment-Inventory Integration Test Runner${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""

# Function to check if a service is running
check_service() {
    local service_name=$1
    local port=$2
    
    echo -ne "${YELLOW}Checking ${service_name}...${NC}"
    
    if nc -z localhost $port 2>/dev/null; then
        echo -e " ${GREEN}✓ Running${NC}"
        return 0
    else
        echo -e " ${RED}✗ Not running${NC}"
        return 1
    fi
}

# Function to wait for service to be ready
wait_for_service() {
    local service_name=$1
    local port=$2
    local max_attempts=30
    local attempt=1
    
    echo -ne "${YELLOW}Waiting for ${service_name}...${NC}"
    
    while [ $attempt -le $max_attempts ]; do
        if nc -z localhost $port 2>/dev/null; then
            echo -e " ${GREEN}✓ Ready${NC}"
            return 0
        fi
        sleep 1
        attempt=$((attempt + 1))
    done
    
    echo -e " ${RED}✗ Timeout${NC}"
    return 1
}

# Check prerequisites
echo -e "${BLUE}Checking prerequisites...${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    exit 1
fi

# Check if required services are running
SERVICES_OK=true

check_service "Customer Service" 22102 || SERVICES_OK=false
check_service "Inventory Service" 22103 || SERVICES_OK=false
check_service "Order Service" 22104 || SERVICES_OK=false
check_service "Payment Service" 22105 || SERVICES_OK=false

echo ""

if [ "$SERVICES_OK" = false ]; then
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}Warning: Not all services are running${NC}"
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo "To start services, run:"
    echo "  cd $PROJECT_ROOT/simple-servers/customer && ./customer-service &"
    echo "  cd $PROJECT_ROOT/simple-servers/inventory && ./inventory-service &"
    echo "  cd $PROJECT_ROOT/simple-servers/order && ./order-server &"
    echo "  cd $PROJECT_ROOT/simple-servers/payment && ./payment-server &"
    echo ""
    read -p "Continue anyway? (y/N) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Navigate to test directory
cd "$SCRIPT_DIR"

# Initialize Go module if needed
if [ ! -f "go.mod" ]; then
    echo -e "${YELLOW}Initializing Go module...${NC}"
    go mod init github.com/makoto-developer/go_microservice_example/tests/integration/order_flow
    go mod tidy
    echo -e "${GREEN}✓ Go module initialized${NC}"
    echo ""
fi

# Download dependencies
echo -e "${BLUE}Downloading dependencies...${NC}"
go mod download
echo ""

# Run tests
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}   Running Integration Tests${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════════${NC}"
echo ""

# Parse command line arguments
TEST_FLAGS=""
VERBOSE=false
BENCHMARK=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -v|--verbose)
            VERBOSE=true
            TEST_FLAGS="$TEST_FLAGS -v"
            shift
            ;;
        -b|--benchmark)
            BENCHMARK=true
            shift
            ;;
        -t|--timeout)
            TEST_FLAGS="$TEST_FLAGS -timeout $2"
            shift 2
            ;;
        *)
            TEST_FLAGS="$TEST_FLAGS $1"
            shift
            ;;
    esac
done

# Run unit tests
if [ "$BENCHMARK" = false ]; then
    echo -e "${YELLOW}Running integration tests...${NC}"
    if go test $TEST_FLAGS .; then
        echo ""
        echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
        echo -e "${GREEN}   ✨ All tests passed!${NC}"
        echo -e "${GREEN}═══════════════════════════════════════════════════════════${NC}"
    else
        echo ""
        echo -e "${RED}═══════════════════════════════════════════════════════════${NC}"
        echo -e "${RED}   ✗ Tests failed${NC}"
        echo -e "${RED}═══════════════════════════════════════════════════════════${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}Running benchmarks...${NC}"
    go test -bench=. -benchmem $TEST_FLAGS .
fi

echo ""

# Generate coverage report if requested
if [[ "$TEST_FLAGS" == *"-cover"* ]]; then
    echo -e "${YELLOW}Generating coverage report...${NC}"
    go test -coverprofile=coverage.out .
    go tool cover -html=coverage.out -o coverage.html
    echo -e "${GREEN}✓ Coverage report generated: coverage.html${NC}"
fi

echo ""
echo -e "${BLUE}Test execution completed at $(date)${NC}"
