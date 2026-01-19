package domain

// ReportStatus represents ReportStatus type
type ReportStatus string

const (
	ReportStatusPending ReportStatus = "PENDING"
	ReportStatusReviewed ReportStatus = "REVIEWED"
	ReportStatusActionTaken ReportStatus = "ACTION_TAKEN"
	ReportStatusDismissed ReportStatus = "DISMISSED"
)

// ReportStatusValues returns all possible values
func ReportStatusValues() []ReportStatus {
	return []ReportStatus{
		ReportStatusPending,
		ReportStatusReviewed,
		ReportStatusActionTaken,
		ReportStatusDismissed,
	}
}

// IsValid checks if the value is valid
func (e ReportStatus) IsValid() bool {
	switch e {
	case ReportStatusPending:
	case ReportStatusReviewed:
	case ReportStatusActionTaken:
	case ReportStatusDismissed:
		return true
	}
	return false
}
