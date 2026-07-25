// Package client は他マイクロサービスへの gRPC クライアントを提供する。
package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderpb "github.com/makoto-developer/go_microservice_example/microservices/order/proto"
)

// OrderClient は注文サービスへの呼び出しを抽象化する(注文履歴系 RPC の委譲先)。
type OrderClient interface {
	// ListOrders は顧客の注文一覧を取得する。
	ListOrders(ctx context.Context, customerID string) ([]*orderpb.Order, error)
	// GetOrderDetail は注文の詳細を取得する。
	GetOrderDetail(ctx context.Context, orderID string) (*orderpb.Order, error)
	// CancelOrder は注文をキャンセルする(決済済みなら返金される)。
	CancelOrder(ctx context.Context, orderID, customerID, reason string) error
	// Reorder は過去の注文と同じ内容で再注文する。
	Reorder(ctx context.Context, orderID, customerID string) (string, error)
	Close() error
}

type grpcOrderClient struct {
	conn *grpc.ClientConn
	stub orderpb.OrderServiceClient
}

func NewOrderClient(addr string) (OrderClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create order client for %s: %w", addr, err)
	}
	return &grpcOrderClient{conn: conn, stub: orderpb.NewOrderServiceClient(conn)}, nil
}

func (c *grpcOrderClient) ListOrders(ctx context.Context, customerID string) ([]*orderpb.Order, error) {
	resp, err := c.stub.ListOrders(ctx, &orderpb.ListOrdersRequest{CustomerId: customerID})
	if err != nil {
		return nil, fmt.Errorf("list orders: %w", err)
	}
	return resp.GetOrders(), nil
}

func (c *grpcOrderClient) GetOrderDetail(ctx context.Context, orderID string) (*orderpb.Order, error) {
	resp, err := c.stub.GetOrderDetail(ctx, &orderpb.GetOrderDetailRequest{OrderId: orderID})
	if err != nil {
		return nil, fmt.Errorf("get order detail: %w", err)
	}
	return resp.GetOrder(), nil
}

func (c *grpcOrderClient) CancelOrder(ctx context.Context, orderID, customerID, reason string) error {
	_, err := c.stub.CancelOrder(ctx, &orderpb.CancelOrderRequest{
		OrderId:     orderID,
		CancelledBy: customerID,
		CancelNote:  reason,
	})
	if err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}
	return nil
}

func (c *grpcOrderClient) Reorder(ctx context.Context, orderID, customerID string) (string, error) {
	resp, err := c.stub.CreateReorder(ctx, &orderpb.CreateReorderRequest{
		CustomerId:      customerID,
		OriginalOrderId: orderID,
	})
	if err != nil {
		return "", fmt.Errorf("reorder: %w", err)
	}
	return resp.GetOrderId(), nil
}

func (c *grpcOrderClient) Close() error {
	return c.conn.Close()
}
