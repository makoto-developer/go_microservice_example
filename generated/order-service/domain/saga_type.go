package domain

// SagaType represents SagaType type
type SagaType string

const (
	SagaTypeCreateOrder SagaType = "CREATE_ORDER"
	SagaTypeCancelOrder SagaType = "CANCEL_ORDER"
)

// SagaTypeValues returns all possible values
func SagaTypeValues() []SagaType {
	return []SagaType{
		SagaTypeCreateOrder,
		SagaTypeCancelOrder,
	}
}

// IsValid checks if the value is valid
func (e SagaType) IsValid() bool {
	switch e {
	case SagaTypeCreateOrder:
	case SagaTypeCancelOrder:
		return true
	}
	return false
}
