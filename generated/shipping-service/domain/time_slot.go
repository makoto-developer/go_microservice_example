package domain

// TimeSlot represents TimeSlot type
type TimeSlot string

const (
	TimeSlotMorning TimeSlot = "MORNING"
	TimeSlotNoonTo2pm TimeSlot = "NOON_TO_2PM"
	TimeSlotPm2To4 TimeSlot = "PM_2_TO_4"
	TimeSlotPm4To6 TimeSlot = "PM_4_TO_6"
	TimeSlotPm6To8 TimeSlot = "PM_6_TO_8"
	TimeSlotPm7To9 TimeSlot = "PM_7_TO_9"
)

// TimeSlotValues returns all possible values
func TimeSlotValues() []TimeSlot {
	return []TimeSlot{
		TimeSlotMorning,
		TimeSlotNoonTo2pm,
		TimeSlotPm2To4,
		TimeSlotPm4To6,
		TimeSlotPm6To8,
		TimeSlotPm7To9,
	}
}

// IsValid checks if the value is valid
func (e TimeSlot) IsValid() bool {
	switch e {
	case TimeSlotMorning:
	case TimeSlotNoonTo2pm:
	case TimeSlotPm2To4:
	case TimeSlotPm4To6:
	case TimeSlotPm6To8:
	case TimeSlotPm7To9:
		return true
	}
	return false
}
