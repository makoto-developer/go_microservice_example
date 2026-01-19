package domain

// DeliveryOption represents DeliveryOption type
type DeliveryOption string

const (
	DeliveryOptionStandard DeliveryOption = "STANDARD"
	DeliveryOptionFrontDoor DeliveryOption = "FRONT_DOOR"
	DeliveryOptionDeliveryBox DeliveryOption = "DELIVERY_BOX"
	DeliveryOptionGasMeterBox DeliveryOption = "GAS_METER_BOX"
	DeliveryOptionGiftWrapping DeliveryOption = "GIFT_WRAPPING"
)

// DeliveryOptionValues returns all possible values
func DeliveryOptionValues() []DeliveryOption {
	return []DeliveryOption{
		DeliveryOptionStandard,
		DeliveryOptionFrontDoor,
		DeliveryOptionDeliveryBox,
		DeliveryOptionGasMeterBox,
		DeliveryOptionGiftWrapping,
	}
}

// IsValid checks if the value is valid
func (e DeliveryOption) IsValid() bool {
	switch e {
	case DeliveryOptionStandard:
	case DeliveryOptionFrontDoor:
	case DeliveryOptionDeliveryBox:
	case DeliveryOptionGasMeterBox:
	case DeliveryOptionGiftWrapping:
		return true
	}
	return false
}
