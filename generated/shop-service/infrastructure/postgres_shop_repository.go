package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type PostgresShopRepository struct {
	db *sql.DB
}

func NewPostgresShopRepository(db *sql.DB) *PostgresShopRepository {
	return &PostgresShopRepository{db: db}
}

func (r *PostgresShopRepository) Create(ctx context.Context, shop *domain.Shop) error {
	query := `
		INSERT INTO shops (
			id, owner_id, name, description, logo_url,
			owner_name, phone_number, business_hours, return_policy,
			status, published, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.ExecContext(ctx, query,
		shop.Id,
		shop.OwnerId,
		shop.Name,
		shop.Description,
		shop.LogoUrl,
		shop.OwnerName,
		shop.PhoneNumber,
		shop.BusinessHours,
		shop.ReturnPolicy,
		shop.Status,
		shop.Published,
		shop.CreatedAt,
		shop.UpdatedAt,
	)

	return err
}

func (r *PostgresShopRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Shop, error) {
	query := `
		SELECT id, owner_id, name, description, logo_url,
			owner_name, phone_number, business_hours, return_policy,
			status, published, created_at, updated_at
		FROM shops
		WHERE id = $1
	`

	var shop domain.Shop
	var logoUrl sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&shop.Id,
		&shop.OwnerId,
		&shop.Name,
		&shop.Description,
		&logoUrl,
		&shop.OwnerName,
		&shop.PhoneNumber,
		&shop.BusinessHours,
		&shop.ReturnPolicy,
		&shop.Status,
		&shop.Published,
		&shop.CreatedAt,
		&shop.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("shop not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	if logoUrl.Valid {
		shop.LogoUrl = &logoUrl.String
	}

	return &shop, nil
}

func (r *PostgresShopRepository) Update(ctx context.Context, shop *domain.Shop) error {
	query := `
		UPDATE shops
		SET owner_id = $1,
			name = $2,
			description = $3,
			logo_url = $4,
			owner_name = $5,
			phone_number = $6,
			business_hours = $7,
			return_policy = $8,
			status = $9,
			published = $10,
			updated_at = $11
		WHERE id = $12
	`

	_, err := r.db.ExecContext(ctx, query,
		shop.OwnerId,
		shop.Name,
		shop.Description,
		shop.LogoUrl,
		shop.OwnerName,
		shop.PhoneNumber,
		shop.BusinessHours,
		shop.ReturnPolicy,
		shop.Status,
		shop.Published,
		time.Now(),
		shop.Id,
	)

	return err
}

func (r *PostgresShopRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM shops WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresShopRepository) List(ctx context.Context, limit, offset int) ([]*domain.Shop, error) {
	query := `
		SELECT id, owner_id, name, description, logo_url,
			owner_name, phone_number, business_hours, return_policy,
			status, published, created_at, updated_at
		FROM shops
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shops []*domain.Shop
	for rows.Next() {
		var shop domain.Shop
		var logoUrl sql.NullString

		err := rows.Scan(
			&shop.Id,
			&shop.OwnerId,
			&shop.Name,
			&shop.Description,
			&logoUrl,
			&shop.OwnerName,
			&shop.PhoneNumber,
			&shop.BusinessHours,
			&shop.ReturnPolicy,
			&shop.Status,
			&shop.Published,
			&shop.CreatedAt,
			&shop.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if logoUrl.Valid {
			shop.LogoUrl = &logoUrl.String
		}

		shops = append(shops, &shop)
	}

	return shops, rows.Err()
}

// BeginTx starts a transaction
func (r *PostgresShopRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
