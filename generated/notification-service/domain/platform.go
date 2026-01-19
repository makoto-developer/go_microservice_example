package domain

// Platform represents Platform type
type Platform string

const (
	PlatformIos Platform = "IOS"
	PlatformAndroid Platform = "ANDROID"
	PlatformWeb Platform = "WEB"
)

// PlatformValues returns all possible values
func PlatformValues() []Platform {
	return []Platform{
		PlatformIos,
		PlatformAndroid,
		PlatformWeb,
	}
}

// IsValid checks if the value is valid
func (e Platform) IsValid() bool {
	switch e {
	case PlatformIos:
	case PlatformAndroid:
	case PlatformWeb:
		return true
	}
	return false
}
