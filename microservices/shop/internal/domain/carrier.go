package domain

type Carrier string

const (
	CarrierYamato    Carrier = "YAMATO"
	CarrierSagawa    Carrier = "SAGAWA"
	CarrierJapanPost Carrier = "JAPAN_POST"
)

func (c Carrier) String() string {
	return string(c)
}

func (c Carrier) IsValid() bool {
	switch c {
	case CarrierYamato, CarrierSagawa, CarrierJapanPost:
		return true
	}
	return false
}
