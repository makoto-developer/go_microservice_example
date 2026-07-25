package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/customer/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CustomerServiceHandler) RegisterPaymentMethod(ctx context.Context, req *pb.RegisterPaymentMethodRequest) (*pb.RegisterPaymentMethodResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	input := usecase.AddPaymentMethodInput{
		CustomerID:            customerID,
		StripePaymentMethodID: req.StripePaymentMethodId,
		IsDefault:             req.IsDefault,
	}

	output, err := h.addPaymentMethodUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.RegisterPaymentMethodResponse{
		Message: "Payment method registered successfully",
		PaymentMethod: &pb.PaymentMethod{
			Id:                    output.PaymentMethodID.String(),
			CustomerId:            customerID.String(),
			StripePaymentMethodId: req.StripePaymentMethodId,
			IsDefault:             req.IsDefault,
		},
	}, nil
}

func (h *CustomerServiceHandler) DeletePaymentMethod(ctx context.Context, req *pb.DeletePaymentMethodRequest) (*pb.DeletePaymentMethodResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	paymentMethodID, err := uuid.Parse(req.PaymentMethodId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payment method ID")
	}

	input := usecase.DeletePaymentMethodInput{
		CustomerID:      customerID,
		PaymentMethodID: paymentMethodID,
	}

	_, err = h.deletePaymentMethodUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.DeletePaymentMethodResponse{
		Success: true,
		Message: "Payment method deleted successfully",
	}, nil
}
