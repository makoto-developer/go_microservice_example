package domain

// Role represents Role type
type Role string

const (
	RoleCustomer Role = "CUSTOMER"
	RoleShopOwner Role = "SHOP_OWNER"
	RoleAdmin Role = "ADMIN"
)

// RoleValues returns all possible values
func RoleValues() []Role {
	return []Role{
		RoleCustomer,
		RoleShopOwner,
		RoleAdmin,
	}
}

// IsValid checks if the value is valid
func (e Role) IsValid() bool {
	switch e {
	case RoleCustomer:
	case RoleShopOwner:
	case RoleAdmin:
		return true
	}
	return false
}
