#!/bin/bash

echo "🛑 Stopping All Microservices"
echo "=============================="
echo ""

# Stop all Go microservices
services=(
  "auth-server"
  "customer-service"
  "inventory-service"
  "order-server"
  "payment-server"
  "notification-service"
  "review-service"
  "shipping"
  "chat-service"
  "search-service"
  "admin-service"
)

for service in "${services[@]}"; do
  pids=$(pgrep -f "$service" | grep -v grep)
  
  if [ -n "$pids" ]; then
    echo "Stopping $service (PIDs: $pids)..."
    kill $pids 2>/dev/null
    echo "  ✅ Stopped"
  else
    echo "⚠️  $service not running"
  fi
done

# Stop Shop Service (Phoenix)
shop_pid=$(pgrep -f "shop-server")
if [ -n "$shop_pid" ]; then
  echo "Stopping Shop Service (PID: $shop_pid)..."
  kill $shop_pid 2>/dev/null
  echo "  ✅ Stopped"
fi

echo ""
echo "=============================="
echo "✅ All services stopped"
