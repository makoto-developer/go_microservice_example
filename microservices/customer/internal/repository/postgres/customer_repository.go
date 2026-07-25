package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/repository"
)

type customerRepository struct {
	db *sql.DB
}

// NewCustomerRepository creates a new PostgreSQL customer repository
func NewCustomerRepository(db *sql.DB) repository.CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	query := `
		INSERT INTO customers (id, user_id, first_name, last_name, phone_number, birth_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		customer.ID,
		customer.UserID,
		customer.FirstName,
		customer.LastName,
		customer.PhoneNumber,
		customer.BirthDate,
		customer.CreatedAt,
		customer.UpdatedAt,
	)
	return err
}

func (r *customerRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Customer, error) {
	customer := &domain.Customer{}
	query := `
		SELECT id, user_id, first_name, last_name, phone_number, birth_date, created_at, updated_at
		FROM customers
		WHERE id = $1
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&customer.ID,
		&customer.UserID,
		&customer.FirstName,
		&customer.LastName,
		&customer.PhoneNumber,
		&customer.BirthDate,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}
	return customer, nil
}

func (r *customerRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*domain.Customer, error) {
	customer := &domain.Customer{}
	query := `
		SELECT id, user_id, first_name, last_name, phone_number, birth_date, created_at, updated_at
		FROM customers
		WHERE user_id = $1
	`
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&customer.ID,
		&customer.UserID,
		&customer.FirstName,
		&customer.LastName,
		&customer.PhoneNumber,
		&customer.BirthDate,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}
	return customer, nil
}

func (r *customerRepository) Update(ctx context.Context, customer *domain.Customer) error {
	query := `
		UPDATE customers
		SET first_name = $2, last_name = $3, phone_number = $4, birth_date = $5, updated_at = $6
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		customer.ID,
		customer.FirstName,
		customer.LastName,
		customer.PhoneNumber,
		customer.BirthDate,
		customer.UpdatedAt,
	)
	return err
}

func (r *customerRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM customers WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *customerRepository) List(ctx context.Context, limit, offset int) ([]*domain.Customer, error) {
	query := `
		SELECT id, user_id, first_name, last_name, phone_number, birth_date, created_at, updated_at
		FROM customers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []*domain.Customer
	for rows.Next() {
		customer := &domain.Customer{}
		err := rows.Scan(
			&customer.ID,
			&customer.UserID,
			&customer.FirstName,
			&customer.LastName,
			&customer.PhoneNumber,
			&customer.BirthDate,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	return customers, rows.Err()
}
