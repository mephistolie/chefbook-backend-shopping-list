package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	api "github.com/mephistolie/chefbook-backend-shopping-list/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/transport/dependencies/service"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/transport/grpc/dto"
)

type ShoppingListServer struct {
	service service.Service
	api.UnsafeShoppingListServiceServer
}

func NewServer(service service.Service) *ShoppingListServer {
	return &ShoppingListServer{
		service: service,
	}
}

func (s *ShoppingListServer) GetShoppingList(_ context.Context, req *api.GetShoppingListRequest) (*api.GetShoppingListResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	shoppingList, err := s.service.GetShoppingList(userId)
	if err != nil {
		return nil, err
	}
	return dto.NewGetShoppingListResponse(shoppingList), nil
}

func (s *ShoppingListServer) SetShoppingList(_ context.Context, req *api.SetShoppingListRequest) (*api.SetShoppingListResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	purchases := dto.ParsePurchases(req.Purchases)
	var lastVersion *int32 = nil
	if req.LastVersion > 0 {
		lastVersion = &req.LastVersion
	}

	version, err := s.service.SetShoppingList(userId, purchases, lastVersion)
	if err != nil {
		return nil, err
	}
	return &api.SetShoppingListResponse{Version: version}, nil
}

func (s *ShoppingListServer) AddToShoppingList(_ context.Context, req *api.AddToShoppingListRequest) (*api.AddToShoppingListResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	purchases := dto.ParsePurchases(req.Purchases)
	var lastVersion *int32 = nil
	if req.LastVersion > 0 {
		lastVersion = &req.LastVersion
	}

	version, err := s.service.AddToShoppingList(userId, purchases, lastVersion)
	if err != nil {
		return nil, err
	}
	return &api.AddToShoppingListResponse{Version: version}, nil
}
