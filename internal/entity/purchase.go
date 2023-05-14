package entity

import "github.com/google/uuid"

type Purchase struct {
	Id          uuid.UUID
	Name        string
	Multiplier  *int
	Purchased   bool
	Amount      *int
	MeasureUnit *string
	RecipeId    *uuid.UUID
}

type RecipeNames map[uuid.UUID]string
