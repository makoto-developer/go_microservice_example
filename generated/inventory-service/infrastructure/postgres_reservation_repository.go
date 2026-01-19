package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/inventory-service/domain"
)

type PostgresReservationRepository struct {
	db *sql.DB
}

func NewPostgresReservationRepository(db *sql.DB) *PostgresReservationRepository {
	return &PostgresReservationRepository{db: db}
}

func (r *PostgresReservationRepository) Create(ctx context.Context, reservation *domain.Reservation) error {
	query := `
		INSERT INTO reservations (
			id, inventory_id, order_id, quantity, status,
			expires_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		reservation.ID,
		reservation.InventoryID,
		reservation.OrderID,
		reservation.Quantity,
		reservation.Status,
		reservation.ExpiresAt,
		reservation.CreatedAt,
		reservation.UpdatedAt,
	)

	return err
}

func (r *PostgresReservationRepository) CreateWithTx(ctx context.Context, tx *sql.Tx, reservation *domain.Reservation) error {
	query := `
		INSERT INTO reservations (
			id, inventory_id, order_id, quantity, status,
			expires_at, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := tx.ExecContext(ctx, query,
		reservation.ID,
		reservation.InventoryID,
		reservation.OrderID,
		reservation.Quantity,
		reservation.Status,
		reservation.ExpiresAt,
		reservation.CreatedAt,
		reservation.UpdatedAt,
	)

	return err
}

func (r *PostgresReservationRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Reservation, error) {
	query := `
		SELECT id, inventory_id, order_id, quantity, status,
			expires_at, confirmed_at, created_at, updated_at
		FROM reservations
		WHERE id = $1
	`

	var reservation domain.Reservation
	var confirmedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&reservation.ID,
		&reservation.InventoryID,
		&reservation.OrderID,
		&reservation.Quantity,
		&reservation.Status,
		&reservation.ExpiresAt,
		&confirmedAt,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("reservation not found: %s", id)
	}
	if err != nil {
		return nil, err
	}

	if confirmedAt.Valid {
		reservation.ConfirmedAt = &confirmedAt.Time
	}

	return &reservation, nil
}

func (r *PostgresReservationRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.Reservation, error) {
	query := `
		SELECT id, inventory_id, order_id, quantity, status,
			expires_at, confirmed_at, created_at, updated_at
		FROM reservations
		WHERE order_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*domain.Reservation
	for rows.Next() {
		var reservation domain.Reservation
		var confirmedAt sql.NullTime

		err := rows.Scan(
			&reservation.ID,
			&reservation.InventoryID,
			&reservation.OrderID,
			&reservation.Quantity,
			&reservation.Status,
			&reservation.ExpiresAt,
			&confirmedAt,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if confirmedAt.Valid {
			reservation.ConfirmedAt = &confirmedAt.Time
		}

		reservations = append(reservations, &reservation)
	}

	return reservations, rows.Err()
}

func (r *PostgresReservationRepository) FindExpired(ctx context.Context) ([]*domain.Reservation, error) {
	query := `
		SELECT id, inventory_id, order_id, quantity, status,
			expires_at, confirmed_at, created_at, updated_at
		FROM reservations
		WHERE status = 'RESERVED' AND expires_at < NOW()
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []*domain.Reservation
	for rows.Next() {
		var reservation domain.Reservation
		var confirmedAt sql.NullTime

		err := rows.Scan(
			&reservation.ID,
			&reservation.InventoryID,
			&reservation.OrderID,
			&reservation.Quantity,
			&reservation.Status,
			&reservation.ExpiresAt,
			&confirmedAt,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if confirmedAt.Valid {
			reservation.ConfirmedAt = &confirmedAt.Time
		}

		reservations = append(reservations, &reservation)
	}

	return reservations, rows.Err()
}

func (r *PostgresReservationRepository) Update(ctx context.Context, reservation *domain.Reservation) error {
	query := `
		UPDATE reservations
		SET status = $1,
			confirmed_at = $2,
			updated_at = $3
		WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query,
		reservation.Status,
		reservation.ConfirmedAt,
		time.Now(),
		reservation.ID,
	)

	return err
}

func (r *PostgresReservationRepository) UpdateWithTx(ctx context.Context, tx *sql.Tx, reservation *domain.Reservation) error {
	query := `
		UPDATE reservations
		SET status = $1,
			confirmed_at = $2,
			updated_at = $3
		WHERE id = $4
	`

	_, err := tx.ExecContext(ctx, query,
		reservation.Status,
		reservation.ConfirmedAt,
		time.Now(),
		reservation.ID,
	)

	return err
}

func (r *PostgresReservationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM reservations WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *PostgresReservationRepository) DeleteWithTx(ctx context.Context, tx *sql.Tx, id uuid.UUID) error {
	query := `DELETE FROM reservations WHERE id = $1`
	_, err := tx.ExecContext(ctx, query, id)
	return err
}
