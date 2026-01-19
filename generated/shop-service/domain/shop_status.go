package domain

// ShopStatus represents ShopStatus type
type ShopStatus string

const (
	ShopStatusPendingApproval ShopStatus = "PENDING_APPROVAL"
	ShopStatusApproved        ShopStatus = "APPROVED"
	ShopStatusRejected        ShopStatus = "REJECTED"
	ShopStatusSuspended       ShopStatus = "SUSPENDED"
)

// ShopStatusValues returns all possible values
func ShopStatusValues() []ShopStatus {
	return []ShopStatus{
		ShopStatusPendingApproval,
		ShopStatusApproved,
		ShopStatusRejected,
		ShopStatusSuspended,
	}
}

// IsValid checks if the value is valid
func (e ShopStatus) IsValid() bool {
	switch e {
	case ShopStatusPendingApproval:
	case ShopStatusApproved:
	case ShopStatusRejected:
	case ShopStatusSuspended:
		return true
	}
	return false
}
