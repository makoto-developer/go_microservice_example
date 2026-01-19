package domain

// HealthStatus represents HealthStatus type
type HealthStatus string

const (
	HealthStatusHealthy HealthStatus = "HEALTHY"
	HealthStatusDegraded HealthStatus = "DEGRADED"
	HealthStatusUnhealthy HealthStatus = "UNHEALTHY"
)

// HealthStatusValues returns all possible values
func HealthStatusValues() []HealthStatus {
	return []HealthStatus{
		HealthStatusHealthy,
		HealthStatusDegraded,
		HealthStatusUnhealthy,
	}
}

// IsValid checks if the value is valid
func (e HealthStatus) IsValid() bool {
	switch e {
	case HealthStatusHealthy:
	case HealthStatusDegraded:
	case HealthStatusUnhealthy:
		return true
	}
	return false
}
