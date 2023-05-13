package service

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/firebase"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/config"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/service/dependencies/repository"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/service/shopping_list"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/service/users"
)

type Service struct {
	ShoppingList
	Users
}

type ShoppingList interface {
	GetShoppingList(userId uuid.UUID) (entity.ShoppingList, error)
	SetShoppingList(userId uuid.UUID, purchases []entity.Purchase, lastVersion *int32) (int32, error)
	AddToShoppingList(userId uuid.UUID, purchases []entity.Purchase, lastVersion *int32) (int32, error)
}

type Users interface {
	AddUser(userId uuid.UUID, messageId uuid.UUID) error
	ImportFirebaseData(userId uuid.UUID, firebaseId string, messageId uuid.UUID) error
	DeleteUser(userId uuid.UUID, messageId uuid.UUID) error
}

func New(
	cfg *config.Config,
	repo repository.ShoppingList,
) (*Service, error) {
	var err error = nil
	var client *firebase.Client = nil
	if len(*cfg.Firebase.Credentials) > 0 {
		credentials := []byte(*cfg.Firebase.Credentials)
		client, err = firebase.NewClient(credentials, "")
		if err != nil {
			return nil, err
		}
		log.Info("Firebase client initialized")
	}

	return &Service{
		ShoppingList: shopping_list.NewService(repo),
		Users:        users.NewService(repo, client),
	}, nil
}
