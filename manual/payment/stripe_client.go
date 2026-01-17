package payment

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// StripeClient はStripe APIのモック実装
type StripeClient struct {
	apiKey string
}

// NewStripeClient はStripeClientを初期化
func NewStripeClient(apiKey string) *StripeClient {
	return &StripeClient{
		apiKey: apiKey,
	}
}

// PaymentIntentRequest はStripe Payment Intent作成リクエスト
type PaymentIntentRequest struct {
	Amount      int64
	Currency    string
	Description string
	Metadata    map[string]string
}

// PaymentIntentResponse はStripe Payment Intentレスポンス
type PaymentIntentResponse struct {
	ID           string
	Amount       int64
	Currency     string
	Status       string
	ClientSecret string
	CreatedAt    time.Time
}

// CreatePaymentIntent はPayment Intentを作成（モック）
func (c *StripeClient) CreatePaymentIntent(ctx context.Context, req PaymentIntentRequest) (*PaymentIntentResponse, error) {
	// モック実装: 実際のStripe APIは呼ばず、ダミーレスポンスを返す
	intentID := fmt.Sprintf("pi_mock_%s", uuid.New().String()[:8])

	resp := &PaymentIntentResponse{
		ID:           intentID,
		Amount:       req.Amount,
		Currency:     req.Currency,
		Status:       "requires_payment_method",
		ClientSecret: fmt.Sprintf("%s_secret_mock", intentID),
		CreatedAt:    time.Now(),
	}

	return resp, nil
}

// ConfirmPaymentIntent はPayment Intentを確定（モック）
func (c *StripeClient) ConfirmPaymentIntent(ctx context.Context, intentID string) (*PaymentIntentResponse, error) {
	// モック実装: 成功ステータスを返す
	resp := &PaymentIntentResponse{
		ID:       intentID,
		Status:   "succeeded",
		Amount:   10000, // ダミー金額
		Currency: "jpy",
	}

	return resp, nil
}

// CancelPaymentIntent はPayment Intentをキャンセル（モック）
func (c *StripeClient) CancelPaymentIntent(ctx context.Context, intentID string) (*PaymentIntentResponse, error) {
	// モック実装: キャンセル成功を返す
	resp := &PaymentIntentResponse{
		ID:       intentID,
		Status:   "canceled",
		Amount:   10000,
		Currency: "jpy",
	}

	return resp, nil
}

// RefundRequest はStripe返金リクエスト
type RefundRequest struct {
	PaymentIntentID string
	Amount          int64
	Reason          string
}

// RefundResponse はStripe返金レスポンス
type RefundResponse struct {
	ID              string
	PaymentIntentID string
	Amount          int64
	Status          string
	Reason          string
	CreatedAt       time.Time
}

// CreateRefund は返金を作成（モック）
func (c *StripeClient) CreateRefund(ctx context.Context, req RefundRequest) (*RefundResponse, error) {
	// モック実装: 返金成功を返す
	refundID := fmt.Sprintf("re_mock_%s", uuid.New().String()[:8])

	resp := &RefundResponse{
		ID:              refundID,
		PaymentIntentID: req.PaymentIntentID,
		Amount:          req.Amount,
		Status:          "succeeded",
		Reason:          req.Reason,
		CreatedAt:       time.Now(),
	}

	return resp, nil
}

// WebhookEvent はStripe Webhookイベント
type WebhookEvent struct {
	ID      string
	Type    string
	Data    map[string]interface{}
	Created time.Time
}

// VerifyWebhookSignature はWebhookの署名検証（モック）
func (c *StripeClient) VerifyWebhookSignature(payload []byte, signature string) (*WebhookEvent, error) {
	// モック実装: 常に成功とする
	event := &WebhookEvent{
		ID:      fmt.Sprintf("evt_mock_%s", uuid.New().String()[:8]),
		Type:    "payment_intent.succeeded",
		Data:    map[string]interface{}{},
		Created: time.Now(),
	}

	return event, nil
}
