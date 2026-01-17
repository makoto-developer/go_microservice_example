package notification

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// EmailClient はSendGrid APIのモック実装
type EmailClient struct {
	apiKey string
}

// NewEmailClient はEmailClientを初期化
func NewEmailClient(apiKey string) *EmailClient {
	return &EmailClient{
		apiKey: apiKey,
	}
}

// EmailRequest はメール送信リクエスト
type EmailRequest struct {
	To          []string
	From        string
	Subject     string
	TextContent string
	HTMLContent string
	TemplateID  string
	TemplateData map[string]interface{}
}

// EmailResponse はメール送信レスポンス
type EmailResponse struct {
	MessageID string
	Status    string
	SentAt    time.Time
}

// SendEmail はメールを送信（モック）
func (c *EmailClient) SendEmail(ctx context.Context, req EmailRequest) (*EmailResponse, error) {
	// モック実装: 実際のSendGrid APIは呼ばず、ダミーレスポンスを返す
	messageID := fmt.Sprintf("msg_mock_%s", uuid.New().String()[:8])

	// ログ出力（実際はメール送信）
	fmt.Printf("[EMAIL MOCK] Sending email to: %v\n", req.To)
	fmt.Printf("[EMAIL MOCK] Subject: %s\n", req.Subject)
	fmt.Printf("[EMAIL MOCK] Template: %s\n", req.TemplateID)

	resp := &EmailResponse{
		MessageID: messageID,
		Status:    "sent",
		SentAt:    time.Now(),
	}

	return resp, nil
}

// SendBulkEmail は一括メール送信（モック）
func (c *EmailClient) SendBulkEmail(ctx context.Context, requests []EmailRequest) ([]*EmailResponse, error) {
	// モック実装: 各メールを順次送信
	responses := make([]*EmailResponse, 0, len(requests))

	for _, req := range requests {
		resp, err := c.SendEmail(ctx, req)
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
	}

	return responses, nil
}

// TemplateRenderer はテンプレートレンダリング
type TemplateRenderer struct{}

// NewTemplateRenderer はTemplateRendererを初期化
func NewTemplateRenderer() *TemplateRenderer {
	return &TemplateRenderer{}
}

// RenderTemplate はテンプレートをレンダリング
func (r *TemplateRenderer) RenderTemplate(templateID string, data map[string]interface{}) (string, string, error) {
	// モック実装: 簡易的なテンプレート変換
	// 実際はHTMLテンプレートエンジンを使用

	var textContent, htmlContent string

	switch templateID {
	case "user_registration":
		name := data["name"].(string)
		confirmURL := data["confirm_url"].(string)
		textContent = fmt.Sprintf("こんにちは %s さん\n\n以下のURLからメールアドレスを確認してください：\n%s", name, confirmURL)
		htmlContent = fmt.Sprintf("<p>こんにちは %s さん</p><p><a href=\"%s\">メールアドレスを確認する</a></p>", name, confirmURL)

	case "order_confirmed":
		orderNumber := data["order_number"].(string)
		totalAmount := data["total_amount"].(int64)
		textContent = fmt.Sprintf("ご注文ありがとうございます。\n\n注文番号: %s\n合計金額: ¥%d", orderNumber, totalAmount)
		htmlContent = fmt.Sprintf("<p>ご注文ありがとうございます。</p><p>注文番号: %s<br>合計金額: ¥%d</p>", orderNumber, totalAmount)

	case "payment_completed":
		orderNumber := data["order_number"].(string)
		textContent = fmt.Sprintf("お支払いが完了しました。\n\n注文番号: %s", orderNumber)
		htmlContent = fmt.Sprintf("<p>お支払いが完了しました。</p><p>注文番号: %s</p>", orderNumber)

	default:
		textContent = "通知メールです"
		htmlContent = "<p>通知メールです</p>"
	}

	return textContent, htmlContent, nil
}
