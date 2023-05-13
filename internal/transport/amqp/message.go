package amqp

import "github.com/google/uuid"

type MessageData struct {
	Id   uuid.UUID
	Type string
	Body []byte
}
