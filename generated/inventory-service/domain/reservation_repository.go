package domain

import (
	"context"

	"github.com/google/uuid"
)

// ReservationRepository defines repository interface for Reservation
type ReservationRepository interface {
	// Create creates a new Reservation
	Create(ctx context.Context, reservation *Reservation) error

	// FindByID finds Reservation by ID
	FindByID(ctx context.Context, id uuid.UUID) (*Reservation, error)

	// Update updates Reservation
	Update(ctx context.Context, reservation *Reservation) error

	// Delete deletes Reservation by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all Reservation
	List(ctx context.Context, limit, offset int) ([]*Reservation, error)
}
