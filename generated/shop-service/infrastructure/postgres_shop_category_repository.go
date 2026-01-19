package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type PostgresShopCategoryRepository struct {
	db *sql.DB
}

func NewPostgresShopCategoryRepository(db *sql.DB) *PostgresShopCategoryRepository {
	return &PostgresShopCategoryRepository{db: db}
}

func (r *PostgresShopCategoryRepository) Create(ctx context.Context, shopCategory *domain.ShopCategory) error {
	query := `
		INSERT INTO shop_categories (
			id, shop_id, category_name, created_at
		) VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query,
		shopCategory.Id,
		shopCategory.ShopId,
		shopCategory.CategoryName,
		shopCategory.CreatedAt,
	)

	return err
}

func (r *PostgresShopCategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ShopCategory, error) {
	query := `
		SELECT id, shop_id, category_name, created_at
		FROM shop_categories
		WHERE id = $1
	`

	var shopCategory domain.ShopCategory

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&shopCategory.Id,
		&shopCategory.ShopId,
		&shopCategory.CategoryName,
		&shopCategory.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("shop category not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	return &shopCategory, nil
}

func (r *PostgresShopCategoryRepository) Update(ctx context.Context, shopCategory *domain.ShopCategory) error {
	query := `
		UPDATE shop_categories
		SET shop_id = $1,
			category_name = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query,
		shopCategory.ShopId,
		shopCategory.CategoryName,
		shopCategory.Id,
	)

	return err
}

func (r *PostgresShopCategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM shop_categories WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresShopCategoryRepository) List(ctx context.Context, limit, offset int) ([]*domain.ShopCategory, error) {
	query := `
		SELECT id, shop_id, category_name, created_at
		FROM shop_categories
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.ShopCategory
	for rows.Next() {
		var category domain.ShopCategory

		err := rows.Scan(
			&category.Id,
			&category.ShopId,
			&category.CategoryName,
			&category.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}

	return categories, rows.Err()
}

// BeginTx starts a transaction
func (r *PostgresShopCategoryRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
