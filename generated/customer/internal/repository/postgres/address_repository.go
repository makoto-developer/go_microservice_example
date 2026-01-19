package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/repository"
)

type addressRepository struct {
	db *sql.DB
}

func NewAddressRepository(db *sql.DB) repository.AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(ctx context.Context, address *domain.Address) error {
	query := `
		INSERT INTO addresses (id, customer_id, address_name, postal_code, prefecture, city,
		                      address_line1, address_line2, recipient_name, recipient_phone,
		                      is_default, deleted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.ExecContext(ctx, query,
		address.ID, address.CustomerID, address.AddressName, address.PostalCode,
		address.Prefecture, address.City, address.AddressLine1, address.AddressLine2,
		address.RecipientName, address.RecipientPhone, address.IsDefault, address.Deleted,
		address.CreatedAt, address.UpdatedAt,
	)
	return err
}

func (r *addressRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Address, error) {
	query := `
		SELECT id, customer_id, address_name, postal_code, prefecture, city,
		       address_line1, address_line2, recipient_name, recipient_phone,
		       is_default, deleted, created_at, updated_at
		FROM addresses WHERE id = $1
	`
	address := &domain.Address{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&address.ID, &address.CustomerID, &address.AddressName, &address.PostalCode,
		&address.Prefecture, &address.City, &address.AddressLine1, &address.AddressLine2,
		&address.RecipientName, &address.RecipientPhone, &address.IsDefault, &address.Deleted,
		&address.CreatedAt, &address.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrAddressNotFound
	}
	return address, err
}

func (r *addressRepository) List(ctx context.Context, customerID uuid.UUID) ([]*domain.Address, error) {
	query := `
		SELECT id, customer_id, address_name, postal_code, prefecture, city,
		       address_line1, address_line2, recipient_name, recipient_phone,
		       is_default, deleted, created_at, updated_at
		FROM addresses WHERE customer_id = $1 AND deleted = false
		ORDER BY is_default DESC, created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []*domain.Address
	for rows.Next() {
		addr := &domain.Address{}
		if err := rows.Scan(
			&addr.ID, &addr.CustomerID, &addr.AddressName, &addr.PostalCode,
			&addr.Prefecture, &addr.City, &addr.AddressLine1, &addr.AddressLine2,
			&addr.RecipientName, &addr.RecipientPhone, &addr.IsDefault, &addr.Deleted,
			&addr.CreatedAt, &addr.UpdatedAt,
		); err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}
	return addresses, rows.Err()
}

func (r *addressRepository) Update(ctx context.Context, address *domain.Address) error {
	query := `
		UPDATE addresses SET address_name = $2, postal_code = $3, prefecture = $4,
		                    city = $5, address_line1 = $6, address_line2 = $7,
		                    recipient_name = $8, recipient_phone = $9, is_default = $10,
		                    updated_at = $11
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		address.ID, address.AddressName, address.PostalCode, address.Prefecture,
		address.City, address.AddressLine1, address.AddressLine2, address.RecipientName,
		address.RecipientPhone, address.IsDefault, address.UpdatedAt,
	)
	return err
}

func (r *addressRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE addresses SET deleted = true, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *addressRepository) GetDefault(ctx context.Context, customerID uuid.UUID) (*domain.Address, error) {
	query := `
		SELECT id, customer_id, address_name, postal_code, prefecture, city,
		       address_line1, address_line2, recipient_name, recipient_phone,
		       is_default, deleted, created_at, updated_at
		FROM addresses WHERE customer_id = $1 AND is_default = true AND deleted = false
	`
	address := &domain.Address{}
	err := r.db.QueryRowContext(ctx, query, customerID).Scan(
		&address.ID, &address.CustomerID, &address.AddressName, &address.PostalCode,
		&address.Prefecture, &address.City, &address.AddressLine1, &address.AddressLine2,
		&address.RecipientName, &address.RecipientPhone, &address.IsDefault, &address.Deleted,
		&address.CreatedAt, &address.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrAddressNotFound
	}
	return address, err
}

func (r *addressRepository) SetDefault(ctx context.Context, customerID, addressID uuid.UUID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `UPDATE addresses SET is_default = false WHERE customer_id = $1`, customerID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `UPDATE addresses SET is_default = true WHERE id = $1`, addressID)
	if err != nil {
		return err
	}

	return tx.Commit()
}
