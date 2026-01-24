package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func main() {
	from := "noreply@shopmall.local"
	to := []string{"test@example.com"}
	subject := "Test Email"
	body := "This is a test email from Go"

	message := fmt.Sprintf("From: %s\r\n", from) +
		fmt.Sprintf("To: %s\r\n", to[0]) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		body

	addr := "localhost:22102"

	err := smtp.SendMail(addr, nil, from, to, []byte(message))
	if err != nil {
		log.Fatalf("Failed to send email: %v", err)
	}

	fmt.Println("✅ Test email sent successfully!")
	fmt.Println("Check Mailhog UI at http://localhost:22103")
}
