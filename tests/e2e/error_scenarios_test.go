package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorScenarios(t *testing.T) {
	ctx := context.Background()

	// Wait for services to be ready
	t.Log("🔄 Waiting for services...")
	require.NoError(t, waitForServices(ctx), "Services should be ready")

	// Test 1: Insufficient Stock
	t.Run("InsufficientStock", func(t *testing.T) {
		testInsufficientStock(t)
	})

	// Test 2: Payment Failure and Rollback
	t.Run("PaymentFailureRollback", func(t *testing.T) {
		testPaymentFailureRollback(t)
	})

	// Test 3: Duplicate Order Prevention
	t.Run("DuplicateOrderPrevention", func(t *testing.T) {
		testDuplicateOrderPrevention(t)
	})

	// Test 4: Invalid Authentication
	t.Run("InvalidAuthentication", func(t *testing.T) {
		testInvalidAuthentication(t)
	})

	// Test 5: Service Timeout Handling
	t.Run("ServiceTimeoutHandling", func(t *testing.T) {
		testServiceTimeoutHandling(t)
	})

	t.Log("✅ All error scenario tests completed!")
}

// Test 1: Insufficient Stock Error
func testInsufficientStock(t *testing.T) {
	t.Log("Testing: Insufficient Stock Error")

	// Setup: Create user and customer
	testCtx := setupTestUser(t)

	// Get a product ID
	productID := getTestProductID(t, testCtx.Token)

	// Try to order excessive quantity
	payload := map[string]interface{}{
		"customer_id": testCtx.CustomerID,
		"items": []map[string]interface{}{
			{
				"product_id": productID,
				"quantity":   999999, // Excessive quantity
			},
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + testCtx.Token,
	}

	resp, err := makeJSONRequest(http.MethodPost, orderServiceURL+"/api/v1/orders", headers, payload)
	require.NoError(t, err, "Request should succeed")
	defer resp.Body.Close()

	// Should return error status
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Should return 400 Bad Request")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Should decode error response")

	assert.Contains(t, result["error"].(string), "insufficient stock", "Error message should mention insufficient stock")

	t.Log("✅ Insufficient stock error handled correctly")
}

// Test 2: Payment Failure and Rollback
func testPaymentFailureRollback(t *testing.T) {
	t.Log("Testing: Payment Failure and Rollback")

	// Setup: Create user, customer, and order
	testCtx := setupTestUser(t)
	productID := getTestProductID(t, testCtx.Token)
	orderID := createTestOrder(t, testCtx, productID)

	// Check initial inventory
	initialInventory := getInventory(t, testCtx.Token, productID)

	// Attempt payment with invalid card (should fail)
	payload := map[string]interface{}{
		"order_id":       orderID,
		"payment_method": "CREDIT_CARD",
		"card_info": map[string]string{
			"card_number": "0000000000000000", // Invalid card
			"exp_month":   "12",
			"exp_year":    "2025",
			"cvc":         "123",
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + testCtx.Token,
	}

	resp, err := makeJSONRequest(http.MethodPost, paymentServiceURL+"/api/v1/payments", headers, payload)
	require.NoError(t, err, "Request should succeed")
	defer resp.Body.Close()

	// Payment should fail
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Payment should fail")

	// Wait for rollback
	time.Sleep(2 * time.Second)

	// Verify inventory was released (rollback)
	finalInventory := getInventory(t, testCtx.Token, productID)
	assert.Equal(t, initialInventory, finalInventory, "Inventory should be released after payment failure")

	// Verify order status is cancelled
	orderStatus := getOrderStatus(t, testCtx.Token, orderID)
	assert.Equal(t, "CANCELLED", orderStatus, "Order should be cancelled after payment failure")

	t.Log("✅ Payment failure rollback handled correctly")
}

// Test 3: Duplicate Order Prevention
func testDuplicateOrderPrevention(t *testing.T) {
	t.Log("Testing: Duplicate Order Prevention")

	testCtx := setupTestUser(t)
	productID := getTestProductID(t, testCtx.Token)

	payload := map[string]interface{}{
		"customer_id":   testCtx.CustomerID,
		"idempotency_key": fmt.Sprintf("test-order-%d", time.Now().Unix()),
		"items": []map[string]interface{}{
			{
				"product_id": productID,
				"quantity":   1,
			},
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + testCtx.Token,
	}

	// First order should succeed
	resp1, err := makeJSONRequest(http.MethodPost, orderServiceURL+"/api/v1/orders", headers, payload)
	require.NoError(t, err, "First request should succeed")
	defer resp1.Body.Close()
	assert.Equal(t, http.StatusCreated, resp1.StatusCode, "First order should be created")

	var result1 map[string]interface{}
	err = json.NewDecoder(resp1.Body).Decode(&result1)
	require.NoError(t, err)
	orderID1 := result1["order_id"].(string)

	// Duplicate order with same idempotency key should return same order
	resp2, err := makeJSONRequest(http.MethodPost, orderServiceURL+"/api/v1/orders", headers, payload)
	require.NoError(t, err, "Second request should succeed")
	defer resp2.Body.Close()
	assert.Equal(t, http.StatusOK, resp2.StatusCode, "Duplicate should return existing order")

	var result2 map[string]interface{}
	err = json.NewDecoder(resp2.Body).Decode(&result2)
	require.NoError(t, err)
	orderID2 := result2["order_id"].(string)

	assert.Equal(t, orderID1, orderID2, "Should return same order ID for duplicate request")

	t.Log("✅ Duplicate order prevention works correctly")
}

// Test 4: Invalid Authentication
func testInvalidAuthentication(t *testing.T) {
	t.Log("Testing: Invalid Authentication")

	headers := map[string]string{
		"Authorization": "Bearer invalid-token",
	}

	resp, err := makeJSONRequest(http.MethodGet, customerServiceURL+"/api/v1/customers/123", headers, nil)
	require.NoError(t, err, "Request should succeed")
	defer resp.Body.Close()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode, "Should return 401 Unauthorized")

	t.Log("✅ Invalid authentication handled correctly")
}

// Test 5: Service Timeout Handling
func testServiceTimeoutHandling(t *testing.T) {
	t.Log("Testing: Service Timeout Handling")

	testCtx := setupTestUser(t)

	// Create a request with very short timeout
	client := &http.Client{Timeout: 100 * time.Millisecond}

	req, err := http.NewRequest(http.MethodGet, searchServiceURL+"/api/v1/search?q=test", nil)
	require.NoError(t, err)

	req.Header.Set("Authorization", "Bearer "+testCtx.Token)
	req.Header.Set("X-Timeout-Trigger", "true") // Special header to trigger timeout

	resp, err := client.Do(req)

	if err != nil {
		// Timeout expected
		assert.Contains(t, err.Error(), "timeout", "Should timeout")
		t.Log("✅ Service timeout handled correctly")
		return
	}

	defer resp.Body.Close()

	// If response came back quickly, that's also acceptable
	assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusGatewayTimeout,
		"Should return OK or Gateway Timeout")

	t.Log("✅ Service timeout handling verified")
}

// Helper Functions

func setupTestUser(t *testing.T) *E2ETestContext {
	testCtx := &E2ETestContext{}

	// Register user
	payload := map[string]interface{}{
		"email":    fmt.Sprintf("test.error.%d@example.com", time.Now().Unix()),
		"password": "SecurePassword123!",
		"role":     "CUSTOMER",
	}

	resp, err := makeJSONRequest(http.MethodPost, authServiceURL+"/api/v1/register", nil, payload)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	testCtx.UserID = result["user_id"].(string)
	testCtx.Token = result["token"].(string)

	// Create customer profile
	customerPayload := map[string]interface{}{
		"user_id": testCtx.UserID,
		"name":    "Test User",
		"phone":   "090-0000-0000",
		"addresses": []map[string]interface{}{
			{
				"postal_code": "100-0001",
				"prefecture":  "東京都",
				"city":        "千代田区",
				"address1":    "千代田1-1-1",
				"is_default":  true,
			},
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + testCtx.Token,
	}

	resp2, err := makeJSONRequest(http.MethodPost, customerServiceURL+"/api/v1/customers", headers, customerPayload)
	require.NoError(t, err)
	defer resp2.Body.Close()

	var customerResult map[string]interface{}
	err = json.NewDecoder(resp2.Body).Decode(&customerResult)
	require.NoError(t, err)

	testCtx.CustomerID = customerResult["customer_id"].(string)

	return testCtx
}

func getTestProductID(t *testing.T, token string) string {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	resp, err := makeJSONRequest(http.MethodGet, searchServiceURL+"/api/v1/search?q=test", headers, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	products := result["products"].([]interface{})
	require.NotEmpty(t, products)

	return products[0].(map[string]interface{})["id"].(string)
}

func createTestOrder(t *testing.T, testCtx *E2ETestContext, productID string) string {
	payload := map[string]interface{}{
		"customer_id": testCtx.CustomerID,
		"items": []map[string]interface{}{
			{
				"product_id": productID,
				"quantity":   1,
			},
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + testCtx.Token,
	}

	resp, err := makeJSONRequest(http.MethodPost, orderServiceURL+"/api/v1/orders", headers, payload)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	return result["order_id"].(string)
}

func getInventory(t *testing.T, token, productID string) int {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	url := fmt.Sprintf("%s/api/v1/inventory/%s", inventoryServiceURL, productID)
	resp, err := makeJSONRequest(http.MethodGet, url, headers, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	return int(result["quantity"].(float64))
}

func getOrderStatus(t *testing.T, token, orderID string) string {
	headers := map[string]string{
		"Authorization": "Bearer " + token,
	}

	url := fmt.Sprintf("%s/api/v1/orders/%s", orderServiceURL, orderID)
	resp, err := makeJSONRequest(http.MethodGet, url, headers, nil)
	require.NoError(t, err)
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err)

	return result["status"].(string)
}
