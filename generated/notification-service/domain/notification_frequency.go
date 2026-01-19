package domain

// NotificationFrequency represents NotificationFrequency type
type NotificationFrequency string

const (
	NotificationFrequencyImmediate NotificationFrequency = "IMMEDIATE"
	NotificationFrequencyDailyDigest NotificationFrequency = "DAILY_DIGEST"
	NotificationFrequencyDisabled NotificationFrequency = "DISABLED"
)

// NotificationFrequencyValues returns all possible values
func NotificationFrequencyValues() []NotificationFrequency {
	return []NotificationFrequency{
		NotificationFrequencyImmediate,
		NotificationFrequencyDailyDigest,
		NotificationFrequencyDisabled,
	}
}

// IsValid checks if the value is valid
func (e NotificationFrequency) IsValid() bool {
	switch e {
	case NotificationFrequencyImmediate:
	case NotificationFrequencyDailyDigest:
	case NotificationFrequencyDisabled:
		return true
	}
	return false
}
