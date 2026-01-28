package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	authServiceURL     = "http://localhost:8081"
	customerServiceURL = "http://localhost:8083"
	searchServiceURL   = "http://localhost:8092"
	shopServiceURL     = "http://localhost:8082"
	inventoryServiceURL = "http://localhost:8084"
	orderServiceURL    = "http://localhost:8085"
	paymentServiceURL  = "http://localhost:8086"
	shippingServiceURL = "http://localhost:8089"
	notificationServiceURL = "http://localhost:8088"
	reviewServiceURL   = "http://localhost:8090"
)

type E2ETestContext struct {
	UserID          string
	Token           string
	CustomerID      string
	ProductID       string
	OrderID         string
	PaymentID       string
	ShipmentID      string
	TrackingNumber  string
	ReviewID        string
}

func TestCompletePurchaseFlow(t *testing.T) {
	ctx := context.Background()
	testCtx := &E2ETestContext{}

	// Wait for services to be ready
	t.Log("🔄 Waiting for all services to be ready...")
	require.NoError(t, waitForServices(ctx), "Services should be ready")

	// Step 1: User Registration
	t.Run("Step1_UserRegistration", func(t *testing.T) {
		testUserRegistration(t, testCtx)
	})

	// Step 2: Customer Profile Creation
	t.Run("Step2_CustomerProfileCreation", func(t *testing.T) {
		testCustomerProfileCreation(t, testCtx)
	})

	// Step 3: Product Search
	t.Run("Step3_ProductSearch", func(t *testing.T) {
		testProductSearch(t, testCtx)
	})

	// Step 4: Product Detail View
	t.Run("Step4_ProductDetailView", func(t *testing.T) {
		testProductDetailView(t, testCtx)
	})

	// Step 5: Inventory Check
	t.Run("Step5_InventoryCheck", func(t *testing.T) {
		testInventoryCheck(t, testCtx)
	})

	// Step 6: Order Creation
	t.Run("Step6_OrderCreation", func(t *testing.T) {
		testOrderCreation(t, testCtx)
	})

	// Step 7: Payment Processing
	t.Run("Step7_PaymentProcessing", func(t *testing.T) {
		testPaymentProcessing(t, testCtx)
	})

	// Step 8: Shipping Arrangement
	t.Run("Step8_ShippingArrangement", func(t *testing.T) {
		testShippingArrangement(t, testCtx)
	})

	// Step 9: Notification Verification
	t.Run("Step9_NotificationVerification", func(t *testing.T) {
		testNotificationVerification(t, testCtx)
	})

	// Step 10: Review Posting
	t.Run("Step10_ReviewPosting", func(t *testing.T) {
		testReviewPosting(t, testCtx)
	})

	t.Log("✅ Complete purchase flow test succeeded!")
}

// Step 1: User Registration
func testUserRegistration(t *testing.T, testCtx *E2ETestContext) {
	t.Log("Step 1: User Registration")

	payload := map[string]interface{}{
		"email":    fmt.Sprintf("test.user.%d@example.com", time.Now().Unix()),
		"password": "SecurePassword123!",
		"role":     "CUSTOMER",
	}

	resp, err := makeJSONRequest(http.MethodPost, authServiceURL+"/api/v1/register", nil, payload)
	require.NoError(t, err, "Registration request should succeed")
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "Should return 201 Created")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Should decode response")

	testCtx.UserID = result["user_id"].(string)
	testCtx.Token = result["token"].(string)

	assert.NotEmpty(t, testCtx.UserID, "Should receive user ID")
	assert.NotEmpty(t, testCtx.Token, "Should receive JWT token")

	t.Logf("✅ User registered: user_id=%s", testCtx.UserID)
}

// Step 2: Customer Profile Creation
func testCustomerProfileCreation(t *testing.T, testCtx *E2ETestContext) {
	t.Log("Step 2: Customer Profile Creation")

	payload := map[string]interface{}{
		"user_id": testCtx.UserID,
		"name":    "山田太郎",
		"phone":   "090-1234-5678",
		"addresses": []map[string]interface{}{
			{
				"postal_code": "100-0001",
				"prefecture":  "東京都",
				"city":        "千代田区",
				"address1":    "千代田1-1-1",
				"address2":    "マンション101",
				"is_default":  true,
			},
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + testCtx.Token,
	}

	resp, err := makeJSONRequest(http.MethodPost, customerServiceURL+"/api/v1/customers", headers, payload)
	require.NoError(t, err, "Customer creation request should succeed")
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "Should return 201 Created")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Should decode response")

	testCtx.CustomerID = result["customer_id"].(string)
	assert.NotEmpty(t, testCtx.CustomerID, "Should receive customer ID")

	t.Logf("✅ Customer profile created: customer_id=%s", testCtx.CustomerID)
}

// Step 3: Product Search
func testProductSearch(t *testing.T, testCtx *E2ETestContext) {
	t.Log("Step 3: Product Search")

	headers := map[string]string{
		"Authorization": "Bearer " + testCtx.Token,
	}

	resp, err := makeJSONRequest(http.MethodGet, searchServiceURL+"/api/v1/search?q=ノートPC", headers, nil)
	require.NoError(t, err, "Search request should succeed")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Should decode response")

	products := result["products"].([]interface{})
	require.NotEmpty(t, products, "Should return at least one product")

	firstProduct := products[0].(map[string]interface{})
	testCtx.ProductID = firstProduct["id"].(string)

	t.Logf("✅ Products found: count=%d, first_product_id=%s", len(products), testCtx.ProductID)
}

// Step 4: Product Detail View
func testProductDetailView(t *testing.T, testCtx *E2ETestContext) {
	t.Log("Step 4: Product Detail View")

	headers := map[string]string{
		"Authorization": "Bearer " + testCtx.Token,
	}

	url := fmt.Sprintf("%s/api/v1/products/%s", shopServiceURL, testCtx.ProductID)
	resp, err := makeJSONRequest(http.MethodGet, url, headers, nil)
	require.NoError(t, err, "Product detail request should succeed")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")

	var product map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&product)
	require.NoError(t, err, "Should decode response")

	assert.Equal(t, testCtx.ProductID, product["id"].(string), "Product ID should match")
	assert.NotEmpty(t, product["name"], "Product should have a name")
	assert.NotEmpty(t, product["price"], "Product should have a price")

	t.Logf("✅ Product details retrieved: name=%s, price=%v", product["name"], product["price"])
}

// Step 5: Inventory Check
func testInventoryCheck(t *testing.T, testCtx *E2ETestContext) {
	t.Log("Step 5: Inventory Check")

	headers := map[string]string{
		"Authorization": "Bearer " + testCtx.Token,
	}

	url := fmt.Sprintf("%s/api/v1/inventory/%s", inventoryServiceURL, testCtx.ProductID)
	resp, err := makeJSONRequest(http.MethodGet, url, headers, nil)
	require.NoError(t, err, "Inventory check request should succeed")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")

	var inventory map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&inventory)
	require.NoError(t, err, "Should decode response")

	quantity := int(inventory["quantity"].(float64))
	assert.Greater(t, quantity, 0, "Product should have stock available")

	t.Logf("✅ Inventory checked: available_quantity=%d", quantity)
}

// Step 6: Order Creation
func testOrderCreation(t *testing.T, testCtx *E2ETestContext) {
	t.Log("Step 6: Order Creation")

	payload := map[string]interface{}{
		"customer_id": testCtx.CustomerID,
		"items": []map[string]interface{}{
			{
				"product_id": testCtx.ProductID,
				"quantity":   1,
			},
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + testCtx.Token,
	}

	resp, err := makeJSONRequest(http.MethodPost, orderServiceURL+"/api/v1/orders", headers, payload)
	require.NoError(t, err, "Order creation request should succeed")
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "Should return 201 Created")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Should decode response")

	testCtx.OrderID = result["order_id"].(string)
	assert.NotEmpty(t, testCtx.OrderID, "Should receive order ID")

	t.Logf("✅ Order created: order_id=%s", testCtx.OrderID)
}

// Step 7: Payment Processing
func testPaymentProcessing(t *testing.T, testCtx *E2ETestContext) {
	t.Log("Step 7: Payment Processing")

	payload := map[string]interface{}{
		"order_id":       testCtx.OrderID,
		"payment_method": "CREDIT_CARD",
		"card_info": map[string]string{
			"card_number": "4242424242424242",
			"exp_month":   "12",
			"exp_year":    "2025",
			"cvc":         "123",
		},
	}

	headers := map[string]string{
		"Authorization": "Bearer " + testCtx.Token,
	}

	resp, err := makeJSONRequest(http.MethodPost, paymentServiceURL+"/api/v1/payments", headers, payload)
	require.NoError(t, err, "Payment request should succeed")
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "Should return 201 Created")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Should decode response")

	testCtx.PaymentID = result["payment_id"].(string)
	status := result["status"].(string)

	assert.NotEmpty(t, testCtx.PaymentID, "Should receive payment ID")
	assert.Equal(t, "COMPLETED", status, "Payment status should be COMPLETED")

	t.Logf("✅ Payment processed: payment_id=%s, status=%s", testCtx.PaymentID, status)
}

// Step 8: Shipping Arrangement
func testShippingArrangement(t *testing.T, testCtx *E2ETestContext) {
	t.Log("Step 8: Shipping Arrangement")

	// Wait for order to be confirmed
	time.Sleep(2 * time.Second)

	headers := map[string]string{
		"Authorization": "Bearer " + testCtx.Token,
	}

	url := fmt.Sprintf("%s/api/v1/shipments?order_id=%s", shippingServiceURL, testCtx.OrderID)
	resp, err := makeJSONRequest(http.MethodGet, url, headers, nil)
	require.NoError(t, err, "Shipment check request should succeed")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Should decode response")

	shipments := result["shipments"].([]interface{})
	require.NotEmpty(t, shipments, "Should have at least one shipment")

	shipment := shipments[0].(map[string]interface{})
	testCtx.ShipmentID = shipment["id"].(string)
	testCtx.TrackingNumber = shipment["tracking_number"].(string)

	assert.NotEmpty(t, testCtx.ShipmentID, "Should receive shipment ID")
	assert.NotEmpty(t, testCtx.TrackingNumber, "Should receive tracking number")

	t.Logf("✅ Shipment arranged: shipment_id=%s, tracking_number=%s", testCtx.ShipmentID, testCtx.TrackingNumber)
}

// Step 9: Notification Verification
func testNotificationVerification(t *testing.T, testCtx *E2ETestContext) {
	t.Log("Step 9: Notification Verification")

	// Wait for notification to be sent
	time.Sleep(2 * time.Second)

	headers := map[string]string{
		"Authorization": "Bearer " + testCtx.Token,
	}

	url := fmt.Sprintf("%s/api/v1/notifications?user_id=%s", notificationServiceURL, testCtx.UserID)
	resp, err := makeJSONRequest(http.MethodGet, url, headers, nil)
	require.NoError(t, err, "Notification check request should succeed")
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "Should return 200 OK")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Should decode response")

	notifications := result["notifications"].([]interface{})
	assert.NotEmpty(t, notifications, "Should have received notifications")

	t.Logf("✅ Notifications verified: count=%d", len(notifications))
}

// Step 10: Review Posting
func testReviewPosting(t *testing.T, testCtx *E2ETestContext) {
	t.Log("Step 10: Review Posting")

	payload := map[string]interface{}{
		"product_id":  testCtx.ProductID,
		"customer_id": testCtx.CustomerID,
		"order_id":    testCtx.OrderID,
		"rating":      5,
		"comment":     "とても良い商品でした！迅速な配送もありがとうございました。",
	}

	headers := map[string]string{
		"Authorization": "Bearer " + testCtx.Token,
	}

	resp, err := makeJSONRequest(http.MethodPost, reviewServiceURL+"/api/v1/reviews", headers, payload)
	require.NoError(t, err, "Review posting request should succeed")
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode, "Should return 201 Created")

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	require.NoError(t, err, "Should decode response")

	testCtx.ReviewID = result["review_id"].(string)
	assert.NotEmpty(t, testCtx.ReviewID, "Should receive review ID")

	t.Logf("✅ Review posted: review_id=%s", testCtx.ReviewID)
}

// Helper Functions

func waitForServices(ctx context.Context) error {
	services := []string{
		authServiceURL + "/health",
		customerServiceURL + "/health",
		searchServiceURL + "/health",
		shopServiceURL + "/health",
		inventoryServiceURL + "/health",
		orderServiceURL + "/health",
		paymentServiceURL + "/health",
		shippingServiceURL + "/health",
		notificationServiceURL + "/health",
		reviewServiceURL + "/health",
	}

	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for services to be ready")
		case <-ticker.C:
			allReady := true
			for _, service := range services {
				if !isServiceReady(service) {
					allReady = false
					break
				}
			}
			if allReady {
				return nil
			}
		}
	}
}

func isServiceReady(url string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func makeJSONRequest(method, url string, headers map[string]string, payload interface{}) (*http.Response, error) {
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	return client.Do(req)
}
