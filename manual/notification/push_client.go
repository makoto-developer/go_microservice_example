package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// PushClient はFCM/APNs APIのモック実装
type PushClient struct {
	fcmAPIKey   string
	apnsKeyPath string
}

// NewPushClient はPushClientを初期化
func NewPushClient(fcmAPIKey string, apnsKeyPath string) *PushClient {
	return &PushClient{
		fcmAPIKey:   fcmAPIKey,
		apnsKeyPath: apnsKeyPath,
	}
}

// Platform はプッシュ通知プラットフォーム
type Platform string

const (
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
)

// PushRequest はプッシュ通知リクエスト
type PushRequest struct {
	DeviceTokens []string
	Platform     Platform
	Title        string
	Body         string
	Data         map[string]string
	Badge        int
	Sound        string
}

// PushResponse はプッシュ通知レスポンス
type PushResponse struct {
	MessageID     string
	SuccessCount  int
	FailureCount  int
	FailedTokens  []string
	SentAt        time.Time
}

// SendPush はプッシュ通知を送信（モック）
func (c *PushClient) SendPush(ctx context.Context, req PushRequest) (*PushResponse, error) {
	// モック実装: 実際のFCM/APNs APIは呼ばず、ダミーレスポンスを返す
	messageID := fmt.Sprintf("push_mock_%s", uuid.New().String()[:8])

	// ログ出力（実際はプッシュ通知送信）
	fmt.Printf("[PUSH MOCK] Platform: %s\n", req.Platform)
	fmt.Printf("[PUSH MOCK] Tokens: %v\n", req.DeviceTokens)
	fmt.Printf("[PUSH MOCK] Title: %s\n", req.Title)
	fmt.Printf("[PUSH MOCK] Body: %s\n", req.Body)

	resp := &PushResponse{
		MessageID:    messageID,
		SuccessCount: len(req.DeviceTokens),
		FailureCount: 0,
		FailedTokens: []string{},
		SentAt:       time.Now(),
	}

	return resp, nil
}

// SendTopicPush はトピック配信（モック）
func (c *PushClient) SendTopicPush(ctx context.Context, topic string, req PushRequest) (*PushResponse, error) {
	// モック実装: トピック配信
	messageID := fmt.Sprintf("topic_mock_%s", uuid.New().String()[:8])

	fmt.Printf("[PUSH MOCK] Topic: %s\n", topic)
	fmt.Printf("[PUSH MOCK] Title: %s\n", req.Title)

	resp := &PushResponse{
		MessageID:    messageID,
		SuccessCount: 1,
		FailureCount: 0,
		FailedTokens: []string{},
		SentAt:       time.Now(),
	}

	return resp, nil
}

// ValidateDeviceToken はデバイストークンを検証
func (c *PushClient) ValidateDeviceToken(ctx context.Context, token string, platform Platform) (bool, error) {
	// モック実装: 常に有効とする
	// 実際はFCM/APNsに問い合わせる
	return true, nil
}

// SubscribeToTopic はトピックを購読
func (c *PushClient) SubscribeToTopic(ctx context.Context, tokens []string, topic string) error {
	// モック実装: 購読成功を返す
	fmt.Printf("[PUSH MOCK] Subscribing tokens to topic: %s\n", topic)
	return nil
}

// UnsubscribeFromTopic はトピックを購読解除
func (c *PushClient) UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) error {
	// モック実装: 購読解除成功を返す
	fmt.Printf("[PUSH MOCK] Unsubscribing tokens from topic: %s\n", topic)
	return nil
}
