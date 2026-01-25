package domain

import (
	"time"

	"github.com/google/uuid"
)

// ProductStatus represents the status of a product
type ProductStatus string

const (
	ProductStatusDraft     ProductStatus = "DRAFT"
	ProductStatusPublished ProductStatus = "PUBLISHED"
	ProductStatusArchived  ProductStatus = "ARCHIVED"
)

// Product represents a product entity
type Product struct {
	ID          uuid.UUID     `db:"id" json:"id"`
	ShopID      uuid.UUID     `db:"shop_id" json:"shop_id"`
	Name          string        `db:"name" json:"name"`
	Description   string        `db:"description" json:"description"`
	Price         string        `db:"price" json:"price"`
	Category      string        `db:"category" json:"category"`
	Weight        *float64      `db:"weight" json:"weight,omitempty"`
	Size          *string       `db:"size" json:"size,omitempty"`
	JANCode       *string       `db:"jan_code" json:"jan_code,omitempty"`
	StockQuantity int           `db:"stock_quantity" json:"stock_quantity"`
	Status        ProductStatus `db:"status" json:"status"`
	Published     bool          `db:"published" json:"published"`
	Deleted       bool          `db:"deleted" json:"deleted"`
	CreatedAt   time.Time     `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at" json:"updated_at"`
	DeletedAt   *time.Time    `db:"deleted_at" json:"deleted_at,omitempty"`
}

// ProductCategory represents a product category
type ProductCategory struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Slug      string    `db:"slug" json:"slug"`
	ParentID  *uuid.UUID `db:"parent_id" json:"parent_id,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
