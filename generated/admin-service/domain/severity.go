package domain

// Severity represents Severity type
type Severity string

const (
	SeverityLow Severity = "LOW"
	SeverityMedium Severity = "MEDIUM"
	SeverityHigh Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

// SeverityValues returns all possible values
func SeverityValues() []Severity {
	return []Severity{
		SeverityLow,
		SeverityMedium,
		SeverityHigh,
		SeverityCritical,
	}
}

// IsValid checks if the value is valid
func (e Severity) IsValid() bool {
	switch e {
	case SeverityLow:
	case SeverityMedium:
	case SeverityHigh:
	case SeverityCritical:
		return true
	}
	return false
}
