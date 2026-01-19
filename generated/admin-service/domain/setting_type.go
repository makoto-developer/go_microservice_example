package domain

// SettingType represents SettingType type
type SettingType string

const (
	SettingTypeString SettingType = "STRING"
	SettingTypeNumber SettingType = "NUMBER"
	SettingTypeBoolean SettingType = "BOOLEAN"
	SettingTypeJson SettingType = "JSON"
)

// SettingTypeValues returns all possible values
func SettingTypeValues() []SettingType {
	return []SettingType{
		SettingTypeString,
		SettingTypeNumber,
		SettingTypeBoolean,
		SettingTypeJson,
	}
}

// IsValid checks if the value is valid
func (e SettingType) IsValid() bool {
	switch e {
	case SettingTypeString:
	case SettingTypeNumber:
	case SettingTypeBoolean:
	case SettingTypeJson:
		return true
	}
	return false
}
