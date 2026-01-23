package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/proto/customer_service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CustomerServiceHandler) AddToCart(ctx context.Context, req *pb.AddToCartRequest) (*pb.AddToCartResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	variationID, err := parseOptionalUUID(req.VariationId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid variation ID")
	}

	input := usecase.AddToCartInput{
		CustomerID:  customerID,
		ProductID:   productID,
		VariationID: variationID,
		Quantity:    int(req.Quantity),
	}

	output, err := h.addToCartUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.AddToCartResponse{
		Message: "Item added to cart successfully",
		CartItem: &pb.CartItem{
			Id:         output.CartItemID.String(),
			CustomerId: customerID.String(),
			ProductId:  productID.String(),
			Quantity:   int32(output.TotalQuantity),
		},
	}, nil
}

func (h *CustomerServiceHandler) GetCart(ctx context.Context, req *pb.GetCartRequest) (*pb.GetCartResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	input := usecase.GetCartInput{CustomerID: customerID}
	output, err := h.getCartUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	var items []*pb.CartItem
	totalQty := int32(0)
	for _, item := range output.CartItems {
		items = append(items, &pb.CartItem{
			Id:         item.ID.String(),
			CustomerId: item.CustomerID.String(),
			ProductId:  item.ProductID.String(),
			Quantity:   int32(item.Quantity),
			CreatedAt:  timestampProto(item.CreatedAt),
		})
		totalQty += int32(item.Quantity)
	}

	return &pb.GetCartResponse{
		CartItems:     items,
		TotalQuantity: totalQty,
	}, nil
}

func (h *CustomerServiceHandler) UpdateCartItemQuantity(ctx context.Context, req *pb.UpdateCartItemQuantityRequest) (*pb.UpdateCartItemQuantityResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	cartItemID, err := uuid.Parse(req.CartItemId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid cart item ID")
	}

	input := usecase.UpdateCartItemQuantityInput{
		CustomerID: customerID,
		CartItemID: cartItemID,
		Quantity:   int(req.Quantity),
	}

	output, err := h.updateCartItemQuantityUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.UpdateCartItemQuantityResponse{
		Message: "Cart item quantity updated successfully",
		CartItem: &pb.CartItem{
			Id:       output.CartItem.ID.String(),
			Quantity: int32(output.CartItem.Quantity),
		},
	}, nil
}

func (h *CustomerServiceHandler) RemoveFromCart(ctx context.Context, req *pb.RemoveFromCartRequest) (*pb.RemoveFromCartResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	cartItemID, err := uuid.Parse(req.CartItemId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid cart item ID")
	}

	input := usecase.RemoveFromCartInput{
		CustomerID: customerID,
		CartItemID: cartItemID,
	}

	_, err = h.removeFromCartUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.RemoveFromCartResponse{
		Success: true,
		Message: "Item removed from cart successfully",
	}, nil
}

func (h *CustomerServiceHandler) MergeGuestCart(ctx context.Context, req *pb.MergeGuestCartRequest) (*pb.MergeGuestCartResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	input := usecase.MergeGuestCartInput{
		CustomerID: customerID,
		SessionID:  req.SessionId,
	}

	_, err = h.mergeGuestCartUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.MergeGuestCartResponse{
		Message: "Guest cart merged successfully",
	}, nil
}
