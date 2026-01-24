package usecase

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

type EmailService struct {
	smtpHost string
	smtpPort string
	from     string
}

func NewEmailService() *EmailService {
	return &EmailService{
		smtpHost: getEnv("SMTP_HOST", "localhost"),
		smtpPort: getEnv("SMTP_PORT", "22102"),
		from:     getEnv("SMTP_FROM", "noreply@shopmall.local"),
	}
}

func (s *EmailService) SendPasswordResetEmail(toEmail, resetToken string) error {
	// パスワードリセットURL（フロントエンドのURL）
	resetURL := fmt.Sprintf("http://localhost:22200/auth/reset-password?token=%s", resetToken)

	htmlTemplate := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>パスワードリセット</title>
</head>
<body style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px;">
    <h2 style="color: #333;">パスワードリセットのリクエスト</h2>
    <p>以下のリンクをクリックして、新しいパスワードを設定してください。</p>
    <p style="margin: 30px 0;">
        <a href="{{.ResetURL}}"
           style="background-color: #4CAF50; color: white; padding: 14px 20px; text-decoration: none; border-radius: 4px; display: inline-block;">
            パスワードをリセット
        </a>
    </p>
    <p style="color: #666; font-size: 14px;">
        このリンクは24時間有効です。<br>
        もしパスワードリセットをリクエストしていない場合は、このメールを無視してください。
    </p>
    <hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
    <p style="color: #999; font-size: 12px;">
        Shop Mall - Online Shopping Platform
    </p>
</body>
</html>
`

	tmpl, err := template.New("password_reset").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse email template: %w", err)
	}

	var body bytes.Buffer
	data := struct {
		ResetURL string
	}{
		ResetURL: resetURL,
	}

	if err := tmpl.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	subject := "パスワードリセットのリクエスト"
	message := fmt.Sprintf("From: %s\r\n", s.from) +
		fmt.Sprintf("To: %s\r\n", toEmail) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body.String()

	addr := fmt.Sprintf("%s:%s", s.smtpHost, s.smtpPort)

	// Mailhogは認証不要
	err = smtp.SendMail(addr, nil, s.from, []string{toEmail}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
