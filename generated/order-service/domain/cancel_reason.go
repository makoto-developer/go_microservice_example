package domain

// CancelReason represents CancelReason type
type CancelReason string

const (
	CancelReasonCustomerNoLongerNeeded CancelReason = "CUSTOMER_NO_LONGER_NEEDED"
	CancelReasonCustomerOrderedByMistake CancelReason = "CUSTOMER_ORDERED_BY_MISTAKE"
	CancelReasonCustomerDeliveryTimeIssue CancelReason = "CUSTOMER_DELIVERY_TIME_ISSUE"
	CancelReasonCustomerOther CancelReason = "CUSTOMER_OTHER"
	CancelReasonShopOutOfStock CancelReason = "SHOP_OUT_OF_STOCK"
	CancelReasonShopDefectiveProduct CancelReason = "SHOP_DEFECTIVE_PRODUCT"
	CancelReasonShopOther CancelReason = "SHOP_OTHER"
)

// CancelReasonValues returns all possible values
func CancelReasonValues() []CancelReason {
	return []CancelReason{
		CancelReasonCustomerNoLongerNeeded,
		CancelReasonCustomerOrderedByMistake,
		CancelReasonCustomerDeliveryTimeIssue,
		CancelReasonCustomerOther,
		CancelReasonShopOutOfStock,
		CancelReasonShopDefectiveProduct,
		CancelReasonShopOther,
	}
}

// IsValid checks if the value is valid
func (e CancelReason) IsValid() bool {
	switch e {
	case CancelReasonCustomerNoLongerNeeded:
	case CancelReasonCustomerOrderedByMistake:
	case CancelReasonCustomerDeliveryTimeIssue:
	case CancelReasonCustomerOther:
	case CancelReasonShopOutOfStock:
	case CancelReasonShopDefectiveProduct:
	case CancelReasonShopOther:
		return true
	}
	return false
}
