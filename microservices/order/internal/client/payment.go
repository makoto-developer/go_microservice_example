// Package client は他マイクロサービスへの gRPC クライアントを提供する。
package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	paymentpb "github.com/makoto-developer/go_microservice_example/microservices/payment/proto"
)

// PaymentClient は決済サービスへの呼び出しを抽象化する。テストではフェイクに差し替える。
type PaymentClient interface {
	// ProcessPayment は決済インテントの作成と確定までを行い、決済 ID を返す。
	ProcessPayment(ctx context.Context, in ProcessPaymentInput) (*PaymentResult, error)
	// CreateCODPayment は代金引換の決済(配達時集金)を作成する。決済は配達完了まで支払い待ちのまま。
	CreateCODPayment(ctx context.Context, in CODPaymentInput) (*PaymentResult, error)
	// RefundByOrder は注文に紐づく決済を全額返金する。決済が存在しない場合はエラーにしない。
	RefundByOrder(ctx context.Context, orderID string, reason string) error
	Close() error
}

type ProcessPaymentInput struct {
	OrderID         string
	CustomerID      string
	PaymentMethodID string
	Amount          int64 // 送料込みの合計金額(円)
	Currency        string
}

type CODPaymentInput struct {
	OrderID string
	Amount  int64 // 送料・代引き手数料込みの合計金額(円)
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

func (c *grpcPaymentClient) CreateCODPayment(ctx context.Context, in CODPaymentInput) (*PaymentResult, error) {
	resp, err := c.stub.CreateCODPayment(ctx, &paymentpb.CreateCODPaymentRequest{
		OrderId: in.OrderID,
		Amount:  fmt.Sprintf("%d", in.Amount),
	})
	if err != nil {
		return nil, fmt.Errorf("create COD payment: %w", err)
	}
	return &PaymentResult{PaymentID: resp.GetPaymentId()}, nil
}

func (c *grpcPaymentClient) RefundByOrder(ctx context.Context, orderID string, reason string) error {
	_, err := c.stub.CreateRefund(ctx, &paymentpb.CreateRefundRequest{
		OrderId: orderID,
		Reason:  reason,
	})
	if err != nil {
		// 決済前の注文キャンセルでは返金対象がないので、NotFound は成功扱いにする
		if status.Code(err) == codes.NotFound {
			return nil
		}
		return fmt.Errorf("refund payment: %w", err)
	}
	return nil
}

func (c *grpcPaymentClient) Close() error {
	return c.conn.Close()
}
