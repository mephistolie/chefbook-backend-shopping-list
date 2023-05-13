package dto

import (
	api "github.com/mephistolie/chefbook-backend-shopping-list/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
)

func NewGetShoppingListResponse(shoppingList entity.ShoppingList) *api.GetShoppingListResponse {
	purchases := make([]*api.Purchase, len(shoppingList.Purchases))
	for i, purchase := range shoppingList.Purchases {
		rawPurchase := api.Purchase{
			Id: purchase.Id.String(),
		}
		purchases[i] = &rawPurchase
	}
	return &api.GetShoppingListResponse{
		Purchases: purchases,
		Version:   shoppingList.Version,
	}
}
