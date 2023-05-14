package dto

import (
	"github.com/google/uuid"
	api "github.com/mephistolie/chefbook-backend-shopping-list/api/v2/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
)

func parsePurchases(rawPurchases []*api.Purchase) ([]entity.Purchase, []uuid.UUID) {
	var purchases []entity.Purchase
	var recipeIds []uuid.UUID

	for _, rawPurchase := range rawPurchases {
		id, err := uuid.Parse(rawPurchase.Id)
		if err != nil {
			continue
		}

		var multiplierPtr *int = nil
		if rawPurchase.Multiplier > 0 {
			multiplier := int(rawPurchase.Multiplier)
			multiplierPtr = &multiplier
		}

		var amountPtr *int = nil
		if rawPurchase.Amount > 0 {
			amount := int(rawPurchase.Amount)
			amountPtr = &amount
		}

		var measureUnitPtr *string = nil
		if len(rawPurchase.MeasureUnit) > 0 {
			measureUnit := rawPurchase.MeasureUnit
			measureUnitPtr = &measureUnit
		}

		var recipeIdPtr *uuid.UUID = nil
		if recipeId, err := uuid.Parse(rawPurchase.RecipeId); err == nil {
			recipeIdPtr = &recipeId
			recipeIds = append(recipeIds, recipeId)
		}

		purchase := entity.Purchase{
			Id:          id,
			Name:        rawPurchase.Name,
			Multiplier:  multiplierPtr,
			Purchased:   rawPurchase.Purchased,
			Amount:      amountPtr,
			MeasureUnit: measureUnitPtr,
			RecipeId:    recipeIdPtr,
		}

		purchases = append(purchases, purchase)
	}

	return purchases, recipeIds
}
