package domain

// NotificationStatus represents NotificationStatus type
type NotificationStatus string

const (
	NotificationStatusPending NotificationStatus = "PENDING"
	NotificationStatusSent NotificationStatus = "SENT"
	NotificationStatusFailed NotificationStatus = "FAILED"
	NotificationStatusCancelled NotificationStatus = "CANCELLED"
)

// NotificationStatusValues returns all possible values
func NotificationStatusValues() []NotificationStatus {
	return []NotificationStatus{
		NotificationStatusPending,
		NotificationStatusSent,
		NotificationStatusFailed,
		NotificationStatusCancelled,
	}
}

// IsValid checks if the value is valid
func (e NotificationStatus) IsValid() bool {
	switch e {
	case NotificationStatusPending:
	case NotificationStatusSent:
	case NotificationStatusFailed:
	case NotificationStatusCancelled:
		return true
	}
	return false
}
