package order_flow

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// PaymentClient wraps the gRPC client for Payment Service
type PaymentClient struct {
	conn    *grpc.ClientConn
	address string
}

// NewPaymentClient creates a new payment client
func NewPaymentClient(address string) (*PaymentClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to payment service: %w", err)
	}

	return &PaymentClient{
		conn:    conn,
		address: address,
	}, nil
}

// Close closes the gRPC connection
func (c *PaymentClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Payment represents a payment record
type Payment struct {
	ID              string
	OrderID         string
	Amount          float64
	Status          string
	Method          string
	TransactionID   string
	CreatedAt       time.Time
	CompletedAt     *time.Time
	FailureReason   string
}

// ProcessPayment processes a payment
func (c *PaymentClient) ProcessPayment(ctx context.Context, req *PaymentRequest) (*Payment, error) {
	// In a real implementation, this would call the gRPC method
	// For now, we'll simulate payment processing
	
	time.Sleep(100 * time.Millisecond) // Simulate payment gateway delay
	
	now := time.Now()
	payment := &Payment{
		ID:            fmt.Sprintf("pay_%s", generateID()),
		OrderID:       req.OrderID,
		Amount:        req.Amount,
		Status:        "completed",
		Method:        req.Method,
		TransactionID: fmt.Sprintf("txn_%s", generateID()),
		CreatedAt:     now,
		CompletedAt:   &now,
	}
	
	return payment, nil
}

// GetPayment retrieves a payment by ID
func (c *PaymentClient) GetPayment(ctx context.Context, paymentID string) (*Payment, error) {
	// In a real implementation, this would call the gRPC method
	
	time.Sleep(50 * time.Millisecond)
	
	now := time.Now()
	payment := &Payment{
		ID:            paymentID,
		OrderID:       "ord_123",
		Amount:        5000.0,
		Status:        "completed",
		Method:        "credit_card",
		TransactionID: "txn_123",
		CreatedAt:     now.Add(-1 * time.Minute),
		CompletedAt:   &now,
	}
	
	return payment, nil
}

// RefundPayment refunds a payment
func (c *PaymentClient) RefundPayment(ctx context.Context, paymentID, reason string) (*Payment, error) {
	// In a real implementation, this would call the gRPC method
	
	time.Sleep(100 * time.Millisecond)
	
	now := time.Now()
	payment := &Payment{
		ID:            paymentID,
		OrderID:       "ord_123",
		Amount:        5000.0,
		Status:        "refunded",
		Method:        "credit_card",
		TransactionID: "txn_123",
		CreatedAt:     now.Add(-1 * time.Hour),
		CompletedAt:   &now,
	}
	
	return payment, nil
}

// GetPaymentByOrderID retrieves a payment by order ID
func (c *PaymentClient) GetPaymentByOrderID(ctx context.Context, orderID string) (*Payment, error) {
	// In a real implementation, this would call the gRPC method
	
	time.Sleep(50 * time.Millisecond)
	
	now := time.Now()
	payment := &Payment{
		ID:            "pay_123",
		OrderID:       orderID,
		Amount:        5000.0,
		Status:        "completed",
		Method:        "credit_card",
		TransactionID: "txn_123",
		CreatedAt:     now.Add(-1 * time.Minute),
		CompletedAt:   &now,
	}
	
	return payment, nil
}

// CancelPayment cancels a pending payment
func (c *PaymentClient) CancelPayment(ctx context.Context, paymentID string) error {
	// In a real implementation, this would call the gRPC method
	
	time.Sleep(50 * time.Millisecond)
	
	return nil
}
