package domain

// MessageType represents MessageType type
type MessageType string

const (
	MessageTypeText MessageType = "TEXT"
	MessageTypeImage MessageType = "IMAGE"
	MessageTypeFile MessageType = "FILE"
)

// MessageTypeValues returns all possible values
func MessageTypeValues() []MessageType {
	return []MessageType{
		MessageTypeText,
		MessageTypeImage,
		MessageTypeFile,
	}
}

// IsValid checks if the value is valid
func (e MessageType) IsValid() bool {
	switch e {
	case MessageTypeText:
	case MessageTypeImage:
	case MessageTypeFile:
		return true
	}
	return false
}
