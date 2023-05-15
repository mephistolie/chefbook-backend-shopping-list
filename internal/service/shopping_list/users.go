package shopping_list

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	shoppingListFail "github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity/fail"
	"time"
)

func (s *Service) GetShoppingListUsers(shoppingListId, requesterId uuid.UUID) ([]uuid.UUID, error) {
	if err := s.checkUserIsShoppingListOwner(requesterId, shoppingListId); err != nil {
		return nil, err
	}
	return s.repo.GetShoppingListUsers(shoppingListId)
}

func (s *Service) GetShoppingListLink(shoppingListId, requesterId uuid.UUID, linkPattern string) (string, time.Time, error) {
	if err := s.checkUserIsShoppingListOwner(requesterId, shoppingListId); err != nil {
		return "", time.Time{}, err
	}
	key, expiresAt, err := s.repo.GetShoppingListKey(shoppingListId)
	if err != nil {
		return "", time.Time{}, err
	}
	return fmt.Sprintf(linkPattern, shoppingListId.String(), key.String()), expiresAt, nil
}

func (s *Service) JoinShoppingList(shoppingListId, userId, key uuid.UUID) error {
	if valid, err := s.repo.IsShoppingListKeyValid(shoppingListId, key); err != nil || !valid {
		return shoppingListFail.GrpcInvalidShoppingListKey
	}
	return s.repo.AddUserToShoppingList(userId, shoppingListId)
}

func (s *Service) DeleteUserFromShoppingList(userId, shoppingListId, requesterId uuid.UUID) error {
	if userId != requesterId {
		if err := s.checkUserIsShoppingListOwner(requesterId, shoppingListId); err != nil {
			return err
		}
	}
	return s.repo.DeleteUserFromShoppingList(userId, shoppingListId)
}

func (s *Service) checkUserHasAccessToShoppingList(userId, shoppingListId uuid.UUID, checkOwnership bool) error {
	if checkOwnership {
		ownerId, err := s.repo.GetShoppingListOwner(shoppingListId)
		if err != nil {
			return err
		}
		if userId == ownerId {
			return nil
		}
	}

	users, err := s.repo.GetShoppingListUsers(shoppingListId)
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

func (s *Service) checkUserIsShoppingListOwner(userId, shoppingListId uuid.UUID) error {
	ownerId, err := s.repo.GetShoppingListOwner(shoppingListId)
	if err != nil {
		return err
	}
	if userId != ownerId {
		return fail.GrpcAccessDenied
	}
	return nil
}
