# Auth Service Integration Tests

This directory contains integration tests for the Auth Service, testing the complete authentication flow including customer and owner authentication.

## Overview

The integration tests verify:
- Customer registration, login, logout, and token refresh
- Owner registration, login, and business verification status
- Token validation and expiration
- Database state after operations
- Error handling for invalid credentials and duplicate registrations

## Prerequisites

### 1. Running Services

Before running the tests, ensure the following services are running:

```bash
# Auth Service (localhost:22100)
cd microservices/auth
./auth-server

# PostgreSQL Database (localhost:22010)
docker ps | grep postgres_auth_dev
```

### 2. Environment Variables

The tests use the following default configuration:

```bash
# Auth Service
AUTH_SERVICE_ADDR=localhost:22100

# Database
AUTH_DB_HOST=localhost
AUTH_DB_PORT=22010
AUTH_DB_USER=auth_user
AUTH_DB_PASSWORD=auth_password
AUTH_DB_NAME=auth_service

# JWT Secret (must match auth service config)
JWT_SECRET=your-secret-key-here-change-in-production
```

You can override these by setting environment variables.

## Test Structure

```
tests/integration/auth/
├── auth_flow_test.go    # Main integration tests
├── client.go            # gRPC client wrapper
├── helpers.go           # Test utilities (DB, tokens, data generation)
├── cleanup.sh           # Database cleanup script
├── run_test.sh          # Test execution script
├── go.mod               # Go module definition
└── README.md            # This file
```

## Running Tests

### Quick Start

```bash
cd tests/integration/auth
./run_test.sh
```

This script will:
1. Check if Auth Service is running
2. Check if database is running
3. Install dependencies
4. Clean up previous test data
5. Run all tests
6. Clean up test data after completion

### Manual Execution

```bash
# Install dependencies
go mod download

# Clean up test data
./cleanup.sh

# Run tests
go test -v -timeout 30s ./...

# Clean up after tests
./cleanup.sh
```

### Run Specific Tests

```bash
# Customer auth flow only
go test -v -run TestCustomerAuthFlow

# Owner auth flow only
go test -v -run TestOwnerAuthFlow

# Invalid credentials tests only
go test -v -run TestInvalidCredentials

# Duplicate registration test only
go test -v -run TestDuplicateRegistration
```

## Test Scenarios

### 1. Customer Auth Flow

Tests the complete customer authentication flow:

1. **Registration**
   - Create new customer account
   - Verify user_id and tokens are returned
   - Verify user exists in database
   - Validate access token

2. **Login**
   - Login with registered credentials
   - Verify tokens are returned
   - Validate both access and refresh tokens

3. **Token Refresh**
   - Refresh access token using refresh token
   - Verify new tokens are different from old ones
   - Validate new tokens

4. **Logout**
   - Logout with refresh token
   - Verify refresh token is invalidated
   - Verify token cannot be used after logout

### 2. Owner Auth Flow

Tests the complete owner authentication flow:

1. **Registration**
   - Create new owner account
   - Verify business verification status is "pending"
   - Verify user exists in database
   - Validate access token

2. **Login**
   - Login with registered credentials
   - Verify business verification status
   - Validate tokens

3. **Business Verification Status**
   - Get business verification status
   - Verify status is "pending"

### 3. Invalid Credentials

Tests error handling:

1. **Non-existent Email**
   - Attempt login with non-existent email
   - Verify login fails

2. **Wrong Password**
   - Register user
   - Attempt login with wrong password
   - Verify login fails

### 4. Duplicate Registration

Tests duplicate email handling:

1. Register user with email
2. Attempt to register again with same email
3. Verify registration fails

## Database Cleanup

### Automatic Cleanup

The `run_test.sh` script automatically cleans up test data before and after tests.

### Manual Cleanup

```bash
./cleanup.sh
```

This removes all users with emails starting with `test_` from both `customer_auth` and `owner_auth` tables.

## Test Data

Test data is automatically generated using:
- Random email addresses (test_XXXXXXXX@example.com)
- Strong passwords (Test@Password123)
- Unique UUIDs for user IDs

All test data is prefixed with `test_` for easy identification and cleanup.

## Troubleshooting

### Auth Service Not Running

```
❌ Error: Auth Service is not running at localhost:22100
```

**Solution**: Start the Auth Service:
```bash
cd microservices/auth
./auth-server
```

### Database Not Running

```
❌ Error: Database is not running at localhost:22010
```

**Solution**: Start the database using docker-compose:
```bash
docker-compose up -d postgres_auth_dev
```

### Connection Timeout

```
context deadline exceeded
```

**Solution**:
- Check if services are accessible
- Increase timeout in test code
- Check network connectivity

### Database Connection Failed

```
Failed to connect to database
```

**Solution**:
- Verify database credentials in `.env`
- Check if database is running
- Verify port 22010 is accessible

### Token Validation Failed

```
invalid token
```

**Solution**:
- Verify JWT_SECRET matches auth service config
- Check if token is expired
- Verify token format

## CI/CD Integration

### GitHub Actions Example

```yaml
- name: Run Auth Integration Tests
  run: |
    cd tests/integration/auth
    ./run_test.sh
  env:
    AUTH_SERVICE_ADDR: localhost:22100
    AUTH_DB_HOST: localhost
    AUTH_DB_PORT: 22010
```

## Expected Output

```
🧪 Auth Service Integration Tests
==================================

Configuration:
  - Auth Service: localhost:22100
  - Database: localhost:22010

🔍 Checking Auth Service availability...
✅ Auth Service is running

🔍 Checking database availability...
✅ Database is running

📦 Installing dependencies...
✅ Dependencies installed

🧹 Cleaning up previous test data...
✅ Cleanup complete!

🚀 Running integration tests...

=== RUN   TestCustomerAuthFlow
=== RUN   TestCustomerAuthFlow/1._Customer_Registration
✅ Registered customer with user_id: 123e4567-e89b-12d3-a456-426614174000
✅ Access token is valid and contains user_id: 123e4567-e89b-12d3-a456-426614174000
=== RUN   TestCustomerAuthFlow/2._Customer_Login
✅ Logged in customer with user_id: 123e4567-e89b-12d3-a456-426614174000
✅ Both access and refresh tokens are valid
=== RUN   TestCustomerAuthFlow/3._Token_Refresh
✅ Token refreshed successfully
✅ New tokens are valid
=== RUN   TestCustomerAuthFlow/4._Customer_Logout
✅ Logged out successfully
✅ Refresh token is invalidated after logout
--- PASS: TestCustomerAuthFlow (2.34s)

=== RUN   TestOwnerAuthFlow
=== RUN   TestOwnerAuthFlow/1._Owner_Registration
✅ Registered owner with user_id: 234e5678-e89b-12d3-a456-426614174001, verification status: pending
✅ Access token is valid and contains user_id: 234e5678-e89b-12d3-a456-426614174001
=== RUN   TestOwnerAuthFlow/2._Owner_Login
✅ Logged in owner with user_id: 234e5678-e89b-12d3-a456-426614174001, business_verified: false
✅ Both access and refresh tokens are valid
=== RUN   TestOwnerAuthFlow/3._Get_Business_Verification_Status
✅ Business verification status: pending, verified: false
--- PASS: TestOwnerAuthFlow (1.89s)

=== RUN   TestInvalidCredentials
=== RUN   TestInvalidCredentials/Login_with_non-existent_email
✅ Login with non-existent email correctly failed: rpc error: code = NotFound desc = user not found
=== RUN   TestInvalidCredentials/Login_with_wrong_password
✅ Login with wrong password correctly failed: rpc error: code = Unauthenticated desc = invalid credentials
--- PASS: TestInvalidCredentials (1.12s)

=== RUN   TestDuplicateRegistration
✅ Duplicate registration correctly failed: rpc error: code = AlreadyExists desc = email already exists
--- PASS: TestDuplicateRegistration (0.67s)

PASS
ok      github.com/makoto-developer/go_microservice_example/tests/integration/auth     6.021s

✅ All tests passed!

🧹 Cleaning up test data...
✅ Cleanup complete!
```

## Contributing

When adding new tests:

1. Follow the existing test structure
2. Use descriptive test names
3. Clean up test data after each test
4. Add appropriate assertions
5. Update this README with new test scenarios
