package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// AddForbiddenWordInput represents input for AddForbiddenWord
type AddForbiddenWordInput struct {
	AdminId uuid.UUID
	Word string
	Context WordContext
	Severity Severity
}

// AddForbiddenWordUsecase defines the interface for AddForbiddenWord
type AddForbiddenWordUsecase interface {
	Execute(ctx context.Context, input AddForbiddenWordInput) error
}

type add_forbidden_wordUsecaseImpl struct {
	// TODO: Add repository dependencies
}

// NewAddForbiddenWordUsecase creates a new instance
func NewAddForbiddenWordUsecase() AddForbiddenWordUsecase {
	return &add_forbidden_wordUsecaseImpl{}
}

// Execute executes AddForbiddenWord
func (u *add_forbidden_wordUsecaseImpl) Execute(ctx context.Context, input AddForbiddenWordInput) error {
	// TODO: Implement business logic

	return nil
}
