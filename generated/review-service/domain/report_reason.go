package domain

// ReportReason represents ReportReason type
type ReportReason string

const (
	ReportReasonInappropriateContent ReportReason = "INAPPROPRIATE_CONTENT"
	ReportReasonSpam ReportReason = "SPAM"
	ReportReasonOffTopic ReportReason = "OFF_TOPIC"
	ReportReasonFalseInformation ReportReason = "FALSE_INFORMATION"
	ReportReasonPersonalInformation ReportReason = "PERSONAL_INFORMATION"
	ReportReasonOther ReportReason = "OTHER"
)

// ReportReasonValues returns all possible values
func ReportReasonValues() []ReportReason {
	return []ReportReason{
		ReportReasonInappropriateContent,
		ReportReasonSpam,
		ReportReasonOffTopic,
		ReportReasonFalseInformation,
		ReportReasonPersonalInformation,
		ReportReasonOther,
	}
}

// IsValid checks if the value is valid
func (e ReportReason) IsValid() bool {
	switch e {
	case ReportReasonInappropriateContent:
	case ReportReasonSpam:
	case ReportReasonOffTopic:
	case ReportReasonFalseInformation:
	case ReportReasonPersonalInformation:
	case ReportReasonOther:
		return true
	}
	return false
}
