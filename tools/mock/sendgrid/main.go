package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type SendEmailRequest struct {
	Personalizations []Personalization `json:"personalizations"`
	From             Email             `json:"from"`
	Subject          string            `json:"subject"`
	Content          []Content         `json:"content"`
}

type Personalization struct {
	To      []Email `json:"to"`
	Subject string  `json:"subject,omitempty"`
}

type Email struct {
	Email string `json:"email"`
	Name  string `json:"name,omitempty"`
}

type Content struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type SendEmailResponse struct {
	MessageID string `json:"message_id"`
	Status    string `json:"status"`
}

type EmailRecord struct {
	MessageID string             `json:"message_id"`
	Request   SendEmailRequest   `json:"request"`
	SentAt    time.Time          `json:"sent_at"`
}

var sentEmails = make(map[string]*EmailRecord)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/v3/mail/send", sendMailHandler)
	http.HandleFunc("/v3/mail/", mailDetailHandler)

	log.Printf("SendGrid Mock Server listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func sendMailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// Validation
	if len(req.Personalizations) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Personalizations required"})
		return
	}

	if req.From.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "From email required"})
		return
	}

	messageID := "msg_" + uuid.New().String()

	// Store email record
	record := &EmailRecord{
		MessageID: messageID,
		Request:   req,
		SentAt:    time.Now(),
	}
	sentEmails[messageID] = record

	log.Printf("[SendGrid Mock] Sent email: %s\n", messageID)
	log.Printf("  From: %s <%s>\n", req.From.Name, req.From.Email)
	log.Printf("  To: %s <%s>\n", req.Personalizations[0].To[0].Name, req.Personalizations[0].To[0].Email)
	log.Printf("  Subject: %s\n", req.Subject)

	resp := SendEmailResponse{
		MessageID: messageID,
		Status:    "sent",
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

func mailDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// GET /v3/mail/list - List all sent emails (for testing)
	if r.URL.Path == "/v3/mail/list" {
		emails := make([]*EmailRecord, 0, len(sentEmails))
		for _, record := range sentEmails {
			emails = append(emails, record)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"total":  len(emails),
			"emails": emails,
		})
		return
	}

	http.Error(w, "Not found", http.StatusNotFound)
}
