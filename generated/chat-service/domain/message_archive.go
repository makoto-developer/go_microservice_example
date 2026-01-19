package domain

import (
	"time"
	"github.com/google/uuid"
)

// MessageArchive represents MessageArchive
type MessageArchive struct {
	Id uuid.UUID `db:"id" json:"id"`
	RoomId uuid.UUID `db:"room_id" json:"room_id"`
	MessageCount int `db:"message_count" json:"message_count"`
	ArchiveUrl string `db:"archive_url" json:"archive_url"`
	ArchivedFrom time.Time `db:"archived_from" json:"archived_from"`
	ArchivedTo time.Time `db:"archived_to" json:"archived_to"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// NewMessageArchive creates a new MessageArchive instance
func NewMessageArchive() *MessageArchive {
	return &MessageArchive{}
}
