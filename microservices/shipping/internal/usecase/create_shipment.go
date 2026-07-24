package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/shipping/internal/repository"
)

type CreateShipmentInput struct {
	OrderID         uuid.UUID
	CustomerID      uuid.UUID
	ShippingAddress string
	Carrier         string
}

type CreateShipmentOutput struct {
	ShipmentID uuid.UUID
	Status     domain.ShipmentStatus
}

type CreateShipmentUsecase interface {
	Execute(ctx context.Context, input CreateShipmentInput) (CreateShipmentOutput, error)
}

type createShipmentUsecaseImpl struct {
	shipmentRepo repository.ShipmentRepository
}

func NewCreateShipmentUsecase(shipmentRepo repository.ShipmentRepository) CreateShipmentUsecase {
	return &createShipmentUsecaseImpl{
		shipmentRepo: shipmentRepo,
	}
}

func (u *createShipmentUsecaseImpl) Execute(ctx context.Context, input CreateShipmentInput) (CreateShipmentOutput, error) {
	shipment := domain.NewShipment(input.OrderID, input.CustomerID, input.ShippingAddress, input.Carrier)

	err := u.shipmentRepo.Create(ctx, shipment)
	if err != nil {
		return CreateShipmentOutput{}, err
	}

	return CreateShipmentOutput{
		ShipmentID: shipment.ID,
		Status:     shipment.Status,
	}, nil
}
