package shipping

import (
	"testing"
)

func TestAddressNormalizer_NormalizePostalCode(t *testing.T) {
	normalizer := NewAddressNormalizer()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"With hyphen", "100-0001", "100-0001", false},
		{"Without hyphen", "1000001", "100-0001", false},
		{"With full-width hyphen", "100ー0001", "100-0001", false},
		{"Invalid format", "12345", "", true},
		{"Too long", "12345678", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizer.NormalizePostalCode(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizePostalCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NormalizePostalCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddressNormalizer_Normalize(t *testing.T) {
	normalizer := NewAddressNormalizer()

	addr, err := normalizer.Normalize("1000001", "東京", "千代田区", "千代田1-1", "")
	if err != nil {
		t.Fatalf("Normalize failed: %v", err)
	}

	if addr.PostalCode != "100-0001" {
		t.Errorf("Expected postal code 100-0001, got %s", addr.PostalCode)
	}

	if addr.Prefecture != "東京都" {
		t.Errorf("Expected prefecture 東京都, got %s", addr.Prefecture)
	}
}

func TestAddressNormalizer_ValidateAddress(t *testing.T) {
	normalizer := NewAddressNormalizer()

	tests := []struct {
		name    string
		addr    *NormalizedAddress
		wantErr bool
	}{
		{
			"Valid address",
			&NormalizedAddress{
				PostalCode: "100-0001",
				Prefecture: "東京都",
				City:       "千代田区",
				Address1:   "千代田1-1",
			},
			false,
		},
		{
			"Missing postal code",
			&NormalizedAddress{
				Prefecture: "東京都",
				City:       "千代田区",
				Address1:   "千代田1-1",
			},
			true,
		},
		{
			"Missing address1",
			&NormalizedAddress{
				PostalCode: "100-0001",
				Prefecture: "東京都",
				City:       "千代田区",
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := normalizer.ValidateAddress(tt.addr)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
