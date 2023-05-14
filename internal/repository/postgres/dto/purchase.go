package dto

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
)

type Purchase struct {
	Id          uuid.UUID  `json:"purchaseId" binding:"required"`
	Name        string     `json:"name" binding:"required"`
	Multiplier  *int       `json:"multiplier,omitempty"`
	Purchased   bool       `json:"purchased"`
	Amount      *int       `json:"amount,omitempty"`
	MeasureUnit *string    `json:"measureUnit,omitempty"`
	RecipeId    *uuid.UUID `json:"recipeId,omitempty"`
}

func newPurchase(purchase entity.Purchase) Purchase {
	return Purchase{
		Id:          purchase.Id,
		Name:        purchase.Name,
		Multiplier:  purchase.Multiplier,
		Purchased:   purchase.Purchased,
		Amount:      purchase.Amount,
		MeasureUnit: purchase.MeasureUnit,
		RecipeId:    purchase.RecipeId,
	}
}

func (l *Purchase) Entity() entity.Purchase {
	return entity.Purchase{
		Id:          l.Id,
		Name:        l.Name,
		Multiplier:  l.Multiplier,
		Purchased:   l.Purchased,
		Amount:      l.Amount,
		MeasureUnit: l.MeasureUnit,
		RecipeId:    l.RecipeId,
	}
}
