package users

import (
	"errors"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/firebase"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/service/dependencies/repository"
	"time"
)

type Service struct {
	repo     repository.ShoppingList
	firebase *firebase.Client
}

func (s *Service) AddUser(userId uuid.UUID) error {
	return s.repo.AddUser(userId)
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

func (s *Service) ImportFirebaseData(userId uuid.UUID, firebaseId string) error {
	if s.firebase == nil {
		log.Warnf("try to import firebase profile with firebase import disabled")
		return errors.New("firebase import disabled")
	}

	log.Infof("importing Firebase data for user %s...", userId)
	firebasePurchases, err := s.firebase.GetShoppingList(firebaseId)
	if err != nil {
		log.Warnf("unable to get firebase shopping list for user %s: %s", userId, err)
		return err
	}

	var purchases []entity.Purchase
	for _, firebasePurchase := range *firebasePurchases {
		purchase := entity.Purchase{
			Id:         uuid.New(),
			Name:       firebasePurchase,
			Multiplier: 1,
		}
		purchases = append(purchases, purchase)
	}

	shoppingList := entity.ShoppingList{
		Purchases: purchases,
		Timestamp: time.Now(),
	}

	return s.repo.SetShoppingList(userId, shoppingList)
}

func (s *Service) DeleteUser(userId uuid.UUID) error {
	return s.repo.DeleteUser(userId)
}
