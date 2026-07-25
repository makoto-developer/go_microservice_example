package domain

import "errors"

var (
	// Customer errors
	ErrCustomerNotFound   = errors.New("customer not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrImageTooLarge      = errors.New("image too large")
	ErrInvalidImageFormat = errors.New("invalid image format")

	// Address errors
	ErrAddressNotFound            = errors.New("address not found")
	ErrInvalidPostalCode          = errors.New("invalid postal code")
	ErrCannotDeleteDefaultAddress = errors.New("cannot delete default address")
	ErrPostalCodeNotFound         = errors.New("postal code not found")
	ErrExternalAPIError           = errors.New("external API error")
	ErrUnauthorizedAccess         = errors.New("unauthorized access")

	// Cart errors
	ErrCartItemNotFound  = errors.New("cart item not found")
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrInvalidQuantity   = errors.New("invalid quantity")

	// Favorite errors
	ErrFavoriteNotFound = errors.New("favorite not found")
	ErrAlreadyFavorited = errors.New("already favorited")

	// Payment Method errors
	ErrPaymentMethodNotFound = errors.New("payment method not found")
	ErrInvalidCardData       = errors.New("invalid card data")

	// Review errors
	ErrReviewNotFound    = errors.New("review not found")
	ErrReviewNotEditable = errors.New("review not editable")
	ErrInvalidRating     = errors.New("invalid rating")
	ErrAlreadyReviewed   = errors.New("already reviewed")
)
