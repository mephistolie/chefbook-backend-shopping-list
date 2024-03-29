package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-common/subscription"
	api "github.com/mephistolie/chefbook-backend-shopping-list/api/v2/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/transport/grpc/dto"
)

func (s *ShoppingListServer) GetShoppingLists(_ context.Context, req *api.GetShoppingListsRequest) (*api.GetShoppingListsResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	shoppingLists, err := s.service.GetShoppingLists(userId)
	if err != nil {
		return nil, err
	}
	return &api.GetShoppingListsResponse{ShoppingLists: dto.NewShoppingListsResponse(shoppingLists)}, nil
}

func (s *ShoppingListServer) CreateSharedShoppingList(_ context.Context, req *api.CreateSharedShoppingListRequest) (*api.CreateSharedShoppingListResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	var shoppingListId *uuid.UUID = nil
	if req.ShoppingListId != nil {
		if id, err := uuid.Parse(*req.ShoppingListId); err == nil {
			shoppingListId = &id
		}
	}
	var name *string = nil
	if len(req.Name) > 0 {
		if len(req.Name) > 64 {
			return nil, fail.CreateGrpcClient(fail.TypeInvalidBody, "name max length is 64")
		}
		name = &req.Name
	}

	if !subscription.IsPremium(req.SubscriptionPlan) && s.checkSubscription {
		return nil, fail.GrpcPremiumRequired
	}

	resultId, err := s.service.CreateSharedShoppingList(userId, shoppingListId, name)
	if err != nil {
		return nil, err
	}

	return &api.CreateSharedShoppingListResponse{ShoppingListId: resultId.String()}, nil
}

func (s *ShoppingListServer) GetShoppingList(_ context.Context, req *api.GetShoppingListRequest) (*api.GetShoppingListResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	var shoppingListIdPtr *uuid.UUID = nil
	if shoppingListId, err := uuid.Parse(req.ShoppingListId); err == nil {
		shoppingListIdPtr = &shoppingListId
	}

	shoppingList, err := s.service.GetShoppingList(shoppingListIdPtr, userId)
	if err != nil {
		return nil, err
	}
	return dto.NewGetShoppingListResponse(shoppingList), nil
}

func (s *ShoppingListServer) SetShoppingListName(_ context.Context, req *api.SetShoppingListNameRequest) (*api.SetShoppingListNameResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	shoppingListId, err := uuid.Parse(req.ShoppingListId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	if req.Name != nil && len(*req.Name) > 64 {
		name := (*req.Name)[0:64]
		req.Name = &name
	}

	if err := s.service.SetShoppingListName(shoppingListId, userId, req.Name); err != nil {
		return nil, err
	}
	return &api.SetShoppingListNameResponse{Message: "shopping list name updated"}, nil
}

func (s *ShoppingListServer) SetShoppingList(_ context.Context, req *api.SetShoppingListRequest) (*api.SetShoppingListResponse, error) {
	input, err := dto.BindSetShoppingListRequest(req)
	if err != nil {
		return nil, err
	}

	version, err := s.service.SetShoppingList(input)
	if err != nil {
		return nil, err
	}
	return &api.SetShoppingListResponse{Version: version}, nil
}

func (s *ShoppingListServer) AddPurchasesToShoppingList(_ context.Context, req *api.SetShoppingListRequest) (*api.SetShoppingListResponse, error) {
	input, err := dto.BindSetShoppingListRequest(req)
	if err != nil {
		return nil, err
	}

	version, err := s.service.AddPurchasesToShoppingList(input)
	if err != nil {
		return nil, err
	}
	return &api.SetShoppingListResponse{Version: version}, nil
}

func (s *ShoppingListServer) DeleteSharedShoppingList(_ context.Context, req *api.DeleteSharedShoppingListRequest) (*api.DeleteSharedShoppingListResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	shoppingListId, err := uuid.Parse(req.ShoppingListId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	err = s.service.DeleteSharedShoppingList(shoppingListId, userId)
	if err != nil {
		return nil, err
	}

	return &api.DeleteSharedShoppingListResponse{Message: "shopping list deleted"}, nil
}
