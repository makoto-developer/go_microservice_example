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
	// ProcessPayment は決済インテントの作成と確定までを行い、決済 ID を返す。
	ProcessPayment(ctx context.Context, in ProcessPaymentInput) (*PaymentResult, error)
	Close() error
}

type ProcessPaymentInput struct {
	OrderID         string
	CustomerID      string
	PaymentMethodID string
	Amount          int64 // 送料込みの合計金額(円)
	Currency        string
}

type PaymentResult struct {
	PaymentID string
}

type grpcPaymentClient struct {
	conn *grpc.ClientConn
	stub paymentpb.PaymentServiceClient
}

// NewPaymentClient は決済サービスへの gRPC クライアントを作る。
// 接続は遅延確立されるため、この時点で決済サービスが起動している必要はない。
func NewPaymentClient(addr string) (PaymentClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create payment client for %s: %w", addr, err)
	}
	return &grpcPaymentClient{conn: conn, stub: paymentpb.NewPaymentServiceClient(conn)}, nil
}

func (c *grpcPaymentClient) ProcessPayment(ctx context.Context, in ProcessPaymentInput) (*PaymentResult, error) {
	intent, err := c.stub.CreatePaymentIntent(ctx, &paymentpb.CreatePaymentIntentRequest{
		OrderId:         in.OrderID,
		Amount:          fmt.Sprintf("%d", in.Amount),
		Currency:        in.Currency,
		PaymentMethodId: in.PaymentMethodID,
		CustomerId:      in.CustomerID,
	})
	if err != nil {
		return nil, fmt.Errorf("create payment intent: %w", err)
	}

	confirm, err := c.stub.ConfirmPayment(ctx, &paymentpb.ConfirmPaymentRequest{
		PaymentId: intent.GetPaymentId(),
	})
	if err != nil {
		return nil, fmt.Errorf("confirm payment: %w", err)
	}
	if !confirm.GetSuccess() {
		return nil, fmt.Errorf("payment was not confirmed: %s", confirm.GetMessage())
	}

	return &PaymentResult{PaymentID: intent.GetPaymentId()}, nil
}

func (c *grpcPaymentClient) Close() error {
	return c.conn.Close()
}
