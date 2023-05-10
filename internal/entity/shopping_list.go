package entity

import "time"

type ShoppingList struct {
	Purchases []Purchase
	Timestamp time.Time
}
