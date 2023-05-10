package dto

import (
	"github.com/google/uuid"
	api "github.com/mephistolie/chefbook-backend-shopping-list/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
)

func ParsePurchases(rawPurchases []*api.Purchase) []entity.Purchase {
	var purchases []entity.Purchase
	for _, rawPurchase := range rawPurchases {
		id, err := uuid.Parse(rawPurchase.Id)
		if err != nil {
			continue
		}

		multiplier := 1
		if rawPurchase.Multiplier > 1 {
			multiplier = int(rawPurchase.Multiplier)
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
		var recipeNamePtr *string = nil
		if recipeId, err := uuid.Parse(rawPurchase.RecipeId); err == nil {
			recipeIdPtr = &recipeId
			if len(rawPurchase.RecipeName) > 0 {
				recipeName := rawPurchase.RecipeName
				recipeNamePtr = &recipeName
			}
		}

		purchase := entity.Purchase{
			Id:          id,
			Name:        rawPurchase.Name,
			Multiplier:  multiplier,
			Purchased:   rawPurchase.Purchased,
			Amount:      amountPtr,
			MeasureUnit: measureUnitPtr,
			RecipeId:    recipeIdPtr,
			RecipeName:  recipeNamePtr,
		}

		purchases = append(purchases, purchase)
	}

	return purchases
}
