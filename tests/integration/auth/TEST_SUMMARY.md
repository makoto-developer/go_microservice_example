# Auth Service Integration Test Summary

## Implementation Complete

Integration tests for the Auth Service have been successfully implemented.

## Files Created

```
tests/integration/auth/
├── auth_flow_test.go       # Main integration tests (4 test suites, 10 test cases)
├── client.go               # gRPC client wrapper
├── helpers.go              # Test utilities (DB, tokens, data generation)
├── go.mod                  # Go module with dependencies
├── cleanup.sh              # Database cleanup script (psql-based)
├── cleanup_docker.sh       # Database cleanup script (docker-based)
├── check_services.sh       # Service availability checker
├── run_test.sh             # Test execution script
├── README.md               # Comprehensive documentation
└── TEST_SUMMARY.md         # This file
```

## Test Coverage

### Test Suites

1. **TestCustomerAuthFlow** - Complete customer authentication flow
   - Registration
   - Login
   - Token refresh
   - Logout

2. **TestOwnerAuthFlow** - Complete owner authentication flow
   - Registration with business verification
   - Login
   - Get business verification status

3. **TestInvalidCredentials** - Error handling
   - Login with non-existent email
   - Login with wrong password

4. **TestDuplicateRegistration** - Duplicate email handling
   - Attempt duplicate registration

### Test Scenarios (10 Total)

#### Customer Auth Flow (4 scenarios)
1. Customer Registration
   - Creates new customer account
   - Verifies user_id and tokens are returned
   - Confirms user exists in database
   - Validates JWT access token

2. Customer Login
   - Authenticates with registered credentials
   - Verifies both access and refresh tokens
   - Validates token contents

3. Token Refresh
   - Refreshes access token using refresh token
   - Verifies new tokens are different from old
   - Validates new tokens

4. Customer Logout
   - Logs out with refresh token
   - Verifies refresh token is invalidated
   - Confirms token cannot be used after logout

#### Owner Auth Flow (3 scenarios)
1. Owner Registration
   - Creates new owner account
   - Verifies business verification status is "pending"
   - Confirms user exists in database
   - Validates JWT access token

2. Owner Login
   - Authenticates with registered credentials
   - Verifies business verification status
   - Validates tokens

3. Get Business Verification Status
   - Retrieves business verification status
   - Confirms status is "pending"

#### Error Handling (3 scenarios)
1. Login with Non-existent Email
   - Verifies proper error response

2. Login with Wrong Password
   - Verifies authentication failure

3. Duplicate Registration
   - Verifies duplicate email is rejected

## Configuration

### Service Endpoints
- **Auth Service**: localhost:22100 (gRPC)
- **PostgreSQL**: localhost:22010

### Database
- **Host**: localhost
- **Port**: 22010
- **User**: postgres
- **Password**: postgres
- **Database**: auth_service

### Tables
- `customer_users` - Customer account data
- `owner_users` - Owner account data
- `customer_refresh_tokens` - Customer refresh tokens
- `owner_refresh_tokens` - Owner refresh tokens

## Running Tests

### Quick Start
```bash
cd tests/integration/auth
./run_test.sh
```

### Manual Execution
```bash
# Check services
./check_services.sh

# Clean up
./cleanup_docker.sh

# Run tests
go test -v -timeout 30s ./...
```

### Run Specific Tests
```bash
go test -v -run TestCustomerAuthFlow
go test -v -run TestOwnerAuthFlow
go test -v -run TestInvalidCredentials
go test -v -run TestDuplicateRegistration
```

## Dependencies

### Go Modules
- `google.golang.org/grpc` - gRPC client
- `google.golang.org/protobuf` - Protocol Buffers
- `github.com/stretchr/testify` - Test assertions
- `github.com/golang-jwt/jwt/v5` - JWT validation
- `github.com/google/uuid` - UUID generation
- `github.com/lib/pq` - PostgreSQL driver

### External Services
- Auth Service (microservices/auth)
- PostgreSQL (Docker container: go_microservice_postgres_auth_dev)

## Key Features

### 1. Automatic Test Data Management
- Unique test emails generated for each test
- All test data prefixed with `test_`
- Automatic cleanup before and after tests

### 2. Database Verification
- Confirms users are created in database
- Verifies user data matches expectations
- Checks token invalidation after logout

### 3. JWT Token Validation
- Validates token signature
- Verifies token claims (user_id, expiration)
- Checks token format and structure

### 4. Service Health Checks
- Verifies Auth Service is running
- Confirms database connectivity
- Validates gRPC connection

### 5. Error Handling Tests
- Tests invalid credentials
- Verifies duplicate email detection
- Confirms proper error responses

## Test Data

### Email Format
```
test_XXXXXXXX@example.com
```
Where XXXXXXXX is a random 8-character UUID prefix.

### Password Format
```
Test@Password123
```
Strong password meeting security requirements.

## Success Criteria

All tests pass when:
1. Auth Service is running on localhost:22100
2. PostgreSQL is running on localhost:22010
3. Database tables exist (customer_users, owner_users)
4. JWT secret matches service configuration

## Cleanup

### Automatic
- `run_test.sh` cleans up before and after tests

### Manual
```bash
# Using Docker
./cleanup_docker.sh

# Using psql (if installed)
./cleanup.sh
```

## Integration with CI/CD

### GitHub Actions Example
```yaml
- name: Start Auth Service
  run: |
    cd microservices/auth
    ./auth-server &
    sleep 5

- name: Run Integration Tests
  run: |
    cd tests/integration/auth
    ./run_test.sh
```

## Troubleshooting

### Common Issues

1. **Auth Service not running**
   ```
   cd microservices/auth && ./auth-server
   ```

2. **Database not running**
   ```
   docker-compose up -d postgres_auth_dev
   ```

3. **Connection timeout**
   - Check service ports
   - Verify network connectivity
   - Review firewall settings

4. **JWT validation failed**
   - Verify JWT_SECRET matches auth service config
   - Check token expiration
   - Confirm token format

## Next Steps

### Potential Enhancements

1. **Email Verification Tests**
   - Add tests for email verification flow
   - Verify token generation and validation

2. **Password Reset Tests**
   - Test password reset request
   - Verify reset token validation
   - Confirm password update

3. **Token Expiration Tests**
   - Test expired access tokens
   - Verify refresh token expiration
   - Check automatic token refresh

4. **Concurrent Access Tests**
   - Multiple login attempts
   - Concurrent token refresh
   - Race condition testing

5. **Performance Tests**
   - Load testing with multiple users
   - Concurrent registration stress test
   - Token refresh performance

## Contact

For issues or questions, please refer to:
- `README.md` - Detailed documentation
- Auth Service implementation: `microservices/auth/`
- Proto definitions: `microservices/auth/proto/`
