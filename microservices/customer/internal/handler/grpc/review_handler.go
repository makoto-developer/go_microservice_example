package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/customer/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/microservices/customer/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CustomerServiceHandler) PostReview(ctx context.Context, req *pb.PostReviewRequest) (*pb.PostReviewResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	orderID, err := uuid.Parse(req.OrderId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid order ID")
	}

	input := usecase.CreateReviewInput{
		CustomerID: customerID,
		ProductID:  productID,
		OrderID:    orderID,
		Rating:     int(req.Rating),
		ReviewText: req.ReviewText,
	}

	output, err := h.createReviewUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.PostReviewResponse{
		Message: "Review posted successfully",
		Review: &pb.Review{
			Id:         output.ReviewID.String(),
			CustomerId: customerID.String(),
			ProductId:  productID.String(),
			OrderId:    orderID.String(),
			Rating:     req.Rating,
			ReviewText: req.ReviewText,
		},
	}, nil
}

func (h *CustomerServiceHandler) UpdateReview(ctx context.Context, req *pb.UpdateReviewRequest) (*pb.UpdateReviewResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	reviewID, err := uuid.Parse(req.ReviewId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid review ID")
	}

	input := usecase.UpdateReviewInput{
		CustomerID: customerID,
		ReviewID:   reviewID,
		Rating:     int(req.Rating),
		ReviewText: req.ReviewText,
	}

	output, err := h.updateReviewUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.UpdateReviewResponse{
		Message: "Review updated successfully",
		Review: &pb.Review{
			Id:         output.Review.ID.String(),
			Rating:     int32(output.Review.Rating),
			ReviewText: output.Review.ReviewText,
			UpdatedAt:  timestampProto(output.Review.UpdatedAt),
		},
	}, nil
}
