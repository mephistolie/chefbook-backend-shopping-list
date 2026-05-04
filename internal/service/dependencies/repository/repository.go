package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
	"time"
)

type ShoppingList interface {
	CreatePersonalShoppingList(ctx context.Context, userId uuid.UUID, messageId uuid.UUID) error
	ImportFirebaseShoppingList(ctx context.Context, shoppingListId uuid.UUID, purchases []entity.Purchase, messageId uuid.UUID) error
	DeleteUserShoppingLists(ctx context.Context, userId uuid.UUID, messageId uuid.UUID) error

	GetShoppingLists(ctx context.Context, userId uuid.UUID) ([]entity.ShoppingListInfo, error)
	SetShoppingListName(ctx context.Context, shoppingListId, userId uuid.UUID, name *string) error
	CreateSharedShoppingList(ctx context.Context, userId uuid.UUID, shoppingListId *uuid.UUID, name *string) (uuid.UUID, error)
	GetShoppingList(ctx context.Context, shoppingListId, userId uuid.UUID) (entity.ShoppingList, error)
	GetPersonalShoppingListId(ctx context.Context, userId uuid.UUID) (uuid.UUID, error)
	GetShoppingListType(ctx context.Context, shoppingListId uuid.UUID) (string, error)
	GetShoppingListOwner(ctx context.Context, shoppingListId uuid.UUID) (uuid.UUID, error)
	SetShoppingList(ctx context.Context, input entity.ShoppingListInput) (int32, error)
	DeleteSharedShoppingList(ctx context.Context, shoppingListId uuid.UUID) error

	GetShoppingListUsers(ctx context.Context, shoppingListId uuid.UUID) ([]uuid.UUID, error)
	GetShoppingListKey(ctx context.Context, shoppingListId uuid.UUID) (uuid.UUID, time.Time, error)
	IsShoppingListKeyValid(ctx context.Context, shoppingListId, key uuid.UUID) (bool, error)
	AddUserToShoppingList(ctx context.Context, userId, shoppingListId uuid.UUID) error
	DeleteUserFromShoppingList(ctx context.Context, userId, shoppingListId uuid.UUID) error
}
