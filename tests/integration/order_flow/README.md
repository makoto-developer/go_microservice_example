# Order-Payment-Inventory Integration Tests

This directory contains integration tests for the complete order flow across multiple microservices.

## Overview

These tests verify the end-to-end functionality of the order processing system, including:

- Order creation
- Inventory reservation and confirmation
- Payment processing
- Rollback scenarios (payment failure, insufficient stock)
- Order cancellation

## Architecture

The tests integrate the following microservices:

```
┌─────────────┐     ┌──────────────┐     ┌─────────────┐
│  Customer   │────▶│    Order     │────▶│   Payment   │
│  Service    │     │   Service    │     │   Service   │
│  (22102)    │     │   (22104)    │     │   (22105)   │
└─────────────┘     └──────┬───────┘     └─────────────┘
                           │
                           ▼
                    ┌──────────────┐
                    │  Inventory   │
                    │   Service    │
                    │   (22103)    │
                    └──────────────┘
```

## Test Scenarios

### 1. Happy Path - Complete Order Flow

Tests the successful completion of an order:

1. ✅ Check inventory availability
2. 🔒 Reserve inventory
3. 📝 Create order (status: pending)
4. 💳 Process payment
5. ✅ Confirm order (status: confirmed)
6. 📦 Confirm inventory

**Expected Result**: Order is confirmed and inventory is decremented.

### 2. Payment Failure Rollback

Tests the rollback process when payment fails:

1. Reserve inventory
2. Create order
3. ❌ Payment fails
4. ↩️ Cancel order
5. ↩️ Release inventory reservation

**Expected Result**: Order is cancelled and inventory is released.

### 3. Insufficient Stock

Tests the behavior when stock is insufficient:

1. Check stock for large quantity
2. ❌ Stock check fails
3. Order is not created

**Expected Result**: System prevents order creation when stock is unavailable.

### 4. Order Cancellation

Tests customer-initiated order cancellation:

1. Reserve inventory
2. Create order
3. 🚫 Customer cancels order
4. ↩️ Release inventory

**Expected Result**: Order is cancelled and inventory is released.

## Prerequisites

### Services Running

Ensure the following services are running:

```bash
# Customer Service
cd simple-servers/customer
./customer-service

# Inventory Service
cd simple-servers/inventory
./inventory-service

# Order Service
cd simple-servers/order
./order-server

# Payment Service
cd simple-servers/payment
./payment-server
```

### Database Setup

Each service requires its own PostgreSQL database:

- Customer Service: `localhost:22012/customer_service`
- Inventory Service: `localhost:22013/inventory_service`
- Order Service: `localhost:22014/order_service`
- Payment Service: `localhost:22015/payment_service`

## Running Tests

### Basic Test Execution

```bash
# Run all tests
./run_test.sh

# Run with verbose output
./run_test.sh -v

# Run specific test
./run_test.sh -run TestOrderPaymentFlow

# Run with timeout
./run_test.sh -timeout 5m
```

### Benchmark Tests

```bash
# Run benchmarks
./run_test.sh -b

# Run benchmarks with memory stats
./run_test.sh -b -benchmem
```

### Coverage Report

```bash
# Generate coverage report
./run_test.sh -cover

# View coverage in browser
open coverage.html
```

## File Structure

```
order_flow/
├── README.md                  # This file
├── run_test.sh               # Test execution script
├── order_payment_test.go     # Main test file
├── test_data.go              # Test fixtures and helpers
├── inventory_client.go       # Inventory Service gRPC client
├── order_client.go           # Order Service gRPC client
└── payment_client.go         # Payment Service gRPC client
```

## Test Data

### Sample Order

```go
OrderRequest{
    CustomerID: "customer_uuid",
    Items: []OrderItem{
        {
            ProductID: "product_uuid",
            Quantity:  5,
            Price:     1000.0,
        },
    },
    TotalAmount: 5000.0,
}
```

### Sample Payment

```go
PaymentRequest{
    OrderID:     "order_uuid",
    Amount:      5000.0,
    Method:      "credit_card",
    CardNumber:  "4242424242424242", // Stripe test card
    ExpiryMonth: 12,
    ExpiryYear:  2025,
    CVV:         "123",
}
```

## Implementation Status

### Current Implementation

- ✅ Test structure and test cases
- ✅ Mock gRPC clients for all services
- ✅ Test scenarios (happy path, rollback, cancellation)
- ✅ Test execution script
- ✅ Documentation

### Future Enhancements

- ⏳ Real gRPC proto definitions
- ⏳ Actual service implementation calls
- ⏳ Database verification
- ⏳ Concurrent order tests
- ⏳ Reservation expiry tests
- ⏳ Performance benchmarks

## Troubleshooting

### Services Not Running

**Error**: `Failed to connect to inventory service`

**Solution**:
```bash
# Check if services are running
lsof -i :22102  # Customer
lsof -i :22103  # Inventory
lsof -i :22104  # Order
lsof -i :22105  # Payment

# Start missing services
cd simple-servers/<service>
./<service>-server
```

### Database Connection Issues

**Error**: `Failed to connect to database`

**Solution**:
```bash
# Check PostgreSQL is running
psql -h localhost -p 22012 -U postgres -d customer_service

# Verify connection strings in .env
cat .env | grep DATABASE_URL
```

### Test Timeout

**Error**: `test timed out after 10m0s`

**Solution**:
```bash
# Increase timeout
./run_test.sh -timeout 30m
```

## Development

### Adding New Tests

1. Create test function in `order_payment_test.go`:
```go
func TestNewScenario(t *testing.T) {
    // Test implementation
}
```

2. Add test data if needed in `test_data.go`
3. Update documentation
4. Run tests: `./run_test.sh -v`

### Adding New Service Client

1. Create `<service>_client.go`
2. Implement client struct and methods
3. Add to test setup
4. Update documentation

## References

- [Auth Service](../../../simple-servers/auth/)
- [Customer Service](../../../simple-servers/customer/)
- [Inventory Service](../../../simple-servers/inventory/)
- [Order Service](../../../simple-servers/order/)
- [Payment Service](../../../simple-servers/payment/)

## Contributing

When adding tests:

1. Follow existing test structure
2. Use descriptive test names
3. Add proper logging with emoji indicators
4. Include rollback scenarios
5. Update this README

## License

Part of the Go Microservice Example project.
