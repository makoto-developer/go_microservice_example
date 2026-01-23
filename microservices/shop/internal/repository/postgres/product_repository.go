package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
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
			id, shop_id, name, description, price, category_id,
			tags, weight, dimensions, jan_code, stock_count,
			status, is_public, is_deleted, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)`

	_, err := r.db.ExecContext(ctx, query,
		product.ID, product.ShopID, product.Name, product.Description,
		product.Price, product.CategoryID, pq.Array(product.Tags),
		product.Weight, product.Dimensions, product.JANCode,
		product.StockCount, product.Status, product.IsPublic,
		product.IsDeleted, product.CreatedAt, product.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	return nil
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `
		SELECT id, shop_id, name, description, price, category_id,
			   tags, weight, dimensions, jan_code, stock_count,
			   status, is_public, is_deleted, created_at, updated_at, deleted_at
		FROM products
		WHERE id = $1`

	var product domain.Product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.ShopID, &product.Name, &product.Description,
		&product.Price, &product.CategoryID, pq.Array(&product.Tags),
		&product.Weight, &product.Dimensions, &product.JANCode,
		&product.StockCount, &product.Status, &product.IsPublic,
		&product.IsDeleted, &product.CreatedAt, &product.UpdatedAt,
		&product.DeletedAt,
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
		SELECT id, shop_id, name, description, price, category_id,
			   tags, weight, dimensions, jan_code, stock_count,
			   status, is_public, is_deleted, created_at, updated_at, deleted_at
		FROM products
		WHERE shop_id = $1`

	if !includeDeleted {
		query += " AND is_deleted = FALSE"
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
			&product.Price, &product.CategoryID, pq.Array(&product.Tags),
			&product.Weight, &product.Dimensions, &product.JANCode,
			&product.StockCount, &product.Status, &product.IsPublic,
			&product.IsDeleted, &product.CreatedAt, &product.UpdatedAt,
			&product.DeletedAt,
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
		SET name = $2, description = $3, price = $4, category_id = $5,
			tags = $6, weight = $7, dimensions = $8, jan_code = $9,
			stock_count = $10, updated_at = $11
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		product.ID, product.Name, product.Description, product.Price,
		product.CategoryID, pq.Array(product.Tags), product.Weight,
		product.Dimensions, product.JANCode, product.StockCount,
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
		SET is_deleted = TRUE, deleted_at = NOW(), updated_at = NOW()
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
	query := `UPDATE products SET is_public = $2, updated_at = NOW() WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, isPublic)
	if err != nil {
		return fmt.Errorf("failed to update product is_public: %w", err)
	}

	return nil
}

func (r *productRepository) UpdateStock(ctx context.Context, id uuid.UUID, stockCount int) error {
	query := `UPDATE products SET stock_count = $2, updated_at = NOW() WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, stockCount)
	if err != nil {
		return fmt.Errorf("failed to update product stock: %w", err)
	}

	return nil
}

func (r *productRepository) List(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	query := `
		SELECT id, shop_id, name, description, price, category_id,
			   tags, weight, dimensions, jan_code, stock_count,
			   status, is_public, is_deleted, created_at, updated_at, deleted_at
		FROM products
		WHERE is_deleted = FALSE
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
			&product.Price, &product.CategoryID, pq.Array(&product.Tags),
			&product.Weight, &product.Dimensions, &product.JANCode,
			&product.StockCount, &product.Status, &product.IsPublic,
			&product.IsDeleted, &product.CreatedAt, &product.UpdatedAt,
			&product.DeletedAt,
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
