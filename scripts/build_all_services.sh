#!/bin/bash

echo "🔨 Building All Microservices"
echo "=============================="
echo ""

BASE_DIR="/Users/user/work/repositories/github.com/makoto-developer/go_microservice_example"

# Build Auth Service
echo "1. Building Auth Service..."
cd "$BASE_DIR/microservices/auth"
go build -o auth-server ./cmd/server/main.go
if [ $? -eq 0 ]; then
  echo "   ✅ Built successfully"
else
  echo "   ❌ Build failed"
fi

# Build simple Go services
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

counter=2
for service in "${services[@]}"; do
  IFS=':' read -r name dir binary <<< "$service"
  
  echo "$counter. Building $name Service..."
  cd "$dir"
  
  if [ -f "go.mod" ]; then
    go build -o "$binary" .
    if [ $? -eq 0 ]; then
      size=$(ls -lh "$binary" | awk '{print $5}')
      echo "   ✅ Built successfully ($size)"
    else
      echo "   ❌ Build failed"
    fi
  else
    echo "   ⚠️  No go.mod found"
  fi
  
  counter=$((counter + 1))
done

echo ""
echo "=============================="
echo "✅ Build complete"
