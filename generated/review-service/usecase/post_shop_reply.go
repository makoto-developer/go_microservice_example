package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// PostShopReplyInput represents input for PostShopReply
type PostShopReplyInput struct {
	ReviewId uuid.UUID
	ShopId uuid.UUID
	Content string
}

// PostShopReplyUsecase defines the interface for PostShopReply
type PostShopReplyUsecase interface {
	Execute(ctx context.Context, input PostShopReplyInput) error
}

type post_shop_replyUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewPostShopReplyUsecase creates a new instance
func NewPostShopReplyUsecase() PostShopReplyUsecase {
	return &post_shop_replyUsecaseImpl{}
}

// Execute executes PostShopReply
func (u *post_shop_replyUsecaseImpl) Execute(ctx context.Context, input PostShopReplyInput) error {
	// TODO: Implement business logic

	return nil
}
