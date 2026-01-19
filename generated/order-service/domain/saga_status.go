package domain

// SagaStatus represents SagaStatus type
type SagaStatus string

const (
	SagaStatusStarted SagaStatus = "STARTED"
	SagaStatusCompensating SagaStatus = "COMPENSATING"
	SagaStatusCompleted SagaStatus = "COMPLETED"
	SagaStatusFailed SagaStatus = "FAILED"
)

// SagaStatusValues returns all possible values
func SagaStatusValues() []SagaStatus {
	return []SagaStatus{
		SagaStatusStarted,
		SagaStatusCompensating,
		SagaStatusCompleted,
		SagaStatusFailed,
	}
}

// IsValid checks if the value is valid
func (e SagaStatus) IsValid() bool {
	switch e {
	case SagaStatusStarted:
	case SagaStatusCompensating:
	case SagaStatusCompleted:
	case SagaStatusFailed:
		return true
	}
	return false
}
