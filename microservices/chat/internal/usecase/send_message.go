package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/microservices/chat/internal/domain"
	"github.com/makoto-developer/go_microservice_example/microservices/chat/internal/repository"
)

type SendMessageInput struct {
	CustomerID uuid.UUID
	ShopID     uuid.UUID
	SenderID   uuid.UUID
	SenderType string
	Content    string
}

type SendMessageOutput struct {
	MessageID uuid.UUID
	RoomID    uuid.UUID
}

type SendMessageUsecase interface {
	Execute(ctx context.Context, input SendMessageInput) (SendMessageOutput, error)
}

type sendMessageUsecaseImpl struct {
	chatRepo repository.ChatRepository
}

func NewSendMessageUsecase(chatRepo repository.ChatRepository) SendMessageUsecase {
	return &sendMessageUsecaseImpl{
		chatRepo: chatRepo,
	}
}

func (u *sendMessageUsecaseImpl) Execute(ctx context.Context, input SendMessageInput) (SendMessageOutput, error) {
	// Get or create chat room
	room, err := u.chatRepo.GetRoomByCustomerAndShop(ctx, input.CustomerID, input.ShopID)
	if err != nil {
		// Create new room if not exists
		room = domain.NewChatRoom(input.CustomerID, input.ShopID)
		err = u.chatRepo.CreateRoom(ctx, room)
		if err != nil {
			return SendMessageOutput{}, err
		}
	}

	// Create message
	message := domain.NewMessage(room.ID, input.SenderID, input.SenderType, domain.MessageTypeText, input.Content)

	err = u.chatRepo.CreateMessage(ctx, message)
	if err != nil {
		return SendMessageOutput{}, err
	}

	return SendMessageOutput{
		MessageID: message.ID,
		RoomID:    room.ID,
	}, nil
}
