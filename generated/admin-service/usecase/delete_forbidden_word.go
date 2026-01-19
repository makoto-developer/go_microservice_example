package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// DeleteForbiddenWordInput represents input for DeleteForbiddenWord
type DeleteForbiddenWordInput struct {
	AdminId uuid.UUID
	WordId uuid.UUID
}

// DeleteForbiddenWordUsecase defines the interface for DeleteForbiddenWord
type DeleteForbiddenWordUsecase interface {
	Execute(ctx context.Context, input DeleteForbiddenWordInput) error
}

type delete_forbidden_wordUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewDeleteForbiddenWordUsecase creates a new instance
func NewDeleteForbiddenWordUsecase() DeleteForbiddenWordUsecase {
	return &delete_forbidden_wordUsecaseImpl{}
}

// Execute executes DeleteForbiddenWord
func (u *delete_forbidden_wordUsecaseImpl) Execute(ctx context.Context, input DeleteForbiddenWordInput) error {
	// TODO: Implement business logic

	return nil
}
