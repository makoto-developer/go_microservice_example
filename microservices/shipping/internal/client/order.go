package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderpb "github.com/makoto-developer/go_microservice_example/microservices/order/proto"
)

// OrderClient は注文サービスへの呼び出しを抽象化する。テストではフェイクに差し替える。
type OrderClient interface {
	// MarkOrderDelivered は配達完了を注文サービスへ反映する。
	MarkOrderDelivered(ctx context.Context, orderID, trackingNumber, carrier string) error
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

func (c *grpcOrderClient) MarkOrderDelivered(ctx context.Context, orderID, trackingNumber, carrier string) error {
	_, err := c.stub.UpdateOrderStatus(ctx, &orderpb.UpdateOrderStatusRequest{
		OrderId:        orderID,
		NewStatus:      orderpb.OrderStatus_DELIVERED,
		ChangedBy:      "shipping-service",
		ChangeReason:   "carrier reported delivery",
		TrackingNumber: trackingNumber,
		Carrier:        carrier,
	})
	if err != nil {
		return fmt.Errorf("update order status: %w", err)
	}
	return nil
}

func (c *grpcOrderClient) Close() error {
	return c.conn.Close()
}
