package shipping

import (
	"context"
	"testing"
)

func TestCarrierClient_CreateShipment(t *testing.T) {
	tests := []struct {
		name    string
		carrier CarrierType
	}{
		{"Yamato", CarrierYamato},
		{"Sagawa", CarrierSagawa},
		{"JapanPost", CarrierJapanPost},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewCarrierClient(tt.carrier, "test_api_key")
			ctx := context.Background()

			req := CreateShipmentRequest{
				SenderName:      "Test Sender",
				SenderAddress:   "Tokyo",
				ReceiverName:    "Test Receiver",
				ReceiverAddress: "Osaka",
				ReceiverPhone:   "090-1234-5678",
				Weight:          1000,
				Size:            "medium",
				COD:             false,
			}

			resp, err := client.CreateShipment(ctx, req)
			if err != nil {
				t.Fatalf("CreateShipment failed: %v", err)
			}

			if resp.TrackingNumber == "" {
				t.Error("Expected non-empty tracking number")
			}

			if len(resp.Label) == 0 {
				t.Error("Expected non-empty label")
			}
		})
	}
}

func TestCarrierClient_GetTracking(t *testing.T) {
	client := NewCarrierClient(CarrierYamato, "test_api_key")
	ctx := context.Background()

	info, err := client.GetTracking(ctx, "YM123456789")
	if err != nil {
		t.Fatalf("GetTracking failed: %v", err)
	}

	if info.TrackingNumber != "YM123456789" {
		t.Errorf("Expected tracking number YM123456789, got %s", info.TrackingNumber)
	}

	if info.Status != StatusInTransit {
		t.Errorf("Expected status in_transit, got %s", info.Status)
	}

	if len(info.History) == 0 {
		t.Error("Expected non-empty history")
	}
}

func TestCarrierClient_CancelShipment(t *testing.T) {
	client := NewCarrierClient(CarrierYamato, "test_api_key")
	ctx := context.Background()

	err := client.CancelShipment(ctx, "YM123456789")
	if err != nil {
		t.Fatalf("CancelShipment failed: %v", err)
	}
}

func TestCarrierClient_CalculateShippingFee(t *testing.T) {
	client := NewCarrierClient(CarrierYamato, "test_api_key")
	ctx := context.Background()

	tests := []struct {
		name   string
		weight int
		size   string
	}{
		{"Small package", 500, "small"},
		{"Medium package", 1500, "medium"},
		{"Large package", 3000, "large"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fee, err := client.CalculateShippingFee(ctx, tt.weight, tt.size, "100-0001", "530-0001")
			if err != nil {
				t.Fatalf("CalculateShippingFee failed: %v", err)
			}

			if fee <= 0 {
				t.Errorf("Expected positive fee, got %d", fee)
			}
		})
	}
}
