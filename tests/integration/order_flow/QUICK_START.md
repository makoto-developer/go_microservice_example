# Quick Start Guide

Get the Order-Payment-Inventory integration tests running in 5 minutes.

## 1. Prerequisites Check

```bash
# Check Go version (requires 1.25+)
go version

# Check if nc (netcat) is available
which nc
```

## 2. Start Services

### Option A: Start All Services (Recommended)

```bash
# From project root
cd /Users/user/work/repositories/github.com/makoto-developer/go_microservice_example

# Start all required services
cd simple-servers/customer && ./customer-service &
cd ../inventory && ./inventory-service &
cd ../order && ./order-server &
cd ../payment && ./payment-server &

# Wait a few seconds for services to start
sleep 5
```

### Option B: Verify Services Are Running

```bash
# Check service ports
lsof -i :22102  # Customer Service
lsof -i :22103  # Inventory Service
lsof -i :22104  # Order Service
lsof -i :22105  # Payment Service
```

## 3. Run Tests

```bash
# Navigate to test directory
cd tests/integration/order_flow

# Run tests
./run_test.sh

# Or run with verbose output
./run_test.sh -v
```

## 4. Expected Output

### Successful Run

```
═══════════════════════════════════════════════════════════
   Order-Payment-Inventory Integration Test Runner
═══════════════════════════════════════════════════════════

Checking prerequisites...

Checking Customer Service... ✓ Running
Checking Inventory Service... ✓ Running
Checking Order Service... ✓ Running
Checking Payment Service... ✓ Running

═══════════════════════════════════════════════════════════
   Running Integration Tests
═══════════════════════════════════════════════════════════

Running integration tests...
PASS
ok      github.com/makoto-developer/go_microservice_example/tests/integration/order_flow       0.500s

═══════════════════════════════════════════════════════════
   ✨ All tests passed!
═══════════════════════════════════════════════════════════
```

### Verbose Output Example

```
=== RUN   TestOrderPaymentFlow
=== RUN   TestOrderPaymentFlow/Happy_Path_-_Complete_Order_Flow
    order_payment_test.go:XX: 🎯 Starting Complete Order Flow Test
    order_payment_test.go:XX: 📦 Step 1: Checking inventory availability...
    order_payment_test.go:XX: ✅ Stock is available
    order_payment_test.go:XX: 🔒 Step 2: Reserving inventory...
    order_payment_test.go:XX: ✅ Inventory reserved: res_1234567890
    order_payment_test.go:XX: 📝 Step 3: Creating order...
    order_payment_test.go:XX: ✅ Order created: ord_1234567890 (status: pending)
    order_payment_test.go:XX: 💳 Step 4: Processing payment...
    order_payment_test.go:XX: ✅ Payment completed: pay_1234567890
    order_payment_test.go:XX: ✅ Step 5: Confirming order...
    order_payment_test.go:XX: ✅ Order confirmed
    order_payment_test.go:XX: 📦 Step 6: Confirming inventory...
    order_payment_test.go:XX: ✅ Inventory confirmed
    order_payment_test.go:XX: ✨ Complete Order Flow Test PASSED
--- PASS: TestOrderPaymentFlow (0.30s)
    --- PASS: TestOrderPaymentFlow/Happy_Path_-_Complete_Order_Flow (0.30s)
PASS
```

## 5. Running Specific Tests

```bash
# Run only happy path test
./run_test.sh -run TestOrderPaymentFlow/Happy_Path

# Run rollback tests
./run_test.sh -run Rollback

# Run with timeout
./run_test.sh -timeout 5m
```

## 6. Benchmarks

```bash
# Run performance benchmarks
./run_test.sh -b

# Run benchmarks with memory stats
./run_test.sh -b -benchmem
```

Expected benchmark output:
```
goos: darwin
goarch: arm64
BenchmarkOrderFlow-10      1000      1000000 ns/op
PASS
```

## 7. Coverage Report

```bash
# Generate coverage report
./run_test.sh -cover

# Open coverage report in browser
open coverage.html
```

## Common Issues

### Issue: Services Not Running

**Symptom**: `✗ Not running` for one or more services

**Solution**:
```bash
# Start the missing service
cd simple-servers/<service>
./<service>-server
```

### Issue: Port Already in Use

**Symptom**: `Failed to listen: address already in use`

**Solution**:
```bash
# Find process using the port
lsof -i :22102

# Kill the process
kill -9 <PID>
```

### Issue: Database Connection Error

**Symptom**: `Failed to connect to database`

**Solution**:
```bash
# Check if PostgreSQL is running
brew services list | grep postgresql

# Start PostgreSQL if needed
brew services start postgresql
```

## Next Steps

- Read [README.md](README.md) for detailed documentation
- Explore test scenarios in [order_payment_test.go](order_payment_test.go)
- Add your own test cases
- Integrate with CI/CD pipeline

## Testing Checklist

- [ ] Go 1.25+ installed
- [ ] All 4 services running
- [ ] PostgreSQL databases available
- [ ] Tests pass with `./run_test.sh`
- [ ] Tests pass with `./run_test.sh -v`
- [ ] Coverage report generated

## Support

For issues or questions:

1. Check [README.md](README.md) Troubleshooting section
2. Review service logs in `simple-servers/<service>/`
3. Verify database connections
4. Check service health endpoints

---

**Ready to test?** Run `./run_test.sh` and watch the magic happen! ✨
