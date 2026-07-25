package domain

type Gender string

const (
	GenderMale           Gender = "MALE"
	GenderFemale         Gender = "FEMALE"
	GenderOther          Gender = "OTHER"
	GenderPreferNotToSay Gender = "PREFER_NOT_TO_SAY"
)

func (g Gender) String() string {
	return string(g)
}

func (g Gender) IsValid() bool {
	switch g {
	case GenderMale, GenderFemale, GenderOther, GenderPreferNotToSay:
		return true
	}
	return false
}
