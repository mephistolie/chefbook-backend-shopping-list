package mq

import (
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

func (s *Service) CreatePersonalShoppingList(userId uuid.UUID, messageId uuid.UUID) error {
	return s.repo.CreatePersonalShoppingList(userId, messageId)
}

func (s *Service) ImportFirebaseShoppingList(userId uuid.UUID, firebaseId string, messageId uuid.UUID) error {
	if s.firebase == nil {
		log.Warnf("try to import firebase profile with firebase import disabled")
		return errors.New("firebase import disabled")
	}

	firebasePurchases, err := s.firebase.GetShoppingList(firebaseId)
	if err != nil {
		log.Warnf("unable to get firebase shopping list for user %s: %s", userId, err)
		return err
	}

	log.Infof("found %d Firebase purchases for user %s...", len(*firebasePurchases), userId)
	var purchases []entity.Purchase
	for _, firebasePurchase := range *firebasePurchases {
		purchase := entity.Purchase{
			Id:   uuid.New(),
			Name: firebasePurchase,
		}
		purchases = append(purchases, purchase)
	}

	shoppingListId, err := s.repo.GetPersonalShoppingListId(userId)
	if err != nil {
		return err
	}
	return s.repo.ImportFirebaseShoppingList(shoppingListId, purchases, messageId)
}

func (s *Service) DeletePersonalShoppingList(userId uuid.UUID, messageId uuid.UUID) error {
	return s.repo.DeletePersonalShoppingList(userId, messageId)
}
