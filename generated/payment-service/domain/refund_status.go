package domain

// RefundStatus represents RefundStatus type
type RefundStatus string

const (
	RefundStatusPending RefundStatus = "PENDING"
	RefundStatusProcessing RefundStatus = "PROCESSING"
	RefundStatusSucceeded RefundStatus = "SUCCEEDED"
	RefundStatusFailed RefundStatus = "FAILED"
)

// RefundStatusValues returns all possible values
func RefundStatusValues() []RefundStatus {
	return []RefundStatus{
		RefundStatusPending,
		RefundStatusProcessing,
		RefundStatusSucceeded,
		RefundStatusFailed,
	}
}

// IsValid checks if the value is valid
func (e RefundStatus) IsValid() bool {
	switch e {
	case RefundStatusPending:
	case RefundStatusProcessing:
	case RefundStatusSucceeded:
	case RefundStatusFailed:
		return true
	}
	return false
}
