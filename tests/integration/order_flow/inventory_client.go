package order_flow

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// InventoryClient wraps the gRPC client for Inventory Service
type InventoryClient struct {
	conn    *grpc.ClientConn
	address string
}

// NewInventoryClient creates a new inventory client
func NewInventoryClient(address string) (*InventoryClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to inventory service: %w", err)
	}

	return &InventoryClient{
		conn:    conn,
		address: address,
	}, nil
}

// Close closes the gRPC connection
func (c *InventoryClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// CheckStock checks if stock is available
func (c *InventoryClient) CheckStock(ctx context.Context, productID string, quantity int32) (bool, error) {
	// In a real implementation, this would call the gRPC method
	// For now, we'll simulate the check
	
	// Simulate network delay
	time.Sleep(50 * time.Millisecond)
	
	// For testing purposes, assume stock is available
	return true, nil
}

// ReserveStock reserves inventory for an order
func (c *InventoryClient) ReserveStock(ctx context.Context, productID string, quantity int32) (*InventoryReservation, error) {
	// In a real implementation, this would call the gRPC method
	// For now, we'll simulate the reservation
	
	time.Sleep(50 * time.Millisecond)
	
	reservation := &InventoryReservation{
		ProductID:    productID,
		Quantity:     quantity,
		ReservationID: fmt.Sprintf("res_%s", generateID()),
		ExpiresAt:    time.Now().Add(15 * time.Minute),
	}
	
	return reservation, nil
}

// ReleaseStock releases a reserved inventory
func (c *InventoryClient) ReleaseStock(ctx context.Context, reservationID string) error {
	// In a real implementation, this would call the gRPC method
	// For now, we'll simulate the release
	
	time.Sleep(50 * time.Millisecond)
	
	return nil
}

// ConfirmStock confirms a reservation and decrements actual stock
func (c *InventoryClient) ConfirmStock(ctx context.Context, reservationID string) error {
	// In a real implementation, this would call the gRPC method
	// For now, we'll simulate the confirmation
	
	time.Sleep(50 * time.Millisecond)
	
	return nil
}

// GetStock retrieves current stock level
func (c *InventoryClient) GetStock(ctx context.Context, productID string) (int32, int32, error) {
	// Returns: available, reserved, error
	// In a real implementation, this would call the gRPC method
	
	time.Sleep(50 * time.Millisecond)
	
	// Return sample stock levels
	return 100, 0, nil
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
