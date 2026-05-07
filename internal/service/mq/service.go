package mq

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/firebase"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/service/dependencies/repository"
)

type Service struct {
	repo     repository.ShoppingList
	firebase *firebase.Client
}

func NewService(
	repo repository.ShoppingList,
	firebase *firebase.Client,
) *Service {
	return &Service{
		repo:     repo,
		firebase: firebase,
	}
}

func (s *Service) CreatePersonalShoppingList(ctx context.Context, userId uuid.UUID, messageId uuid.UUID) error {
	return s.repo.CreatePersonalShoppingList(ctx, userId, messageId)
}

func (s *Service) ImportFirebaseShoppingList(ctx context.Context, userId uuid.UUID, firebaseId string, messageId uuid.UUID) error {
	if s.firebase == nil {
		log.LogWarn(ctx, log.Event{
			Event:     "firebase.import.disabled",
			Message:   "try to import firebase profile with firebase import disabled",
			Component: log.ComponentFirebase,
			UserID:    userId.String(),
		})
		return errors.New("firebase import disabled")
	}

	firebasePurchases, err := s.firebase.GetShoppingList(firebaseId)
	if err != nil {
		log.LogWarnError(ctx, log.Event{
			Event:     "firebase.shopping_list.load_failed",
			Message:   "unable to get firebase shopping list",
			Component: log.ComponentFirebase,
			UserID:    userId.String(),
		}, err)
		return err
	}

	log.Log(ctx, log.Event{
		Event:     "firebase.shopping_list.loaded",
		Message:   "firebase shopping list loaded",
		Component: log.ComponentFirebase,
		UserID:    userId.String(),
		Payload: map[string]any{
			"purchases_count": len(*firebasePurchases),
		},
	})
	var purchases []entity.Purchase
	for _, firebasePurchase := range *firebasePurchases {
		purchase := entity.Purchase{
			Id:   uuid.New(),
			Name: firebasePurchase,
		}
		purchases = append(purchases, purchase)
	}

	shoppingListId, err := s.repo.GetPersonalShoppingListId(ctx, userId)
	if err != nil {
		return err
	}
	return s.repo.ImportFirebaseShoppingList(ctx, shoppingListId, purchases, messageId)
}

func (s *Service) DeleteUserShoppingLists(ctx context.Context, userId uuid.UUID, messageId uuid.UUID) error {
	return s.repo.DeleteUserShoppingLists(ctx, userId, messageId)
}
