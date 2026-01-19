package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// DeleteShopReplyInput represents input for DeleteShopReply
type DeleteShopReplyInput struct {
	ReplyId uuid.UUID
	ShopId uuid.UUID
}

// DeleteShopReplyUsecase defines the interface for DeleteShopReply
type DeleteShopReplyUsecase interface {
	Execute(ctx context.Context, input DeleteShopReplyInput) error
}

type delete_shop_replyUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewDeleteShopReplyUsecase creates a new instance
func NewDeleteShopReplyUsecase() DeleteShopReplyUsecase {
	return &delete_shop_replyUsecaseImpl{}
}

// Execute executes DeleteShopReply
func (u *delete_shop_replyUsecaseImpl) Execute(ctx context.Context, input DeleteShopReplyInput) error {
	// TODO: Implement business logic

	return nil
}
