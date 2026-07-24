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
	// NotifyOrderConfirmed は注文確認メールを送る。
	NotifyOrderConfirmed(ctx context.Context, in OrderNotificationInput) error
	// NotifyOrderCancelled はキャンセル確認メールを送る。
	NotifyOrderCancelled(ctx context.Context, in OrderNotificationInput) error
	Close() error
}

type OrderNotificationInput struct {
	CustomerID    string
	CustomerEmail string
	OrderNumber   string
	TotalAmount   int64
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

func (c *grpcNotificationClient) NotifyOrderConfirmed(ctx context.Context, in OrderNotificationInput) error {
	return c.sendEmail(ctx, in, "order_confirmed", notificationpb.NotificationType_ORDER_CONFIRMED)
}

func (c *grpcNotificationClient) NotifyOrderCancelled(ctx context.Context, in OrderNotificationInput) error {
	return c.sendEmail(ctx, in, "order_cancelled", notificationpb.NotificationType_ORDER_CANCELLED)
}

func (c *grpcNotificationClient) sendEmail(ctx context.Context, in OrderNotificationInput, templateKey string, nType notificationpb.NotificationType) error {
	_, err := c.stub.SendEmail(ctx, &notificationpb.SendEmailRequest{
		UserId:      in.CustomerID,
		Email:       in.CustomerEmail,
		TemplateKey: templateKey,
		Variables: map[string]string{
			"order_number": in.OrderNumber,
			"total_amount": fmt.Sprintf("%d", in.TotalAmount),
		},
		NotificationType: nType,
	})
	if err != nil {
		return fmt.Errorf("send %s email: %w", templateKey, err)
	}
	return nil
}

func (c *grpcNotificationClient) Close() error {
	return c.conn.Close()
}
