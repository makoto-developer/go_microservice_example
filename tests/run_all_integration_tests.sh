#!/bin/bash

# Master Integration Test Runner
# Executes all integration and E2E tests in the correct order

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 Master Integration Test Runner"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Project: Go Microservice Example - Online Shop Mall"
echo "Location: $PROJECT_ROOT"
echo "Started: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0
FAILED_TEST_NAMES=()

# Function to run a test suite
run_test_suite() {
    local suite_name=$1
    local test_command=$2
    local suite_dir=$3

    TOTAL_TESTS=$((TOTAL_TESTS + 1))

    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${BLUE}🧪 Test Suite $TOTAL_TESTS: $suite_name${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""

    if [ -n "$suite_dir" ]; then
        cd "$suite_dir" || {
            echo -e "${RED}❌ Failed to change directory to $suite_dir${NC}"
            FAILED_TESTS=$((FAILED_TESTS + 1))
            FAILED_TEST_NAMES+=("$suite_name")
            return 1
        }
    fi

    if eval "$test_command"; then
        echo ""
        echo -e "${GREEN}✅ $suite_name - PASSED${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
        return 0
    else
        echo ""
        echo -e "${RED}❌ $suite_name - FAILED${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
        FAILED_TEST_NAMES+=("$suite_name")
        return 1
    fi
}

# Start time
START_TIME=$(date +%s)

# ==========================================
# Phase 0: Environment Check
# ==========================================

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📋 Phase 0: Environment Check"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${RED}❌ Docker is not running. Please start Docker first.${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Docker is running${NC}"

# Check if services are up
echo ""
echo "Checking if services are running..."
if ! docker-compose -f "$PROJECT_ROOT/docker-compose.yml" ps | grep -q "Up"; then
    echo -e "${YELLOW}⚠️  Services are not running${NC}"
    echo ""
    read -p "Start all services? [y/N] " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        echo "Starting services..."
        cd "$PROJECT_ROOT"
        docker-compose up -d
        echo "Waiting for services to be ready (30 seconds)..."
        sleep 30
    else
        echo -e "${RED}❌ Cannot proceed without running services${NC}"
        exit 1
    fi
fi
echo -e "${GREEN}✅ Services are running${NC}"

# ==========================================
# Phase 1: Health Checks
# ==========================================

run_test_suite \
    "All Services Health Check" \
    "bash all_services_health_check.sh" \
    "$SCRIPT_DIR/e2e"

# ==========================================
# Phase 2: Unit Tests (Individual Services)
# ==========================================

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📦 Phase 2: Unit Tests (skipped - run separately if needed)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "💡 Tip: Run unit tests individually for each service:"
echo "   cd simple-servers/[service-name]"
echo "   go test ./..."
echo ""

# ==========================================
# Phase 3: Integration Tests
# ==========================================

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔗 Phase 3: Integration Tests"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Auth Flow Tests
if [ -f "$SCRIPT_DIR/integration/auth/run_test.sh" ]; then
    run_test_suite \
        "Auth Service Integration Tests" \
        "bash run_test.sh" \
        "$SCRIPT_DIR/integration/auth"
else
    echo -e "${YELLOW}⚠️  Auth integration tests not found - skipping${NC}"
fi

# Order Flow Tests
if [ -f "$SCRIPT_DIR/integration/order_flow/run_test.sh" ]; then
    run_test_suite \
        "Order-Payment Flow Tests" \
        "bash run_test.sh" \
        "$SCRIPT_DIR/integration/order_flow"
else
    echo -e "${YELLOW}⚠️  Order flow integration tests not found - skipping${NC}"
fi

# Notification Flow Tests
if [ -f "$SCRIPT_DIR/integration/notification_flow/run_all_tests.sh" ]; then
    run_test_suite \
        "Notification Flow Tests" \
        "bash run_all_tests.sh" \
        "$SCRIPT_DIR/integration/notification_flow"
else
    echo -e "${YELLOW}⚠️  Notification flow integration tests not found - skipping${NC}"
fi

# ==========================================
# Phase 4: E2E Tests
# ==========================================

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🌐 Phase 4: End-to-End Tests"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# E2E Complete Purchase Flow
run_test_suite \
    "E2E: Complete Purchase Flow" \
    "go test -v -run TestCompletePurchaseFlow -timeout 5m" \
    "$SCRIPT_DIR/e2e" \
    || true  # Continue even if failed

# E2E Error Scenarios
run_test_suite \
    "E2E: Error Scenarios" \
    "go test -v -run TestErrorScenarios -timeout 5m" \
    "$SCRIPT_DIR/e2e" \
    || true  # Continue even if failed

# E2E Performance Tests (Optional)
echo ""
echo -e "${YELLOW}⚠️  Performance tests can take 10+ minutes${NC}"
read -p "Run E2E performance tests? [y/N] " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Yy]$ ]]; then
    run_test_suite \
        "E2E: Performance Tests" \
        "go test -v -run TestPerformanceScenarios -timeout 15m" \
        "$SCRIPT_DIR/e2e" \
        || true  # Continue even if failed
else
    echo "⏭️  Skipping performance tests"
fi

# ==========================================
# Final Summary
# ==========================================

END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Test Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "Total Test Suites: $TOTAL_TESTS"
echo -e "${GREEN}Passed: $PASSED_TESTS${NC}"
echo -e "${RED}Failed: $FAILED_TESTS${NC}"
echo ""
echo "Total Duration: ${DURATION}s ($(date -u -d @${DURATION} +%H:%M:%S 2>/dev/null || echo $DURATION seconds))"
echo "Completed: $(date '+%Y-%m-%d %H:%M:%S')"
echo ""

if [ $FAILED_TESTS -eq 0 ]; then
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${GREEN}✅ ALL TESTS PASSED! 🎉${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "🎊 Congratulations! Your microservices are working correctly!"
    echo ""
    exit 0
else
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${RED}❌ SOME TESTS FAILED${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "Failed test suites:"
    for test_name in "${FAILED_TEST_NAMES[@]}"; do
        echo -e "  ${RED}✗${NC} $test_name"
    done
    echo ""
    echo "💡 Troubleshooting tips:"
    echo "  1. Check service logs: docker-compose logs [service-name]"
    echo "  2. Verify database connections and schemas"
    echo "  3. Check network connectivity between services"
    echo "  4. Review test output above for specific errors"
    echo "  5. Run failed tests individually for detailed debugging"
    echo ""
    echo "To re-run a specific test suite:"
    echo "  cd tests/[integration|e2e]"
    echo "  ./[test-script].sh"
    echo ""
    exit 1
fi
