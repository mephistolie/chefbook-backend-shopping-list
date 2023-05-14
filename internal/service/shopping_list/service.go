package shopping_list

import (
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/service/dependencies/repository"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/service/dependencies/services"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/service/mail"
)

type Service struct {
	repo repository.ShoppingList
	auth *services.Auth
	mail *mail.Service
}

func NewService(
	repo repository.ShoppingList,
	auth *services.Auth,
	mail *mail.Service,
) *Service {
	return &Service{
		repo: repo,
		auth: auth,
		mail: mail,
	}
}
