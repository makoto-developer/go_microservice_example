package order_flow

import (
	"time"

	"github.com/google/uuid"
)

// TestData contains test fixtures and helper data
type TestData struct {
	CustomerID  string
	ProductID   string
	ShopID      string
	OrderID     string
	PaymentID   string
	InitialStock int32
}

// NewTestData creates a new test data instance with random IDs
func NewTestData() *TestData {
	return &TestData{
		CustomerID:   uuid.New().String(),
		ProductID:    uuid.New().String(),
		ShopID:       uuid.New().String(),
		InitialStock: 100,
	}
}

// OrderRequest represents a test order request
type OrderRequest struct {
	CustomerID string
	Items      []OrderItem
	TotalAmount float64
}

// OrderItem represents an item in an order
type OrderItem struct {
	ProductID string
	Quantity  int32
	Price     float64
}

// PaymentRequest represents a payment request
type PaymentRequest struct {
	OrderID     string
	Amount      float64
	Method      string
	CardNumber  string
	ExpiryMonth int32
	ExpiryYear  int32
	CVV         string
}

// InventoryReservation represents an inventory reservation
type InventoryReservation struct {
	ProductID    string
	Quantity     int32
	ReservationID string
	ExpiresAt    time.Time
}

// NewSampleOrderRequest creates a sample order request for testing
func NewSampleOrderRequest(td *TestData) *OrderRequest {
	return &OrderRequest{
		CustomerID: td.CustomerID,
		Items: []OrderItem{
			{
				ProductID: td.ProductID,
				Quantity:  5,
				Price:     1000.0,
			},
		},
		TotalAmount: 5000.0,
	}
}

// NewSamplePaymentRequest creates a sample payment request for testing
func NewSamplePaymentRequest(orderID string, amount float64) *PaymentRequest {
	return &PaymentRequest{
		OrderID:     orderID,
		Amount:      amount,
		Method:      "credit_card",
		CardNumber:  "4242424242424242", // Stripe test card
		ExpiryMonth: 12,
		ExpiryYear:  2025,
		CVV:         "123",
	}
}
