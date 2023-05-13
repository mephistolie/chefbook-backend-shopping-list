package dto

import (
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
)

func NewPurchasesDto(entities []entity.Purchase) []Purchase {
	purchases := make([]Purchase, len(entities))
	for i, purchase := range entities {
		purchases[i] = newPurchase(purchase)
	}
	return purchases
}

func NewPurchasesEntity(dto []Purchase) []entity.Purchase {
	purchases := make([]entity.Purchase, len(dto))
	for i, purchase := range dto {
		purchases[i] = purchase.Entity()
	}
	return purchases
}
