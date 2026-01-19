package domain

// ReviewStatus represents ReviewStatus type
type ReviewStatus string

const (
	ReviewStatusPending ReviewStatus = "PENDING"
	ReviewStatusApproved ReviewStatus = "APPROVED"
	ReviewStatusRejected ReviewStatus = "REJECTED"
	ReviewStatusDeleted ReviewStatus = "DELETED"
)

// ReviewStatusValues returns all possible values
func ReviewStatusValues() []ReviewStatus {
	return []ReviewStatus{
		ReviewStatusPending,
		ReviewStatusApproved,
		ReviewStatusRejected,
		ReviewStatusDeleted,
	}
}

// IsValid checks if the value is valid
func (e ReviewStatus) IsValid() bool {
	switch e {
	case ReviewStatusPending:
	case ReviewStatusApproved:
	case ReviewStatusRejected:
	case ReviewStatusDeleted:
		return true
	}
	return false
}
