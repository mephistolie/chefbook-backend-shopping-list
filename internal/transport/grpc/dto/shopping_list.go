package dto

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	api "github.com/mephistolie/chefbook-backend-shopping-list/api/v2/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
)

func BindSetShoppingListRequest(req *api.SetShoppingListRequest) (entity.ShoppingListInput, error) {
	editorId, err := uuid.Parse(req.EditorId)
	if err != nil {
		return entity.ShoppingListInput{}, fail.GrpcInvalidBody
	}
	var shoppingListIdPtr *uuid.UUID = nil
	if shoppingListId, err := uuid.Parse(req.ShoppingListId); err == nil {
		shoppingListIdPtr = &shoppingListId
	}

	var lastVersion *int32 = nil
	if req.LastVersion > 0 {
		lastVersion = &req.LastVersion
	}

	return entity.ShoppingListInput{
		ShoppingListId: shoppingListIdPtr,
		EditorId:       editorId,
		Purchases:      parsePurchases(req.Purchases),
		LastVersion:    lastVersion,
	}, nil
}

func NewShoppingListsResponse(shoppingLists []entity.ShoppingListInfo) []*api.ShoppingListInfo {
	response := make([]*api.ShoppingListInfo, len(shoppingLists))
	for i, shoppingList := range shoppingLists {
		rawShoppingList := api.ShoppingListInfo{
			Id:      shoppingList.Id.String(),
			Type:    string(shoppingList.Type),
			OwnerId: shoppingList.OwnerId.String(),
		}
		response[i] = &rawShoppingList
	}
	return response
}

func NewGetShoppingListResponse(shoppingList entity.ShoppingList) *api.GetShoppingListResponse {
	name := ""
	if shoppingList.Name != nil {
		name = *shoppingList.Name
	}

	purchases := make([]*api.Purchase, len(shoppingList.Purchases))
	for i, purchase := range shoppingList.Purchases {
		var recipeIdPtr *string
		if purchase.RecipeId != nil {
			recipeId := purchase.RecipeId.String()
			recipeIdPtr = &recipeId
		}

		rawPurchase := api.Purchase{
			Id:          purchase.Id.String(),
			Name:        purchase.Name,
			Multiplier:  purchase.Multiplier,
			Purchased:   purchase.Purchased,
			Amount:      purchase.Amount,
			MeasureUnit: purchase.MeasureUnit,
			RecipeId:    recipeIdPtr,
		}
		purchases[i] = &rawPurchase
	}

	return &api.GetShoppingListResponse{
		Id:          shoppingList.Id.String(),
		Name:        name,
		Type:        string(shoppingList.Type),
		OwnerId:     shoppingList.OwnerId.String(),
		Purchases:   purchases,
		RecipeNames: shoppingList.RecipeNames,
		Version:     shoppingList.Version,
	}
}
