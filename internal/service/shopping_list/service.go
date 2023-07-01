package shopping_list

import (
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/repository/grpc"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/service/dependencies/repository"
)

type Service struct {
	repo repository.ShoppingList
	grpc *grpc.Repository
}

func NewService(
	repo repository.ShoppingList,
	grpc *grpc.Repository,
) *Service {
	return &Service{
		repo: repo,
		grpc: grpc,
	}
}
