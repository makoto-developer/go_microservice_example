package clients

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type ReviewClient struct {
	db *sql.DB
}

type Review struct {
	ID            uuid.UUID
	CustomerID    uuid.UUID
	ProductID     uuid.UUID
	OrderID       uuid.UUID
	Rating        int
	ReviewText    string
	EditableUntil time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewReviewClient(databaseURL string) (*ReviewClient, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to review database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping review database: %w", err)
	}

	return &ReviewClient{db: db}, nil
}

func (c *ReviewClient) Close() error {
	return c.db.Close()
}

func (c *ReviewClient) CreateReview(customerID, productID, orderID uuid.UUID, rating int, reviewText string) (*Review, error) {
	if rating < 1 || rating > 5 {
		return nil, fmt.Errorf("rating must be between 1 and 5")
	}

	// Editable until 30 days from now
	editableUntil := time.Now().Add(30 * 24 * time.Hour)

	query := `
		INSERT INTO reviews (id, customer_id, product_id, order_id, rating, review_text, editable_until)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, customer_id, product_id, order_id, rating, review_text, editable_until, created_at, updated_at
	`

	var review Review
	reviewID := uuid.New()

	err := c.db.QueryRow(query, reviewID, customerID, productID, orderID, rating, reviewText, editableUntil).Scan(
		&review.ID,
		&review.CustomerID,
		&review.ProductID,
		&review.OrderID,
		&review.Rating,
		&review.ReviewText,
		&review.EditableUntil,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return &review, nil
}

func (c *ReviewClient) GetReview(id uuid.UUID) (*Review, error) {
	query := `
		SELECT id, customer_id, product_id, order_id, rating, review_text, editable_until, created_at, updated_at
		FROM reviews
		WHERE id = $1
	`

	var review Review
	err := c.db.QueryRow(query, id).Scan(
		&review.ID,
		&review.CustomerID,
		&review.ProductID,
		&review.OrderID,
		&review.Rating,
		&review.ReviewText,
		&review.EditableUntil,
		&review.CreatedAt,
		&review.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	return &review, nil
}

func (c *ReviewClient) GetReviewsByProduct(productID uuid.UUID) ([]Review, error) {
	query := `
		SELECT id, customer_id, product_id, order_id, rating, review_text, editable_until, created_at, updated_at
		FROM reviews
		WHERE product_id = $1
		ORDER BY created_at DESC
	`

	rows, err := c.db.Query(query, productID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}
	defer rows.Close()

	var reviews []Review
	for rows.Next() {
		var review Review
		err := rows.Scan(
			&review.ID,
			&review.CustomerID,
			&review.ProductID,
			&review.OrderID,
			&review.Rating,
			&review.ReviewText,
			&review.EditableUntil,
			&review.CreatedAt,
			&review.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan review: %w", err)
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

func (c *ReviewClient) GetAverageRating(productID uuid.UUID) (float64, error) {
	query := `
		SELECT COALESCE(AVG(rating), 0) as avg_rating
		FROM reviews
		WHERE product_id = $1
	`

	var avgRating float64
	err := c.db.QueryRow(query, productID).Scan(&avgRating)
	if err != nil {
		return 0, fmt.Errorf("failed to get average rating: %w", err)
	}

	return avgRating, nil
}

func (c *ReviewClient) UpdateReview(id uuid.UUID, rating int, reviewText string) error {
	if rating < 1 || rating > 5 {
		return fmt.Errorf("rating must be between 1 and 5")
	}

	// Check if review is still editable
	var editableUntil time.Time
	err := c.db.QueryRow("SELECT editable_until FROM reviews WHERE id = $1", id).Scan(&editableUntil)
	if err != nil {
		return fmt.Errorf("failed to get review: %w", err)
	}

	if time.Now().After(editableUntil) {
		return fmt.Errorf("review is no longer editable")
	}

	query := `
		UPDATE reviews
		SET rating = $1, review_text = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	_, err = c.db.Exec(query, rating, reviewText, id)
	if err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}

	return nil
}

func (c *ReviewClient) DeleteReview(id uuid.UUID) error {
	query := `DELETE FROM reviews WHERE id = $1`
	_, err := c.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}
	return nil
}
