package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

type PaymentIntent struct {
	ID            string `json:"id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	Status        string `json:"status"` // requires_payment_method, succeeded, failed
	ClientSecret  string `json:"client_secret"`
	PaymentMethod string `json:"payment_method,omitempty"`
	Created       int64  `json:"created"`
}

type Refund struct {
	ID              string `json:"id"`
	Amount          int64  `json:"amount"`
	PaymentIntentID string `json:"payment_intent"`
	Status          string `json:"status"` // succeeded, failed
	Created         int64  `json:"created"`
}

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

var (
	paymentIntents = make(map[string]*PaymentIntent)
	refunds        = make(map[string]*Refund)
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/v1/payment_intents", paymentIntentsHandler)
	http.HandleFunc("/v1/payment_intents/", paymentIntentDetailHandler)
	http.HandleFunc("/v1/refunds", refundsHandler)

	log.Printf("Stripe Mock Server listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func paymentIntentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Amount   int64  `json:"amount"`
		Currency string `json:"currency"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "invalid_request_error", "Invalid request body", http.StatusBadRequest)
		return
	}

	pi := &PaymentIntent{
		ID:           "pi_" + uuid.New().String(),
		Amount:       req.Amount,
		Currency:     req.Currency,
		Status:       "requires_payment_method",
		ClientSecret: "pi_secret_" + uuid.New().String(),
		Created:      time.Now().Unix(),
	}

	paymentIntents[pi.ID] = pi

	log.Printf("[Stripe Mock] Created Payment Intent: %s, Amount: %d %s\n", pi.ID, pi.Amount, pi.Currency)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pi)
}

func paymentIntentDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract ID from path: /v1/payment_intents/{id} or /v1/payment_intents/{id}/confirm
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/v1/payment_intents/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		sendError(w, "invalid_request_error", "Payment Intent ID required", http.StatusBadRequest)
		return
	}

	piID := parts[0]
	pi, exists := paymentIntents[piID]

	if !exists {
		sendError(w, "resource_missing_error", "Payment Intent not found", http.StatusNotFound)
		return
	}

	// GET /v1/payment_intents/{id}
	if r.Method == http.MethodGet && len(parts) == 1 {
		log.Printf("[Stripe Mock] Retrieved Payment Intent: %s\n", piID)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pi)
		return
	}

	// POST /v1/payment_intents/{id}/confirm
	if r.Method == http.MethodPost && len(parts) == 2 && parts[1] == "confirm" {
		var req struct {
			PaymentMethod string `json:"payment_method"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendError(w, "invalid_request_error", "Invalid request body", http.StatusBadRequest)
			return
		}

		pi.PaymentMethod = req.PaymentMethod
		pi.Status = "succeeded" // Mock: always succeed

		log.Printf("[Stripe Mock] Confirmed Payment Intent: %s with method: %s\n", piID, req.PaymentMethod)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pi)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func refundsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Amount          int64  `json:"amount"`
		PaymentIntentID string `json:"payment_intent"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "invalid_request_error", "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate payment intent exists
	pi, exists := paymentIntents[req.PaymentIntentID]
	if !exists {
		sendError(w, "resource_missing_error", "Payment Intent not found", http.StatusNotFound)
		return
	}

	// Validate amount
	if req.Amount > pi.Amount {
		sendError(w, "invalid_request_error", "Refund amount exceeds payment amount", http.StatusBadRequest)
		return
	}

	refund := &Refund{
		ID:              "re_" + uuid.New().String(),
		Amount:          req.Amount,
		PaymentIntentID: req.PaymentIntentID,
		Status:          "succeeded",
		Created:         time.Now().Unix(),
	}

	refunds[refund.ID] = refund

	log.Printf("[Stripe Mock] Created Refund: %s, Amount: %d for PI: %s\n", refund.ID, refund.Amount, req.PaymentIntentID)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(refund)
}

func sendError(w http.ResponseWriter, errType, message string, statusCode int) {
	resp := ErrorResponse{}
	resp.Error.Type = errType
	resp.Error.Message = message

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(resp)
}
