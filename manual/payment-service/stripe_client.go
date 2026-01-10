package custom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type StripeClient struct {
	BaseURL   string
	APIKey    string
	WebhookSecret string
	httpClient *http.Client
}

type ChargeRequest struct {
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	Source      string `json:"source"`
	Description string `json:"description"`
}

type ChargeResponse struct {
	ID            string `json:"id"`
	Amount        int64  `json:"amount"`
	Currency      string `json:"currency"`
	Status        string `json:"status"`
	Created       int64  `json:"created"`
	Description   string `json:"description"`
	FailureCode   string `json:"failure_code,omitempty"`
	FailureMessage string `json:"failure_message,omitempty"`
}

type RefundRequest struct {
	ChargeID string `json:"charge_id"`
	Amount   int64  `json:"amount,omitempty"`
	Reason   string `json:"reason,omitempty"`
}

type RefundResponse struct {
	ID       string `json:"id"`
	ChargeID string `json:"charge_id"`
	Amount   int64  `json:"amount"`
	Status   string `json:"status"`
	Created  int64  `json:"created"`
	Reason   string `json:"reason"`
}

func NewStripeClient() *StripeClient {
	baseURL := os.Getenv("STRIPE_API_URL")
	if baseURL == "" {
		baseURL = "http://localhost:20010"
	}

	apiKey := os.Getenv("STRIPE_API_KEY")
	if apiKey == "" {
		apiKey = "sk_test_mock"
	}

	webhookSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	return &StripeClient{
		BaseURL:       baseURL,
		APIKey:        apiKey,
		WebhookSecret: webhookSecret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *StripeClient) CreateCharge(req ChargeRequest) (*ChargeResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/v1/charges", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stripe API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var chargeResp ChargeResponse
	if err := json.Unmarshal(body, &chargeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &chargeResp, nil
}

func (c *StripeClient) CreateRefund(req RefundRequest) (*RefundResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL+"/v1/refunds", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("stripe API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	var refundResp RefundResponse
	if err := json.Unmarshal(body, &refundResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &refundResp, nil
}

func (c *StripeClient) ValidateAmount(amount int64, currency string) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be positive: got %d", amount)
	}

	switch currency {
	case "JPY":
		if amount < 50 {
			return fmt.Errorf("amount must be at least 50 JPY")
		}
	case "USD":
		if amount < 50 {
			return fmt.Errorf("amount must be at least 0.50 USD (50 cents)")
		}
	default:
		return fmt.Errorf("unsupported currency: %s", currency)
	}

	return nil
}
