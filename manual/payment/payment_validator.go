package payment

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidAmount      = errors.New("invalid amount: must be positive")
	ErrInvalidCurrency    = errors.New("invalid currency: only JPY supported")
	ErrAmountTooLarge     = errors.New("amount too large: exceeds maximum limit")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
)

const (
	MaxPaymentAmount = 10000000 // 1000万円
	MinPaymentAmount = 100      // 100円
)

// PaymentValidator は決済検証ロジック
type PaymentValidator struct{}

// NewPaymentValidator はPaymentValidatorを初期化
func NewPaymentValidator() *PaymentValidator {
	return &PaymentValidator{}
}

// ValidateAmount は金額を検証
func (v *PaymentValidator) ValidateAmount(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	if amount < MinPaymentAmount {
		return fmt.Errorf("amount too small: minimum is %d", MinPaymentAmount)
	}

	if amount > MaxPaymentAmount {
		return ErrAmountTooLarge
	}

	return nil
}

// ValidateCurrency は通貨コードを検証
func (v *PaymentValidator) ValidateCurrency(currency string) error {
	// 現状はJPYのみサポート
	if currency != "jpy" && currency != "JPY" {
		return ErrInvalidCurrency
	}

	return nil
}

// PaymentMethod は決済方法
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodCOD        PaymentMethod = "cash_on_delivery"
)

// ValidatePaymentMethod は決済方法を検証
func (v *PaymentValidator) ValidatePaymentMethod(method string) error {
	switch PaymentMethod(method) {
	case PaymentMethodCreditCard, PaymentMethodCOD:
		return nil
	default:
		return ErrInvalidPaymentMethod
	}
}

// CODFeeCalculator は代引き手数料計算
type CODFeeCalculator struct{}

// NewCODFeeCalculator はCODFeeCalculatorを初期化
func NewCODFeeCalculator() *CODFeeCalculator {
	return &CODFeeCalculator{}
}

// CalculateFee は代引き手数料を計算
func (c *CODFeeCalculator) CalculateFee(amount int64) int64 {
	// 代引き手数料のテーブル
	switch {
	case amount < 10000:
		return 330
	case amount < 30000:
		return 440
	case amount < 100000:
		return 660
	case amount < 300000:
		return 1100
	default:
		return 1650
	}
}
