package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/proto/customer-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CustomerServiceHandler struct {
	pb.UnimplementedCustomerServiceServer
	getCustomerProfileUsecase       usecase.GetCustomerProfileUsecase
	updateCustomerProfileUsecase    usecase.UpdateCustomerProfileUsecase
	registerAddressUsecase          usecase.RegisterAddressUsecase
	updateAddressUsecase            usecase.UpdateAddressUsecase
	deleteAddressUsecase            usecase.DeleteAddressUsecase
	addToCartUsecase                usecase.AddToCartUsecase
	getCartUsecase                  usecase.GetCartUsecase
	updateCartItemQuantityUsecase   usecase.UpdateCartItemQuantityUsecase
	removeFromCartUsecase           usecase.RemoveFromCartUsecase
	mergeGuestCartUsecase           usecase.MergeGuestCartUsecase
	addToFavoriteUsecase            usecase.AddToFavoriteUsecase
	listFavoritesUsecase            usecase.ListFavoritesUsecase
	removeFromFavoriteUsecase       usecase.RemoveFromFavoriteUsecase
	addPaymentMethodUsecase         usecase.AddPaymentMethodUsecase
	deletePaymentMethodUsecase      usecase.DeletePaymentMethodUsecase
	createReviewUsecase             usecase.CreateReviewUsecase
	updateReviewUsecase             usecase.UpdateReviewUsecase
}

func NewCustomerServiceHandler(
	getCustomerProfileUsecase usecase.GetCustomerProfileUsecase,
	updateCustomerProfileUsecase usecase.UpdateCustomerProfileUsecase,
	registerAddressUsecase usecase.RegisterAddressUsecase,
	updateAddressUsecase usecase.UpdateAddressUsecase,
	deleteAddressUsecase usecase.DeleteAddressUsecase,
	addToCartUsecase usecase.AddToCartUsecase,
	getCartUsecase usecase.GetCartUsecase,
	updateCartItemQuantityUsecase usecase.UpdateCartItemQuantityUsecase,
	removeFromCartUsecase usecase.RemoveFromCartUsecase,
	mergeGuestCartUsecase usecase.MergeGuestCartUsecase,
	addToFavoriteUsecase usecase.AddToFavoriteUsecase,
	listFavoritesUsecase usecase.ListFavoritesUsecase,
	removeFromFavoriteUsecase usecase.RemoveFromFavoriteUsecase,
	addPaymentMethodUsecase usecase.AddPaymentMethodUsecase,
	deletePaymentMethodUsecase usecase.DeletePaymentMethodUsecase,
	createReviewUsecase usecase.CreateReviewUsecase,
	updateReviewUsecase usecase.UpdateReviewUsecase,
) *CustomerServiceHandler {
	return &CustomerServiceHandler{
		getCustomerProfileUsecase:     getCustomerProfileUsecase,
		updateCustomerProfileUsecase:  updateCustomerProfileUsecase,
		registerAddressUsecase:        registerAddressUsecase,
		updateAddressUsecase:          updateAddressUsecase,
		deleteAddressUsecase:          deleteAddressUsecase,
		addToCartUsecase:              addToCartUsecase,
		getCartUsecase:                getCartUsecase,
		updateCartItemQuantityUsecase: updateCartItemQuantityUsecase,
		removeFromCartUsecase:         removeFromCartUsecase,
		mergeGuestCartUsecase:         mergeGuestCartUsecase,
		addToFavoriteUsecase:          addToFavoriteUsecase,
		listFavoritesUsecase:          listFavoritesUsecase,
		removeFromFavoriteUsecase:     removeFromFavoriteUsecase,
		addPaymentMethodUsecase:       addPaymentMethodUsecase,
		deletePaymentMethodUsecase:    deletePaymentMethodUsecase,
		createReviewUsecase:           createReviewUsecase,
		updateReviewUsecase:           updateReviewUsecase,
	}
}

func (h *CustomerServiceHandler) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	input := usecase.GetCustomerProfileInput{CustomerID: customerID}
	output, err := h.getCustomerProfileUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.GetProfileResponse{
		Customer: &pb.Customer{
			Id:                     output.Customer.ID.String(),
			UserId:                 output.Customer.UserID.String(),
			FirstName:              output.Customer.FirstName,
			LastName:               output.Customer.LastName,
			Phone:                  output.Customer.Phone,
			BirthDate:              formatDate(output.Customer.BirthDate),
			Gender:                 domainGenderToProto(output.Customer.Gender),
			ProfileImageUrl:        formatOptionalString(output.Customer.ProfileImageURL),
			ProfileThumbnail_100Url: formatOptionalString(output.Customer.ProfileThumbnail100URL),
			ProfileThumbnail_200Url: formatOptionalString(output.Customer.ProfileThumbnail200URL),
			CreatedAt:              timestampProto(output.Customer.CreatedAt),
			UpdatedAt:              timestampProto(output.Customer.UpdatedAt),
		},
		Message: "Profile retrieved successfully",
	}, nil
}

func (h *CustomerServiceHandler) UpdateProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	customerID, err := uuid.Parse(req.CustomerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid customer ID")
	}

	birthDate, err := parseDate(req.BirthDate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birth date")
	}

	input := usecase.UpdateCustomerProfileInput{
		CustomerID: customerID,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Phone:      req.Phone,
		BirthDate:  birthDate,
		Gender:     protoGenderToDomain(req.Gender),
	}

	output, err := h.updateCustomerProfileUsecase.Execute(ctx, input)
	if err != nil {
		return nil, mapDomainError(err)
	}

	return &pb.UpdateProfileResponse{
		Customer: &pb.Customer{
			Id:        output.Customer.ID.String(),
			UserId:    output.Customer.UserID.String(),
			FirstName: output.Customer.FirstName,
			LastName:  output.Customer.LastName,
			Phone:     output.Customer.Phone,
			UpdatedAt: timestampProto(output.Customer.UpdatedAt),
		},
		Message: "Profile updated successfully",
	}, nil
}

func mapDomainError(err error) error {
	switch err {
	case domain.ErrCustomerNotFound:
		return status.Error(codes.NotFound, "customer not found")
	case domain.ErrAddressNotFound:
		return status.Error(codes.NotFound, "address not found")
	case domain.ErrCartItemNotFound:
		return status.Error(codes.NotFound, "cart item not found")
	case domain.ErrAlreadyFavorited:
		return status.Error(codes.AlreadyExists, "already favorited")
	case domain.ErrPaymentMethodNotFound:
		return status.Error(codes.NotFound, "payment method not found")
	case domain.ErrReviewNotFound:
		return status.Error(codes.NotFound, "review not found")
	case domain.ErrInvalidQuantity:
		return status.Error(codes.InvalidArgument, "invalid quantity")
	case domain.ErrInvalidRating:
		return status.Error(codes.InvalidArgument, "invalid rating")
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
