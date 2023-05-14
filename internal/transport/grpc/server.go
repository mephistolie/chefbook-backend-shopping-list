package grpc

import (
	api "github.com/mephistolie/chefbook-backend-shopping-list/api/v2/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/transport/dependencies/service"
)

type ShoppingListServer struct {
	api.UnsafeShoppingListServiceServer
	service           service.Service
	checkSubscription bool
}

func NewServer(service service.Service, checkSubscription bool) *ShoppingListServer {
	return &ShoppingListServer{
		service:           service,
		checkSubscription: checkSubscription,
	}
}
