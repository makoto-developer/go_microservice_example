package payment

import (
	"testing"
)

func TestPaymentValidator_ValidateAmount(t *testing.T) {
	validator := NewPaymentValidator()

	tests := []struct {
		name    string
		amount  int64
		wantErr bool
	}{
		{"Valid amount", 1000, false},
		{"Zero amount", 0, true},
		{"Negative amount", -100, true},
		{"Too small", 50, true},
		{"Minimum amount", MinPaymentAmount, false},
		{"Maximum amount", MaxPaymentAmount, false},
		{"Too large", MaxPaymentAmount + 1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateAmount(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAmount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPaymentValidator_ValidateCurrency(t *testing.T) {
	validator := NewPaymentValidator()

	tests := []struct {
		name     string
		currency string
		wantErr  bool
	}{
		{"JPY lowercase", "jpy", false},
		{"JPY uppercase", "JPY", false},
		{"USD", "usd", true},
		{"EUR", "eur", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCurrency(tt.currency)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateCurrency() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPaymentValidator_ValidatePaymentMethod(t *testing.T) {
	validator := NewPaymentValidator()

	tests := []struct {
		name    string
		method  string
		wantErr bool
	}{
		{"Credit card", "credit_card", false},
		{"COD", "cash_on_delivery", false},
		{"Invalid method", "bank_transfer", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidatePaymentMethod(tt.method)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePaymentMethod() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCODFeeCalculator_CalculateFee(t *testing.T) {
	calculator := NewCODFeeCalculator()

	tests := []struct {
		name   string
		amount int64
		want   int64
	}{
		{"Under 10000", 5000, 330},
		{"10000-30000", 15000, 440},
		{"30000-100000", 50000, 660},
		{"100000-300000", 150000, 1100},
		{"Over 300000", 400000, 1650},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculator.CalculateFee(tt.amount)
			if got != tt.want {
				t.Errorf("CalculateFee() = %v, want %v", got, tt.want)
			}
		})
	}
}
