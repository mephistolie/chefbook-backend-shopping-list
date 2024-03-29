package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	api "github.com/mephistolie/chefbook-backend-shopping-list/api/v2/proto/implementation/v1"
	shoppingListFail "github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity/fail"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/transport/grpc/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *ShoppingListServer) GetShoppingListUsers(_ context.Context, req *api.GetShoppingListUsersRequest) (*api.GetShoppingListUsersResponse, error) {
	requesterId, err := uuid.Parse(req.RequesterId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	shoppingListId, err := uuid.Parse(req.ShoppingListId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	users, err := s.service.GetShoppingListUsers(shoppingListId, requesterId)
	if err != nil {
		return nil, err
	}

	return &api.GetShoppingListUsersResponse{Users: dto.NewShoppingListUsersResponse(users)}, nil
}

func (s *ShoppingListServer) GetSharedShoppingListLink(_ context.Context, req *api.GetSharedShoppingListLinkRequest) (*api.GetSharedShoppingListLinkResponse, error) {
	requesterId, err := uuid.Parse(req.RequesterId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	shoppingListId, err := uuid.Parse(req.ShoppingListId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	link, expiresAt, err := s.service.GetShoppingListLink(shoppingListId, requesterId, req.LinkPattern)
	if err != nil {
		return nil, err
	}

	return &api.GetSharedShoppingListLinkResponse{Link: link, ExpiresAt: timestamppb.New(expiresAt)}, nil
}

func (s *ShoppingListServer) JoinShoppingList(_ context.Context, req *api.JoinShoppingListRequest) (*api.JoinShoppingListResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	shoppingListId, err := uuid.Parse(req.ShoppingListId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	key, err := uuid.Parse(req.Key)
	if err != nil {
		return nil, shoppingListFail.GrpcInvalidShoppingListKey
	}

	err = s.service.JoinShoppingList(shoppingListId, userId, key)
	if err != nil {
		return nil, err
	}

	return &api.JoinShoppingListResponse{Message: "joined"}, nil
}

func (s *ShoppingListServer) DeleteUserFromShoppingList(_ context.Context, req *api.DeleteUserFromShoppingListRequest) (*api.DeleteUserFromShoppingListResponse, error) {
	requesterId, err := uuid.Parse(req.RequesterId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	shoppingListId, err := uuid.Parse(req.ShoppingListId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	err = s.service.DeleteUserFromShoppingList(userId, shoppingListId, requesterId)
	if err != nil {
		return nil, err
	}

	return &api.DeleteUserFromShoppingListResponse{Message: "user excluded from shopping list"}, nil
}
