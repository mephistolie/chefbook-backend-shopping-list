package repository

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
)

type ShoppingList interface {
	AddUser(userId uuid.UUID, messageId uuid.UUID) error
	ImportFirebaseProfile(userId uuid.UUID, purchases []entity.Purchase, messageId uuid.UUID) error
	DeleteUser(userId uuid.UUID, messageId uuid.UUID) error

	GetShoppingList(userId uuid.UUID) (entity.ShoppingList, error)
	SetShoppingList(userId uuid.UUID, purchases []entity.Purchase, lastVersion *int32) (int32, error)
}
