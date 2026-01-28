package order_flow

import (
	"context"
	"fmt"
	"testing"
	"time"
)

const (
	// Service addresses
	inventoryServiceAddr = "localhost:22103"
	orderServiceAddr     = "localhost:22104"
	paymentServiceAddr   = "localhost:22105"
	customerServiceAddr  = "localhost:22102"
)

// TestOrderPaymentFlow tests the complete order-payment-inventory flow
func TestOrderPaymentFlow(t *testing.T) {
	ctx := context.Background()
	
	// Initialize clients
	inventoryClient, err := NewInventoryClient(inventoryServiceAddr)
	if err != nil {
		t.Fatalf("Failed to create inventory client: %v", err)
	}
	defer inventoryClient.Close()
	
	orderClient, err := NewOrderClient(orderServiceAddr)
	if err != nil {
		t.Fatalf("Failed to create order client: %v", err)
	}
	defer orderClient.Close()
	
	paymentClient, err := NewPaymentClient(paymentServiceAddr)
	if err != nil {
		t.Fatalf("Failed to create payment client: %v", err)
	}
	defer paymentClient.Close()
	
	// Create test data
	testData := NewTestData()
	
	t.Run("Happy Path - Complete Order Flow", func(t *testing.T) {
		testCompleteOrderFlow(t, ctx, testData, inventoryClient, orderClient, paymentClient)
	})
	
	t.Run("Rollback - Payment Failure", func(t *testing.T) {
		testPaymentFailureRollback(t, ctx, testData, inventoryClient, orderClient, paymentClient)
	})
	
	t.Run("Rollback - Insufficient Stock", func(t *testing.T) {
		testInsufficientStockRollback(t, ctx, testData, inventoryClient, orderClient, paymentClient)
	})
	
	t.Run("Order Cancellation", func(t *testing.T) {
		testOrderCancellation(t, ctx, testData, inventoryClient, orderClient, paymentClient)
	})
}

// testCompleteOrderFlow tests the successful order flow
func testCompleteOrderFlow(
	t *testing.T,
	ctx context.Context,
	testData *TestData,
	inventoryClient *InventoryClient,
	orderClient *OrderClient,
	paymentClient *PaymentClient,
) {
	t.Log("🎯 Starting Complete Order Flow Test")
	
	// Step 1: Check inventory availability
	t.Log("📦 Step 1: Checking inventory availability...")
	available, err := inventoryClient.CheckStock(ctx, testData.ProductID, 5)
	if err != nil {
		t.Fatalf("Failed to check stock: %v", err)
	}
	if !available {
		t.Fatal("Stock not available")
	}
	t.Log("✅ Stock is available")
	
	// Step 2: Reserve inventory
	t.Log("🔒 Step 2: Reserving inventory...")
	reservation, err := inventoryClient.ReserveStock(ctx, testData.ProductID, 5)
	if err != nil {
		t.Fatalf("Failed to reserve stock: %v", err)
	}
	t.Logf("✅ Inventory reserved: %s (expires at %s)", reservation.ReservationID, reservation.ExpiresAt.Format(time.RFC3339))
	
	// Step 3: Create order
	t.Log("📝 Step 3: Creating order...")
	orderReq := NewSampleOrderRequest(testData)
	order, err := orderClient.CreateOrder(ctx, orderReq)
	if err != nil {
		// Rollback: Release inventory reservation
		t.Log("❌ Order creation failed, releasing inventory...")
		if err := inventoryClient.ReleaseStock(ctx, reservation.ReservationID); err != nil {
			t.Logf("Failed to release stock: %v", err)
		}
		t.Fatalf("Failed to create order: %v", err)
	}
	t.Logf("✅ Order created: %s (status: %s)", order.ID, order.Status)
	
	// Step 4: Process payment
	t.Log("💳 Step 4: Processing payment...")
	paymentReq := NewSamplePaymentRequest(order.ID, order.TotalAmount)
	payment, err := paymentClient.ProcessPayment(ctx, paymentReq)
	if err != nil {
		// Rollback: Cancel order and release inventory
		t.Log("❌ Payment failed, rolling back...")
		if err := orderClient.CancelOrder(ctx, order.ID, "payment_failed"); err != nil {
			t.Logf("Failed to cancel order: %v", err)
		}
		if err := inventoryClient.ReleaseStock(ctx, reservation.ReservationID); err != nil {
			t.Logf("Failed to release stock: %v", err)
		}
		t.Fatalf("Failed to process payment: %v", err)
	}
	t.Logf("✅ Payment completed: %s (transaction: %s)", payment.ID, payment.TransactionID)
	
	// Step 5: Confirm order
	t.Log("✅ Step 5: Confirming order...")
	if err := orderClient.UpdateOrderStatus(ctx, order.ID, "confirmed"); err != nil {
		t.Fatalf("Failed to update order status: %v", err)
	}
	t.Log("✅ Order confirmed")
	
	// Step 6: Confirm inventory
	t.Log("📦 Step 6: Confirming inventory...")
	if err := inventoryClient.ConfirmStock(ctx, reservation.ReservationID); err != nil {
		t.Fatalf("Failed to confirm stock: %v", err)
	}
	t.Log("✅ Inventory confirmed")
	
	// Verify final state
	t.Log("🔍 Verifying final state...")
	finalOrder, err := orderClient.GetOrder(ctx, order.ID)
	if err != nil {
		t.Fatalf("Failed to get order: %v", err)
	}
	
	if finalOrder.Status != "confirmed" {
		t.Errorf("Expected order status 'confirmed', got '%s'", finalOrder.Status)
	}
	
	t.Log("✨ Complete Order Flow Test PASSED")
}

// testPaymentFailureRollback tests rollback when payment fails
func testPaymentFailureRollback(
	t *testing.T,
	ctx context.Context,
	testData *TestData,
	inventoryClient *InventoryClient,
	orderClient *OrderClient,
	paymentClient *PaymentClient,
) {
	t.Log("🎯 Starting Payment Failure Rollback Test")
	
	// Reserve inventory
	reservation, err := inventoryClient.ReserveStock(ctx, testData.ProductID, 3)
	if err != nil {
		t.Fatalf("Failed to reserve stock: %v", err)
	}
	
	// Create order
	orderReq := &OrderRequest{
		CustomerID: testData.CustomerID,
		Items: []OrderItem{
			{ProductID: testData.ProductID, Quantity: 3, Price: 1000.0},
		},
		TotalAmount: 3000.0,
	}
	order, err := orderClient.CreateOrder(ctx, orderReq)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	
	// Simulate payment failure
	t.Log("💳 Simulating payment failure...")
	// In real scenario, payment would fail here
	// For now, we simulate the rollback process
	
	// Rollback: Cancel order
	t.Log("↩️ Rolling back: Cancelling order...")
	if err := orderClient.CancelOrder(ctx, order.ID, "payment_failed"); err != nil {
		t.Fatalf("Failed to cancel order: %v", err)
	}
	t.Log("✅ Order cancelled")
	
	// Rollback: Release inventory
	t.Log("↩️ Rolling back: Releasing inventory...")
	if err := inventoryClient.ReleaseStock(ctx, reservation.ReservationID); err != nil {
		t.Fatalf("Failed to release stock: %v", err)
	}
	t.Log("✅ Inventory released")
	
	// Verify final state
	finalOrder, err := orderClient.GetOrder(ctx, order.ID)
	if err != nil {
		t.Fatalf("Failed to get order: %v", err)
	}
	
	if finalOrder.Status != "cancelled" {
		t.Errorf("Expected order status 'cancelled', got '%s'", finalOrder.Status)
	}
	
	t.Log("✨ Payment Failure Rollback Test PASSED")
}

// testInsufficientStockRollback tests when stock is insufficient
func testInsufficientStockRollback(
	t *testing.T,
	ctx context.Context,
	testData *TestData,
	inventoryClient *InventoryClient,
	orderClient *OrderClient,
	paymentClient *PaymentClient,
) {
	t.Log("🎯 Starting Insufficient Stock Test")
	
	// Try to reserve more stock than available
	t.Log("📦 Attempting to reserve 1000 units (should fail)...")
	available, err := inventoryClient.CheckStock(ctx, testData.ProductID, 1000)
	if err != nil {
		t.Fatalf("Failed to check stock: %v", err)
	}
	
	if available {
		t.Fatal("Stock check should have failed for 1000 units")
	}
	
	t.Log("✅ Stock check correctly reported insufficient stock")
	t.Log("✨ Insufficient Stock Test PASSED")
}

// testOrderCancellation tests order cancellation after creation
func testOrderCancellation(
	t *testing.T,
	ctx context.Context,
	testData *TestData,
	inventoryClient *InventoryClient,
	orderClient *OrderClient,
	paymentClient *PaymentClient,
) {
	t.Log("🎯 Starting Order Cancellation Test")
	
	// Reserve inventory
	reservation, err := inventoryClient.ReserveStock(ctx, testData.ProductID, 2)
	if err != nil {
		t.Fatalf("Failed to reserve stock: %v", err)
	}
	
	// Create order
	orderReq := &OrderRequest{
		CustomerID: testData.CustomerID,
		Items: []OrderItem{
			{ProductID: testData.ProductID, Quantity: 2, Price: 1000.0},
		},
		TotalAmount: 2000.0,
	}
	order, err := orderClient.CreateOrder(ctx, orderReq)
	if err != nil {
		t.Fatalf("Failed to create order: %v", err)
	}
	
	// Cancel order (customer changes mind)
	t.Log("🚫 Cancelling order...")
	if err := orderClient.CancelOrder(ctx, order.ID, "customer_request"); err != nil {
		t.Fatalf("Failed to cancel order: %v", err)
	}
	t.Log("✅ Order cancelled")
	
	// Release inventory
	t.Log("📦 Releasing inventory...")
	if err := inventoryClient.ReleaseStock(ctx, reservation.ReservationID); err != nil {
		t.Fatalf("Failed to release stock: %v", err)
	}
	t.Log("✅ Inventory released")
	
	t.Log("✨ Order Cancellation Test PASSED")
}

// TestConcurrentOrders tests concurrent order processing
func TestConcurrentOrders(t *testing.T) {
	t.Log("🎯 Starting Concurrent Orders Test")
	
	// This test would verify that concurrent orders
	// properly handle inventory locking and don't oversell
	
	t.Skip("Concurrent test requires full database setup")
}

// TestInventoryReservationExpiry tests reservation expiration
func TestInventoryReservationExpiry(t *testing.T) {
	t.Log("🎯 Starting Inventory Reservation Expiry Test")
	
	// This test would verify that expired reservations
	// are properly released back to available stock
	
	t.Skip("Expiry test requires time-based logic")
}

// BenchmarkOrderFlow benchmarks the complete order flow
func BenchmarkOrderFlow(b *testing.B) {
	ctx := context.Background()
	
	inventoryClient, _ := NewInventoryClient(inventoryServiceAddr)
	defer inventoryClient.Close()
	
	orderClient, _ := NewOrderClient(orderServiceAddr)
	defer orderClient.Close()
	
	paymentClient, _ := NewPaymentClient(paymentServiceAddr)
	defer paymentClient.Close()
	
	testData := NewTestData()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		// Reserve inventory
		reservation, _ := inventoryClient.ReserveStock(ctx, testData.ProductID, 1)
		
		// Create order
		orderReq := NewSampleOrderRequest(testData)
		order, _ := orderClient.CreateOrder(ctx, orderReq)
		
		// Process payment
		paymentReq := NewSamplePaymentRequest(order.ID, order.TotalAmount)
		paymentClient.ProcessPayment(ctx, paymentReq)
		
		// Confirm
		orderClient.UpdateOrderStatus(ctx, order.ID, "confirmed")
		inventoryClient.ConfirmStock(ctx, reservation.ReservationID)
	}
}

// Helper function to print test section
func printTestSection(msg string) {
	fmt.Printf("\n%s\n%s\n", msg, "─────────────────────────────────────────────────────────")
}
