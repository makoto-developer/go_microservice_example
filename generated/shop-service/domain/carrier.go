package domain

// Carrier represents Carrier type
type Carrier string

const (
	CarrierYamato    Carrier = "YAMATO"
	CarrierSagawa    Carrier = "SAGAWA"
	CarrierJapanPost Carrier = "JAPAN_POST"
)

// CarrierValues returns all possible values
func CarrierValues() []Carrier {
	return []Carrier{
		CarrierYamato,
		CarrierSagawa,
		CarrierJapanPost,
	}
}

// IsValid checks if the value is valid
func (e Carrier) IsValid() bool {
	switch e {
	case CarrierYamato:
	case CarrierSagawa:
	case CarrierJapanPost:
		return true
	}
	return false
}
