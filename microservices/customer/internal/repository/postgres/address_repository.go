package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository"
)

type addressRepository struct {
	db *sql.DB
}

// NewAddressRepository creates a new PostgreSQL address repository
func NewAddressRepository(db *sql.DB) repository.AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(ctx context.Context, address *domain.Address) error {
	query := `
		INSERT INTO addresses (id, customer_id, label, postal_code, prefecture, city, address_line1, address_line2,
			recipient_name, recipient_phone, is_default, is_deleted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.ExecContext(ctx, query,
		address.ID,
		address.CustomerID,
		address.Label,
		address.PostalCode,
		address.Prefecture,
		address.City,
		address.AddressLine1,
		address.AddressLine2,
		address.RecipientName,
		address.RecipientPhone,
		address.IsDefault,
		address.IsDeleted,
		address.CreatedAt,
		address.UpdatedAt,
	)
	return err
}

func (r *addressRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Address, error) {
	address := &domain.Address{}
	query := `
		SELECT id, customer_id, label, postal_code, prefecture, city, address_line1, address_line2,
			recipient_name, recipient_phone, is_default, is_deleted, created_at, updated_at, deleted_at
		FROM addresses
		WHERE id = $1 AND is_deleted = false
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&address.ID,
		&address.CustomerID,
		&address.Label,
		&address.PostalCode,
		&address.Prefecture,
		&address.City,
		&address.AddressLine1,
		&address.AddressLine2,
		&address.RecipientName,
		&address.RecipientPhone,
		&address.IsDefault,
		&address.IsDeleted,
		&address.CreatedAt,
		&address.UpdatedAt,
		&address.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("address not found")
		}
		return nil, err
	}
	return address, nil
}

func (r *addressRepository) List(ctx context.Context, customerID uuid.UUID) ([]*domain.Address, error) {
	query := `
		SELECT id, customer_id, label, postal_code, prefecture, city, address_line1, address_line2,
			recipient_name, recipient_phone, is_default, is_deleted, created_at, updated_at, deleted_at
		FROM addresses
		WHERE customer_id = $1 AND is_deleted = false
		ORDER BY is_default DESC, created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []*domain.Address
	for rows.Next() {
		address := &domain.Address{}
		err := rows.Scan(
			&address.ID,
			&address.CustomerID,
			&address.Label,
			&address.PostalCode,
			&address.Prefecture,
			&address.City,
			&address.AddressLine1,
			&address.AddressLine2,
			&address.RecipientName,
			&address.RecipientPhone,
			&address.IsDefault,
			&address.IsDeleted,
			&address.CreatedAt,
			&address.UpdatedAt,
			&address.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}

	return addresses, rows.Err()
}

func (r *addressRepository) Update(ctx context.Context, address *domain.Address) error {
	query := `
		UPDATE addresses
		SET label = $2, postal_code = $3, prefecture = $4, city = $5, address_line1 = $6, address_line2 = $7,
			recipient_name = $8, recipient_phone = $9, is_default = $10, updated_at = $11
		WHERE id = $1 AND is_deleted = false
	`
	_, err := r.db.ExecContext(ctx, query,
		address.ID,
		address.Label,
		address.PostalCode,
		address.Prefecture,
		address.City,
		address.AddressLine1,
		address.AddressLine2,
		address.RecipientName,
		address.RecipientPhone,
		address.IsDefault,
		address.UpdatedAt,
	)
	return err
}

func (r *addressRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE addresses
		SET is_deleted = true, deleted_at = NOW()
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *addressRepository) GetDefault(ctx context.Context, customerID uuid.UUID) (*domain.Address, error) {
	address := &domain.Address{}
	query := `
		SELECT id, customer_id, label, postal_code, prefecture, city, address_line1, address_line2,
			recipient_name, recipient_phone, is_default, is_deleted, created_at, updated_at, deleted_at
		FROM addresses
		WHERE customer_id = $1 AND is_default = true AND is_deleted = false
		LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(
		&address.ID,
		&address.CustomerID,
		&address.Label,
		&address.PostalCode,
		&address.Prefecture,
		&address.City,
		&address.AddressLine1,
		&address.AddressLine2,
		&address.RecipientName,
		&address.RecipientPhone,
		&address.IsDefault,
		&address.IsDeleted,
		&address.CreatedAt,
		&address.UpdatedAt,
		&address.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No default address is ok
		}
		return nil, err
	}
	return address, nil
}

func (r *addressRepository) SetDefault(ctx context.Context, customerID, addressID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Unset all current default addresses for this customer
	query1 := `UPDATE addresses SET is_default = false WHERE customer_id = $1 AND is_deleted = false`
	if _, err := tx.ExecContext(ctx, query1, customerID); err != nil {
		return err
	}

	// Set the new default address
	query2 := `UPDATE addresses SET is_default = true WHERE id = $1 AND customer_id = $2 AND is_deleted = false`
	if _, err := tx.ExecContext(ctx, query2, addressID, customerID); err != nil {
		return err
	}

	return tx.Commit()
}
