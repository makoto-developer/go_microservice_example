package notification_flow

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/makoto-developer/go_microservice_example/tests/integration/notification_flow/clients"
)

// TestCreateShipment tests shipment creation flow
func TestCreateShipment(t *testing.T) {
	// Connect to shipping database
	shippingClient, err := clients.NewShippingClient(
		"postgresql://postgres:postgres_password@localhost:22016/shipping_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer shippingClient.Close()

	// Test data
	orderID := uuid.New()
	recipientName := "John Doe"
	recipientPhone := "+81-90-1234-5678"
	shippingAddress := "123 Main St, Tokyo, Japan"

	// 1. Create shipment
	shipment, err := shippingClient.CreateShipment(orderID, recipientName, recipientPhone, shippingAddress)
	require.NoError(t, err)
	assert.NotNil(t, shipment)
	assert.Equal(t, orderID, shipment.OrderID)
	assert.Equal(t, "preparing", shipment.Status)
	assert.Equal(t, recipientName, shipment.RecipientName)
	assert.Equal(t, recipientPhone, shipment.RecipientPhone)
	assert.Equal(t, shippingAddress, shipment.ShippingAddress)

	// 2. Verify shipment can be retrieved
	retrievedShipment, err := shippingClient.GetShipment(shipment.ID)
	require.NoError(t, err)
	assert.Equal(t, shipment.ID, retrievedShipment.ID)
	assert.Equal(t, "preparing", retrievedShipment.Status)

	// 3. Verify shipment can be retrieved by order ID
	shipmentByOrder, err := shippingClient.GetShipmentByOrderID(orderID)
	require.NoError(t, err)
	assert.Equal(t, shipment.ID, shipmentByOrder.ID)
}

// TestShipmentStatusFlow tests complete shipment status lifecycle
func TestShipmentStatusFlow(t *testing.T) {
	shippingClient, err := clients.NewShippingClient(
		"postgresql://postgres:postgres_password@localhost:22016/shipping_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer shippingClient.Close()

	// 1. Create shipment (status: preparing)
	orderID := uuid.New()
	shipment, err := shippingClient.CreateShipment(
		orderID,
		"Jane Smith",
		"+81-80-9876-5432",
		"456 Oak Ave, Osaka, Japan",
	)
	require.NoError(t, err)
	assert.Equal(t, "preparing", shipment.Status)

	// 2. Update to shipped status
	trackingNumber := "TRK" + uuid.New().String()[:8]
	carrier := "YamatoTransport"
	err = shippingClient.UpdateShipmentStatus(shipment.ID, "shipped", &trackingNumber, &carrier)
	require.NoError(t, err)

	// Verify shipped status
	shipment, err = shippingClient.GetShipment(shipment.ID)
	require.NoError(t, err)
	assert.Equal(t, "shipped", shipment.Status)
	assert.NotNil(t, shipment.TrackingNumber)
	assert.Equal(t, trackingNumber, *shipment.TrackingNumber)
	assert.NotNil(t, shipment.Carrier)
	assert.Equal(t, carrier, *shipment.Carrier)
	assert.NotNil(t, shipment.ShippedAt)

	// 3. Update to in_transit status
	err = shippingClient.UpdateShipmentStatus(shipment.ID, "in_transit", nil, nil)
	require.NoError(t, err)

	shipment, err = shippingClient.GetShipment(shipment.ID)
	require.NoError(t, err)
	assert.Equal(t, "in_transit", shipment.Status)

	// 4. Update to delivered status
	err = shippingClient.UpdateShipmentStatus(shipment.ID, "delivered", nil, nil)
	require.NoError(t, err)

	shipment, err = shippingClient.GetShipment(shipment.ID)
	require.NoError(t, err)
	assert.Equal(t, "delivered", shipment.Status)
	assert.NotNil(t, shipment.DeliveredAt)
}

// TestTrackingEvents tests tracking event creation and retrieval
func TestTrackingEvents(t *testing.T) {
	shippingClient, err := clients.NewShippingClient(
		"postgresql://postgres:postgres_password@localhost:22016/shipping_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer shippingClient.Close()

	// 1. Create shipment
	shipment, err := shippingClient.CreateShipment(
		uuid.New(),
		"Test User",
		"+81-90-0000-0000",
		"Test Address",
	)
	require.NoError(t, err)

	// 2. Add tracking events
	location1 := "Tokyo Distribution Center"
	description1 := "Package received at distribution center"
	event1, err := shippingClient.AddTrackingEvent(shipment.ID, "preparing", &location1, &description1)
	require.NoError(t, err)
	assert.Equal(t, shipment.ID, event1.ShipmentID)
	assert.Equal(t, "preparing", event1.Status)

	time.Sleep(100 * time.Millisecond) // Ensure different timestamps

	location2 := "In Transit"
	description2 := "Package is on the way"
	event2, err := shippingClient.AddTrackingEvent(shipment.ID, "in_transit", &location2, &description2)
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	location3 := "Osaka Delivery Hub"
	description3 := "Package arrived at delivery hub"
	event3, err := shippingClient.AddTrackingEvent(shipment.ID, "out_for_delivery", &location3, &description3)
	require.NoError(t, err)

	// 3. Get all tracking events
	events, err := shippingClient.GetTrackingEvents(shipment.ID)
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(events), 3)

	// 4. Verify events are ordered by event_time ASC
	assert.True(t, events[0].EventTime.Before(events[1].EventTime) ||
		events[0].EventTime.Equal(events[1].EventTime))
	assert.True(t, events[1].EventTime.Before(events[2].EventTime) ||
		events[1].EventTime.Equal(events[2].EventTime))

	// 5. Verify event details
	assert.Equal(t, "preparing", events[0].Status)
	assert.Equal(t, location1, *events[0].Location)
	assert.Equal(t, description1, *events[0].Description)
}

// TestShippingWithNotification tests integration between shipping and notification
func TestShippingWithNotification(t *testing.T) {
	// Connect to both services
	shippingClient, err := clients.NewShippingClient(
		"postgresql://postgres:postgres_password@localhost:22016/shipping_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer shippingClient.Close()

	notificationClient, err := clients.NewNotificationClient(
		"postgresql://postgres:postgres_password@localhost:22017/notification_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer notificationClient.Close()

	// 1. Create shipment
	orderID := uuid.New()
	shipment, err := shippingClient.CreateShipment(
		orderID,
		"Alice Johnson",
		"+81-70-1111-2222",
		"789 Pine Rd, Kyoto, Japan",
	)
	require.NoError(t, err)

	// 2. Update shipment to shipped
	trackingNumber := "TRK123456789"
	carrier := "JapanPost"
	err = shippingClient.UpdateShipmentStatus(shipment.ID, "shipped", &trackingNumber, &carrier)
	require.NoError(t, err)

	// 3. Create shipping update notification
	template, err := notificationClient.GetTemplate("shipping_update")
	require.NoError(t, err)

	recipient := "alice@example.com"
	subject := "Shipping Update - Order #" + orderID.String()[:8]
	body := "Your order has been shipped.\n\nTracking Number: " + trackingNumber

	notification, err := notificationClient.CreateNotification(
		template.ID,
		recipient,
		template.Channel,
		subject,
		body,
	)
	require.NoError(t, err)
	assert.Equal(t, "pending", notification.Status)

	// 4. Mark notification as sent
	sentTime := time.Now()
	err = notificationClient.UpdateNotificationStatus(notification.ID, "sent", &sentTime, nil)
	require.NoError(t, err)

	// 5. Verify notification was sent
	updatedNotification, err := notificationClient.GetNotification(notification.ID)
	require.NoError(t, err)
	assert.Equal(t, "sent", updatedNotification.Status)
}

// TestMultipleTrackingEventsChronology tests tracking events are properly chronological
func TestMultipleTrackingEventsChronology(t *testing.T) {
	shippingClient, err := clients.NewShippingClient(
		"postgresql://postgres:postgres_password@localhost:22016/shipping_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer shippingClient.Close()

	// Create shipment
	shipment, err := shippingClient.CreateShipment(
		uuid.New(),
		"Test User",
		"+81-90-0000-0000",
		"Test Address",
	)
	require.NoError(t, err)

	// Add events with clear time progression
	statuses := []string{"preparing", "shipped", "in_transit", "out_for_delivery", "delivered"}
	for _, status := range statuses {
		location := "Location for " + status
		description := "Description for " + status
		_, err := shippingClient.AddTrackingEvent(shipment.ID, status, &location, &description)
		require.NoError(t, err)
		time.Sleep(50 * time.Millisecond) // Ensure different timestamps
	}

	// Get events and verify chronology
	events, err := shippingClient.GetTrackingEvents(shipment.ID)
	require.NoError(t, err)
	assert.Equal(t, len(statuses), len(events))

	// Verify order
	for i := 0; i < len(events)-1; i++ {
		assert.True(t, events[i].EventTime.Before(events[i+1].EventTime))
	}

	// Verify status progression
	for i, status := range statuses {
		assert.Equal(t, status, events[i].Status)
	}
}

// TestShipmentUniqueOrderID tests that order_id is unique (one shipment per order)
func TestShipmentUniqueOrderID(t *testing.T) {
	shippingClient, err := clients.NewShippingClient(
		"postgresql://postgres:postgres_password@localhost:22016/shipping_service?sslmode=disable",
	)
	require.NoError(t, err)
	defer shippingClient.Close()

	orderID := uuid.New()

	// Create first shipment
	_, err = shippingClient.CreateShipment(orderID, "User 1", "+81-90-1111-1111", "Address 1")
	require.NoError(t, err)

	// Try to create duplicate shipment for same order
	_, err = shippingClient.CreateShipment(orderID, "User 2", "+81-90-2222-2222", "Address 2")
	assert.Error(t, err) // Should fail due to unique constraint on order_id
}
