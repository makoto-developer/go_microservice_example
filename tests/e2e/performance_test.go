package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerformanceScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	ctx := context.Background()

	// Wait for services
	t.Log("🔄 Waiting for services...")
	require.NoError(t, waitForServices(ctx), "Services should be ready")

	// Test 1: Concurrent User Registration
	t.Run("ConcurrentUserRegistration", func(t *testing.T) {
		testConcurrentUserRegistration(t)
	})

	// Test 2: Concurrent Order Creation
	t.Run("ConcurrentOrderCreation", func(t *testing.T) {
		testConcurrentOrderCreation(t)
	})

	// Test 3: Search Performance
	t.Run("SearchPerformance", func(t *testing.T) {
		testSearchPerformance(t)
	})

	// Test 4: Database Connection Pool
	t.Run("DatabaseConnectionPool", func(t *testing.T) {
		testDatabaseConnectionPool(t)
	})

	t.Log("✅ Performance tests completed!")
}

// Test 1: Concurrent User Registration
func testConcurrentUserRegistration(t *testing.T) {
	t.Log("Testing: Concurrent User Registration (100 users)")

	concurrentUsers := 100
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	start := time.Now()

	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			payload := map[string]interface{}{
				"email":    fmt.Sprintf("perf.test.%d.%d@example.com", time.Now().Unix(), index),
				"password": "SecurePassword123!",
				"role":     "CUSTOMER",
			}

			resp, err := makeJSONRequest(http.MethodPost, authServiceURL+"/api/v1/register", nil, payload)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusCreated {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	t.Logf("✅ Concurrent registrations: %d/%d succeeded in %s", successCount, concurrentUsers, elapsed)
	t.Logf("   Average: %.2f registrations/second", float64(successCount)/elapsed.Seconds())

	assert.GreaterOrEqual(t, successCount, concurrentUsers*90/100, "At least 90% should succeed")
	assert.Less(t, elapsed, 30*time.Second, "Should complete within 30 seconds")
}

// Test 2: Concurrent Order Creation
func testConcurrentOrderCreation(t *testing.T) {
	t.Log("Testing: Concurrent Order Creation (50 orders)")

	// Setup: Create test users
	numOrders := 50
	testUsers := make([]*E2ETestContext, numOrders)

	t.Log("Setting up test users...")
	for i := 0; i < numOrders; i++ {
		testUsers[i] = setupTestUser(t)
	}

	productID := getTestProductID(t, testUsers[0].Token)

	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	start := time.Now()

	for i := 0; i < numOrders; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			payload := map[string]interface{}{
				"customer_id": testUsers[index].CustomerID,
				"items": []map[string]interface{}{
					{
						"product_id": productID,
						"quantity":   1,
					},
				},
			}

			headers := map[string]string{
				"Authorization": "Bearer " + testUsers[index].Token,
			}

			resp, err := makeJSONRequest(http.MethodPost, orderServiceURL+"/api/v1/orders", headers, payload)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusCreated {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	elapsed := time.Since(start)

	t.Logf("✅ Concurrent orders: %d/%d succeeded in %s", successCount, numOrders, elapsed)
	t.Logf("   Average: %.2f orders/second", float64(successCount)/elapsed.Seconds())

	assert.GreaterOrEqual(t, successCount, numOrders*80/100, "At least 80% should succeed")
	assert.Less(t, elapsed, 60*time.Second, "Should complete within 60 seconds")
}

// Test 3: Search Performance
func testSearchPerformance(t *testing.T) {
	t.Log("Testing: Search Performance (100 concurrent searches)")

	testCtx := setupTestUser(t)
	queries := []string{"ノートPC", "スマートフォン", "タブレット", "カメラ", "時計"}

	concurrentSearches := 100
	var wg sync.WaitGroup
	responseTimes := make([]time.Duration, 0, concurrentSearches)
	var mu sync.Mutex

	start := time.Now()

	for i := 0; i < concurrentSearches; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			query := queries[index%len(queries)]
			headers := map[string]string{
				"Authorization": "Bearer " + testCtx.Token,
			}

			searchStart := time.Now()
			resp, err := makeJSONRequest(http.MethodGet,
				fmt.Sprintf("%s/api/v1/search?q=%s", searchServiceURL, query),
				headers, nil)
			searchElapsed := time.Since(searchStart)

			if err != nil {
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				mu.Lock()
				responseTimes = append(responseTimes, searchElapsed)
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	totalElapsed := time.Since(start)

	// Calculate statistics
	if len(responseTimes) > 0 {
		var total time.Duration
		min := responseTimes[0]
		max := responseTimes[0]

		for _, rt := range responseTimes {
			total += rt
			if rt < min {
				min = rt
			}
			if rt > max {
				max = rt
			}
		}

		avg := total / time.Duration(len(responseTimes))

		t.Logf("✅ Search performance:")
		t.Logf("   Successful searches: %d/%d", len(responseTimes), concurrentSearches)
		t.Logf("   Total time: %s", totalElapsed)
		t.Logf("   Average response time: %s", avg)
		t.Logf("   Min response time: %s", min)
		t.Logf("   Max response time: %s", max)
		t.Logf("   Throughput: %.2f searches/second", float64(len(responseTimes))/totalElapsed.Seconds())

		assert.Less(t, avg, 500*time.Millisecond, "Average response time should be under 500ms")
		assert.Less(t, max, 2*time.Second, "Max response time should be under 2s")
	}
}

// Test 4: Database Connection Pool
func testDatabaseConnectionPool(t *testing.T) {
	t.Log("Testing: Database Connection Pool Handling")

	testCtx := setupTestUser(t)
	productID := getTestProductID(t, testCtx.Token)

	// Simulate many concurrent database operations
	concurrentOps := 200
	var wg sync.WaitGroup
	successCount := 0
	var mu sync.Mutex

	start := time.Now()

	for i := 0; i < concurrentOps; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			headers := map[string]string{
				"Authorization": "Bearer " + testCtx.Token,
			}

			// Product detail retrieval (database operation)
			url := fmt.Sprintf("%s/api/v1/products/%s", shopServiceURL, productID)
			resp, err := makeJSONRequest(http.MethodGet, url, headers, nil)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	t.Logf("✅ Database operations: %d/%d succeeded in %s", successCount, concurrentOps, elapsed)
	t.Logf("   Throughput: %.2f ops/second", float64(successCount)/elapsed.Seconds())

	assert.GreaterOrEqual(t, successCount, concurrentOps*95/100, "At least 95% should succeed")
	assert.Less(t, elapsed, 30*time.Second, "Should complete within 30 seconds")
}

// Benchmark: End-to-End Purchase Flow
func BenchmarkEndToEndPurchaseFlow(b *testing.B) {
	ctx := context.Background()
	if err := waitForServices(ctx); err != nil {
		b.Fatalf("Services not ready: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		testCtx := &E2ETestContext{}

		// User Registration
		payload := map[string]interface{}{
			"email":    fmt.Sprintf("bench.%d.%d@example.com", time.Now().Unix(), i),
			"password": "SecurePassword123!",
			"role":     "CUSTOMER",
		}

		resp, err := makeJSONRequest(http.MethodPost, authServiceURL+"/api/v1/register", nil, payload)
		if err != nil {
			b.Fatal(err)
		}

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		testCtx.UserID = result["user_id"].(string)
		testCtx.Token = result["token"].(string)

		// Customer Creation
		customerPayload := map[string]interface{}{
			"user_id": testCtx.UserID,
			"name":    "Bench User",
			"phone":   "090-0000-0000",
		}

		headers := map[string]string{
			"Authorization": "Bearer " + testCtx.Token,
		}

		resp, _ = makeJSONRequest(http.MethodPost, customerServiceURL+"/api/v1/customers", headers, customerPayload)
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		testCtx.CustomerID = result["customer_id"].(string)

		// Search Product
		resp, _ = makeJSONRequest(http.MethodGet, searchServiceURL+"/api/v1/search?q=test", headers, nil)
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		products := result["products"].([]interface{})
		productID := products[0].(map[string]interface{})["id"].(string)

		// Create Order
		orderPayload := map[string]interface{}{
			"customer_id": testCtx.CustomerID,
			"items": []map[string]interface{}{
				{
					"product_id": productID,
					"quantity":   1,
				},
			},
		}

		resp, _ = makeJSONRequest(http.MethodPost, orderServiceURL+"/api/v1/orders", headers, orderPayload)
		resp.Body.Close()
	}
}
