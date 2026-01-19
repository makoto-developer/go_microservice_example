package domain

// PaymentMethodType represents PaymentMethodType type
type PaymentMethodType string

const (
	PaymentMethodTypeCreditCard PaymentMethodType = "CREDIT_CARD"
	PaymentMethodTypeCashOnDelivery PaymentMethodType = "CASH_ON_DELIVERY"
)

// PaymentMethodTypeValues returns all possible values
func PaymentMethodTypeValues() []PaymentMethodType {
	return []PaymentMethodType{
		PaymentMethodTypeCreditCard,
		PaymentMethodTypeCashOnDelivery,
	}
}

// IsValid checks if the value is valid
func (e PaymentMethodType) IsValid() bool {
	switch e {
	case PaymentMethodTypeCreditCard:
	case PaymentMethodTypeCashOnDelivery:
		return true
	}
	return false
}
