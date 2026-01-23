package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository"
)

type reviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) repository.ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(ctx context.Context, review *domain.Review) error {
	query := `
		INSERT INTO reviews (id, customer_id, product_id, order_id, rating, review_text, editable_until, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		review.ID,
		review.CustomerID,
		review.ProductID,
		review.OrderID,
		review.Rating,
		review.ReviewText,
		review.EditableUntil,
		review.CreatedAt,
		review.UpdatedAt,
	)
	return err
}

func (r *reviewRepository) Update(ctx context.Context, review *domain.Review) error {
	query := `
		UPDATE reviews
		SET rating = $1, review_text = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query,
		review.Rating,
		review.ReviewText,
		review.UpdatedAt,
		review.ID,
	)
	return err
}

func (r *reviewRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM reviews WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *reviewRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Review, error) {
	query := `
		SELECT id, customer_id, product_id, order_id, rating, review_text, editable_until, created_at, updated_at
		FROM reviews WHERE id = $1
	`
	rev := &domain.Review{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rev.ID,
		&rev.CustomerID,
		&rev.ProductID,
		&rev.OrderID,
		&rev.Rating,
		&rev.ReviewText,
		&rev.EditableUntil,
		&rev.CreatedAt,
		&rev.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrReviewNotFound
	}
	return rev, err
}

func (r *reviewRepository) GetByOrderAndProduct(ctx context.Context, orderID, productID uuid.UUID) (*domain.Review, error) {
	query := `
		SELECT id, customer_id, product_id, order_id, rating, review_text, editable_until, created_at, updated_at
		FROM reviews WHERE order_id = $1 AND product_id = $2
	`
	rev := &domain.Review{}
	err := r.db.QueryRowContext(ctx, query, orderID, productID).Scan(
		&rev.ID,
		&rev.CustomerID,
		&rev.ProductID,
		&rev.OrderID,
		&rev.Rating,
		&rev.ReviewText,
		&rev.EditableUntil,
		&rev.CreatedAt,
		&rev.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrReviewNotFound
	}
	return rev, err
}

func (r *reviewRepository) ListByCustomer(ctx context.Context, customerID uuid.UUID) ([]*domain.Review, error) {
	query := `
		SELECT id, customer_id, product_id, order_id, rating, review_text, editable_until, created_at, updated_at
		FROM reviews WHERE customer_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*domain.Review
	for rows.Next() {
		rev := &domain.Review{}
		if err := rows.Scan(
			&rev.ID,
			&rev.CustomerID,
			&rev.ProductID,
			&rev.OrderID,
			&rev.Rating,
			&rev.ReviewText,
			&rev.EditableUntil,
			&rev.CreatedAt,
			&rev.UpdatedAt,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, rows.Err()
}

func (r *reviewRepository) ListByProduct(ctx context.Context, productID uuid.UUID) ([]*domain.Review, error) {
	query := `
		SELECT id, customer_id, product_id, order_id, rating, review_text, editable_until, created_at, updated_at
		FROM reviews WHERE product_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reviews []*domain.Review
	for rows.Next() {
		rev := &domain.Review{}
		if err := rows.Scan(
			&rev.ID,
			&rev.CustomerID,
			&rev.ProductID,
			&rev.OrderID,
			&rev.Rating,
			&rev.ReviewText,
			&rev.EditableUntil,
			&rev.CreatedAt,
			&rev.UpdatedAt,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, rev)
	}
	return reviews, rows.Err()
}

func (r *reviewRepository) AddImage(ctx context.Context, image *domain.ReviewImage) error {
	query := `INSERT INTO review_images (id, review_id, image_url, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, image.ID, image.ReviewID, image.ImageURL, image.CreatedAt)
	return err
}

func (r *reviewRepository) GetImages(ctx context.Context, reviewID uuid.UUID) ([]*domain.ReviewImage, error) {
	query := `SELECT id, review_id, image_url, created_at FROM review_images WHERE review_id = $1`
	rows, err := r.db.QueryContext(ctx, query, reviewID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*domain.ReviewImage
	for rows.Next() {
		img := &domain.ReviewImage{}
		if err := rows.Scan(&img.ID, &img.ReviewID, &img.ImageURL, &img.CreatedAt); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, rows.Err()
}
