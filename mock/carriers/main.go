package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

type TrackingRequest struct {
	TrackingNumber string `json:"tracking_number"`
	Carrier        string `json:"carrier"` // yamato, sagawa, japan_post
}

type TrackingResponse struct{
	TrackingNumber string            `json:"tracking_number"`
	Carrier        string            `json:"carrier"`
	Status         string            `json:"status"`
	StatusHistory  []TrackingHistory `json:"status_history"`
	EstimatedDelivery time.Time      `json:"estimated_delivery,omitempty"`
}

type TrackingHistory struct {
	Status    string    `json:"status"`
	Location  string    `json:"location"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/api/tracking", trackingHandler)

	log.Printf("Carriers Mock API Server listening on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func trackingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req TrackingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	log.Printf("[Carriers Mock] Tracking request for: %s (%s)\n", req.TrackingNumber, req.Carrier)

	// Mock tracking data
	statuses := []string{"picked_up", "in_transit", "out_for_delivery", "delivered"}
	currentStatus := statuses[rand.Intn(len(statuses))]

	history := []TrackingHistory{
		{
			Status:    "picked_up",
			Location:  "Tokyo Distribution Center",
			Timestamp: time.Now().Add(-48 * time.Hour),
		},
		{
			Status:    "in_transit",
			Location:  "Osaka Hub",
			Timestamp: time.Now().Add(-24 * time.Hour),
		},
	}

	if currentStatus == "out_for_delivery" || currentStatus == "delivered" {
		history = append(history, TrackingHistory{
			Status:    "out_for_delivery",
			Location:  "Local Delivery Center",
			Timestamp: time.Now().Add(-2 * time.Hour),
		})
	}

	if currentStatus == "delivered" {
		history = append(history, TrackingHistory{
			Status:    "delivered",
			Location:  "Customer Address",
			Timestamp: time.Now(),
		})
	}

	resp := TrackingResponse{
		TrackingNumber:    req.TrackingNumber,
		Carrier:           req.Carrier,
		Status:            currentStatus,
		StatusHistory:     history,
		EstimatedDelivery: time.Now().Add(24 * time.Hour),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
