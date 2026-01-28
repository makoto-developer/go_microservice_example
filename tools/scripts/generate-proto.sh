#!/bin/bash

set -e

PROTO_DIR="proto"
OUT_DIR="proto"

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed"
    echo "Install with: brew install protobuf"
    exit 1
fi

# Check if protoc-gen-go is installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Error: protoc-gen-go is not installed"
    echo "Install with: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest"
    exit 1
fi

# Check if protoc-gen-go-grpc is installed
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Error: protoc-gen-go-grpc is not installed"
    echo "Install with: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
    exit 1
fi

# Generate proto for all services
SERVICES=(
    "auth-service"
    "shop-service"
    "customer-service"
    "inventory-service"
    "order-service"
    "payment-service"
    "shipping-service"
    "notification-service"
    "review-service"
    "chat-service"
    "search-service"
    "admin-service"
)

echo "Generating Protocol Buffers code..."

for service in "${SERVICES[@]}"; do
    echo "Processing $service..."
    
    proto_file="${PROTO_DIR}/${service}/v1/${service}.proto"
    
    if [ ! -f "$proto_file" ]; then
        echo "Warning: $proto_file not found, skipping"
        continue
    fi
    
    protoc \
        --go_out="${OUT_DIR}" \
        --go_opt=module=github.com/makoto-developer/go_microservice_example/proto \
        --go-grpc_out="${OUT_DIR}" \
        --go-grpc_opt=module=github.com/makoto-developer/go_microservice_example/proto \
        --proto_path="${PROTO_DIR}" \
        "${proto_file}"
    
    echo "✓ Generated code for $service"
done

echo ""
echo "All Protocol Buffers code generated successfully!"
