package auth

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TestClient wraps gRPC connections for auth service
type TestClient struct {
	conn         *grpc.ClientConn
	customerConn *grpc.ClientConn
	ownerConn    *grpc.ClientConn
}

// NewTestClient creates a new test client
func NewTestClient(authServiceAddr string) (*TestClient, error) {
	// Set up connection options
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to auth service
	conn, err := grpc.DialContext(ctx, authServiceAddr, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service at %s: %w", authServiceAddr, err)
	}

	return &TestClient{
		conn:         conn,
		customerConn: conn,
		ownerConn:    conn,
	}, nil
}

// Close closes all connections
func (c *TestClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetCustomerConn returns the customer auth connection
func (c *TestClient) GetCustomerConn() *grpc.ClientConn {
	return c.customerConn
}

// GetOwnerConn returns the owner auth connection
func (c *TestClient) GetOwnerConn() *grpc.ClientConn {
	return c.ownerConn
}
