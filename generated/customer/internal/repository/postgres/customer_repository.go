package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository"
)

type customerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) repository.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	query := `
		INSERT INTO customers (id, user_id, first_name, last_name, phone, birth_date, gender,
		                      profile_image_url, profile_thumbnail_100_url, profile_thumbnail_200_url,
		                      created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		customer.ID, customer.UserID, customer.FirstName, customer.LastName, customer.Phone,
		customer.BirthDate, customer.Gender, customer.ProfileImageURL,
		customer.ProfileThumbnail100URL, customer.ProfileThumbnail200URL,
		customer.CreatedAt, customer.UpdatedAt,
	)
	return err
}

func (r *customerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	query := `
		SELECT id, user_id, first_name, last_name, phone, birth_date, gender,
		       profile_image_url, profile_thumbnail_100_url, profile_thumbnail_200_url,
		       created_at, updated_at
		FROM customers WHERE id = $1
	`
	customer := &domain.Customer{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&customer.ID, &customer.UserID, &customer.FirstName, &customer.LastName, &customer.Phone,
		&customer.BirthDate, &customer.Gender, &customer.ProfileImageURL,
		&customer.ProfileThumbnail100URL, &customer.ProfileThumbnail200URL,
		&customer.CreatedAt, &customer.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *customerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Customer, error) {
	query := `
		SELECT id, user_id, first_name, last_name, phone, birth_date, gender,
		       profile_image_url, profile_thumbnail_100_url, profile_thumbnail_200_url,
		       created_at, updated_at
		FROM customers WHERE user_id = $1
	`
	customer := &domain.Customer{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&customer.ID, &customer.UserID, &customer.FirstName, &customer.LastName, &customer.Phone,
		&customer.BirthDate, &customer.Gender, &customer.ProfileImageURL,
		&customer.ProfileThumbnail100URL, &customer.ProfileThumbnail200URL,
		&customer.CreatedAt, &customer.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrCustomerNotFound
	}
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	query := `
		UPDATE customers SET first_name = $2, last_name = $3, phone = $4, birth_date = $5,
		                    gender = $6, profile_image_url = $7, profile_thumbnail_100_url = $8,
		                    profile_thumbnail_200_url = $9, updated_at = $10
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		customer.ID, customer.FirstName, customer.LastName, customer.Phone, customer.BirthDate,
		customer.Gender, customer.ProfileImageURL, customer.ProfileThumbnail100URL,
		customer.ProfileThumbnail200URL, customer.UpdatedAt,
	)
	return err
}
