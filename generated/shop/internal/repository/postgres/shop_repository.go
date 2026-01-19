package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type shopRepository struct {
	db *sql.DB
}

func NewShopRepository(db *sql.DB) repository.ShopRepository {
	return &shopRepository{db: db}
}

func (r *shopRepository) Create(ctx context.Context, shop *domain.Shop) error {
	query := `
		INSERT INTO shops (id, owner_id, name, description, logo_url, owner_name, phone_number,
		                   business_hours, return_policy, status, published, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	_, err := r.db.ExecContext(ctx, query,
		shop.ID, shop.OwnerID, shop.Name, shop.Description, shop.LogoURL, shop.OwnerName,
		shop.PhoneNumber, shop.BusinessHours, shop.ReturnPolicy, shop.Status, shop.Published,
		shop.CreatedAt, shop.UpdatedAt,
	)
	return err
}

func (r *shopRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Shop, error) {
	query := `
		SELECT id, owner_id, name, description, logo_url, owner_name, phone_number,
		       business_hours, return_policy, status, published, created_at, updated_at
		FROM shops WHERE id = $1
	`
	shop := &domain.Shop{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&shop.ID, &shop.OwnerID, &shop.Name, &shop.Description, &shop.LogoURL, &shop.OwnerName,
		&shop.PhoneNumber, &shop.BusinessHours, &shop.ReturnPolicy, &shop.Status, &shop.Published,
		&shop.CreatedAt, &shop.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrShopNotFound
	}
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (r *shopRepository) GetByOwnerID(ctx context.Context, ownerID uuid.UUID) (*domain.Shop, error) {
	query := `
		SELECT id, owner_id, name, description, logo_url, owner_name, phone_number,
		       business_hours, return_policy, status, published, created_at, updated_at
		FROM shops WHERE owner_id = $1
	`
	shop := &domain.Shop{}
	err := r.db.QueryRowContext(ctx, query, ownerID).Scan(
		&shop.ID, &shop.OwnerID, &shop.Name, &shop.Description, &shop.LogoURL, &shop.OwnerName,
		&shop.PhoneNumber, &shop.BusinessHours, &shop.ReturnPolicy, &shop.Status, &shop.Published,
		&shop.CreatedAt, &shop.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrShopNotFound
	}
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (r *shopRepository) Update(ctx context.Context, shop *domain.Shop) error {
	query := `
		UPDATE shops SET name = $2, description = $3, logo_url = $4, owner_name = $5,
		                phone_number = $6, business_hours = $7, return_policy = $8, updated_at = $9
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		shop.ID, shop.Name, shop.Description, shop.LogoURL, shop.OwnerName,
		shop.PhoneNumber, shop.BusinessHours, shop.ReturnPolicy, shop.UpdatedAt,
	)
	return err
}

func (r *shopRepository) UpdateStatus(ctx context.Context, shopID uuid.UUID, status domain.ShopStatus) error {
	query := `UPDATE shops SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, shopID, status)
	return err
}

func (r *shopRepository) UpdatePublished(ctx context.Context, shopID uuid.UUID, published bool) error {
	query := `UPDATE shops SET published = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, shopID, published)
	return err
}

func (r *shopRepository) AddCategory(ctx context.Context, category *domain.ShopCategory) error {
	query := `
		INSERT INTO shop_categories (id, shop_id, category_name, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.ExecContext(ctx, query, category.ID, category.ShopID, category.CategoryName, category.CreatedAt)
	return err
}

func (r *shopRepository) GetCategories(ctx context.Context, shopID uuid.UUID) ([]*domain.ShopCategory, error) {
	query := `SELECT id, shop_id, category_name, created_at FROM shop_categories WHERE shop_id = $1`
	rows, err := r.db.QueryContext(ctx, query, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*domain.ShopCategory
	for rows.Next() {
		cat := &domain.ShopCategory{}
		if err := rows.Scan(&cat.ID, &cat.ShopID, &cat.CategoryName, &cat.CreatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, rows.Err()
}

func (r *shopRepository) DeleteCategories(ctx context.Context, shopID uuid.UUID) error {
	query := `DELETE FROM shop_categories WHERE shop_id = $1`
	_, err := r.db.ExecContext(ctx, query, shopID)
	return err
}
