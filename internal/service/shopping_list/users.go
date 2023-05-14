package shopping_list

import (
	"context"
	"github.com/google/uuid"
	api "github.com/mephistolie/chefbook-backend-auth/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
)

func (s *Service) GetShoppingListInvites(userId uuid.UUID) ([]entity.ShoppingListInfo, error) {
	return s.repo.GetShoppingLists(userId, true)
}

func (s *Service) GetShoppingListUsers(shoppingListId, requesterId uuid.UUID) ([]uuid.UUID, error) {
	if err := s.checkUserIsShoppingListOwner(requesterId, shoppingListId); err != nil {
		return nil, err
	}
	return s.repo.GetShoppingListUsers(shoppingListId)
}

func (s *Service) InviteShoppingListUser(userId, shoppingListId, requesterId uuid.UUID) error {
	if err := s.checkUserIsShoppingListOwner(requesterId, shoppingListId); err != nil {
		return err
	}

	err := s.repo.InviteUserToShoppingList(userId, shoppingListId)
	if err == nil {
		go s.processShoppingListUserInvite(shoppingListId, userId)
	}

	return err
}

func (s *Service) processShoppingListUserInvite(shoppingListId, userId uuid.UUID) {
	user, err := s.auth.GetAuthInfo(context.Background(), &api.GetAuthInfoRequest{Id: userId.String()})
	if err != nil {
		return
	}
	if user.IsBlocked {
		_ = s.repo.DeleteUserFromShoppingList(userId, shoppingListId)
		return
	}
	if user.IsActivated {
		s.mail.SendShoppingListInviteMail(user.Email)
	}
}

func (s *Service) AcceptShoppingListInvite(userId, shoppingListId uuid.UUID) error {
	return s.repo.AcceptShoppingListInvite(userId, shoppingListId)
}

func (s *Service) DeclineShoppingListInvite(userId, shoppingListId uuid.UUID) error {
	return s.repo.DeleteUserFromShoppingList(userId, shoppingListId)
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
