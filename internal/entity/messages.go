package entity

import "github.com/google/uuid"

type MessageData struct {
	EventId uuid.UUID
	Type    string
	Body    []byte
}
