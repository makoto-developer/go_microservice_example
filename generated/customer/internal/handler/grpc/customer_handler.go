package grpc

import (
	"errors"

	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/customer/internal/usecase"
	pb "github.com/makoto-developer/go_microservice_example/proto/customer-service/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CustomerServiceHandler implements the gRPC customer service
type CustomerServiceHandler struct {
	pb.UnimplementedCustomerServiceServer

	// Profile usecases
	getProfileUsecase    usecase.GetCustomerProfileUsecase
	updateProfileUsecase usecase.UpdateCustomerProfileUsecase

	// Address usecases
	registerAddressUsecase usecase.RegisterAddressUsecase
	updateAddressUsecase   usecase.UpdateAddressUsecase
	deleteAddressUsecase   usecase.DeleteAddressUsecase

	// Cart usecases
	addToCartUsecase              usecase.AddToCartUsecase
	getCartUsecase                usecase.GetCartUsecase
	updateCartItemQuantityUsecase usecase.UpdateCartItemQuantityUsecase
	removeFromCartUsecase         usecase.RemoveFromCartUsecase
	mergeGuestCartUsecase         usecase.MergeGuestCartUsecase

	// Favorite usecases
	addToFavoriteUsecase     usecase.AddToFavoriteUsecase
	listFavoritesUsecase     usecase.ListFavoritesUsecase
	removeFromFavoriteUsecase usecase.RemoveFromFavoriteUsecase

	// Payment method usecases
	addPaymentMethodUsecase    usecase.AddPaymentMethodUsecase
	deletePaymentMethodUsecase usecase.DeletePaymentMethodUsecase

	// Review usecases
	createReviewUsecase usecase.CreateReviewUsecase
	updateReviewUsecase usecase.UpdateReviewUsecase
}

// NewCustomerServiceHandler creates a new customer service handler
func NewCustomerServiceHandler(
	getProfileUsecase usecase.GetCustomerProfileUsecase,
	updateProfileUsecase usecase.UpdateCustomerProfileUsecase,
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
		getProfileUsecase:             getProfileUsecase,
		updateProfileUsecase:          updateProfileUsecase,
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

// mapDomainError maps domain errors to gRPC status errors
func mapDomainError(err error) error {
	if errors.Is(err, domain.ErrCustomerNotFound) {
		return status.Error(codes.NotFound, "customer not found")
	}
	if errors.Is(err, domain.ErrAddressNotFound) {
		return status.Error(codes.NotFound, "address not found")
	}
	if errors.Is(err, domain.ErrAlreadyFavorited) {
		return status.Error(codes.AlreadyExists, "already favorited")
	}
	return status.Error(codes.Internal, err.Error())
}
