package domain

// PaymentMethod represents PaymentMethod type
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "CREDIT_CARD"
	PaymentMethodCashOnDelivery PaymentMethod = "CASH_ON_DELIVERY"
)

// PaymentMethodValues returns all possible values
func PaymentMethodValues() []PaymentMethod {
	return []PaymentMethod{
		PaymentMethodCreditCard,
		PaymentMethodCashOnDelivery,
	}
}

// IsValid checks if the value is valid
func (e PaymentMethod) IsValid() bool {
	switch e {
	case PaymentMethodCreditCard:
	case PaymentMethodCashOnDelivery:
		return true
	}
	return false
}
