package domain

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID            uuid.UUID `db:"id" json:"id"`
	ShopID        uuid.UUID `db:"shop_id" json:"shop_id"`
	Name          string    `db:"name" json:"name"`
	Description   string    `db:"description" json:"description"`
	Price         float64   `db:"price" json:"price"`
	Category      string    `db:"category" json:"category"`
	StockQuantity int       `db:"stock_quantity" json:"stock_quantity"`
	Weight        *float64  `db:"weight" json:"weight,omitempty"`
	Size          *string   `db:"size" json:"size,omitempty"`
	JANCode       *string   `db:"jan_code" json:"jan_code,omitempty"`
	Published     bool      `db:"published" json:"published"`
	Deleted       bool      `db:"deleted" json:"deleted"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at"`
}

func NewProduct(shopID uuid.UUID, name, description string, price float64, category string, stockQuantity int) *Product {
	now := time.Now()
	return &Product{
		ID:            uuid.New(),
		ShopID:        shopID,
		Name:          name,
		Description:   description,
		Price:         price,
		Category:      category,
		StockQuantity: stockQuantity,
		Published:     false,
		Deleted:       false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

func (p *Product) IsAvailable() bool {
	return p.Published && !p.Deleted && p.StockQuantity > 0
}

func (p *Product) DecrementStock(quantity int) error {
	if p.StockQuantity < quantity {
		return ErrInsufficientStock
	}
	p.StockQuantity -= quantity
	p.UpdatedAt = time.Now()
	return nil
}
