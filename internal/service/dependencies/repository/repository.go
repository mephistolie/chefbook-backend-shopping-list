package repository

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
	"time"
)

type ShoppingList interface {
	CreatePersonalShoppingList(userId uuid.UUID, messageId uuid.UUID) error
	ImportFirebaseShoppingList(shoppingListId uuid.UUID, purchases []entity.Purchase, messageId uuid.UUID) error
	DeleteUserShoppingLists(userId uuid.UUID, messageId uuid.UUID) error

	GetShoppingLists(userId uuid.UUID) ([]entity.ShoppingListInfo, error)
	SetShoppingListName(shoppingListId, userId uuid.UUID, name *string) error
	CreateSharedShoppingList(userId uuid.UUID, shoppingListId *uuid.UUID, name *string) (uuid.UUID, error)
	GetShoppingList(shoppingListId uuid.UUID) (entity.ShoppingList, error)
	GetPersonalShoppingListId(userId uuid.UUID) (uuid.UUID, error)
	GetShoppingListType(shoppingListId uuid.UUID) (string, error)
	GetShoppingListOwner(shoppingListId uuid.UUID) (uuid.UUID, error)
	SetShoppingList(input entity.ShoppingListInput) (int32, error)
	DeleteSharedShoppingList(shoppingListId uuid.UUID) error

	GetShoppingListUsers(shoppingListId uuid.UUID) ([]uuid.UUID, error)
	GetShoppingListKey(shoppingListId uuid.UUID) (uuid.UUID, time.Time, error)
	IsShoppingListKeyValid(shoppingListId, key uuid.UUID) (bool, error)
	AddUserToShoppingList(userId, shoppingListId uuid.UUID) error
	DeleteUserFromShoppingList(userId, shoppingListId uuid.UUID) error
}
