package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
	"github.com/shopspring/decimal"
)

type PostgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (
			id, shop_id, name, description, price, category,
			stock_quantity, weight, size, jan_code,
			published, deleted, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`

	_, err := r.db.ExecContext(ctx, query,
		product.Id,
		product.ShopId,
		product.Name,
		product.Description,
		product.Price,
		product.Category,
		product.StockQuantity,
		product.Weight,
		product.Size,
		product.JanCode,
		product.Published,
		product.Deleted,
		product.CreatedAt,
		product.UpdatedAt,
	)

	return err
}

func (r *PostgresProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `
		SELECT id, shop_id, name, description, price, category,
			stock_quantity, weight, size, jan_code,
			published, deleted, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product domain.Product
	var weight sql.NullString
	var size sql.NullString
	var janCode sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.Id,
		&product.ShopId,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.Category,
		&product.StockQuantity,
		&weight,
		&size,
		&janCode,
		&product.Published,
		&product.Deleted,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	if weight.Valid {
		w, _ := decimal.NewFromString(weight.String)
		product.Weight = &w
	}
	if size.Valid {
		product.Size = &size.String
	}
	if janCode.Valid {
		product.JanCode = &janCode.String
	}

	return &product, nil
}

func (r *PostgresProductRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products
		SET shop_id = $1,
			name = $2,
			description = $3,
			price = $4,
			category = $5,
			stock_quantity = $6,
			weight = $7,
			size = $8,
			jan_code = $9,
			published = $10,
			deleted = $11,
			updated_at = $12
		WHERE id = $13
	`

	_, err := r.db.ExecContext(ctx, query,
		product.ShopId,
		product.Name,
		product.Description,
		product.Price,
		product.Category,
		product.StockQuantity,
		product.Weight,
		product.Size,
		product.JanCode,
		product.Published,
		product.Deleted,
		time.Now(),
		product.Id,
	)

	return err
}

func (r *PostgresProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresProductRepository) List(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	query := `
		SELECT id, shop_id, name, description, price, category,
			stock_quantity, weight, size, jan_code,
			published, deleted, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		var weight sql.NullString
		var size sql.NullString
		var janCode sql.NullString

		err := rows.Scan(
			&product.Id,
			&product.ShopId,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Category,
			&product.StockQuantity,
			&weight,
			&size,
			&janCode,
			&product.Published,
			&product.Deleted,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if weight.Valid {
			w, _ := decimal.NewFromString(weight.String)
			product.Weight = &w
		}
		if size.Valid {
			product.Size = &size.String
		}
		if janCode.Valid {
			product.JanCode = &janCode.String
		}

		products = append(products, &product)
	}

	return products, rows.Err()
}

func (r *PostgresProductRepository) ListByShopID(ctx context.Context, shopID uuid.UUID, category string, publishedOnly bool, limit, offset int) ([]*domain.Product, error) {
	query := `
		SELECT id, shop_id, name, description, price, category,
			stock_quantity, weight, size, jan_code,
			published, deleted, created_at, updated_at
		FROM products
		WHERE deleted = false
	`

	args := []interface{}{}
	argIndex := 1

	// Shop ID filter (empty UUID means all shops)
	if shopID != uuid.Nil {
		query += fmt.Sprintf(" AND shop_id = $%d", argIndex)
		args = append(args, shopID)
		argIndex++
	}

	// Category filter
	if category != "" {
		query += fmt.Sprintf(" AND category = $%d", argIndex)
		args = append(args, category)
		argIndex++
	}

	// Published filter
	if publishedOnly {
		query += " AND published = true"
	}

	query += fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		var weight sql.NullString
		var size sql.NullString
		var janCode sql.NullString

		err := rows.Scan(
			&product.Id,
			&product.ShopId,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Category,
			&product.StockQuantity,
			&weight,
			&size,
			&janCode,
			&product.Published,
			&product.Deleted,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if weight.Valid {
			w, _ := decimal.NewFromString(weight.String)
			product.Weight = &w
		}
		if size.Valid {
			product.Size = &size.String
		}
		if janCode.Valid {
			product.JanCode = &janCode.String
		}

		products = append(products, &product)
	}

	return products, rows.Err()
}

// BeginTx starts a transaction
func (r *PostgresProductRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
