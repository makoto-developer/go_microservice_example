package infrastructure

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type PostgresProductTagRepository struct {
	db *sql.DB
}

func NewPostgresProductTagRepository(db *sql.DB) *PostgresProductTagRepository {
	return &PostgresProductTagRepository{db: db}
}

func (r *PostgresProductTagRepository) Create(ctx context.Context, productTag *domain.ProductTag) error {
	query := `
		INSERT INTO product_tags (
			id, product_id, tag_name, created_at
		) VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query,
		productTag.Id,
		productTag.ProductId,
		productTag.TagName,
		productTag.CreatedAt,
	)

	return err
}

func (r *PostgresProductTagRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.ProductTag, error) {
	query := `
		SELECT id, product_id, tag_name, created_at
		FROM product_tags
		WHERE id = $1
	`

	var productTag domain.ProductTag

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&productTag.Id,
		&productTag.ProductId,
		&productTag.TagName,
		&productTag.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product tag not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	return &productTag, nil
}

func (r *PostgresProductTagRepository) Update(ctx context.Context, productTag *domain.ProductTag) error {
	query := `
		UPDATE product_tags
		SET product_id = $1,
			tag_name = $2
		WHERE id = $3
	`

	_, err := r.db.ExecContext(ctx, query,
		productTag.ProductId,
		productTag.TagName,
		productTag.Id,
	)

	return err
}

func (r *PostgresProductTagRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM product_tags WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresProductTagRepository) List(ctx context.Context, limit, offset int) ([]*domain.ProductTag, error) {
	query := `
		SELECT id, product_id, tag_name, created_at
		FROM product_tags
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*domain.ProductTag
	for rows.Next() {
		var tag domain.ProductTag

		err := rows.Scan(
			&tag.Id,
			&tag.ProductId,
			&tag.TagName,
			&tag.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		tags = append(tags, &tag)
	}

	return tags, rows.Err()
}

// BeginTx starts a transaction
func (r *PostgresProductTagRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}
