package domain

import "errors"

var (
	// Shop errors
	ErrShopNotFound          = errors.New("shop not found")
	ErrShopAlreadyExists     = errors.New("owner already has a shop")
	ErrShopNotApproved       = errors.New("shop not approved")
	ErrInvalidShopData       = errors.New("invalid shop data")
	ErrUnauthorizedAccess    = errors.New("unauthorized access")

	// Product errors
	ErrProductNotFound       = errors.New("product not found")
	ErrInvalidProductData    = errors.New("invalid product data")
	ErrInsufficientStock     = errors.New("insufficient stock")
	ErrMaxImagesExceeded     = errors.New("maximum images exceeded")
	ErrImageTooLarge         = errors.New("image too large")
	ErrInvalidImageFormat    = errors.New("invalid image format")
	ErrDuplicateSKU          = errors.New("duplicate SKU")

	// Order errors
	ErrOrderNotFound         = errors.New("order not found")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrInvalidDateRange      = errors.New("invalid date range")
	ErrNoDataFound           = errors.New("no data found")
)
