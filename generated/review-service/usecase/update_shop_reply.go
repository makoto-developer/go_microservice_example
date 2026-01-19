package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// UpdateShopReplyInput represents input for UpdateShopReply
type UpdateShopReplyInput struct {
	ReplyId uuid.UUID
	ShopId uuid.UUID
	Content string
}

// UpdateShopReplyUsecase defines the interface for UpdateShopReply
type UpdateShopReplyUsecase interface {
	Execute(ctx context.Context, input UpdateShopReplyInput) error
}

type update_shop_replyUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewUpdateShopReplyUsecase creates a new instance
func NewUpdateShopReplyUsecase() UpdateShopReplyUsecase {
	return &update_shop_replyUsecaseImpl{}
}

// Execute executes UpdateShopReply
func (u *update_shop_replyUsecaseImpl) Execute(ctx context.Context, input UpdateShopReplyInput) error {
	// TODO: Implement business logic

	return nil
}
