package entity

import "github.com/google/uuid"

type Purchase struct {
	Id          uuid.UUID
	Name        string
	Multiplier  *int32
	Purchased   bool
	Amount      *float32
	MeasureUnit *string
	RecipeId    *uuid.UUID
}
