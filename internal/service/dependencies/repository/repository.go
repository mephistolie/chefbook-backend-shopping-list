package repository

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
)

type ShoppingList interface {
	AddUser(userId uuid.UUID) error
	DeleteUser(userId uuid.UUID) error

	GetShoppingList(userId uuid.UUID) (entity.ShoppingList, error)
	SetShoppingList(userId uuid.UUID, shoppingList entity.ShoppingList) error
}
