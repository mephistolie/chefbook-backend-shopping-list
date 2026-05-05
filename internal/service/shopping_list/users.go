package shopping_list

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
	shoppingListFail "github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity/fail"
	"time"
)

func (s *Service) GetShoppingListUsers(ctx context.Context, shoppingListId, requesterId uuid.UUID) ([]entity.User, error) {
	if err := s.checkUserIsShoppingListOwner(ctx, requesterId, shoppingListId); err != nil {
		return nil, err
	}
	ids, err := s.repo.GetShoppingListUsers(ctx, shoppingListId)
	if err != nil {
		return nil, err
	}

	var rawIds []string
	var users []entity.User
	for _, id := range ids {
		rawIds = append(rawIds, id.String())
		users = append(users, entity.User{Id: id})
	}

	profiles := s.getProfilesInfo(ctx, rawIds)
	for i := range users {
		if profile, ok := profiles[users[i].Id.String()]; ok {
			users[i].Name = profile.VisibleName
			users[i].Avatar = profile.Avatar
		}
	}

	return users, nil
}

func (s *Service) GetShoppingListLink(ctx context.Context, shoppingListId, requesterId uuid.UUID, linkPattern string) (string, time.Time, error) {
	if shoppingListType, err := s.repo.GetShoppingListType(ctx, shoppingListId); err != nil || shoppingListType != string(entity.ShoppingListTypeShared) {
		return "", time.Time{}, shoppingListFail.GrpcPersonalShoppingList
	}

	if err := s.checkUserIsShoppingListOwner(ctx, requesterId, shoppingListId); err != nil {
		return "", time.Time{}, err
	}

	key, expiresAt, err := s.repo.GetShoppingListKey(ctx, shoppingListId)
	if err != nil {
		return "", time.Time{}, err
	}
	return fmt.Sprintf(linkPattern, shoppingListId.String(), key.String()), expiresAt, nil
}

func (s *Service) JoinShoppingList(ctx context.Context, shoppingListId, userId, key uuid.UUID) error {
	if valid, err := s.repo.IsShoppingListKeyValid(ctx, shoppingListId, key); err != nil || !valid {
		return shoppingListFail.GrpcInvalidShoppingListKey
	}
	return s.repo.AddUserToShoppingList(ctx, userId, shoppingListId)
}

func (s *Service) DeleteUserFromShoppingList(ctx context.Context, userId, shoppingListId, requesterId uuid.UUID) error {
	ownerId, err := s.repo.GetShoppingListOwner(ctx, shoppingListId)
	if err != nil {
		return err
	}
	if requesterId != ownerId {
		return fail.GrpcAccessDenied
	}
	if userId == ownerId {
		return shoppingListFail.GrpcShoppingListOwner
	}
	return s.repo.DeleteUserFromShoppingList(ctx, userId, shoppingListId)
}

func (s *Service) checkUserHasAccessToShoppingList(ctx context.Context, userId, shoppingListId uuid.UUID, checkOwnership bool) error {
	if checkOwnership {
		ownerId, err := s.repo.GetShoppingListOwner(ctx, shoppingListId)
		if err != nil {
			return err
		}
		if userId == ownerId {
			return nil
		}
	}

	users, err := s.repo.GetShoppingListUsers(ctx, shoppingListId)
	if err != nil {
		return err
	}
	accessed := false
	for _, userWithAccess := range users {
		if userId == userWithAccess {
			accessed = true
			break
		}
	}
	if !accessed {
		return fail.GrpcAccessDenied
	}

	return nil
}

func (s *Service) checkUserIsShoppingListOwner(ctx context.Context, userId, shoppingListId uuid.UUID) error {
	ownerId, err := s.repo.GetShoppingListOwner(ctx, shoppingListId)
	if err != nil {
		return err
	}
	if userId != ownerId {
		return fail.GrpcAccessDenied
	}
	return nil
}
