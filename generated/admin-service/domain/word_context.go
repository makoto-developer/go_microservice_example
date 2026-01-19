package domain

// WordContext represents WordContext type
type WordContext string

const (
	WordContextReview WordContext = "REVIEW"
	WordContextChat WordContext = "CHAT"
	WordContextAll WordContext = "ALL"
)

// WordContextValues returns all possible values
func WordContextValues() []WordContext {
	return []WordContext{
		WordContextReview,
		WordContextChat,
		WordContextAll,
	}
}

// IsValid checks if the value is valid
func (e WordContext) IsValid() bool {
	switch e {
	case WordContextReview:
	case WordContextChat:
	case WordContextAll:
		return true
	}
	return false
}
