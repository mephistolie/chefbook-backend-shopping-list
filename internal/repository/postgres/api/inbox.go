package api

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
)

type Inbox interface {
	GetPendingMessages() ([]entity.MessageData, error)
	AddMessage(msg entity.MessageData) error
	CheckMessageProcessed(eventId uuid.UUID) error
}
