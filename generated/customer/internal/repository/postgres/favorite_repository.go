package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository"
)

type favoriteRepository struct {
	db *sql.DB
}

func NewFavoriteRepository(db *sql.DB) repository.FavoriteRepository {
	return &favoriteRepository{db: db}
}

func (r *favoriteRepository) Add(ctx context.Context, favorite *domain.Favorite) error {
	query := `INSERT INTO favorites (id, customer_id, product_id, notify_on_restock, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.ExecContext(ctx, query, favorite.ID, favorite.CustomerID, favorite.ProductID, favorite.NotifyOnRestock, favorite.CreatedAt)
	return err
}

func (r *favoriteRepository) Remove(ctx context.Context, customerID, productID uuid.UUID) error {
	query := `DELETE FROM favorites WHERE customer_id = $1 AND product_id = $2`
	_, err := r.db.ExecContext(ctx, query, customerID, productID)
	return err
}

func (r *favoriteRepository) List(ctx context.Context, customerID uuid.UUID) ([]*domain.Favorite, error) {
	query := `SELECT id, customer_id, product_id, notify_on_restock, created_at FROM favorites WHERE customer_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var favorites []*domain.Favorite
	for rows.Next() {
		fav := &domain.Favorite{}
		if err := rows.Scan(&fav.ID, &fav.CustomerID, &fav.ProductID, &fav.NotifyOnRestock, &fav.CreatedAt); err != nil {
			return nil, err
		}
		favorites = append(favorites, fav)
	}
	return favorites, rows.Err()
}

func (r *favoriteRepository) Exists(ctx context.Context, customerID, productID uuid.UUID) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM favorites WHERE customer_id = $1 AND product_id = $2)`
	var exists bool
	err := r.db.QueryRowContext(ctx, query, customerID, productID).Scan(&exists)
	return exists, err
}
