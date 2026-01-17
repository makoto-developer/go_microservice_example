package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// writeWait はメッセージ書き込みのタイムアウト
	writeWait = 10 * time.Second

	// pongWait はpongメッセージのタイムアウト
	pongWait = 60 * time.Second

	// pingPeriod はpingメッセージの送信間隔
	pingPeriod = (pongWait * 9) / 10

	// maxMessageSize は最大メッセージサイズ
	maxMessageSize = 512 * 1024 // 512KB
)

// Client はWebSocketクライアント
type Client struct {
	ID     string
	UserID string
	RoomID string
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
}

// NewClient はClientを初期化
func NewClient(hub *Hub, conn *websocket.Conn, userID string, roomID string) *Client {
	return &Client{
		ID:     generateClientID(),
		UserID: userID,
		RoomID: roomID,
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
	}
}

// ReadPump はWebSocketからメッセージを読み取る
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, messageData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket error: %v", err)
			}
			break
		}

		// メッセージをパース
		var msg Message
		if err := json.Unmarshal(messageData, &msg); err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		// メッセージの検証
		if msg.RoomID != c.RoomID {
			log.Printf("invalid room_id: expected %s, got %s", c.RoomID, msg.RoomID)
			continue
		}

		msg.UserID = c.UserID
		msg.Timestamp = time.Now()

		// Hubに配信
		c.hub.broadcast <- &msg
	}
}

// WritePump はメッセージをWebSocketに書き込む
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hubがチャネルを閉じた
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// キューに溜まっているメッセージを追加送信
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// generateClientID はクライアントIDを生成
func generateClientID() string {
	// 簡易実装: UUIDを使用
	// 実際はより短いIDを生成することが多い
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString はランダム文字列を生成
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
