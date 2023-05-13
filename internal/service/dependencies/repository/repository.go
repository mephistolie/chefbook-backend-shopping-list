package repository

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
)

type ShoppingList interface {
	AddUser(userId uuid.UUID, messageId uuid.UUID) error
	ImportFirebaseProfile(userId uuid.UUID, shoppingList entity.ShoppingList, messageId uuid.UUID) error
	DeleteUser(userId uuid.UUID, messageId uuid.UUID) error

	GetShoppingList(userId uuid.UUID) (entity.ShoppingList, error)
	SetShoppingList(userId uuid.UUID, shoppingList entity.ShoppingList, lastVersion *int32) error
}
