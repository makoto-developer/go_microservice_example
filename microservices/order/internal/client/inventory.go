package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventorypb "github.com/makoto-developer/go_microservice_example/microservices/inventory/proto"
)

// InventoryClient は在庫サービスへの呼び出しを抽象化する。テストではフェイクに差し替える。
type InventoryClient interface {
	// ReserveOrderStock は注文の全商品を一括で引き当てる(全部成功 or 全部なし)。
	ReserveOrderStock(ctx context.Context, orderID string, items []StockItem) error
	// ConfirmOrderStock は引当を確定する(決済成功後)。
	ConfirmOrderStock(ctx context.Context, orderID string) error
	// ReleaseOrderStock は引当を解放する(決済失敗・キャンセル時の補償)。
	ReleaseOrderStock(ctx context.Context, orderID string) error
	Close() error
}

type StockItem struct {
	ProductID string
	Quantity  int
}

type grpcInventoryClient struct {
	conn *grpc.ClientConn
	stub inventorypb.InventoryServiceClient
}

func NewInventoryClient(addr string) (InventoryClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create inventory client for %s: %w", addr, err)
	}
	return &grpcInventoryClient{conn: conn, stub: inventorypb.NewInventoryServiceClient(conn)}, nil
}

func (c *grpcInventoryClient) ReserveOrderStock(ctx context.Context, orderID string, items []StockItem) error {
	reservations := make([]*inventorypb.ReservationRequest, 0, len(items))
	for _, item := range items {
		reservations = append(reservations, &inventorypb.ReservationRequest{
			// inventory サービスの簡略運用に合わせ、inventory_id に product_id を渡す
			InventoryId: item.ProductID,
			Quantity:    int32(item.Quantity),
		})
	}
	_, err := c.stub.BulkReserveStock(ctx, &inventorypb.BulkReserveStockRequest{
		OrderId:      orderID,
		Reservations: reservations,
	})
	if err != nil {
		return fmt.Errorf("reserve stock: %w", err)
	}
	return nil
}

func (c *grpcInventoryClient) ConfirmOrderStock(ctx context.Context, orderID string) error {
	_, err := c.stub.ConfirmStock(ctx, &inventorypb.ConfirmStockRequest{OrderId: orderID})
	if err != nil {
		return fmt.Errorf("confirm stock: %w", err)
	}
	return nil
}

func (c *grpcInventoryClient) ReleaseOrderStock(ctx context.Context, orderID string) error {
	_, err := c.stub.ReleaseStock(ctx, &inventorypb.ReleaseStockRequest{OrderId: orderID})
	if err != nil {
		return fmt.Errorf("release stock: %w", err)
	}
	return nil
}

func (c *grpcInventoryClient) Close() error {
	return c.conn.Close()
}
