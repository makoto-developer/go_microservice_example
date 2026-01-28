#!/bin/bash

# E2E Test Runner Script
# Executes all E2E tests in order

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 E2E Test Runner - Complete Purchase Flow"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Step 0: Check if services are running
echo "📋 Step 0: Health Check"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if ! bash "$SCRIPT_DIR/all_services_health_check.sh"; then
    echo ""
    echo "❌ Services are not ready. Please start all services first."
    echo ""
    echo "To start services:"
    echo "  cd $PROJECT_ROOT"
    echo "  docker-compose up -d"
    echo ""
    exit 1
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 Running E2E Tests"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Navigate to tests directory
cd "$SCRIPT_DIR"

# Initialize Go module if not exists
if [ ! -f go.mod ]; then
    echo "Initializing Go module..."
    go mod init github.com/makoto-developer/go_microservice_example/tests/e2e
    go mod tidy
fi

# Run tests with verbose output
TEST_FAILED=false

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 Test Suite 1: Complete Purchase Flow"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if ! go test -v -run TestCompletePurchaseFlow -timeout 5m; then
    TEST_FAILED=true
    echo "❌ Complete Purchase Flow test failed"
else
    echo "✅ Complete Purchase Flow test passed"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 Test Suite 2: Error Scenarios"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if ! go test -v -run TestErrorScenarios -timeout 5m; then
    TEST_FAILED=true
    echo "❌ Error Scenarios test failed"
else
    echo "✅ Error Scenarios test passed"
fi

echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 Test Suite 3: Performance Tests (Optional)"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "⚠️  Performance tests can take several minutes"
read -p "Run performance tests? [y/N] " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    if ! go test -v -run TestPerformanceScenarios -timeout 10m; then
        TEST_FAILED=true
        echo "❌ Performance tests failed"
    else
        echo "✅ Performance tests passed"
    fi
else
    echo "⏭️  Skipping performance tests"
fi

# Summary
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 E2E Test Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

if [ "$TEST_FAILED" = true ]; then
    echo "❌ Some tests failed"
    echo ""
    echo "💡 Troubleshooting tips:"
    echo "  1. Check service logs: docker-compose logs [service-name]"
    echo "  2. Verify database connections"
    echo "  3. Check network connectivity between services"
    echo "  4. Review test output above for specific errors"
    echo ""
    exit 1
else
    echo "✅ All E2E tests passed successfully!"
    echo ""
    echo "🎉 Your microservices are working correctly!"
    echo ""
    exit 0
fi
