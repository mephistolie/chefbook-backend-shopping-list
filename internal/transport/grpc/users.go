package grpc

import (
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-common/subscription"
	api "github.com/mephistolie/chefbook-backend-shopping-list/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/transport/grpc/dto"
)

func (s *ShoppingListServer) GetShoppingListUsers(_ context.Context, req *api.GetShoppingListUsersRequest) (*api.GetShoppingListUsersResponse, error) {
	requesterId, err := uuid.Parse(req.RequesterId)
	if err == nil {
		return nil, fail.GrpcInvalidBody
	}
	shoppingListId, err := uuid.Parse(req.ShoppingListId)
	if err == nil {
		return nil, fail.GrpcInvalidBody
	}

	invites, err := s.service.GetShoppingListUsers(shoppingListId, requesterId)
	if err != nil {
		return nil, err
	}

	return &api.GetShoppingListUsersResponse{Users: dto.NewInvitesResponse(invites)}, nil
}

func (s *ShoppingListServer) GetShoppingListInvites(_ context.Context, req *api.GetShoppingListInvitesRequest) (*api.GetShoppingListInvitesResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}

	invites, err := s.service.GetShoppingListInvites(userId)
	if err != nil {
		return nil, err
	}

	return &api.GetShoppingListInvitesResponse{ShoppingLists: dto.NewShoppingListsResponse(invites)}, nil
}

func (s *ShoppingListServer) InviteUserToShoppingList(_ context.Context, req *api.InviteUserToShoppingListRequest) (*api.InviteUserToShoppingListResponse, error) {
	requesterId, err := uuid.Parse(req.RequesterId)
	if err == nil {
		return nil, fail.GrpcInvalidBody
	}
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	shoppingListId, err := uuid.Parse(req.ShoppingListId)
	if err == nil {
		return nil, fail.GrpcInvalidBody
	}

	if !subscription.IsPremium(req.SubscriptionPlan) {
		return nil, fail.GrpcPremiumRequired
	}

	err = s.service.InviteShoppingListUser(userId, shoppingListId, requesterId)
	if err != nil {
		return nil, err
	}

	return &api.InviteUserToShoppingListResponse{Message: "user invited"}, nil
}

func (s *ShoppingListServer) AcceptShoppingListInvite(_ context.Context, req *api.AcceptShoppingListInviteRequest) (*api.AcceptShoppingListInviteResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	shoppingListId, err := uuid.Parse(req.ShoppingListId)
	if err == nil {
		return nil, fail.GrpcInvalidBody
	}

	err = s.service.AcceptShoppingListInvite(userId, shoppingListId)
	if err != nil {
		return nil, err
	}

	return &api.AcceptShoppingListInviteResponse{Message: "invite accepted"}, nil
}

func (s *ShoppingListServer) DeclineShoppingListInvite(_ context.Context, req *api.DeclineShoppingListInviteRequest) (*api.DeclineShoppingListInviteResponse, error) {
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	shoppingListId, err := uuid.Parse(req.ShoppingListId)
	if err == nil {
		return nil, fail.GrpcInvalidBody
	}

	err = s.service.DeclineShoppingListInvite(userId, shoppingListId)
	if err != nil {
		return nil, err
	}

	return &api.DeclineShoppingListInviteResponse{Message: "invite declined"}, nil
}

func (s *ShoppingListServer) DeleteUserFromShoppingList(_ context.Context, req *api.DeleteUserFromShoppingListRequest) (*api.DeleteUserFromShoppingListResponse, error) {
	requesterId, err := uuid.Parse(req.RequesterId)
	if err == nil {
		return nil, fail.GrpcInvalidBody
	}
	userId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, fail.GrpcInvalidBody
	}
	shoppingListId, err := uuid.Parse(req.ShoppingListId)
	if err == nil {
		return nil, fail.GrpcInvalidBody
	}

	err = s.service.DeleteUserFromShoppingList(userId, shoppingListId, requesterId)
	if err != nil {
		return nil, err
	}

	return &api.DeleteUserFromShoppingListResponse{Message: "user excluded from shopping list"}, nil
}
