package domain

import (
	"context"

	"github.com/google/uuid"
)

// CarrierTrackingRepository defines repository interface for CarrierTracking
type CarrierTrackingRepository interface {
	// Create creates a new CarrierTracking
	Create(ctx context.Context, carrier_tracking *CarrierTracking) error

	// FindByID finds CarrierTracking by ID
	FindByID(ctx context.Context, id uuid.UUID) (*CarrierTracking, error)

	// Update updates CarrierTracking
	Update(ctx context.Context, carrier_tracking *CarrierTracking) error

	// Delete deletes CarrierTracking by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all CarrierTracking
	List(ctx context.Context, limit, offset int) ([]*CarrierTracking, error)
}
