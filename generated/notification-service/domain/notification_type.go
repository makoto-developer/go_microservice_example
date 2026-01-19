package domain

// NotificationType represents NotificationType type
type NotificationType string

const (
	NotificationTypeUserRegistered NotificationType = "USER_REGISTERED"
	NotificationTypeEmailVerified NotificationType = "EMAIL_VERIFIED"
	NotificationTypeOrderConfirmed NotificationType = "ORDER_CONFIRMED"
	NotificationTypeOrderShipped NotificationType = "ORDER_SHIPPED"
	NotificationTypeOrderDelivered NotificationType = "ORDER_DELIVERED"
	NotificationTypeOrderCancelled NotificationType = "ORDER_CANCELLED"
	NotificationTypePaymentCompleted NotificationType = "PAYMENT_COMPLETED"
	NotificationTypePaymentFailed NotificationType = "PAYMENT_FAILED"
	NotificationTypeRefundCompleted NotificationType = "REFUND_COMPLETED"
	NotificationTypeShopApproved NotificationType = "SHOP_APPROVED"
	NotificationTypeShopRejected NotificationType = "SHOP_REJECTED"
	NotificationTypeStockLowAlert NotificationType = "STOCK_LOW_ALERT"
	NotificationTypeStockOutAlert NotificationType = "STOCK_OUT_ALERT"
	NotificationTypeStockRestored NotificationType = "STOCK_RESTORED"
	NotificationTypeChatMessageReceived NotificationType = "CHAT_MESSAGE_RECEIVED"
	NotificationTypePasswordReset NotificationType = "PASSWORD_RESET"
	NotificationTypeCampaignStarted NotificationType = "CAMPAIGN_STARTED"
)

// NotificationTypeValues returns all possible values
func NotificationTypeValues() []NotificationType {
	return []NotificationType{
		NotificationTypeUserRegistered,
		NotificationTypeEmailVerified,
		NotificationTypeOrderConfirmed,
		NotificationTypeOrderShipped,
		NotificationTypeOrderDelivered,
		NotificationTypeOrderCancelled,
		NotificationTypePaymentCompleted,
		NotificationTypePaymentFailed,
		NotificationTypeRefundCompleted,
		NotificationTypeShopApproved,
		NotificationTypeShopRejected,
		NotificationTypeStockLowAlert,
		NotificationTypeStockOutAlert,
		NotificationTypeStockRestored,
		NotificationTypeChatMessageReceived,
		NotificationTypePasswordReset,
		NotificationTypeCampaignStarted,
	}
}

// IsValid checks if the value is valid
func (e NotificationType) IsValid() bool {
	switch e {
	case NotificationTypeUserRegistered:
	case NotificationTypeEmailVerified:
	case NotificationTypeOrderConfirmed:
	case NotificationTypeOrderShipped:
	case NotificationTypeOrderDelivered:
	case NotificationTypeOrderCancelled:
	case NotificationTypePaymentCompleted:
	case NotificationTypePaymentFailed:
	case NotificationTypeRefundCompleted:
	case NotificationTypeShopApproved:
	case NotificationTypeShopRejected:
	case NotificationTypeStockLowAlert:
	case NotificationTypeStockOutAlert:
	case NotificationTypeStockRestored:
	case NotificationTypeChatMessageReceived:
	case NotificationTypePasswordReset:
	case NotificationTypeCampaignStarted:
		return true
	}
	return false
}
