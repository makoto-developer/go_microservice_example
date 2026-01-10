package custom

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type EmailClient struct {
	SMTPHost string
	SMTPPort string
	From     string
}

type EmailMessage struct {
	To      []string
	Subject string
	Body    string
	IsHTML  bool
}

func NewEmailClient() *EmailClient {
	smtpHost := os.Getenv("MAILHOG_HOST")
	if smtpHost == "" {
		smtpHost = "mailhog_dev"
	}

	smtpPort := os.Getenv("MAILHOG_SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "1025"
	}

	from := os.Getenv("EMAIL_FROM")
	if from == "" {
		from = "noreply@go-microservice.local"
	}

	return &EmailClient{
		SMTPHost: smtpHost,
		SMTPPort: smtpPort,
		From:     from,
	}
}

func (c *EmailClient) SendEmail(msg EmailMessage) error {
	if len(msg.To) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	if msg.Subject == "" {
		return fmt.Errorf("subject is required")
	}

	if msg.Body == "" {
		return fmt.Errorf("body is required")
	}

	headers := c.buildHeaders(msg)
	body := []byte(headers + "\r\n" + msg.Body)

	addr := fmt.Sprintf("%s:%s", c.SMTPHost, c.SMTPPort)

	err := smtp.SendMail(addr, nil, c.From, msg.To, body)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (c *EmailClient) buildHeaders(msg EmailMessage) string {
	var headers strings.Builder

	headers.WriteString(fmt.Sprintf("From: %s\r\n", c.From))
	headers.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(msg.To, ", ")))
	headers.WriteString(fmt.Sprintf("Subject: %s\r\n", msg.Subject))
	headers.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	headers.WriteString("MIME-Version: 1.0\r\n")

	if msg.IsHTML {
		headers.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	} else {
		headers.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	}

	return headers.String()
}

func (c *EmailClient) SendOrderConfirmation(to string, orderID string, amount int64) error {
	subject := fmt.Sprintf("注文確認 - 注文番号: %s", orderID)
	body := fmt.Sprintf(`
ご注文ありがとうございます。

注文番号: %s
金額: ¥%d

ご不明な点がございましたら、お気軽にお問い合わせください。

よろしくお願いいたします。
`, orderID, amount)

	return c.SendEmail(EmailMessage{
		To:      []string{to},
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	})
}

func (c *EmailClient) SendPaymentSuccess(to string, orderID string, amount int64) error {
	subject := fmt.Sprintf("お支払い完了 - 注文番号: %s", orderID)
	body := fmt.Sprintf(`
お支払いが正常に完了しました。

注文番号: %s
支払金額: ¥%d

商品の発送準備を開始いたします。

よろしくお願いいたします。
`, orderID, amount)

	return c.SendEmail(EmailMessage{
		To:      []string{to},
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	})
}

func (c *EmailClient) SendShipmentNotification(to string, orderID string, trackingNumber string) error {
	subject := fmt.Sprintf("発送完了 - 注文番号: %s", orderID)
	body := fmt.Sprintf(`
ご注文の商品を発送いたしました。

注文番号: %s
追跡番号: %s

配送状況は追跡番号でご確認いただけます。

よろしくお願いいたします。
`, orderID, trackingNumber)

	return c.SendEmail(EmailMessage{
		To:      []string{to},
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	})
}

func (c *EmailClient) SendWelcomeEmail(to string, userName string) error {
	subject := "ご登録ありがとうございます"
	body := fmt.Sprintf(`
%s 様

ご登録いただきありがとうございます。

当サービスをご利用いただけます。
ご不明な点がございましたら、お気軽にお問い合わせください。

今後ともよろしくお願いいたします。
`, userName)

	return c.SendEmail(EmailMessage{
		To:      []string{to},
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	})
}

func (c *EmailClient) SendPasswordReset(to string, resetToken string) error {
	subject := "パスワードリセットのご案内"
	body := fmt.Sprintf(`
パスワードリセットのリクエストを受け付けました。

以下のトークンを使用してパスワードをリセットしてください。
リセットトークン: %s

このトークンは24時間有効です。

もしこのリクエストにお心当たりがない場合は、このメールを無視してください。

よろしくお願いいたします。
`, resetToken)

	return c.SendEmail(EmailMessage{
		To:      []string{to},
		Subject: subject,
		Body:    body,
		IsHTML:  false,
	})
}
