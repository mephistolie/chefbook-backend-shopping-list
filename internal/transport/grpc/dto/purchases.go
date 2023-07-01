package dto

import (
	"github.com/google/uuid"
	api "github.com/mephistolie/chefbook-backend-shopping-list/api/v2/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
)

const (
	maxUnitLength = 10
)

func parsePurchases(rawPurchases []*api.Purchase) []entity.Purchase {
	var purchases []entity.Purchase

	for _, rawPurchase := range rawPurchases {
		id, err := uuid.Parse(rawPurchase.Id)
		if err != nil {
			continue
		}
		if len(rawPurchase.Name) == 0 {
			continue
		}

		if rawPurchase.MeasureUnit != nil && len(*rawPurchase.MeasureUnit) > maxUnitLength {
			measureUnit := (*rawPurchase.MeasureUnit)[0:maxUnitLength]
			rawPurchase.MeasureUnit = &measureUnit
		}

		var recipeIdPtr *uuid.UUID = nil
		if rawPurchase.RecipeId != nil {
			if recipeId, err := uuid.Parse(*rawPurchase.RecipeId); err == nil {
				recipeIdPtr = &recipeId
			}
		}

		purchase := entity.Purchase{
			Id:          id,
			Name:        rawPurchase.Name,
			Multiplier:  rawPurchase.Multiplier,
			Purchased:   rawPurchase.Purchased,
			Amount:      rawPurchase.Multiplier,
			MeasureUnit: rawPurchase.MeasureUnit,
			RecipeId:    recipeIdPtr,
		}

		purchases = append(purchases, purchase)
	}

	return purchases
}
