package repository

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
)

type ShoppingList interface {
	CreatePersonalShoppingList(userId uuid.UUID, messageId uuid.UUID) error
	ImportFirebaseShoppingList(shoppingListId uuid.UUID, purchases []entity.Purchase, messageId uuid.UUID) error
	DeletePersonalShoppingList(userId uuid.UUID, messageId uuid.UUID) error

	GetShoppingLists(userId uuid.UUID, onlyPending bool) ([]entity.ShoppingListInfo, error)
	SetShoppingListName(shoppingListId uuid.UUID, name *string) error
	CreateSharedShoppingList(userId uuid.UUID, shoppingListId *uuid.UUID, name *string) (uuid.UUID, error)
	GetShoppingList(shoppingListId uuid.UUID) (entity.ShoppingList, error)
	GetPersonalShoppingListId(userId uuid.UUID) (uuid.UUID, error)
	GetShoppingListOwner(shoppingListId uuid.UUID) (uuid.UUID, error)
	SetShoppingList(input entity.ShoppingListInput) (int32, error)
	DeleteSharedShoppingList(shoppingListId uuid.UUID) error

	GetShoppingListUsers(shoppingListId uuid.UUID) ([]uuid.UUID, error)
	InviteUserToShoppingList(userId, shoppingListId uuid.UUID) error
	AcceptShoppingListInvite(userId, shoppingListId uuid.UUID) error
	DeleteUserFromShoppingList(userId, shoppingListId uuid.UUID) error
}
