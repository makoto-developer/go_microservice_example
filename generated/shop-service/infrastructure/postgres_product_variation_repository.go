package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type PostgresProductVariationRepository struct {
	db *sql.DB
}

func NewPostgresProductVariationRepository(db *sql.DB) *PostgresProductVariationRepository {
	return &PostgresProductVariationRepository{db: db}
}

func (r *PostgresProductVariationRepository) Create(ctx context.Context, variation *domain.ProductVariation) error {
	query := `
		INSERT INTO product_variations (
			id, product_id, sku, attribute_name, attribute_value,
			price, stock_quantity, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		variation.Id,
		variation.ProductId,
		variation.Sku,
		variation.AttributeName,
		variation.AttributeValue,
		variation.Price,
		variation.StockQuantity,
		variation.CreatedAt,
		variation.UpdatedAt,
	)

	return err
}

func (r *PostgresProductVariationRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductVariation, error) {
	query := `
		SELECT id, product_id, sku, attribute_name, attribute_value,
			price, stock_quantity, created_at, updated_at
		FROM product_variations
		WHERE id = $1
	`

	var variation domain.ProductVariation

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&variation.Id,
		&variation.ProductId,
		&variation.Sku,
		&variation.AttributeName,
		&variation.AttributeValue,
		&variation.Price,
		&variation.StockQuantity,
		&variation.CreatedAt,
		&variation.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product variation not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	return &variation, nil
}

func (r *PostgresProductVariationRepository) FindBySku(ctx context.Context, sku string) (*domain.ProductVariation, error) {
	query := `
		SELECT id, product_id, sku, attribute_name, attribute_value,
			price, stock_quantity, created_at, updated_at
		FROM product_variations
		WHERE sku = $1
	`

	var variation domain.ProductVariation

	err := r.db.QueryRowContext(ctx, query, sku).Scan(
		&variation.Id,
		&variation.ProductId,
		&variation.Sku,
		&variation.AttributeName,
		&variation.AttributeValue,
		&variation.Price,
		&variation.StockQuantity,
		&variation.CreatedAt,
		&variation.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product variation not found with sku: %s", sku)
	}
	if err != nil {
		return nil, err
	}

	return &variation, nil
}

func (r *PostgresProductVariationRepository) Update(ctx context.Context, variation *domain.ProductVariation) error {
	query := `
		UPDATE product_variations
		SET product_id = $1,
			sku = $2,
			attribute_name = $3,
			attribute_value = $4,
			price = $5,
			stock_quantity = $6,
			updated_at = $7
		WHERE id = $8
	`

	_, err := r.db.ExecContext(ctx, query,
		variation.ProductId,
		variation.Sku,
		variation.AttributeName,
		variation.AttributeValue,
		variation.Price,
		variation.StockQuantity,
		time.Now(),
		variation.Id,
	)

	return err
}

func (r *PostgresProductVariationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM product_variations WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresProductVariationRepository) List(ctx context.Context, limit, offset int) ([]*domain.ProductVariation, error) {
	query := `
		SELECT id, product_id, sku, attribute_name, attribute_value,
			price, stock_quantity, created_at, updated_at
		FROM product_variations
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variations []*domain.ProductVariation
	for rows.Next() {
		var variation domain.ProductVariation

		err := rows.Scan(
			&variation.Id,
			&variation.ProductId,
			&variation.Sku,
			&variation.AttributeName,
			&variation.AttributeValue,
			&variation.Price,
			&variation.StockQuantity,
			&variation.CreatedAt,
			&variation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		variations = append(variations, &variation)
	}

	return variations, rows.Err()
}

// BeginTx starts a transaction
func (r *PostgresProductVariationRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
