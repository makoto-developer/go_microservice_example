package domain

// PeriodType represents PeriodType type
type PeriodType string

const (
	PeriodTypeHourly PeriodType = "HOURLY"
	PeriodTypeDaily PeriodType = "DAILY"
	PeriodTypeWeekly PeriodType = "WEEKLY"
	PeriodTypeMonthly PeriodType = "MONTHLY"
)

// PeriodTypeValues returns all possible values
func PeriodTypeValues() []PeriodType {
	return []PeriodType{
		PeriodTypeHourly,
		PeriodTypeDaily,
		PeriodTypeWeekly,
		PeriodTypeMonthly,
	}
}

// IsValid checks if the value is valid
func (e PeriodType) IsValid() bool {
	switch e {
	case PeriodTypeHourly:
	case PeriodTypeDaily:
	case PeriodTypeWeekly:
	case PeriodTypeMonthly:
		return true
	}
	return false
}
