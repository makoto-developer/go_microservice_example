package auth

import (
	"context"
	"fmt"
	"net/smtp"
	"os"
)

// EmailSender はメール送信のインターフェース
type EmailSender interface {
	SendVerificationEmail(ctx context.Context, toEmail, token string) error
	SendPasswordResetEmail(ctx context.Context, toEmail, token string) error
}

// SMTPEmailSender はSMTP経由でメールを送信する実装
type SMTPEmailSender struct {
	host     string
	port     string
	from     string
}

// NewSMTPEmailSender は SMTPEmailSender のインスタンスを生成
func NewSMTPEmailSender() *SMTPEmailSender {
	host := os.Getenv("MAILHOG_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("MAILHOG_SMTP_PORT")
	if port == "" {
		port = "20004"
	}

	from := os.Getenv("EMAIL_FROM")
	if from == "" {
		from = "noreply@shopmall.example.com"
	}

	return &SMTPEmailSender{
		host: host,
		port: port,
		from: from,
	}
}

// SendVerificationEmail はメール認証用のメールを送信
func (s *SMTPEmailSender) SendVerificationEmail(ctx context.Context, toEmail, token string) error {
	subject := "メールアドレスの確認"
	body := fmt.Sprintf(`
こんにちは！

オンラインショップモールへようこそ。
以下のリンクをクリックして、メールアドレスを確認してください。

認証トークン: %s

このリンクは24時間有効です。

よろしくお願いいたします。
オンラインショップモール運営チーム
`, token)

	return s.sendEmail(toEmail, subject, body)
}

// SendPasswordResetEmail はパスワードリセット用のメールを送信
func (s *SMTPEmailSender) SendPasswordResetEmail(ctx context.Context, toEmail, token string) error {
	subject := "パスワードリセット"
	body := fmt.Sprintf(`
こんにちは！

パスワードのリセットリクエストを受け付けました。
以下のトークンを使用してパスワードをリセットしてください。

リセットトークン: %s

このリンクは1時間有効です。

このリクエストに心当たりがない場合は、このメールを無視してください。

よろしくお願いいたします。
オンラインショップモール運営チーム
`, token)

	return s.sendEmail(toEmail, subject, body)
}

// sendEmail は実際のメール送信処理
func (s *SMTPEmailSender) sendEmail(to, subject, body string) error {
	// メッセージを構築
	message := fmt.Sprintf("From: %s\r\n", s.from)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "Content-Type: text/plain; charset=UTF-8\r\n"
	message += "\r\n"
	message += body

	// MailHog は認証不要なので、smtp.SendMail を直接使用
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	// MailHog は認証なしで動作
	err := smtp.SendMail(
		addr,
		nil, // 認証なし
		s.from,
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	fmt.Printf("✉️  Email sent to %s (subject: %s)\n", to, subject)
	return nil
}
