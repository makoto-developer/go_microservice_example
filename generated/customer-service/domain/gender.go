package domain

// Gender represents Gender type
type Gender string

const (
	GenderMale Gender = "MALE"
	GenderFemale Gender = "FEMALE"
	GenderOther Gender = "OTHER"
	GenderPreferNotToSay Gender = "PREFER_NOT_TO_SAY"
)

// GenderValues returns all possible values
func GenderValues() []Gender {
	return []Gender{
		GenderMale,
		GenderFemale,
		GenderOther,
		GenderPreferNotToSay,
	}
}

// IsValid checks if the value is valid
func (e Gender) IsValid() bool {
	switch e {
	case GenderMale:
	case GenderFemale:
	case GenderOther:
	case GenderPreferNotToSay:
		return true
	}
	return false
}
