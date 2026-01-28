#!/bin/bash

# E2E Test - All Services Health Check Script

set -e

echo "🏥 E2E Health Check: Verifying all microservices..."

# Service definitions
declare -A SERVICES=(
    ["auth"]="http://localhost:8081/health"
    ["shop"]="http://localhost:8082/health"
    ["customer"]="http://localhost:8083/health"
    ["inventory"]="http://localhost:8084/health"
    ["order"]="http://localhost:8085/health"
    ["payment"]="http://localhost:8086/health"
    ["shipping"]="http://localhost:8089/health"
    ["notification"]="http://localhost:8088/health"
    ["review"]="http://localhost:8090/health"
    ["search"]="http://localhost:8092/health"
    ["chat"]="http://localhost:8091/health"
    ["admin"]="http://localhost:8093/health"
)

TIMEOUT=30
INTERVAL=2
FAILED_SERVICES=()

# Function to check a single service
check_service() {
    local name=$1
    local url=$2
    local elapsed=0

    echo -n "Checking $name service... "

    while [ $elapsed -lt $TIMEOUT ]; do
        if curl -s -f "$url" > /dev/null 2>&1; then
            echo "✅ OK"
            return 0
        fi
        sleep $INTERVAL
        elapsed=$((elapsed + INTERVAL))
    done

    echo "❌ TIMEOUT"
    FAILED_SERVICES+=("$name")
    return 1
}

# Check all services
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔍 Checking all services (timeout: ${TIMEOUT}s)..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

for service_name in "${!SERVICES[@]}"; do
    check_service "$service_name" "${SERVICES[$service_name]}" || true
done

# Summary
echo ""
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Health Check Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

total_services=${#SERVICES[@]}
failed_count=${#FAILED_SERVICES[@]}
success_count=$((total_services - failed_count))

echo "Total services: $total_services"
echo "Healthy: $success_count"
echo "Failed: $failed_count"
echo ""

if [ $failed_count -eq 0 ]; then
    echo "✅ All services are healthy!"
    echo ""
    exit 0
else
    echo "❌ Some services failed:"
    for service in "${FAILED_SERVICES[@]}"; do
        echo "   - $service"
    done
    echo ""
    echo "💡 Tip: Run 'docker-compose ps' to check container status"
    echo ""
    exit 1
fi
