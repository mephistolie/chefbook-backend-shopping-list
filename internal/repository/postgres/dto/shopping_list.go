package dto

import (
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
	"time"
)

type ShoppingList struct {
	Purchases []Purchase `json:"purchases"`
	Timestamp time.Time  `json:"timestamp"`
}

func NewShoppingList(shoppingList entity.ShoppingList) ShoppingList {
	purchases := make([]Purchase, len(shoppingList.Purchases))
	for i, purchase := range shoppingList.Purchases {
		purchases[i] = newPurchase(purchase)
	}

	return ShoppingList{
		Purchases: purchases,
		Timestamp: shoppingList.Timestamp,
	}
}

func (l *ShoppingList) Entity() entity.ShoppingList {
	purchases := make([]entity.Purchase, len(l.Purchases))
	for i, purchase := range l.Purchases {
		purchases[i] = purchase.Entity()
	}

	return entity.ShoppingList{
		Purchases: purchases,
		Timestamp: l.Timestamp,
	}
}
