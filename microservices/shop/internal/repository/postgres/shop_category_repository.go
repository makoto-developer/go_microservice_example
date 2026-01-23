package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
)

type shopCategoryRepository struct {
	db *sql.DB
}

// NewShopCategoryRepository creates a new shop category repository
func NewShopCategoryRepository(db *sql.DB) *shopCategoryRepository {
	return &shopCategoryRepository{db: db}
}

func (r *shopCategoryRepository) AddCategory(ctx context.Context, shopID, categoryID uuid.UUID) error {
	query := `
		INSERT INTO shop_categories (id, shop_id, category_id, created_at)
		VALUES ($1, $2, $3, NOW())`

	_, err := r.db.ExecContext(ctx, query, uuid.New(), shopID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to add category: %w", err)
	}

	return nil
}

func (r *shopCategoryRepository) RemoveCategory(ctx context.Context, shopID, categoryID uuid.UUID) error {
	query := `DELETE FROM shop_categories WHERE shop_id = $1 AND category_id = $2`

	_, err := r.db.ExecContext(ctx, query, shopID, categoryID)
	if err != nil {
		return fmt.Errorf("failed to remove category: %w", err)
	}

	return nil
}

func (r *shopCategoryRepository) GetByShopID(ctx context.Context, shopID uuid.UUID) ([]*domain.Category, error) {
	query := `
		SELECT c.id, c.name, c.slug, c.created_at, c.updated_at
		FROM categories c
		INNER JOIN shop_categories sc ON c.id = sc.category_id
		WHERE sc.shop_id = $1`

	rows, err := r.db.QueryContext(ctx, query, shopID)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories by shop: %w", err)
	}
	defer rows.Close()

	var categories []*domain.Category
	for rows.Next() {
		var category domain.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Slug, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, &category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}
