package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type SendNotificationRequest struct {
	Token        string                 `json:"token"`
	Notification map[string]string      `json:"notification"`
	Data         map[string]interface{} `json:"data,omitempty"`
}

type SendNotificationResponse struct {
	MessageID string `json:"message_id"`
	Success   bool   `json:"success"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/v1/projects/mock-project/messages:send", sendMessageHandler)

	log.Printf("FCM Mock Server listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	messageID := "fcm_" + uuid.New().String()

	log.Printf("[FCM Mock] Sent push notification: %s\n", messageID)
	log.Printf("  Token: %s\n", req.Token)
	log.Printf("  Title: %s\n", req.Notification["title"])
	log.Printf("  Body: %s\n", req.Notification["body"])

	resp := SendNotificationResponse{
		MessageID: messageID,
		Success:   true,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
