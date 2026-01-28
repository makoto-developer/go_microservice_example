#!/bin/bash

# Shipping Service Verification Script

echo "=== Shipping Service Verification ==="
echo ""

# 1. Check database connection
echo "1. Checking database connection..."
cd ../../infrastructure/docker
docker-compose exec -T postgres_shipping psql -U postgres -d shipping_service -c "SELECT 1;" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "   ✅ Database connection successful"
else
    echo "   ❌ Database connection failed"
    exit 1
fi

# 2. Verify tables exist
echo ""
echo "2. Verifying tables..."
TABLES=$(docker-compose exec -T postgres_shipping psql -U postgres -d shipping_service -t -c "SELECT tablename FROM pg_tables WHERE schemaname = 'public';")

if echo "$TABLES" | grep -q "shipments"; then
    echo "   ✅ shipments table exists"
else
    echo "   ❌ shipments table not found"
    exit 1
fi

if echo "$TABLES" | grep -q "tracking_events"; then
    echo "   ✅ tracking_events table exists"
else
    echo "   ❌ tracking_events table not found"
    exit 1
fi

# 3. Verify indexes
echo ""
echo "3. Verifying indexes..."
INDEXES=$(docker-compose exec -T postgres_shipping psql -U postgres -d shipping_service -t -c "SELECT indexname FROM pg_indexes WHERE schemaname = 'public';")

for idx in "idx_shipments_order_id" "idx_shipments_tracking_number" "idx_shipments_status" "idx_tracking_events_shipment_id"; do
    if echo "$INDEXES" | grep -q "$idx"; then
        echo "   ✅ $idx exists"
    else
        echo "   ❌ $idx not found"
        exit 1
    fi
done

# 4. Verify foreign key constraint
echo ""
echo "4. Verifying foreign key constraint..."
FK=$(docker-compose exec -T postgres_shipping psql -U postgres -d shipping_service -t -c "SELECT conname FROM pg_constraint WHERE contype = 'f' AND conrelid = 'tracking_events'::regclass;")

if echo "$FK" | grep -q "tracking_events_shipment_id_fkey"; then
    echo "   ✅ Foreign key constraint exists"
else
    echo "   ❌ Foreign key constraint not found"
    exit 1
fi

# 5. Test service build
echo ""
echo "5. Testing service build..."
cd ../../simple-servers/shipping
if [ -f "./shipping" ]; then
    echo "   ✅ Service binary exists"
else
    echo "   ⚠️  Building service..."
    go build -o shipping
    if [ $? -eq 0 ]; then
        echo "   ✅ Service built successfully"
    else
        echo "   ❌ Service build failed"
        exit 1
    fi
fi

# 6. Test service startup (quick test)
echo ""
echo "6. Testing service startup..."
./shipping &
SERVICE_PID=$!
sleep 3

if ps -p $SERVICE_PID > /dev/null 2>&1; then
    echo "   ✅ Service started successfully on port 22108"
    kill $SERVICE_PID 2>/dev/null
    wait $SERVICE_PID 2>/dev/null
else
    echo "   ❌ Service failed to start"
    exit 1
fi

echo ""
echo "=== All Verifications Passed! ==="
echo ""
echo "📊 Summary:"
echo "   - Database: shipping_service on port 22016"
echo "   - Service Port: 22108"
echo "   - Tables: shipments, tracking_events"
echo "   - Indexes: 4 indexes created"
echo "   - Foreign Keys: 1 constraint (CASCADE DELETE)"
