package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/chat/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/chat/internal/usecase"
)

type mockChatRepository struct {
	createRoomFunc               func(ctx context.Context, room *domain.ChatRoom) error
	getRoomByCustomerAndShopFunc func(ctx context.Context, customerID, shopID uuid.UUID) (*domain.ChatRoom, error)
	createMessageFunc            func(ctx context.Context, message *domain.Message) error
}

func (m *mockChatRepository) CreateRoom(ctx context.Context, room *domain.ChatRoom) error {
	if m.createRoomFunc != nil {
		return m.createRoomFunc(ctx, room)
	}
	return nil
}

func (m *mockChatRepository) GetRoomByID(ctx context.Context, id uuid.UUID) (*domain.ChatRoom, error) {
	return nil, nil
}

func (m *mockChatRepository) GetRoomByCustomerAndShop(ctx context.Context, customerID, shopID uuid.UUID) (*domain.ChatRoom, error) {
	if m.getRoomByCustomerAndShopFunc != nil {
		return m.getRoomByCustomerAndShopFunc(ctx, customerID, shopID)
	}
	return nil, errors.New("not found")
}

func (m *mockChatRepository) GetRoomsByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*domain.ChatRoom, error) {
	return nil, nil
}

func (m *mockChatRepository) CreateMessage(ctx context.Context, message *domain.Message) error {
	if m.createMessageFunc != nil {
		return m.createMessageFunc(ctx, message)
	}
	return nil
}

func (m *mockChatRepository) GetMessagesByRoomID(ctx context.Context, roomID uuid.UUID) ([]*domain.Message, error) {
	return nil, nil
}

func (m *mockChatRepository) MarkMessageRead(ctx context.Context, messageID uuid.UUID) error {
	return nil
}

func TestSendMessageUsecase_NewRoom(t *testing.T) {
	customerID := uuid.New()
	shopID := uuid.New()
	senderID := customerID

	repo := &mockChatRepository{
		getRoomByCustomerAndShopFunc: func(ctx context.Context, cid, sid uuid.UUID) (*domain.ChatRoom, error) {
			return nil, errors.New("not found") // Room doesn't exist
		},
		createRoomFunc: func(ctx context.Context, room *domain.ChatRoom) error {
			if room.CustomerID != customerID {
				t.Errorf("expected customer ID %v, got %v", customerID, room.CustomerID)
			}
			return nil
		},
		createMessageFunc: func(ctx context.Context, message *domain.Message) error {
			if message.Content != "Hello" {
				t.Errorf("expected content 'Hello', got %v", message.Content)
			}
			return nil
		},
	}

	uc := usecase.NewSendMessageUsecase(repo)

	input := usecase.SendMessageInput{
		CustomerID: customerID,
		ShopID:     shopID,
		SenderID:   senderID,
		SenderType: "customer",
		Content:    "Hello",
	}

	output, err := uc.Execute(context.Background(), input)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if output.MessageID == uuid.Nil {
		t.Error("expected message ID to be generated")
	}
}
