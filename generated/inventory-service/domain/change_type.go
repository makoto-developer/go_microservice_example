package domain

// ChangeType represents ChangeType type
type ChangeType string

const (
	ChangeTypeInitial ChangeType = "INITIAL"
	ChangeTypeRestock ChangeType = "RESTOCK"
	ChangeTypeReturn ChangeType = "RETURN"
	ChangeTypeReservation ChangeType = "RESERVATION"
	ChangeTypeRelease ChangeType = "RELEASE"
	ChangeTypeConfirmation ChangeType = "CONFIRMATION"
	ChangeTypeDamage ChangeType = "DAMAGE"
	ChangeTypeStockTaking ChangeType = "STOCK_TAKING"
	ChangeTypeManualAdjustment ChangeType = "MANUAL_ADJUSTMENT"
)

// ChangeTypeValues returns all possible values
func ChangeTypeValues() []ChangeType {
	return []ChangeType{
		ChangeTypeInitial,
		ChangeTypeRestock,
		ChangeTypeReturn,
		ChangeTypeReservation,
		ChangeTypeRelease,
		ChangeTypeConfirmation,
		ChangeTypeDamage,
		ChangeTypeStockTaking,
		ChangeTypeManualAdjustment,
	}
}

// IsValid checks if the value is valid
func (e ChangeType) IsValid() bool {
	switch e {
	case ChangeTypeInitial:
	case ChangeTypeRestock:
	case ChangeTypeReturn:
	case ChangeTypeReservation:
	case ChangeTypeRelease:
	case ChangeTypeConfirmation:
	case ChangeTypeDamage:
	case ChangeTypeStockTaking:
	case ChangeTypeManualAdjustment:
		return true
	}
	return false
}
