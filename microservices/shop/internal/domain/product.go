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
	Name        string        `db:"name" json:"name"`
	Description string        `db:"description" json:"description"`
	Price       int64         `db:"price" json:"price"`
	CategoryID  uuid.UUID     `db:"category_id" json:"category_id"`
	Tags        []string      `db:"tags" json:"tags"`
	Weight      *float64      `db:"weight" json:"weight,omitempty"`
	Dimensions  *string       `db:"dimensions" json:"dimensions,omitempty"`
	JANCode     *string       `db:"jan_code" json:"jan_code,omitempty"`
	StockCount  int           `db:"stock_count" json:"stock_count"`
	Status      ProductStatus `db:"status" json:"status"`
	IsPublic    bool          `db:"is_public" json:"is_public"`
	IsDeleted   bool          `db:"is_deleted" json:"is_deleted"`
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
