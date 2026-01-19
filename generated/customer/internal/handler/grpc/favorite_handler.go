package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/proto/customer-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CustomerServiceHandler) AddToFavorite(ctx context.Context, req *pb.AddToFavoriteRequest) (*pb.AddToFavoriteResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid product ID")
	}

	input := usecase.AddToFavoriteInput{
		CustomerID:      customerID,
		ProductID:       productID,
		NotifyOnRestock: req.NotifyOnRestock,
	}

	output, err := h.addToFavoriteUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.AddToFavoriteResponse{
		Message: "Product added to favorites successfully",
		Favorite: &pb.Favorite{
			Id:              output.FavoriteID.String(),
			CustomerId:      customerID.String(),
			ProductId:       productID.String(),
			NotifyOnRestock: req.NotifyOnRestock,
		},
	}, nil
}

func (h *CustomerServiceHandler) GetFavorites(ctx context.Context, req *pb.GetFavoritesRequest) (*pb.GetFavoritesResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	input := usecase.ListFavoritesInput{CustomerID: customerID}
	output, err := h.listFavoritesUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	var favorites []*pb.Favorite
	for _, fav := range output.Favorites {
		favorites = append(favorites, &pb.Favorite{
			Id:              fav.ID.String(),
			CustomerId:      fav.CustomerID.String(),
			ProductId:       fav.ProductID.String(),
			NotifyOnRestock: fav.NotifyOnRestock,
			CreatedAt:       timestampProto(fav.CreatedAt),
		})
	}

	return &pb.GetFavoritesResponse{
		Favorites:  favorites,
		TotalCount: int32(len(favorites)),
	}, nil
}

func (h *CustomerServiceHandler) RemoveFromFavorite(ctx context.Context, req *pb.RemoveFromFavoriteRequest) (*pb.RemoveFromFavoriteResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	favoriteID, err := uuid.Parse(req.FavoriteId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid favorite ID")
	}

	input := usecase.RemoveFromFavoriteInput{
		CustomerID: customerID,
		ProductID:  favoriteID,
	}

	_, err = h.removeFromFavoriteUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.RemoveFromFavoriteResponse{
		Success: true,
		Message: "Product removed from favorites successfully",
	}, nil
}
