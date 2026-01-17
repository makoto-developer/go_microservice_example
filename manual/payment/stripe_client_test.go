package payment

import (
	"context"
	"testing"
)

func TestStripeClient_CreatePaymentIntent(t *testing.T) {
	client := NewStripeClient("test_key")
	ctx := context.Background()

	req := PaymentIntentRequest{
		Amount:      10000,
		Currency:    "jpy",
		Description: "Test payment",
		Metadata: map[string]string{
			"order_id": "test_order_123",
		},
	}

	resp, err := client.CreatePaymentIntent(ctx, req)
	if err != nil {
		t.Fatalf("CreatePaymentIntent failed: %v", err)
	}

	if resp.Amount != req.Amount {
		t.Errorf("Expected amount %d, got %d", req.Amount, resp.Amount)
	}

	if resp.Currency != req.Currency {
		t.Errorf("Expected currency %s, got %s", req.Currency, resp.Currency)
	}

	if resp.Status != "requires_payment_method" {
		t.Errorf("Expected status requires_payment_method, got %s", resp.Status)
	}

	if resp.ID == "" {
		t.Error("Expected non-empty payment intent ID")
	}
}

func TestStripeClient_ConfirmPaymentIntent(t *testing.T) {
	client := NewStripeClient("test_key")
	ctx := context.Background()

	resp, err := client.ConfirmPaymentIntent(ctx, "pi_mock_12345")
	if err != nil {
		t.Fatalf("ConfirmPaymentIntent failed: %v", err)
	}

	if resp.Status != "succeeded" {
		t.Errorf("Expected status succeeded, got %s", resp.Status)
	}
}

func TestStripeClient_CancelPaymentIntent(t *testing.T) {
	client := NewStripeClient("test_key")
	ctx := context.Background()

	resp, err := client.CancelPaymentIntent(ctx, "pi_mock_12345")
	if err != nil {
		t.Fatalf("CancelPaymentIntent failed: %v", err)
	}

	if resp.Status != "canceled" {
		t.Errorf("Expected status canceled, got %s", resp.Status)
	}
}

func TestStripeClient_CreateRefund(t *testing.T) {
	client := NewStripeClient("test_key")
	ctx := context.Background()

	req := RefundRequest{
		PaymentIntentID: "pi_mock_12345",
		Amount:          5000,
		Reason:          "requested_by_customer",
	}

	resp, err := client.CreateRefund(ctx, req)
	if err != nil {
		t.Fatalf("CreateRefund failed: %v", err)
	}

	if resp.Amount != req.Amount {
		t.Errorf("Expected amount %d, got %d", req.Amount, resp.Amount)
	}

	if resp.Status != "succeeded" {
		t.Errorf("Expected status succeeded, got %s", resp.Status)
	}
}

func TestStripeClient_VerifyWebhookSignature(t *testing.T) {
	client := NewStripeClient("test_key")

	payload := []byte(`{"type": "payment_intent.succeeded"}`)
	signature := "test_signature"

	event, err := client.VerifyWebhookSignature(payload, signature)
	if err != nil {
		t.Fatalf("VerifyWebhookSignature failed: %v", err)
	}

	if event.Type != "payment_intent.succeeded" {
		t.Errorf("Expected event type payment_intent.succeeded, got %s", event.Type)
	}

	if event.ID == "" {
		t.Error("Expected non-empty event ID")
	}
}
