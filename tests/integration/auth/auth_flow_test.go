package auth

import (
	"context"
	"testing"
	"time"

	customer_authv1 "github.com/makoto-developer/go_microservice_example/microservices/auth/proto/customer_auth/v1"
	owner_authv1 "github.com/makoto-developer/go_microservice_example/microservices/auth/proto/owner_auth/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMain sets up and tears down test environment
func TestMain(m *testing.M) {
	// Run tests
	m.Run()
}

// TestCustomerAuthFlow tests the complete customer authentication flow
func TestCustomerAuthFlow(t *testing.T) {
	// Setup
	client, err := NewTestClient(AuthServiceAddr)
	require.NoError(t, err, "Failed to create test client")
	defer client.Close()

	dbHelper := NewDBHelper(t)
	defer dbHelper.Close()

	tokenHelper := NewTokenHelper()
	dataGen := NewTestDataGenerator()

	// Generate test data
	testEmail := dataGen.GenerateTestEmail()
	testPassword := dataGen.GenerateTestPassword()

	// Cleanup before and after test
	defer dbHelper.CleanupTestData(t, testEmail)
	dbHelper.CleanupTestData(t, testEmail)

	// Create gRPC client
	customerClient := customer_authv1.NewCustomerAuthServiceClient(client.GetCustomerConn())

	t.Run("1. Customer Registration", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Register new customer
		resp, err := customerClient.Register(ctx, &customer_authv1.CustomerRegisterRequest{
			Email:    testEmail,
			Password: testPassword,
		})

		// Assertions
		require.NoError(t, err, "Registration should succeed")
		require.NotNil(t, resp, "Response should not be nil")

		assert.NotEmpty(t, resp.UserId, "User ID should be returned")
		assert.NotEmpty(t, resp.AccessToken, "Access token should be returned")
		assert.NotEmpty(t, resp.RefreshToken, "Refresh token should be returned")

		t.Logf("✅ Registered customer with user_id: %s", resp.UserId)

		// Verify user exists in database
		exists, err := dbHelper.UserExistsInCustomer(testEmail)
		require.NoError(t, err, "Database check should succeed")
		assert.True(t, exists, "User should exist in database")

		// Validate access token
		claims, err := tokenHelper.ValidateToken(resp.AccessToken)
		require.NoError(t, err, "Access token should be valid")
		assert.Equal(t, resp.UserId, claims["user_id"], "Token should contain correct user_id")

		t.Logf("✅ Access token is valid and contains user_id: %s", claims["user_id"])
	})

	t.Run("2. Customer Login", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Login with registered credentials
		resp, err := customerClient.Login(ctx, &customer_authv1.CustomerLoginRequest{
			Email:    testEmail,
			Password: testPassword,
		})

		// Assertions
		require.NoError(t, err, "Login should succeed")
		require.NotNil(t, resp, "Response should not be nil")

		assert.NotEmpty(t, resp.UserId, "User ID should be returned")
		assert.NotEmpty(t, resp.AccessToken, "Access token should be returned")
		assert.NotEmpty(t, resp.RefreshToken, "Refresh token should be returned")

		t.Logf("✅ Logged in customer with user_id: %s", resp.UserId)

		// Validate tokens
		accessClaims, err := tokenHelper.ValidateToken(resp.AccessToken)
		require.NoError(t, err, "Access token should be valid")
		assert.Equal(t, resp.UserId, accessClaims["user_id"], "Access token should contain correct user_id")

		refreshClaims, err := tokenHelper.ValidateToken(resp.RefreshToken)
		require.NoError(t, err, "Refresh token should be valid")
		assert.Equal(t, resp.UserId, refreshClaims["user_id"], "Refresh token should contain correct user_id")

		t.Logf("✅ Both access and refresh tokens are valid")
	})

	t.Run("3. Token Refresh", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// First login to get refresh token
		loginResp, err := customerClient.Login(ctx, &customer_authv1.CustomerLoginRequest{
			Email:    testEmail,
			Password: testPassword,
		})
		require.NoError(t, err, "Login should succeed")

		oldAccessToken := loginResp.AccessToken
		oldRefreshToken := loginResp.RefreshToken

		// Refresh token
		refreshResp, err := customerClient.RefreshToken(ctx, &customer_authv1.CustomerRefreshTokenRequest{
			RefreshToken: oldRefreshToken,
		})

		// Assertions
		require.NoError(t, err, "Token refresh should succeed")
		require.NotNil(t, refreshResp, "Response should not be nil")

		assert.NotEmpty(t, refreshResp.AccessToken, "New access token should be returned")
		assert.NotEmpty(t, refreshResp.RefreshToken, "New refresh token should be returned")
		assert.NotEqual(t, oldAccessToken, refreshResp.AccessToken, "New access token should be different")

		t.Logf("✅ Token refreshed successfully")

		// Validate new tokens
		_, err = tokenHelper.ValidateToken(refreshResp.AccessToken)
		require.NoError(t, err, "New access token should be valid")

		_, err = tokenHelper.ValidateToken(refreshResp.RefreshToken)
		require.NoError(t, err, "New refresh token should be valid")

		t.Logf("✅ New tokens are valid")
	})

	t.Run("4. Customer Logout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// First login to get refresh token
		loginResp, err := customerClient.Login(ctx, &customer_authv1.CustomerLoginRequest{
			Email:    testEmail,
			Password: testPassword,
		})
		require.NoError(t, err, "Login should succeed")

		// Logout
		logoutResp, err := customerClient.Logout(ctx, &customer_authv1.CustomerLogoutRequest{
			RefreshToken: loginResp.RefreshToken,
		})

		// Assertions
		require.NoError(t, err, "Logout should succeed")
		require.NotNil(t, logoutResp, "Response should not be nil")
		assert.True(t, logoutResp.Success, "Logout should be successful")

		t.Logf("✅ Logged out successfully")

		// Try to use refresh token after logout (should fail)
		_, err = customerClient.RefreshToken(ctx, &customer_authv1.CustomerRefreshTokenRequest{
			RefreshToken: loginResp.RefreshToken,
		})
		assert.Error(t, err, "Using refresh token after logout should fail")

		t.Logf("✅ Refresh token is invalidated after logout")
	})
}

// TestOwnerAuthFlow tests the complete owner authentication flow
func TestOwnerAuthFlow(t *testing.T) {
	// Setup
	client, err := NewTestClient(AuthServiceAddr)
	require.NoError(t, err, "Failed to create test client")
	defer client.Close()

	dbHelper := NewDBHelper(t)
	defer dbHelper.Close()

	tokenHelper := NewTokenHelper()
	dataGen := NewTestDataGenerator()

	// Generate test data
	testEmail := dataGen.GenerateTestEmail()
	testPassword := dataGen.GenerateTestPassword()

	// Cleanup before and after test
	defer dbHelper.CleanupTestData(t, testEmail)
	dbHelper.CleanupTestData(t, testEmail)

	// Create gRPC client
	ownerClient := owner_authv1.NewOwnerAuthServiceClient(client.GetOwnerConn())

	t.Run("1. Owner Registration", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Register new owner
		resp, err := ownerClient.Register(ctx, &owner_authv1.OwnerRegisterRequest{
			Email:    testEmail,
			Password: testPassword,
		})

		// Assertions
		require.NoError(t, err, "Registration should succeed")
		require.NotNil(t, resp, "Response should not be nil")

		assert.NotEmpty(t, resp.UserId, "User ID should be returned")
		assert.NotEmpty(t, resp.AccessToken, "Access token should be returned")
		assert.NotEmpty(t, resp.RefreshToken, "Refresh token should be returned")
		assert.Equal(t, "pending", resp.BusinessVerificationStatus, "Business verification status should be pending")

		t.Logf("✅ Registered owner with user_id: %s, verification status: %s",
			resp.UserId, resp.BusinessVerificationStatus)

		// Verify user exists in database
		exists, err := dbHelper.UserExistsInOwner(testEmail)
		require.NoError(t, err, "Database check should succeed")
		assert.True(t, exists, "User should exist in database")

		// Validate access token
		claims, err := tokenHelper.ValidateToken(resp.AccessToken)
		require.NoError(t, err, "Access token should be valid")
		assert.Equal(t, resp.UserId, claims["user_id"], "Token should contain correct user_id")

		t.Logf("✅ Access token is valid and contains user_id: %s", claims["user_id"])
	})

	t.Run("2. Owner Login", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Login with registered credentials
		resp, err := ownerClient.Login(ctx, &owner_authv1.OwnerLoginRequest{
			Email:    testEmail,
			Password: testPassword,
		})

		// Assertions
		require.NoError(t, err, "Login should succeed")
		require.NotNil(t, resp, "Response should not be nil")

		assert.NotEmpty(t, resp.UserId, "User ID should be returned")
		assert.NotEmpty(t, resp.AccessToken, "Access token should be returned")
		assert.NotEmpty(t, resp.RefreshToken, "Refresh token should be returned")
		assert.False(t, resp.BusinessVerified, "Business should not be verified yet")
		assert.Equal(t, "pending", resp.BusinessVerificationStatus, "Business verification status should be pending")

		t.Logf("✅ Logged in owner with user_id: %s, business_verified: %v",
			resp.UserId, resp.BusinessVerified)

		// Validate tokens
		_, err = tokenHelper.ValidateToken(resp.AccessToken)
		require.NoError(t, err, "Access token should be valid")

		_, err = tokenHelper.ValidateToken(resp.RefreshToken)
		require.NoError(t, err, "Refresh token should be valid")

		t.Logf("✅ Both access and refresh tokens are valid")
	})

	t.Run("3. Get Business Verification Status", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// First get user_id from login
		loginResp, err := ownerClient.Login(ctx, &owner_authv1.OwnerLoginRequest{
			Email:    testEmail,
			Password: testPassword,
		})
		require.NoError(t, err, "Login should succeed")

		// Get business verification status
		statusResp, err := ownerClient.GetBusinessVerificationStatus(ctx,
			&owner_authv1.OwnerGetBusinessVerificationStatusRequest{
				UserId: loginResp.UserId,
			})

		// Assertions
		require.NoError(t, err, "Getting verification status should succeed")
		require.NotNil(t, statusResp, "Response should not be nil")

		assert.False(t, statusResp.BusinessVerified, "Business should not be verified")
		assert.Equal(t, "pending", statusResp.BusinessVerificationStatus,
			"Business verification status should be pending")

		t.Logf("✅ Business verification status: %s, verified: %v",
			statusResp.BusinessVerificationStatus, statusResp.BusinessVerified)
	})
}

// TestInvalidCredentials tests authentication with invalid credentials
func TestInvalidCredentials(t *testing.T) {
	// Setup
	client, err := NewTestClient(AuthServiceAddr)
	require.NoError(t, err, "Failed to create test client")
	defer client.Close()

	customerClient := customer_authv1.NewCustomerAuthServiceClient(client.GetCustomerConn())

	t.Run("Login with non-existent email", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		_, err := customerClient.Login(ctx, &customer_authv1.CustomerLoginRequest{
			Email:    "nonexistent@example.com",
			Password: "password123",
		})

		assert.Error(t, err, "Login with non-existent email should fail")
		t.Logf("✅ Login with non-existent email correctly failed: %v", err)
	})

	t.Run("Login with wrong password", func(t *testing.T) {
		// First register a user
		dbHelper := NewDBHelper(t)
		defer dbHelper.Close()

		dataGen := NewTestDataGenerator()
		testEmail := dataGen.GenerateTestEmail()
		testPassword := dataGen.GenerateTestPassword()

		defer dbHelper.CleanupTestData(t, testEmail)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Register
		_, err := customerClient.Register(ctx, &customer_authv1.CustomerRegisterRequest{
			Email:    testEmail,
			Password: testPassword,
		})
		require.NoError(t, err, "Registration should succeed")

		// Try to login with wrong password
		_, err = customerClient.Login(ctx, &customer_authv1.CustomerLoginRequest{
			Email:    testEmail,
			Password: "WrongPassword123!",
		})

		assert.Error(t, err, "Login with wrong password should fail")
		t.Logf("✅ Login with wrong password correctly failed: %v", err)
	})
}

// TestDuplicateRegistration tests duplicate email registration
func TestDuplicateRegistration(t *testing.T) {
	// Setup
	client, err := NewTestClient(AuthServiceAddr)
	require.NoError(t, err, "Failed to create test client")
	defer client.Close()

	dbHelper := NewDBHelper(t)
	defer dbHelper.Close()

	dataGen := NewTestDataGenerator()
	testEmail := dataGen.GenerateTestEmail()
	testPassword := dataGen.GenerateTestPassword()

	defer dbHelper.CleanupTestData(t, testEmail)

	customerClient := customer_authv1.NewCustomerAuthServiceClient(client.GetCustomerConn())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Register first time
	_, err = customerClient.Register(ctx, &customer_authv1.CustomerRegisterRequest{
		Email:    testEmail,
		Password: testPassword,
	})
	require.NoError(t, err, "First registration should succeed")

	// Try to register again with same email
	_, err = customerClient.Register(ctx, &customer_authv1.CustomerRegisterRequest{
		Email:    testEmail,
		Password: testPassword,
	})

	assert.Error(t, err, "Duplicate registration should fail")
	t.Logf("✅ Duplicate registration correctly failed: %v", err)
}
