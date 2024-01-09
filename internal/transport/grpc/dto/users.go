package dto

import (
	api "github.com/mephistolie/chefbook-backend-shopping-list/api/v2/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
)

func NewShoppingListUsersResponse(users []entity.User) []*api.ShoppingListUser {
	dtos := make([]*api.ShoppingListUser, len(users))
	for _, user := range users {
		dtos = append(dtos, &api.ShoppingListUser{
			Id:     user.Id.String(),
			Name:   user.Name,
			Avatar: user.Avatar,
		})
	}
	return dtos
}
