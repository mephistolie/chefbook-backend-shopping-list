package entity

import "github.com/google/uuid"

type ShoppingListType string

var (
	ShoppingListTypePersonal ShoppingListType = "personal"
	ShoppingListTypeShared   ShoppingListType = "shared"
)

type ShoppingListInfo struct {
	Id      uuid.UUID
	Name    *string
	Type    ShoppingListType
	OwnerId uuid.UUID
	Version int32
}

type ShoppingList struct {
	Id          uuid.UUID
	Name        *string
	Purchases   []Purchase
	RecipeNames map[string]string
	Type        ShoppingListType
	OwnerId     uuid.UUID
	Version     int32
}

type ShoppingListInput struct {
	ShoppingListId *uuid.UUID
	EditorId       uuid.UUID
	Purchases      []Purchase
	LastVersion    *int32
}
