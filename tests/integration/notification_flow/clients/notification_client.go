package clients

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type NotificationClient struct {
	db *sql.DB
}

type NotificationTemplate struct {
	ID           uuid.UUID
	Name         string
	Channel      string
	Subject      string
	BodyTemplate string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Notification struct {
	ID           uuid.UUID
	TemplateID   uuid.UUID
	Recipient    string
	Channel      string
	Subject      string
	Body         string
	Status       string
	SentAt       *time.Time
	ErrorMessage *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func NewNotificationClient(databaseURL string) (*NotificationClient, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to notification database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping notification database: %w", err)
	}

	return &NotificationClient{db: db}, nil
}

func (c *NotificationClient) Close() error {
	return c.db.Close()
}

func (c *NotificationClient) GetTemplate(name string) (*NotificationTemplate, error) {
	query := `
		SELECT id, name, channel, subject, body_template, created_at, updated_at
		FROM notification_templates
		WHERE name = $1
	`

	var template NotificationTemplate
	err := c.db.QueryRow(query, name).Scan(
		&template.ID,
		&template.Name,
		&template.Channel,
		&template.Subject,
		&template.BodyTemplate,
		&template.CreatedAt,
		&template.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	return &template, nil
}

func (c *NotificationClient) CreateNotification(templateID uuid.UUID, recipient, channel, subject, body string) (*Notification, error) {
	query := `
		INSERT INTO notifications (template_id, recipient, channel, subject, body, status)
		VALUES ($1, $2, $3, $4, $5, 'pending')
		RETURNING id, template_id, recipient, channel, subject, body, status, sent_at, error_message, created_at, updated_at
	`

	var notification Notification
	err := c.db.QueryRow(query, templateID, recipient, channel, subject, body).Scan(
		&notification.ID,
		&notification.TemplateID,
		&notification.Recipient,
		&notification.Channel,
		&notification.Subject,
		&notification.Body,
		&notification.Status,
		&notification.SentAt,
		&notification.ErrorMessage,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return &notification, nil
}

func (c *NotificationClient) GetNotification(id uuid.UUID) (*Notification, error) {
	query := `
		SELECT id, template_id, recipient, channel, subject, body, status, sent_at, error_message, created_at, updated_at
		FROM notifications
		WHERE id = $1
	`

	var notification Notification
	err := c.db.QueryRow(query, id).Scan(
		&notification.ID,
		&notification.TemplateID,
		&notification.Recipient,
		&notification.Channel,
		&notification.Subject,
		&notification.Body,
		&notification.Status,
		&notification.SentAt,
		&notification.ErrorMessage,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	return &notification, nil
}

func (c *NotificationClient) UpdateNotificationStatus(id uuid.UUID, status string, sentAt *time.Time, errorMsg *string) error {
	query := `
		UPDATE notifications
		SET status = $1, sent_at = $2, error_message = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4
	`

	_, err := c.db.Exec(query, status, sentAt, errorMsg, id)
	if err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}

	return nil
}

func (c *NotificationClient) GetNotificationsByRecipient(recipient string) ([]Notification, error) {
	query := `
		SELECT id, template_id, recipient, channel, subject, body, status, sent_at, error_message, created_at, updated_at
		FROM notifications
		WHERE recipient = $1
		ORDER BY created_at DESC
	`

	rows, err := c.db.Query(query, recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifications []Notification
	for rows.Next() {
		var notification Notification
		err := rows.Scan(
			&notification.ID,
			&notification.TemplateID,
			&notification.Recipient,
			&notification.Channel,
			&notification.Subject,
			&notification.Body,
			&notification.Status,
			&notification.SentAt,
			&notification.ErrorMessage,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}
