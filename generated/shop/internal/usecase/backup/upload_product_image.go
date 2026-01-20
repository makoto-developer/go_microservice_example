package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type UploadProductImageInput struct {
	ProductID    uuid.UUID
	ImageData    []byte
	DisplayOrder int
}

type UploadProductImageOutput struct {
	ImageID    uuid.UUID
	URL        string
	Thumbnails map[string]string
}

type UploadProductImageUsecase interface {
	Execute(ctx context.Context, input UploadProductImageInput) (UploadProductImageOutput, error)
}

type uploadProductImageUsecase struct {
	productRepo repository.ProductRepository
}

func NewUploadProductImageUsecase(productRepo repository.ProductRepository) UploadProductImageUsecase {
	return &uploadProductImageUsecase{productRepo: productRepo}
}

func (u *uploadProductImageUsecase) Execute(ctx context.Context, input UploadProductImageInput) (UploadProductImageOutput, error) {
	if _, err := u.productRepo.GetByID(ctx, input.ProductID); err != nil {
		return UploadProductImageOutput{}, err
	}

	count, err := u.productRepo.CountImages(ctx, input.ProductID)
	if err != nil {
		return UploadProductImageOutput{}, err
	}
	if count >= 5 {
		return UploadProductImageOutput{}, domain.ErrMaxImagesExceeded
	}

	if len(input.ImageData) > 5*1024*1024 {
		return UploadProductImageOutput{}, domain.ErrImageTooLarge
	}

	imageURL := fmt.Sprintf("https://example.com/products/%s/images/%s.jpg", input.ProductID, uuid.New())
	thumb200 := fmt.Sprintf("%s?size=200", imageURL)
	thumb400 := fmt.Sprintf("%s?size=400", imageURL)
	thumb800 := fmt.Sprintf("%s?size=800", imageURL)

	image := domain.NewProductImage(input.ProductID, imageURL, input.DisplayOrder, thumb200, thumb400, thumb800)

	if err := u.productRepo.AddImage(ctx, image); err != nil {
		return UploadProductImageOutput{}, err
	}

	return UploadProductImageOutput{
		ImageID: image.ID,
		URL:     image.URL,
		Thumbnails: map[string]string{
			"200": thumb200,
			"400": thumb400,
			"800": thumb800,
		},
	}, nil
}
