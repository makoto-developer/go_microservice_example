package order_flow

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// OrderClient wraps the gRPC client for Order Service
type OrderClient struct {
	conn    *grpc.ClientConn
	address string
}

// NewOrderClient creates a new order client
func NewOrderClient(address string) (*OrderClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}

	return &OrderClient{
		conn:    conn,
		address: address,
	}, nil
}

// Close closes the gRPC connection
func (c *OrderClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Order represents an order in the system
type Order struct {
	ID          string
	CustomerID  string
	Status      string
	TotalAmount float64
	Items       []OrderItem
	CreatedAt   time.Time
}

// CreateOrder creates a new order
func (c *OrderClient) CreateOrder(ctx context.Context, req *OrderRequest) (*Order, error) {
	// In a real implementation, this would call the gRPC method
	// For now, we'll simulate order creation
	
	time.Sleep(50 * time.Millisecond)
	
	order := &Order{
		ID:          fmt.Sprintf("ord_%s", generateID()),
		CustomerID:  req.CustomerID,
		Status:      "pending",
		TotalAmount: req.TotalAmount,
		Items:       req.Items,
		CreatedAt:   time.Now(),
	}
	
	return order, nil
}

// GetOrder retrieves an order by ID
func (c *OrderClient) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	// In a real implementation, this would call the gRPC method
	
	time.Sleep(50 * time.Millisecond)
	
	order := &Order{
		ID:          orderID,
		CustomerID:  "customer_123",
		Status:      "pending",
		TotalAmount: 5000.0,
		Items: []OrderItem{
			{
				ProductID: "product_123",
				Quantity:  5,
				Price:     1000.0,
			},
		},
		CreatedAt: time.Now().Add(-1 * time.Minute),
	}
	
	return order, nil
}

// UpdateOrderStatus updates the status of an order
func (c *OrderClient) UpdateOrderStatus(ctx context.Context, orderID, status string) error {
	// In a real implementation, this would call the gRPC method
	
	time.Sleep(50 * time.Millisecond)
	
	return nil
}

// CancelOrder cancels an order
func (c *OrderClient) CancelOrder(ctx context.Context, orderID, reason string) error {
	// In a real implementation, this would call the gRPC method
	
	time.Sleep(50 * time.Millisecond)
	
	return nil
}

// ListOrders lists orders for a customer
func (c *OrderClient) ListOrders(ctx context.Context, customerID string) ([]*Order, error) {
	// In a real implementation, this would call the gRPC method
	
	time.Sleep(50 * time.Millisecond)
	
	orders := []*Order{
		{
			ID:          "ord_1",
			CustomerID:  customerID,
			Status:      "confirmed",
			TotalAmount: 5000.0,
			CreatedAt:   time.Now().Add(-1 * time.Hour),
		},
	}
	
	return orders, nil
}
