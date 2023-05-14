package entity

import "github.com/google/uuid"

type ShoppingListType string

var (
	ShoppingListTypePersonal = "personal"
	ShoppingListTypeShared   = "shared"
)

type ShoppingListInfo struct {
	Id      uuid.UUID
	Name    *string
	Type    ShoppingListType
	OwnerId uuid.UUID
}

type ShoppingList struct {
	Id          uuid.UUID
	Name        *string
	Purchases   []Purchase
	RecipeNames RecipeNames
	Type        ShoppingListType
	OwnerId     uuid.UUID
	Version     int32
}

type ShoppingListInput struct {
	ShoppingListId *uuid.UUID
	EditorId       uuid.UUID
	Purchases      []Purchase
	RecipeNames    RecipeNames
	LastVersion    *int32
}
