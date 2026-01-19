package domain

import (
	"context"

	"github.com/google/uuid"
)

// ShopReplyRepository defines repository interface for ShopReply
type ShopReplyRepository interface {
	// Create creates a new ShopReply
	Create(ctx context.Context, shop_reply *ShopReply) error

	// FindByID finds ShopReply by ID
	FindByID(ctx context.Context, id uuid.UUID) (*ShopReply, error)

	// Update updates ShopReply
	Update(ctx context.Context, shop_reply *ShopReply) error

	// Delete deletes ShopReply by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List lists all ShopReply
	List(ctx context.Context, limit, offset int) ([]*ShopReply, error)

	// FindByReviewId finds ShopReply by review_id
	FindByReviewId(ctx context.Context, review_id uuid.UUID) (*ShopReply, error)
}
