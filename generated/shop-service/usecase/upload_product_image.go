package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop-service/domain"
)

type UploadProductImageInput struct {
	ProductID    uuid.UUID
	ImageURL     string
	DisplayOrder int
}

type UploadProductImageOutput struct {
	ImageID uuid.UUID
	Message string
}

type UploadProductImageUsecase interface {
	Execute(ctx context.Context, input UploadProductImageInput) (*UploadProductImageOutput, error)
}

type uploadProductImageUsecaseImpl struct {
	productImageRepo domain.ProductImageRepository
}

func NewUploadProductImageUsecase(productImageRepo domain.ProductImageRepository) UploadProductImageUsecase {
	return &uploadProductImageUsecaseImpl{
		productImageRepo: productImageRepo,
	}
}

func (u *uploadProductImageUsecaseImpl) Execute(ctx context.Context, input UploadProductImageInput) (*UploadProductImageOutput, error) {
	image := &domain.ProductImage{
		Id:              uuid.New(),
		ProductId:       input.ProductID,
		Url:             input.ImageURL,
		DisplayOrder:    input.DisplayOrder,
		Thumbnail200Url: input.ImageURL + "?w=200",
		Thumbnail400Url: input.ImageURL + "?w=400",
		Thumbnail800Url: input.ImageURL + "?w=800",
		CreatedAt:       time.Now(),
	}

	if err := u.productImageRepo.Create(ctx, image); err != nil {
		return nil, fmt.Errorf("failed to create product image: %w", err)
	}

	return &UploadProductImageOutput{
		ImageID: image.Id,
		Message: "Product image uploaded successfully",
	}, nil
}
