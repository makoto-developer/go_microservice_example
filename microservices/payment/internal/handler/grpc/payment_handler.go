package grpc

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/repository"
	"github.com/makoto-developer/go_microservice_example/microservices/payment/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/payment/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentServiceHandler struct {
	pb.UnimplementedPaymentServiceServer
	processPaymentUsecase usecase.ProcessPaymentUsecase
	paymentRepo           repository.PaymentRepository
}

func NewPaymentServiceHandler(
	processPaymentUsecase usecase.ProcessPaymentUsecase,
	paymentRepo repository.PaymentRepository,
) *PaymentServiceHandler {
	return &PaymentServiceHandler{
		processPaymentUsecase: processPaymentUsecase,
		paymentRepo:           paymentRepo,
	}
}

func (h *PaymentServiceHandler) CreatePaymentIntent(ctx context.Context, req *pb.CreatePaymentIntentRequest) (*pb.CreatePaymentIntentResponse, error) {
	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}

	amount, err := strconv.Atoi(req.GetAmount())
	if err != nil || amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid amount")
	}

	input := usecase.ProcessPaymentInput{
		OrderID:       orderID,
		Amount:        amount,
		PaymentMethod: "credit_card", // モック決済のため支払い手段は固定(payment_method_id の解決は未実装)
	}

	output, err := h.processPaymentUsecase.Execute(ctx, input)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreatePaymentIntentResponse{
		PaymentId:    output.PaymentID.String(),
		ClientSecret: "secret_" + output.PaymentID.String(),
		Message:      "Payment intent created successfully",
	}, nil
}

func (h *PaymentServiceHandler) ConfirmPayment(ctx context.Context, req *pb.ConfirmPaymentRequest) (*pb.ConfirmPaymentResponse, error) {
	paymentID, err := uuid.Parse(req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payment ID")
	}

	payment, err := h.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if payment == nil {
		return nil, status.Error(codes.NotFound, "payment not found")
	}

	return &pb.ConfirmPaymentResponse{
		Success: true,
		Status:  pb.PaymentStatus_PAYMENT_STATUS_SUCCEEDED,
		Message: "Payment confirmed successfully",
	}, nil
}

func (h *PaymentServiceHandler) GetPaymentStatus(ctx context.Context, req *pb.GetPaymentStatusRequest) (*pb.GetPaymentStatusResponse, error) {
	paymentID, err := uuid.Parse(req.PaymentId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid payment ID")
	}

	payment, err := h.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if payment == nil {
		return nil, status.Error(codes.NotFound, "payment not found")
	}

	return &pb.GetPaymentStatusResponse{
		Status:        pb.PaymentStatus_PAYMENT_STATUS_SUCCEEDED,
		TransactionId: payment.TransactionID,
	}, nil
}
