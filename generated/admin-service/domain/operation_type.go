package domain

// OperationType represents OperationType type
type OperationType string

const (
	OperationTypeUserCreated OperationType = "USER_CREATED"
	OperationTypeUserUpdated OperationType = "USER_UPDATED"
	OperationTypeUserRoleChanged OperationType = "USER_ROLE_CHANGED"
	OperationTypeUserSuspended OperationType = "USER_SUSPENDED"
	OperationTypeUserActivated OperationType = "USER_ACTIVATED"
	OperationTypeShopApproved OperationType = "SHOP_APPROVED"
	OperationTypeShopRejected OperationType = "SHOP_REJECTED"
	OperationTypeShopSuspended OperationType = "SHOP_SUSPENDED"
	OperationTypeShopActivated OperationType = "SHOP_ACTIVATED"
	OperationTypeSettingUpdated OperationType = "SETTING_UPDATED"
	OperationTypeCategoryCreated OperationType = "CATEGORY_CREATED"
	OperationTypeCategoryUpdated OperationType = "CATEGORY_UPDATED"
	OperationTypeCategoryDeleted OperationType = "CATEGORY_DELETED"
	OperationTypeReviewApproved OperationType = "REVIEW_APPROVED"
	OperationTypeReviewRejected OperationType = "REVIEW_REJECTED"
	OperationTypeReviewDeleted OperationType = "REVIEW_DELETED"
)

// OperationTypeValues returns all possible values
func OperationTypeValues() []OperationType {
	return []OperationType{
		OperationTypeUserCreated,
		OperationTypeUserUpdated,
		OperationTypeUserRoleChanged,
		OperationTypeUserSuspended,
		OperationTypeUserActivated,
		OperationTypeShopApproved,
		OperationTypeShopRejected,
		OperationTypeShopSuspended,
		OperationTypeShopActivated,
		OperationTypeSettingUpdated,
		OperationTypeCategoryCreated,
		OperationTypeCategoryUpdated,
		OperationTypeCategoryDeleted,
		OperationTypeReviewApproved,
		OperationTypeReviewRejected,
		OperationTypeReviewDeleted,
	}
}

// IsValid checks if the value is valid
func (e OperationType) IsValid() bool {
	switch e {
	case OperationTypeUserCreated:
	case OperationTypeUserUpdated:
	case OperationTypeUserRoleChanged:
	case OperationTypeUserSuspended:
	case OperationTypeUserActivated:
	case OperationTypeShopApproved:
	case OperationTypeShopRejected:
	case OperationTypeShopSuspended:
	case OperationTypeShopActivated:
	case OperationTypeSettingUpdated:
	case OperationTypeCategoryCreated:
	case OperationTypeCategoryUpdated:
	case OperationTypeCategoryDeleted:
	case OperationTypeReviewApproved:
	case OperationTypeReviewRejected:
	case OperationTypeReviewDeleted:
		return true
	}
	return false
}
