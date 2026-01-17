package notification

import (
	"context"
	"testing"
)

func TestPushClient_SendPush(t *testing.T) {
	client := NewPushClient("fcm_key", "/path/to/apns.p8")
	ctx := context.Background()

	tests := []struct {
		name     string
		platform Platform
	}{
		{"Android", PlatformAndroid},
		{"iOS", PlatformIOS},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := PushRequest{
				DeviceTokens: []string{"token1", "token2"},
				Platform:     tt.platform,
				Title:        "Test Notification",
				Body:         "This is a test",
				Data: map[string]string{
					"order_id": "12345",
				},
			}

			resp, err := client.SendPush(ctx, req)
			if err != nil {
				t.Fatalf("SendPush failed: %v", err)
			}

			if resp.SuccessCount != len(req.DeviceTokens) {
				t.Errorf("Expected success count %d, got %d", len(req.DeviceTokens), resp.SuccessCount)
			}

			if resp.FailureCount != 0 {
				t.Errorf("Expected failure count 0, got %d", resp.FailureCount)
			}
		})
	}
}

func TestPushClient_SendTopicPush(t *testing.T) {
	client := NewPushClient("fcm_key", "/path/to/apns.p8")
	ctx := context.Background()

	req := PushRequest{
		Title: "Topic Notification",
		Body:  "This is a topic notification",
	}

	resp, err := client.SendTopicPush(ctx, "news", req)
	if err != nil {
		t.Fatalf("SendTopicPush failed: %v", err)
	}

	if resp.MessageID == "" {
		t.Error("Expected non-empty message ID")
	}
}

func TestPushClient_ValidateDeviceToken(t *testing.T) {
	client := NewPushClient("fcm_key", "/path/to/apns.p8")
	ctx := context.Background()

	valid, err := client.ValidateDeviceToken(ctx, "test_token", PlatformAndroid)
	if err != nil {
		t.Fatalf("ValidateDeviceToken failed: %v", err)
	}

	if !valid {
		t.Error("Expected token to be valid")
	}
}

func TestPushClient_SubscribeToTopic(t *testing.T) {
	client := NewPushClient("fcm_key", "/path/to/apns.p8")
	ctx := context.Background()

	tokens := []string{"token1", "token2"}
	err := client.SubscribeToTopic(ctx, tokens, "news")
	if err != nil {
		t.Fatalf("SubscribeToTopic failed: %v", err)
	}
}

func TestPushClient_UnsubscribeFromTopic(t *testing.T) {
	client := NewPushClient("fcm_key", "/path/to/apns.p8")
	ctx := context.Background()

	tokens := []string{"token1", "token2"}
	err := client.UnsubscribeFromTopic(ctx, tokens, "news")
	if err != nil {
		t.Fatalf("UnsubscribeFromTopic failed: %v", err)
	}
}
