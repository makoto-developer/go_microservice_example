package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// RabbitMQEventPublisher RabbitMQイベント発行
type RabbitMQEventPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewRabbitMQEventPublisher コンストラクタ
func NewRabbitMQEventPublisher(amqpURL string) (*RabbitMQEventPublisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Exchange宣言
	err = ch.ExchangeDeclare(
		"inventory_events", // name
		"topic",            // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQEventPublisher{
		conn:    conn,
		channel: ch,
	}, nil
}

// Close 接続クローズ
func (p *RabbitMQEventPublisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

// PublishStockReserved 在庫引き当てイベント発行
func (p *RabbitMQEventPublisher) PublishStockReserved(ctx context.Context, event StockReservedEvent) error {
	return p.publish(ctx, "inventory.reserved", event)
}

// PublishStockReleased 在庫リリースイベント発行
func (p *RabbitMQEventPublisher) PublishStockReleased(ctx context.Context, event StockReleasedEvent) error {
	return p.publish(ctx, "inventory.released", event)
}

// PublishStockConfirmed 在庫確定イベント発行
func (p *RabbitMQEventPublisher) PublishStockConfirmed(ctx context.Context, event StockConfirmedEvent) error {
	return p.publish(ctx, "inventory.confirmed", event)
}

// PublishLowStockAlert 低在庫アラートイベント発行
func (p *RabbitMQEventPublisher) PublishLowStockAlert(ctx context.Context, event LowStockAlertEvent) error {
	return p.publish(ctx, "inventory.low", event)
}

// PublishOutOfStockAlert 在庫切れアラートイベント発行
func (p *RabbitMQEventPublisher) PublishOutOfStockAlert(ctx context.Context, event OutOfStockAlertEvent) error {
	return p.publish(ctx, "inventory.out_of_stock", event)
}

// PublishStockRestored 在庫復活イベント発行
func (p *RabbitMQEventPublisher) PublishStockRestored(ctx context.Context, event StockRestoredEvent) error {
	return p.publish(ctx, "inventory.restored", event)
}

// PublishInventoryUpdated 在庫更新イベント発行
func (p *RabbitMQEventPublisher) PublishInventoryUpdated(ctx context.Context, event InventoryUpdatedEvent) error {
	return p.publish(ctx, "inventory.updated", event)
}

// publish 共通のイベント発行処理
func (p *RabbitMQEventPublisher) publish(ctx context.Context, routingKey string, event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.channel.Publish(
		"inventory_events", // exchange
		routingKey,         // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
			Timestamp:   time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	return nil
}

// Event Structures

type StockReservedEvent struct {
	ReservationID uuid.UUID `json:"reservation_id"`
	InventoryID   uuid.UUID `json:"inventory_id"`
	ProductID     uuid.UUID `json:"product_id"`
	OrderID       uuid.UUID `json:"order_id"`
	Quantity      int       `json:"quantity"`
	ExpiresAt     time.Time `json:"expires_at"`
	CreatedAt     time.Time `json:"created_at"`
}

type StockReleasedEvent struct {
	ReservationID uuid.UUID `json:"reservation_id"`
	InventoryID   uuid.UUID `json:"inventory_id"`
	ProductID     uuid.UUID `json:"product_id"`
	OrderID       uuid.UUID `json:"order_id"`
	Quantity      int       `json:"quantity"`
	ReleasedAt    time.Time `json:"released_at"`
}

type StockConfirmedEvent struct {
	ReservationID uuid.UUID `json:"reservation_id"`
	InventoryID   uuid.UUID `json:"inventory_id"`
	ProductID     uuid.UUID `json:"product_id"`
	OrderID       uuid.UUID `json:"order_id"`
	Quantity      int       `json:"quantity"`
	ConfirmedAt   time.Time `json:"confirmed_at"`
}

type LowStockAlertEvent struct {
	InventoryID     uuid.UUID `json:"inventory_id"`
	ProductID       uuid.UUID `json:"product_id"`
	ShopID          uuid.UUID `json:"shop_id"`
	CurrentQuantity int       `json:"current_quantity"`
	Threshold       int       `json:"threshold"`
	AlertedAt       time.Time `json:"alerted_at"`
}

type OutOfStockAlertEvent struct {
	InventoryID uuid.UUID `json:"inventory_id"`
	ProductID   uuid.UUID `json:"product_id"`
	ShopID      uuid.UUID `json:"shop_id"`
	AlertedAt   time.Time `json:"alerted_at"`
}

type StockRestoredEvent struct {
	InventoryID uuid.UUID `json:"inventory_id"`
	ProductID   uuid.UUID `json:"product_id"`
	ShopID      uuid.UUID `json:"shop_id"`
	NewQuantity int       `json:"new_quantity"`
	RestoredAt  time.Time `json:"restored_at"`
}

type InventoryUpdatedEvent struct {
	InventoryID uuid.UUID `json:"inventory_id"`
	ProductID   uuid.UUID `json:"product_id"`
	ShopID      uuid.UUID `json:"shop_id"`
	OldQuantity int       `json:"old_quantity"`
	NewQuantity int       `json:"new_quantity"`
	ChangeType  string    `json:"change_type"`
	UpdatedAt   time.Time `json:"updated_at"`
}
