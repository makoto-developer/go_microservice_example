package domain

import (
	"github.com/google/uuid"
	"time"
)

// NotificationPreference represents NotificationPreference
type NotificationPreference struct {
	Id uuid.UUID `db:"id" json:"id"`
	UserId uuid.UUID `db:"user_id" json:"user_id"`
	EmailEnabled bool `db:"email_enabled" json:"email_enabled"`
	PushEnabled bool `db:"push_enabled" json:"push_enabled"`
	EmailOrderUpdates bool `db:"email_order_updates" json:"email_order_updates"`
	EmailShopUpdates bool `db:"email_shop_updates" json:"email_shop_updates"`
	EmailChatMessages bool `db:"email_chat_messages" json:"email_chat_messages"`
	PushOrderUpdates bool `db:"push_order_updates" json:"push_order_updates"`
	PushStockRestored bool `db:"push_stock_restored" json:"push_stock_restored"`
	PushCampaigns bool `db:"push_campaigns" json:"push_campaigns"`
	PushChatMessages bool `db:"push_chat_messages" json:"push_chat_messages"`
	Frequency NotificationFrequency `db:"frequency" json:"frequency"`
	QuietHoursStart *time `db:"quiet_hours_start" json:"quiet_hours_start,omitempty"`
	QuietHoursEnd *time `db:"quiet_hours_end" json:"quiet_hours_end,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// NewNotificationPreference creates a new NotificationPreference instance
func NewNotificationPreference() *NotificationPreference {
	return &NotificationPreference{}
}
