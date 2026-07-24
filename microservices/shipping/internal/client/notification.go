package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	notificationpb "github.com/makoto-developer/go_microservice_example/microservices/notification/proto"
)

// NotificationClient は通知サービスへの呼び出しを抽象化する。テストではフェイクに差し替える。
type NotificationClient interface {
	// NotifyDelivered は配達完了の通知を送る。
	NotifyDelivered(ctx context.Context, customerID, orderID, trackingNumber string) error
	Close() error
}

type grpcNotificationClient struct {
	conn *grpc.ClientConn
	stub notificationpb.NotificationServiceClient
}

func NewNotificationClient(addr string) (NotificationClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create notification client for %s: %w", addr, err)
	}
	return &grpcNotificationClient{conn: conn, stub: notificationpb.NewNotificationServiceClient(conn)}, nil
}

func (c *grpcNotificationClient) NotifyDelivered(ctx context.Context, customerID, orderID, trackingNumber string) error {
	_, err := c.stub.SendEmail(ctx, &notificationpb.SendEmailRequest{
		UserId:      customerID,
		TemplateKey: "order_delivered",
		Variables: map[string]string{
			"order_id":        orderID,
			"tracking_number": trackingNumber,
		},
		NotificationType: notificationpb.NotificationType_ORDER_DELIVERED,
	})
	if err != nil {
		return fmt.Errorf("send order_delivered email: %w", err)
	}
	return nil
}

func (c *grpcNotificationClient) Close() error {
	return c.conn.Close()
}
