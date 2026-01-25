package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
)

type shopRepository struct {
	db *sql.DB
}

// NewShopRepository creates a new shop repository
func NewShopRepository(db *sql.DB) *shopRepository {
	return &shopRepository{db: db}
}

func (r *shopRepository) Create(ctx context.Context, shop *domain.Shop) error {
	query := `
		INSERT INTO shops (
			id, owner_id, name, description, logo_url,
			owner_name, phone_number, business_hours, return_policy,
			status, published, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		)`

	_, err := r.db.ExecContext(ctx, query,
		shop.ID, shop.OwnerID, shop.Name, shop.Description, shop.LogoURL,
		shop.OwnerName, shop.PhoneNumber, shop.BusinessHours, shop.ReturnPolicy,
		shop.Status, shop.Published, shop.CreatedAt, shop.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create shop: %w", err)
	}

	return nil
}

func (r *shopRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Shop, error) {
	query := `
		SELECT id, owner_id, name, description, logo_url,
			   owner_name, phone_number, business_hours, return_policy,
			   status, published, created_at, updated_at, approved_at, approved_by
		FROM shops
		WHERE id = $1`

	var shop domain.Shop
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&shop.ID, &shop.OwnerID, &shop.Name, &shop.Description, &shop.LogoURL,
		&shop.OwnerName, &shop.PhoneNumber, &shop.BusinessHours, &shop.ReturnPolicy,
		&shop.Status, &shop.Published, &shop.CreatedAt, &shop.UpdatedAt,
		&shop.ApprovedAt, &shop.ApprovedBy,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("shop not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}

	return &shop, nil
}

func (r *shopRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) ([]*domain.Shop, error) {
	query := `
		SELECT id, owner_id, name, description, logo_url,
			   owner_name, phone_number, business_hours, return_policy,
			   status, published, created_at, updated_at, approved_at, approved_by
		FROM shops
		WHERE owner_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get shops by owner: %w", err)
	}
	defer rows.Close()

	var shops []*domain.Shop
	for rows.Next() {
		var shop domain.Shop
		err := rows.Scan(
			&shop.ID, &shop.OwnerID, &shop.Name, &shop.Description, &shop.LogoURL,
			&shop.OwnerName, &shop.PhoneNumber, &shop.BusinessHours, &shop.ReturnPolicy,
			&shop.Status, &shop.Published, &shop.CreatedAt, &shop.UpdatedAt,
			&shop.ApprovedAt, &shop.ApprovedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan shop: %w", err)
		}
		shops = append(shops, &shop)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating shops: %w", err)
	}

	return shops, nil
}

func (r *shopRepository) Update(ctx context.Context, shop *domain.Shop) error {
	query := `
		UPDATE shops
		SET name = $2, description = $3, logo_url = $4,
			owner_name = $5, phone_number = $6, business_hours = $7,
			return_policy = $8, updated_at = $9
		WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query,
		shop.ID, shop.Name, shop.Description, shop.LogoURL,
		shop.OwnerName, shop.PhoneNumber, shop.BusinessHours,
		shop.ReturnPolicy, shop.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update shop: %w", err)
	}

	return nil
}

func (r *shopRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ShopStatus) error {
	query := `UPDATE shops SET status = $2, updated_at = NOW() WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update shop status: %w", err)
	}

	return nil
}

func (r *shopRepository) UpdateIsPublic(ctx context.Context, id uuid.UUID, isPublic bool) error {
	query := `UPDATE shops SET published = $2, updated_at = NOW() WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id, isPublic)
	if err != nil {
		return fmt.Errorf("failed to update shop published status: %w", err)
	}

	return nil
}

func (r *shopRepository) List(ctx context.Context, limit, offset int) ([]*domain.Shop, error) {
	query := `
		SELECT id, owner_id, name, description, logo_url,
			   owner_name, phone_number, business_hours, return_policy,
			   status, published, created_at, updated_at, approved_at, approved_by
		FROM shops
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list shops: %w", err)
	}
	defer rows.Close()

	var shops []*domain.Shop
	for rows.Next() {
		var shop domain.Shop
		err := rows.Scan(
			&shop.ID, &shop.OwnerID, &shop.Name, &shop.Description, &shop.LogoURL,
			&shop.OwnerName, &shop.PhoneNumber, &shop.BusinessHours, &shop.ReturnPolicy,
			&shop.Status, &shop.Published, &shop.CreatedAt, &shop.UpdatedAt,
			&shop.ApprovedAt, &shop.ApprovedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan shop: %w", err)
		}
		shops = append(shops, &shop)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating shops: %w", err)
	}

	return shops, nil
}
