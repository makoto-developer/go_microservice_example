package domain

// NotificationChannel represents NotificationChannel type
type NotificationChannel string

const (
	NotificationChannelEmail NotificationChannel = "EMAIL"
	NotificationChannelPush NotificationChannel = "PUSH"
	NotificationChannelSms NotificationChannel = "SMS"
)

// NotificationChannelValues returns all possible values
func NotificationChannelValues() []NotificationChannel {
	return []NotificationChannel{
		NotificationChannelEmail,
		NotificationChannelPush,
		NotificationChannelSms,
	}
}

// IsValid checks if the value is valid
func (e NotificationChannel) IsValid() bool {
	switch e {
	case NotificationChannelEmail:
	case NotificationChannelPush:
	case NotificationChannelSms:
		return true
	}
	return false
}
