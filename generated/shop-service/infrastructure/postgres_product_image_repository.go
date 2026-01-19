package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type PostgresProductImageRepository struct {
	db *sql.DB
}

func NewPostgresProductImageRepository(db *sql.DB) *PostgresProductImageRepository {
	return &PostgresProductImageRepository{db: db}
}

func (r *PostgresProductImageRepository) Create(ctx context.Context, productImage *domain.ProductImage) error {
	query := `
		INSERT INTO product_images (
			id, product_id, url, display_order,
			thumbnail_200_url, thumbnail_400_url, thumbnail_800_url,
			created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		productImage.Id,
		productImage.ProductId,
		productImage.Url,
		productImage.DisplayOrder,
		productImage.Thumbnail200Url,
		productImage.Thumbnail400Url,
		productImage.Thumbnail800Url,
		productImage.CreatedAt,
	)

	return err
}

func (r *PostgresProductImageRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductImage, error) {
	query := `
		SELECT id, product_id, url, display_order,
			thumbnail_200_url, thumbnail_400_url, thumbnail_800_url,
			created_at
		FROM product_images
		WHERE id = $1
	`

	var productImage domain.ProductImage

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&productImage.Id,
		&productImage.ProductId,
		&productImage.Url,
		&productImage.DisplayOrder,
		&productImage.Thumbnail200Url,
		&productImage.Thumbnail400Url,
		&productImage.Thumbnail800Url,
		&productImage.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product image not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	return &productImage, nil
}

func (r *PostgresProductImageRepository) Update(ctx context.Context, productImage *domain.ProductImage) error {
	query := `
		UPDATE product_images
		SET product_id = $1,
			url = $2,
			display_order = $3,
			thumbnail_200_url = $4,
			thumbnail_400_url = $5,
			thumbnail_800_url = $6
		WHERE id = $7
	`

	_, err := r.db.ExecContext(ctx, query,
		productImage.ProductId,
		productImage.Url,
		productImage.DisplayOrder,
		productImage.Thumbnail200Url,
		productImage.Thumbnail400Url,
		productImage.Thumbnail800Url,
		productImage.Id,
	)

	return err
}

func (r *PostgresProductImageRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM product_images WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresProductImageRepository) List(ctx context.Context, limit, offset int) ([]*domain.ProductImage, error) {
	query := `
		SELECT id, product_id, url, display_order,
			thumbnail_200_url, thumbnail_400_url, thumbnail_800_url,
			created_at
		FROM product_images
		ORDER BY display_order ASC, created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*domain.ProductImage
	for rows.Next() {
		var image domain.ProductImage

		err := rows.Scan(
			&image.Id,
			&image.ProductId,
			&image.Url,
			&image.DisplayOrder,
			&image.Thumbnail200Url,
			&image.Thumbnail400Url,
			&image.Thumbnail800Url,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		images = append(images, &image)
	}

	return images, rows.Err()
}

// BeginTx starts a transaction
func (r *PostgresProductImageRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
