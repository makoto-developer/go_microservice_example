// Package client は他マイクロサービスへの gRPC クライアントを提供する。
package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	paymentpb "github.com/makoto-developer/go_microservice_example/microservices/payment/proto"
)

// PaymentClient は決済サービスへの呼び出しを抽象化する。テストではフェイクに差し替える。
type PaymentClient interface {
	// ConfirmCODByOrder は注文に紐づく「支払い待ちの代引き決済」を集金確定にする。
	// 代引きでない注文・決済が無い注文では何もしない(エラーにしない)。
	ConfirmCODByOrder(ctx context.Context, orderID string) error
	Close() error
}

type grpcPaymentClient struct {
	conn *grpc.ClientConn
	stub paymentpb.PaymentServiceClient
}

func NewPaymentClient(addr string) (PaymentClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create payment client for %s: %w", addr, err)
	}
	return &grpcPaymentClient{conn: conn, stub: paymentpb.NewPaymentServiceClient(conn)}, nil
}

func (c *grpcPaymentClient) ConfirmCODByOrder(ctx context.Context, orderID string) error {
	list, err := c.stub.ListPayments(ctx, &paymentpb.ListPaymentsRequest{OrderId: orderID})
	if err != nil {
		return fmt.Errorf("list payments for order %s: %w", orderID, err)
	}

	for _, p := range list.GetPayments() {
		if p.GetPaymentMethod() != paymentpb.PaymentMethodType_CASH_ON_DELIVERY {
			continue
		}
		if p.GetStatus() != paymentpb.PaymentStatus_PAYMENT_STATUS_PENDING {
			continue
		}
		if _, err := c.stub.ConfirmCODPayment(ctx, &paymentpb.ConfirmCODPaymentRequest{
			PaymentId: p.GetId(),
			OrderId:   orderID,
		}); err != nil {
			return fmt.Errorf("confirm COD payment %s: %w", p.GetId(), err)
		}
	}
	return nil
}

func (c *grpcPaymentClient) Close() error {
	return c.conn.Close()
}
