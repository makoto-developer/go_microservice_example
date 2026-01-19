package domain

// Carrier represents Carrier type
type Carrier string

const (
	CarrierYamato Carrier = "YAMATO"
	CarrierSagawa Carrier = "SAGAWA"
	CarrierJapanPost Carrier = "JAPAN_POST"
	CarrierClickpost Carrier = "CLICKPOST"
	CarrierNekopos Carrier = "NEKOPOS"
)

// CarrierValues returns all possible values
func CarrierValues() []Carrier {
	return []Carrier{
		CarrierYamato,
		CarrierSagawa,
		CarrierJapanPost,
		CarrierClickpost,
		CarrierNekopos,
	}
}

// IsValid checks if the value is valid
func (e Carrier) IsValid() bool {
	switch e {
	case CarrierYamato:
	case CarrierSagawa:
	case CarrierJapanPost:
	case CarrierClickpost:
	case CarrierNekopos:
		return true
	}
	return false
}
