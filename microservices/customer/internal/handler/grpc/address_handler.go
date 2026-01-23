package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/proto/customer_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CustomerServiceHandler) RegisterAddress(ctx context.Context, req *pb.RegisterAddressRequest) (*pb.RegisterAddressResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	var addressLine2 *string
	if req.AddressLine2 != "" {
		addressLine2 = &req.AddressLine2
	}

	input := usecase.RegisterAddressInput{
		CustomerID:     customerID,
		AddressName:    req.AddressName,
		PostalCode:     req.PostalCode,
		Prefecture:     req.Prefecture,
		City:           req.City,
		AddressLine1:   req.AddressLine1,
		AddressLine2:   addressLine2,
		RecipientName:  req.RecipientName,
		RecipientPhone: req.RecipientPhone,
		IsDefault:      req.IsDefault,
	}

	output, err := h.registerAddressUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.RegisterAddressResponse{
		Message: "Address registered successfully",
		Address: &pb.Address{
			Id:         output.AddressID.String(),
			CustomerId: customerID.String(),
		},
	}, nil
}

func (h *CustomerServiceHandler) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (*pb.UpdateAddressResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	addressID, err := uuid.Parse(req.AddressId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address ID")
	}

	var addressLine2 *string
	if req.AddressLine2 != "" {
		addressLine2 = &req.AddressLine2
	}

	input := usecase.UpdateAddressInput{
		CustomerID:     customerID,
		AddressID:      addressID,
		AddressName:    req.AddressName,
		PostalCode:     req.PostalCode,
		Prefecture:     req.Prefecture,
		City:           req.City,
		AddressLine1:   req.AddressLine1,
		AddressLine2:   addressLine2,
		RecipientName:  req.RecipientName,
		RecipientPhone: req.RecipientPhone,
	}

	output, err := h.updateAddressUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.UpdateAddressResponse{
		Message: "Address updated successfully",
		Address: &pb.Address{
			Id:         output.Address.ID.String(),
			CustomerId: output.Address.CustomerID.String(),
		},
	}, nil
}

func (h *CustomerServiceHandler) DeleteAddress(ctx context.Context, req *pb.DeleteAddressRequest) (*pb.DeleteAddressResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	addressID, err := uuid.Parse(req.AddressId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid address ID")
	}

	input := usecase.DeleteAddressInput{
		CustomerID: customerID,
		AddressID:  addressID,
	}

	_, err = h.deleteAddressUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.DeleteAddressResponse{
		Success: true,
		Message: "Address deleted successfully",
	}, nil
}
