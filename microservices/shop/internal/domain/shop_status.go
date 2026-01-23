package domain

type ShopStatus string

const (
	ShopStatusPendingApproval ShopStatus = "PENDING_APPROVAL"
	ShopStatusApproved        ShopStatus = "APPROVED"
	ShopStatusRejected        ShopStatus = "REJECTED"
	ShopStatusSuspended       ShopStatus = "SUSPENDED"
)

func (s ShopStatus) String() string {
	return string(s)
}

func (s ShopStatus) IsValid() bool {
	switch s {
	case ShopStatusPendingApproval, ShopStatusApproved, ShopStatusRejected, ShopStatusSuspended:
		return true
	}
	return false
}
