package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	shippingpb "github.com/makoto-developer/go_microservice_example/microservices/shipping/proto"
)

// ShippingClient は配送サービスへの呼び出しを抽象化する。テストではフェイクに差し替える。
type ShippingClient interface {
	// CreateShipment は注文の出荷を起票する。
	CreateShipment(ctx context.Context, in CreateShipmentInput) (string, error)
	// CalculateFee は配送方法と重量から送料を見積もる。
	CalculateFee(ctx context.Context, shippingMethod string, totalWeightGrams int32) (int64, error)
	Close() error
}

type CreateShipmentInput struct {
	OrderID         string
	CustomerID      string
	ShippingAddress string
}

type grpcShippingClient struct {
	conn *grpc.ClientConn
	stub shippingpb.ShippingServiceClient
}

func NewShippingClient(addr string) (ShippingClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create shipping client for %s: %w", addr, err)
	}
	return &grpcShippingClient{conn: conn, stub: shippingpb.NewShippingServiceClient(conn)}, nil
}

func (c *grpcShippingClient) CreateShipment(ctx context.Context, in CreateShipmentInput) (string, error) {
	resp, err := c.stub.CreateShipment(ctx, &shippingpb.CreateShipmentRequest{
		OrderId:         in.OrderID,
		CustomerId:      in.CustomerID,
		ShippingAddress: in.ShippingAddress,
	})
	if err != nil {
		return "", fmt.Errorf("create shipment: %w", err)
	}
	return resp.GetShipmentId(), nil
}

func (c *grpcShippingClient) CalculateFee(ctx context.Context, shippingMethod string, totalWeightGrams int32) (int64, error) {
	resp, err := c.stub.CalculateShippingFee(ctx, &shippingpb.CalculateShippingFeeRequest{
		ShippingMethod:   shippingMethod,
		TotalWeightGrams: totalWeightGrams,
	})
	if err != nil {
		return 0, fmt.Errorf("calculate shipping fee: %w", err)
	}
	var fee int64
	if _, err := fmt.Sscanf(resp.GetFee(), "%d", &fee); err != nil {
		return 0, fmt.Errorf("invalid fee %q: %w", resp.GetFee(), err)
	}
	return fee, nil
}

func (c *grpcShippingClient) Close() error {
	return c.conn.Close()
}
