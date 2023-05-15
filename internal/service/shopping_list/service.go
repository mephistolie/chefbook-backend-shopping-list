package shopping_list

import (
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/service/dependencies/repository"
)

type Service struct {
	repo repository.ShoppingList
}

func NewService(
	repo repository.ShoppingList,
) *Service {
	return &Service{repo: repo}
}
