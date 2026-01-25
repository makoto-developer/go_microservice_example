package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
)

type productRepository struct {
	db *sql.DB
}

// NewProductRepository creates a new product repository
func NewProductRepository(db *sql.DB) *productRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (
			id, shop_id, name, description, price, category,
			weight, size, jan_code, stock_quantity,
			published, deleted, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)`

	_, err := r.db.ExecContext(ctx, query,
		product.ID, product.ShopID, product.Name, product.Description,
		product.Price, product.Category,
		product.Weight, product.Size, product.JANCode,
		product.StockQuantity, product.Published,
		product.Deleted, product.CreatedAt, product.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `
		SELECT id, shop_id, name, description, price, category,
			   weight, size, jan_code, stock_quantity,
			   published, deleted, created_at, updated_at
		FROM products
		WHERE id = $1`

	var product domain.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.ShopID, &product.Name, &product.Description,
		&product.Price, &product.Category,
		&product.Weight, &product.Size, &product.JANCode,
		&product.StockQuantity, &product.Published,
		&product.Deleted, &product.CreatedAt, &product.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

func (r *productRepository) GetByShopID(ctx context.Context, shopID uuid.UUID, includeDeleted bool) ([]*domain.Product, error) {
	query := `
		SELECT id, shop_id, name, description, price, category,
			   weight, size, jan_code, stock_quantity,
			   published, deleted, created_at, updated_at
		FROM products
		WHERE shop_id = $1`

	if !includeDeleted {
		query += " AND deleted = FALSE"
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by shop: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID, &product.ShopID, &product.Name, &product.Description,
			&product.Price, &product.Category,
			&product.Weight, &product.Size, &product.JANCode,
			&product.StockQuantity, &product.Published,
			&product.Deleted, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products
		SET name = $2, description = $3, price = $4, category = $5,
			weight = $6, size = $7, jan_code = $8,
			stock_quantity = $9, updated_at = $10
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		product.ID, product.Name, product.Description, product.Price,
		product.Category, product.Weight,
		product.Size, product.JANCode, product.StockQuantity,
		product.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE products
		SET deleted = TRUE, updated_at = NOW()
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	return nil
}

func (r *productRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ProductStatus) error {
	query := `UPDATE products SET status = $2, updated_at = NOW() WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update product status: %w", err)
	}

	return nil
}

func (r *productRepository) UpdateIsPublic(ctx context.Context, id uuid.UUID, isPublic bool) error {
	query := `UPDATE products SET published = $2, updated_at = NOW() WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, isPublic)
	if err != nil {
		return fmt.Errorf("failed to update product published: %w", err)
	}

	return nil
}

func (r *productRepository) UpdateStock(ctx context.Context, id uuid.UUID, stockCount int) error {
	query := `UPDATE products SET stock_quantity = $2, updated_at = NOW() WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, stockCount)
	if err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	return nil
}

func (r *productRepository) List(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	query := `
		SELECT id, shop_id, name, description, price, category,
			   weight, size, jan_code, stock_quantity,
			   published, deleted, created_at, updated_at
		FROM products
		WHERE deleted = FALSE
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID, &product.ShopID, &product.Name, &product.Description,
			&product.Price, &product.Category,
			&product.Weight, &product.Size, &product.JANCode,
			&product.StockQuantity, &product.Published,
			&product.Deleted, &product.CreatedAt, &product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return products, nil
}
