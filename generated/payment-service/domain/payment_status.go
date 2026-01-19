package domain

// PaymentStatus represents PaymentStatus type
type PaymentStatus string

const (
	PaymentStatusPending PaymentStatus = "PENDING"
	PaymentStatusProcessing PaymentStatus = "PROCESSING"
	PaymentStatusRequiresAuthentication PaymentStatus = "REQUIRES_AUTHENTICATION"
	PaymentStatusSucceeded PaymentStatus = "SUCCEEDED"
	PaymentStatusFailed PaymentStatus = "FAILED"
	PaymentStatusRefunding PaymentStatus = "REFUNDING"
	PaymentStatusRefunded PaymentStatus = "REFUNDED"
)

// PaymentStatusValues returns all possible values
func PaymentStatusValues() []PaymentStatus {
	return []PaymentStatus{
		PaymentStatusPending,
		PaymentStatusProcessing,
		PaymentStatusRequiresAuthentication,
		PaymentStatusSucceeded,
		PaymentStatusFailed,
		PaymentStatusRefunding,
		PaymentStatusRefunded,
	}
}

// IsValid checks if the value is valid
func (e PaymentStatus) IsValid() bool {
	switch e {
	case PaymentStatusPending:
	case PaymentStatusProcessing:
	case PaymentStatusRequiresAuthentication:
	case PaymentStatusSucceeded:
	case PaymentStatusFailed:
	case PaymentStatusRefunding:
	case PaymentStatusRefunded:
		return true
	}
	return false
}
