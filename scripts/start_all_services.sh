#!/bin/bash

echo "🚀 Starting All Microservices"
echo "=============================="
echo ""

BASE_DIR="/Users/user/work/repositories/github.com/makoto-developer/go_microservice_example"

# Start Auth Service
echo "1. Starting Auth Service..."
cd "$BASE_DIR/microservices/auth"
if [ -f "./auth-server" ]; then
  ./auth-server > auth-server.log 2>&1 &
  echo "   ✅ Auth Service started (PID: $!)"
else
  echo "   ⚠️  Auth Service binary not found"
fi

# Start Shop Service (Phoenix)
echo "2. Starting Shop Service (Phoenix)..."
cd "$BASE_DIR/simple-servers/admin"
if pgrep -f "shop-server" > /dev/null; then
  echo "   ✅ Shop Service already running"
else
  echo "   ⚠️  Shop Service not running (start manually with: cd simple-servers/admin && mix phx.server)"
fi

# Start simple Go services
services=(
  "customer:$BASE_DIR/simple-servers/customer:customer-service"
  "inventory:$BASE_DIR/simple-servers/inventory:inventory-service"
  "order:$BASE_DIR/simple-servers/order:order-server"
  "payment:$BASE_DIR/simple-servers/payment:payment-server"
  "notification:$BASE_DIR/simple-servers/notification:notification-service"
  "review:$BASE_DIR/simple-servers/review:review-service"
  "shipping:$BASE_DIR/simple-servers/shipping:shipping"
  "chat:$BASE_DIR/simple-servers/chat:chat-service"
  "search:$BASE_DIR/simple-servers/search:search-service"
  "admin:$BASE_DIR/simple-servers/admin:admin-service"
)

counter=3
for service in "${services[@]}"; do
  IFS=':' read -r name dir binary <<< "$service"
  
  echo "$counter. Starting $name Service..."
  cd "$dir"
  
  if [ -f "./$binary" ]; then
    ./$binary > /tmp/$binary.log 2>&1 &
    echo "   ✅ $name Service started (PID: $!)"
  else
    echo "   ⚠️  $name Service binary not found at $dir/$binary"
  fi
  
  counter=$((counter + 1))
  sleep 0.5
done

echo ""
echo "=============================="
echo "✅ Service startup complete"
echo ""
echo "Run 'make status' to check service health"
