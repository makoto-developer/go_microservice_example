package notification

import (
	"context"
	"testing"
)

func TestEmailClient_SendEmail(t *testing.T) {
	client := NewEmailClient("test_api_key")
	ctx := context.Background()

	req := EmailRequest{
		To:          []string{"test@example.com"},
		From:        "noreply@example.com",
		Subject:     "Test Email",
		TextContent: "This is a test email",
		HTMLContent: "<p>This is a test email</p>",
	}

	resp, err := client.SendEmail(ctx, req)
	if err != nil {
		t.Fatalf("SendEmail failed: %v", err)
	}

	if resp.MessageID == "" {
		t.Error("Expected non-empty message ID")
	}

	if resp.Status != "sent" {
		t.Errorf("Expected status sent, got %s", resp.Status)
	}
}

func TestEmailClient_SendBulkEmail(t *testing.T) {
	client := NewEmailClient("test_api_key")
	ctx := context.Background()

	requests := []EmailRequest{
		{
			To:      []string{"user1@example.com"},
			From:    "noreply@example.com",
			Subject: "Test 1",
		},
		{
			To:      []string{"user2@example.com"},
			From:    "noreply@example.com",
			Subject: "Test 2",
		},
	}

	responses, err := client.SendBulkEmail(ctx, requests)
	if err != nil {
		t.Fatalf("SendBulkEmail failed: %v", err)
	}

	if len(responses) != len(requests) {
		t.Errorf("Expected %d responses, got %d", len(requests), len(responses))
	}
}

func TestTemplateRenderer_RenderTemplate(t *testing.T) {
	renderer := NewTemplateRenderer()

	tests := []struct {
		name       string
		templateID string
		data       map[string]interface{}
	}{
		{
			"User registration",
			"user_registration",
			map[string]interface{}{
				"name":        "Test User",
				"confirm_url": "https://example.com/confirm",
			},
		},
		{
			"Order confirmed",
			"order_confirmed",
			map[string]interface{}{
				"order_number": "ORD-12345",
				"total_amount": int64(10000),
			},
		},
		{
			"Payment completed",
			"payment_completed",
			map[string]interface{}{
				"order_number": "ORD-12345",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			text, html, err := renderer.RenderTemplate(tt.templateID, tt.data)
			if err != nil {
				t.Fatalf("RenderTemplate failed: %v", err)
			}

			if text == "" {
				t.Error("Expected non-empty text content")
			}

			if html == "" {
				t.Error("Expected non-empty HTML content")
			}
		})
	}
}
