package domain

// PresenceStatus represents PresenceStatus type
type PresenceStatus string

const (
	PresenceStatusOnline PresenceStatus = "ONLINE"
	PresenceStatusOffline PresenceStatus = "OFFLINE"
	PresenceStatusAway PresenceStatus = "AWAY"
)

// PresenceStatusValues returns all possible values
func PresenceStatusValues() []PresenceStatus {
	return []PresenceStatus{
		PresenceStatusOnline,
		PresenceStatusOffline,
		PresenceStatusAway,
	}
}

// IsValid checks if the value is valid
func (e PresenceStatus) IsValid() bool {
	switch e {
	case PresenceStatusOnline:
	case PresenceStatusOffline:
	case PresenceStatusAway:
		return true
	}
	return false
}
