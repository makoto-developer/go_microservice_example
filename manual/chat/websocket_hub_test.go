package chat

import (
	"testing"
	"time"
)

func TestHub_RegisterClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client := &Client{
		ID:     "test-client-1",
		UserID: "user-1",
		RoomID: "room-1",
		hub:    hub,
		send:   make(chan []byte, 256),
	}

	hub.register <- client

	// Wait for registration
	time.Sleep(100 * time.Millisecond)

	if len(hub.rooms) != 1 {
		t.Errorf("Expected 1 room, got %d", len(hub.rooms))
	}
}

func TestHub_BroadcastMessage(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	client1 := &Client{
		ID:     "test-client-1",
		UserID: "user-1",
		RoomID: "room-1",
		hub:    hub,
		send:   make(chan []byte, 256),
	}

	client2 := &Client{
		ID:     "test-client-2",
		UserID: "user-2",
		RoomID: "room-1",
		hub:    hub,
		send:   make(chan []byte, 256),
	}

	hub.register <- client1
	hub.register <- client2

	// Wait for registration
	time.Sleep(100 * time.Millisecond)

	msg := NewMessage(MessageTypeText, "room-1", "user-1", "Hello, World!")
	hub.broadcast <- msg

	// Wait for broadcast
	time.Sleep(100 * time.Millisecond)

	// Both clients should receive the message
	select {
	case <-client1.send:
		// Message received
	case <-time.After(1 * time.Second):
		t.Error("Client 1 did not receive message")
	}

	select {
	case <-client2.send:
		// Message received
	case <-time.After(1 * time.Second):
		t.Error("Client 2 did not receive message")
	}
}

func TestRoom_GetOnlineUsers(t *testing.T) {
	room := NewRoom("test-room")

	client1 := &Client{
		ID:     "client-1",
		UserID: "user-1",
		RoomID: "test-room",
	}

	client2 := &Client{
		ID:     "client-2",
		UserID: "user-2",
		RoomID: "test-room",
	}

	room.AddClient(client1)
	room.AddClient(client2)

	users := room.GetOnlineUsers()
	if len(users) != 2 {
		t.Errorf("Expected 2 online users, got %d", len(users))
	}

	room.RemoveClient(client1)

	users = room.GetOnlineUsers()
	if len(users) != 1 {
		t.Errorf("Expected 1 online user after removal, got %d", len(users))
	}
}

func TestRoom_IsEmpty(t *testing.T) {
	room := NewRoom("test-room")

	if !room.IsEmpty() {
		t.Error("Expected empty room")
	}

	client := &Client{
		ID:     "client-1",
		UserID: "user-1",
		RoomID: "test-room",
	}

	room.AddClient(client)

	if room.IsEmpty() {
		t.Error("Expected non-empty room")
	}

	room.RemoveClient(client)

	if !room.IsEmpty() {
		t.Error("Expected empty room after removal")
	}
}

func TestNewMessage(t *testing.T) {
	msg := NewMessage(MessageTypeText, "room-1", "user-1", "Test message")

	if msg.ID == "" {
		t.Error("Expected non-empty message ID")
	}

	if msg.Type != MessageTypeText {
		t.Errorf("Expected type text, got %s", msg.Type)
	}

	if msg.Content != "Test message" {
		t.Errorf("Expected content 'Test message', got %s", msg.Content)
	}
}
