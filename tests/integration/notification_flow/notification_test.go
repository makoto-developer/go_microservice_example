package notification_flow

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/makoto-developer/go_microservice_example/tests/integration/notification_flow/clients"
)

// TestOrderConfirmationNotification tests order confirmation notification flow
func TestOrderConfirmationNotification(t *testing.T) {
	// Connect to notification database
	notificationClient, err := clients.NewNotificationClient(
		"postgresql://postgres:postgres_password@localhost:22017/notification_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer notificationClient.Close()

	// 1. Get order_confirmation template
	template, err := notificationClient.GetTemplate("order_confirmation")
	require.NoError(t, err)
	assert.Equal(t, "order_confirmation", template.Name)
	assert.Equal(t, "email", template.Channel)

	// 2. Create notification for order confirmation
	recipient := "customer@example.com"
	subject := "Order Confirmation - Order #12345"
	body := "Dear John Doe,\n\nYour order #12345 has been confirmed.\n\nThank you!"

	notification, err := notificationClient.CreateNotification(
		template.ID,
		recipient,
		template.Channel,
		subject,
		body,
	)
	require.NoError(t, err)
	assert.Equal(t, "pending", notification.Status)
	assert.Equal(t, recipient, notification.Recipient)

	// 3. Verify notification is created
	retrievedNotification, err := notificationClient.GetNotification(notification.ID)
	require.NoError(t, err)
	assert.Equal(t, notification.ID, retrievedNotification.ID)
	assert.Equal(t, "pending", retrievedNotification.Status)

	// 4. Update notification status to sent
	sentTime := time.Now()
	err = notificationClient.UpdateNotificationStatus(notification.ID, "sent", &sentTime, nil)
	require.NoError(t, err)

	// 5. Verify notification status is updated
	updatedNotification, err := notificationClient.GetNotification(notification.ID)
	require.NoError(t, err)
	assert.Equal(t, "sent", updatedNotification.Status)
	assert.NotNil(t, updatedNotification.SentAt)

	// Cleanup
	// Note: In production, implement cleanup logic
}

// TestPaymentSuccessNotification tests payment success notification flow
func TestPaymentSuccessNotification(t *testing.T) {
	notificationClient, err := clients.NewNotificationClient(
		"postgresql://postgres:postgres_password@localhost:22017/notification_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer notificationClient.Close()

	// 1. Get payment_success template
	template, err := notificationClient.GetTemplate("payment_success")
	require.NoError(t, err)
	assert.Equal(t, "payment_success", template.Name)
	assert.Equal(t, "email", template.Channel)

	// 2. Create payment success notification
	recipient := "customer@example.com"
	subject := "Payment Successful - Order #12345"
	body := "Dear John Doe,\n\nYour payment for order #12345 has been successfully processed.\n\nAmount Paid: $150.00"

	notification, err := notificationClient.CreateNotification(
		template.ID,
		recipient,
		template.Channel,
		subject,
		body,
	)
	require.NoError(t, err)
	assert.Equal(t, "pending", notification.Status)

	// 3. Simulate sending and update status
	sentTime := time.Now()
	err = notificationClient.UpdateNotificationStatus(notification.ID, "sent", &sentTime, nil)
	require.NoError(t, err)

	// 4. Verify
	updatedNotification, err := notificationClient.GetNotification(notification.ID)
	require.NoError(t, err)
	assert.Equal(t, "sent", updatedNotification.Status)
}

// TestShippingUpdateNotification tests shipping update notification flow
func TestShippingUpdateNotification(t *testing.T) {
	notificationClient, err := clients.NewNotificationClient(
		"postgresql://postgres:postgres_password@localhost:22017/notification_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer notificationClient.Close()

	// 1. Get shipping_update template
	template, err := notificationClient.GetTemplate("shipping_update")
	require.NoError(t, err)
	assert.Equal(t, "shipping_update", template.Name)

	// 2. Create shipping update notification
	recipient := "customer@example.com"
	subject := "Shipping Update - Order #12345"
	body := "Your order #12345 has been shipped.\n\nTracking Number: TRK123456789\nEstimated Delivery: 2025-02-01"

	notification, err := notificationClient.CreateNotification(
		template.ID,
		recipient,
		template.Channel,
		subject,
		body,
	)
	require.NoError(t, err)

	// 3. Update to sent
	sentTime := time.Now()
	err = notificationClient.UpdateNotificationStatus(notification.ID, "sent", &sentTime, nil)
	require.NoError(t, err)

	// 4. Verify
	updatedNotification, err := notificationClient.GetNotification(notification.ID)
	require.NoError(t, err)
	assert.Equal(t, "sent", updatedNotification.Status)
	assert.NotNil(t, updatedNotification.SentAt)
}

// TestNotificationFailureHandling tests notification failure scenario
func TestNotificationFailureHandling(t *testing.T) {
	notificationClient, err := clients.NewNotificationClient(
		"postgresql://postgres:postgres_password@localhost:22017/notification_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer notificationClient.Close()

	// 1. Get template
	template, err := notificationClient.GetTemplate("order_confirmation")
	require.NoError(t, err)

	// 2. Create notification
	notification, err := notificationClient.CreateNotification(
		template.ID,
		"invalid-email",
		template.Channel,
		"Test Subject",
		"Test Body",
	)
	require.NoError(t, err)

	// 3. Simulate failure
	errorMsg := "Invalid email address"
	err = notificationClient.UpdateNotificationStatus(notification.ID, "failed", nil, &errorMsg)
	require.NoError(t, err)

	// 4. Verify failure status
	updatedNotification, err := notificationClient.GetNotification(notification.ID)
	require.NoError(t, err)
	assert.Equal(t, "failed", updatedNotification.Status)
	assert.NotNil(t, updatedNotification.ErrorMessage)
	assert.Equal(t, errorMsg, *updatedNotification.ErrorMessage)
}

// TestMultipleNotificationsForRecipient tests fetching multiple notifications for a recipient
func TestMultipleNotificationsForRecipient(t *testing.T) {
	notificationClient, err := clients.NewNotificationClient(
		"postgresql://postgres:postgres_password@localhost:22017/notification_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer notificationClient.Close()

	recipient := "multi-test@example.com"

	// Create multiple notifications
	template, err := notificationClient.GetTemplate("order_confirmation")
	require.NoError(t, err)

	for i := 0; i < 3; i++ {
		_, err := notificationClient.CreateNotification(
			template.ID,
			recipient,
			template.Channel,
			"Test Subject",
			"Test Body",
		)
		require.NoError(t, err)
		time.Sleep(100 * time.Millisecond) // Ensure different timestamps
	}

	// Get all notifications for recipient
	notifications, err := notificationClient.GetNotificationsByRecipient(recipient)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(notifications), 3)

	// Verify they are ordered by created_at DESC
	if len(notifications) >= 2 {
		assert.True(t, notifications[0].CreatedAt.After(notifications[1].CreatedAt) ||
			notifications[0].CreatedAt.Equal(notifications[1].CreatedAt))
	}
}
