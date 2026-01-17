package chat

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Hub はWebSocket接続を管理
type Hub struct {
	// rooms はチャットルームのマップ
	rooms map[string]*Room

	// register はクライアント登録チャネル
	register chan *Client

	// unregister はクライアント登録解除チャネル
	unregister chan *Client

	// broadcast はメッセージ配信チャネル
	broadcast chan *Message

	mu sync.RWMutex
}

// NewHub はHubを初期化
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]*Room),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 256),
	}
}

// Run はHubを実行
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)
		}
	}
}

// registerClient はクライアントを登録
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.rooms[client.RoomID]
	if !exists {
		room = NewRoom(client.RoomID)
		h.rooms[client.RoomID] = room
	}

	room.AddClient(client)

	// 参加通知を送信
	joinMsg := &Message{
		Type:      MessageTypeSystem,
		RoomID:    client.RoomID,
		UserID:    client.UserID,
		Content:   "joined the room",
		Timestamp: time.Now(),
	}
	h.broadcast <- joinMsg
}

// unregisterClient はクライアントを登録解除
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, exists := h.rooms[client.RoomID]
	if !exists {
		return
	}

	room.RemoveClient(client)

	// ルームが空になったら削除
	if room.IsEmpty() {
		delete(h.rooms, client.RoomID)
	}

	// 退出通知を送信
	leaveMsg := &Message{
		Type:      MessageTypeSystem,
		RoomID:    client.RoomID,
		UserID:    client.UserID,
		Content:   "left the room",
		Timestamp: time.Now(),
	}
	h.broadcast <- leaveMsg
}

// broadcastMessage はメッセージを配信
func (h *Hub) broadcastMessage(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, exists := h.rooms[message.RoomID]
	if !exists {
		return
	}

	room.Broadcast(message)
}

// Room はチャットルーム
type Room struct {
	ID      string
	clients map[string]*Client
	mu      sync.RWMutex
}

// NewRoom はRoomを初期化
func NewRoom(id string) *Room {
	return &Room{
		ID:      id,
		clients: make(map[string]*Client),
	}
}

// AddClient はクライアントを追加
func (r *Room) AddClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[client.ID] = client
}

// RemoveClient はクライアントを削除
func (r *Room) RemoveClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, client.ID)
}

// Broadcast はメッセージを全クライアントに配信
func (r *Room) Broadcast(message *Message) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	for _, client := range r.clients {
		select {
		case client.send <- data:
		default:
			// 送信バッファが詰まっている場合はクライアントを切断
			close(client.send)
			delete(r.clients, client.ID)
		}
	}
}

// IsEmpty はルームが空かどうか
func (r *Room) IsEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients) == 0
}

// GetOnlineUsers はオンラインユーザー一覧を取得
func (r *Room) GetOnlineUsers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]string, 0, len(r.clients))
	for _, client := range r.clients {
		users = append(users, client.UserID)
	}
	return users
}

// MessageType はメッセージタイプ
type MessageType string

const (
	MessageTypeText     MessageType = "text"
	MessageTypeImage    MessageType = "image"
	MessageTypeFile     MessageType = "file"
	MessageTypeSystem   MessageType = "system"
	MessageTypeTyping   MessageType = "typing"
	MessageTypePresence MessageType = "presence"
)

// Message はチャットメッセージ
type Message struct {
	ID        string      `json:"id"`
	Type      MessageType `json:"type"`
	RoomID    string      `json:"room_id"`
	UserID    string      `json:"user_id"`
	Content   string      `json:"content"`
	FileURL   string      `json:"file_url,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// NewMessage は新しいメッセージを作成
func NewMessage(msgType MessageType, roomID string, userID string, content string) *Message {
	return &Message{
		ID:        uuid.New().String(),
		Type:      msgType,
		RoomID:    roomID,
		UserID:    userID,
		Content:   content,
		Timestamp: time.Now(),
	}
}
