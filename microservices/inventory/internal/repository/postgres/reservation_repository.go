package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/inventory/internal/domain"
)

type reservationRepository struct {
	db *sql.DB
}

func NewReservationRepository(db *sql.DB) *reservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(ctx context.Context, reservation *domain.Reservation) error {
	query := `INSERT INTO reservations (id, inventory_id, order_id, quantity, status, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, reservation.ID, reservation.InventoryID, reservation.OrderID,
		reservation.Quantity, reservation.Status, reservation.ExpiresAt, reservation.CreatedAt, reservation.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}
	return nil
}

func (r *reservationRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Reservation, error) {
	query := `SELECT id, inventory_id, order_id, quantity, status, expires_at, created_at, updated_at
		FROM reservations WHERE id = $1`

	var res domain.Reservation
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&res.ID, &res.InventoryID, &res.OrderID, &res.Quantity, &res.Status, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("reservation not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get reservation: %w", err)
	}
	return &res, nil
}

func (r *reservationRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*domain.Reservation, error) {
	query := `SELECT id, inventory_id, order_id, quantity, status, expires_at, created_at, updated_at
		FROM reservations WHERE order_id = $1`

	rows, err := r.db.QueryContext(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reservations: %w", err)
	}
	defer rows.Close()

	var reservations []*domain.Reservation
	for rows.Next() {
		var res domain.Reservation
		if err := rows.Scan(&res.ID, &res.InventoryID, &res.OrderID, &res.Quantity, &res.Status, &res.ExpiresAt, &res.CreatedAt, &res.UpdatedAt); err != nil {
			return nil, err
		}
		reservations = append(reservations, &res)
	}
	return reservations, nil
}

func (r *reservationRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.ReservationStatus) error {
	query := `UPDATE reservations SET status = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}
	return nil
}

func (r *reservationRepository) DeleteExpired(ctx context.Context) error {
	query := `UPDATE reservations SET status = $1, updated_at = NOW()
		WHERE status = $2 AND expires_at < NOW()`
	_, err := r.db.ExecContext(ctx, query, domain.ReservationStatusExpired, domain.ReservationStatusPending)
	if err != nil {
		return fmt.Errorf("failed to delete expired: %w", err)
	}
	return nil
}
