package dto

import (
	"github.com/google/uuid"
	api "github.com/mephistolie/chefbook-backend-shopping-list/api/v2/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
	"math"
)

const (
	maxNameLength  = 75
	maxAmountCount = 10000
	maxUnitLength  = 15
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
		if len([]rune(rawPurchase.Name)) > maxNameLength {
			measureUnit := string([]rune(rawPurchase.Name)[0:maxNameLength])
			rawPurchase.MeasureUnit = &measureUnit
		}

		if rawPurchase.Amount != nil {
			if *rawPurchase.Amount <= 0 {
				rawPurchase.Amount = nil
			} else if *rawPurchase.Amount > maxAmountCount {
				*rawPurchase.Amount = maxAmountCount
			}
		}
		if rawPurchase.Amount != nil {
			*rawPurchase.Amount = float32(math.Round(float64(*rawPurchase.Amount)*1000) / 1000)
		}

		if rawPurchase.MeasureUnit != nil && len([]rune(*rawPurchase.MeasureUnit)) > maxUnitLength {
			measureUnit := string([]rune(*rawPurchase.MeasureUnit)[0:maxUnitLength])
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
			Amount:      rawPurchase.Amount,
			MeasureUnit: rawPurchase.MeasureUnit,
			RecipeId:    recipeIdPtr,
		}

		purchases = append(purchases, purchase)
	}

	return purchases
}
