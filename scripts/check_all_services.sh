#!/bin/bash

echo ""

# Check services
services=(
  "Auth:auth-server:22100"
  "Shop:shop-server:4000"
  "Customer:customer-service:22102"
  "Inventory:inventory-service:22103"
  "Order:order-server:22104"
  "Payment:payment-server:22105"
  "Notification:notification-service:22106"
  "Review:review-service:22107"
  "Shipping:shipping:22108"
  "Chat:chat-service:22109"
  "Search:search-service:22110"
  "Admin:admin-service:22111"
)

running=0
total=${#services[@]}

for service in "${services[@]}"; do
  IFS=':' read -r name binary port <<< "$service"
  
  if pgrep -f "$binary" > /dev/null 2>&1; then
    pid=$(pgrep -f "$binary" | head -1)
    printf "✅ %-20s Running (PID: %5s, Port: %5s)\n" "$name Service" "$pid" "$port"
    running=$((running + 1))
  else
    printf "❌ %-20s Not running\n" "$name Service"
  fi
done

echo ""
echo "Summary: $running/$total services running"

# Check databases
echo ""
echo "Database Status:"
healthy=$(docker ps --filter "name=postgres" --filter "health=healthy" --format "{{.Names}}" 2>/dev/null | wc -l)
echo "  Healthy Databases: $healthy/12"
