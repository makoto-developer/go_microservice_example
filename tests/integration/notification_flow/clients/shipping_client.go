package clients

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type ShippingClient struct {
	db *sql.DB
}

type Shipment struct {
	ID              uuid.UUID
	OrderID         uuid.UUID
	TrackingNumber  *string
	Carrier         *string
	Status          string
	ShippedAt       *time.Time
	DeliveredAt     *time.Time
	RecipientName   string
	RecipientPhone  string
	ShippingAddress string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type TrackingEvent struct {
	ID          uuid.UUID
	ShipmentID  uuid.UUID
	Status      string
	Location    *string
	Description *string
	EventTime   time.Time
	CreatedAt   time.Time
}

func NewShippingClient(databaseURL string) (*ShippingClient, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to shipping database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping shipping database: %w", err)
	}

	return &ShippingClient{db: db}, nil
}

func (c *ShippingClient) Close() error {
	return c.db.Close()
}

func (c *ShippingClient) CreateShipment(orderID uuid.UUID, recipientName, recipientPhone, shippingAddress string) (*Shipment, error) {
	query := `
		INSERT INTO shipments (order_id, recipient_name, recipient_phone, shipping_address, status)
		VALUES ($1, $2, $3, $4, 'preparing')
		RETURNING id, order_id, tracking_number, carrier, status, shipped_at, delivered_at,
		          recipient_name, recipient_phone, shipping_address, created_at, updated_at
	`

	var shipment Shipment
	err := c.db.QueryRow(query, orderID, recipientName, recipientPhone, shippingAddress).Scan(
		&shipment.ID,
		&shipment.OrderID,
		&shipment.TrackingNumber,
		&shipment.Carrier,
		&shipment.Status,
		&shipment.ShippedAt,
		&shipment.DeliveredAt,
		&shipment.RecipientName,
		&shipment.RecipientPhone,
		&shipment.ShippingAddress,
		&shipment.CreatedAt,
		&shipment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create shipment: %w", err)
	}

	return &shipment, nil
}

func (c *ShippingClient) GetShipment(id uuid.UUID) (*Shipment, error) {
	query := `
		SELECT id, order_id, tracking_number, carrier, status, shipped_at, delivered_at,
		       recipient_name, recipient_phone, shipping_address, created_at, updated_at
		FROM shipments
		WHERE id = $1
	`

	var shipment Shipment
	err := c.db.QueryRow(query, id).Scan(
		&shipment.ID,
		&shipment.OrderID,
		&shipment.TrackingNumber,
		&shipment.Carrier,
		&shipment.Status,
		&shipment.ShippedAt,
		&shipment.DeliveredAt,
		&shipment.RecipientName,
		&shipment.RecipientPhone,
		&shipment.ShippingAddress,
		&shipment.CreatedAt,
		&shipment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}

	return &shipment, nil
}

func (c *ShippingClient) GetShipmentByOrderID(orderID uuid.UUID) (*Shipment, error) {
	query := `
		SELECT id, order_id, tracking_number, carrier, status, shipped_at, delivered_at,
		       recipient_name, recipient_phone, shipping_address, created_at, updated_at
		FROM shipments
		WHERE order_id = $1
	`

	var shipment Shipment
	err := c.db.QueryRow(query, orderID).Scan(
		&shipment.ID,
		&shipment.OrderID,
		&shipment.TrackingNumber,
		&shipment.Carrier,
		&shipment.Status,
		&shipment.ShippedAt,
		&shipment.DeliveredAt,
		&shipment.RecipientName,
		&shipment.RecipientPhone,
		&shipment.ShippingAddress,
		&shipment.CreatedAt,
		&shipment.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get shipment: %w", err)
	}

	return &shipment, nil
}

func (c *ShippingClient) UpdateShipmentStatus(id uuid.UUID, status string, trackingNumber, carrier *string) error {
	var shippedAt *time.Time
	var deliveredAt *time.Time

	if status == "shipped" {
		now := time.Now()
		shippedAt = &now
	} else if status == "delivered" {
		now := time.Now()
		deliveredAt = &now
	}

	query := `
		UPDATE shipments
		SET status = $1, tracking_number = COALESCE($2, tracking_number),
		    carrier = COALESCE($3, carrier), shipped_at = COALESCE($4, shipped_at),
		    delivered_at = COALESCE($5, delivered_at), updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
	`

	_, err := c.db.Exec(query, status, trackingNumber, carrier, shippedAt, deliveredAt, id)
	if err != nil {
		return fmt.Errorf("failed to update shipment status: %w", err)
	}

	return nil
}

func (c *ShippingClient) AddTrackingEvent(shipmentID uuid.UUID, status string, location, description *string) (*TrackingEvent, error) {
	query := `
		INSERT INTO tracking_events (shipment_id, status, location, description, event_time)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, shipment_id, status, location, description, event_time, created_at
	`

	var event TrackingEvent
	eventTime := time.Now()

	err := c.db.QueryRow(query, shipmentID, status, location, description, eventTime).Scan(
		&event.ID,
		&event.ShipmentID,
		&event.Status,
		&event.Location,
		&event.Description,
		&event.EventTime,
		&event.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to add tracking event: %w", err)
	}

	return &event, nil
}

func (c *ShippingClient) GetTrackingEvents(shipmentID uuid.UUID) ([]TrackingEvent, error) {
	query := `
		SELECT id, shipment_id, status, location, description, event_time, created_at
		FROM tracking_events
		WHERE shipment_id = $1
		ORDER BY event_time ASC
	`

	rows, err := c.db.Query(query, shipmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tracking events: %w", err)
	}
	defer rows.Close()

	var events []TrackingEvent
	for rows.Next() {
		var event TrackingEvent
		err := rows.Scan(
			&event.ID,
			&event.ShipmentID,
			&event.Status,
			&event.Location,
			&event.Description,
			&event.EventTime,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tracking event: %w", err)
		}
		events = append(events, event)
	}

	return events, nil
}
