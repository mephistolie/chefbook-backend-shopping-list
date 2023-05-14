package dto

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	api "github.com/mephistolie/chefbook-backend-shopping-list/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
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

	purchases, recipeIds := parsePurchases(req.Purchases)
	recipeNames := parseRecipeNames(req.RecipeNames, recipeIds)
	for i := range purchases {
		if purchases[i].RecipeId != nil {
			if _, ok := recipeNames[*purchases[i].RecipeId]; !ok {
				purchases[i].RecipeId = nil
			}
		}
	}

	var lastVersion *int32 = nil
	if req.LastVersion > 0 {
		lastVersion = &req.LastVersion
	}

	return entity.ShoppingListInput{
		ShoppingListId: shoppingListIdPtr,
		EditorId:       editorId,
		Purchases:      purchases,
		RecipeNames:    recipeNames,
		LastVersion:    lastVersion,
	}, nil
}

func parseRecipeNames(request map[string]string, usedIds []uuid.UUID) map[uuid.UUID]string {
	var names map[uuid.UUID]string
	for _, id := range usedIds {
		if name, ok := request[id.String()]; ok {
			names[id] = name
		}
	}
	return names
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
	purchases := make([]*api.Purchase, len(shoppingList.Purchases))
	for i, purchase := range shoppingList.Purchases {
		multiplier := 0
		if purchase.Multiplier != nil && *purchase.Multiplier > 0 {
			multiplier = *purchase.Multiplier
		}

		amount := 0
		if purchase.Amount != nil && *purchase.Amount > 0 {
			amount = *purchase.Amount
		}

		measureUnit := ""
		if purchase.MeasureUnit != nil {
			measureUnit = *purchase.MeasureUnit
		}

		recipeId := ""
		if purchase.RecipeId != nil {
			recipeId = purchase.RecipeId.String()
		}

		rawPurchase := api.Purchase{
			Id:          purchase.Id.String(),
			Name:        purchase.Name,
			Multiplier:  int32(multiplier),
			Purchased:   purchase.Purchased,
			Amount:      int32(amount),
			MeasureUnit: measureUnit,
			RecipeId:    recipeId,
		}
		purchases[i] = &rawPurchase
	}
	return &api.GetShoppingListResponse{
		Purchases: purchases,
		Version:   shoppingList.Version,
	}
}
