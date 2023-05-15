package service

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/firebase"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/config"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/service/dependencies/repository"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/service/mq"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/service/shopping_list"
	"time"
)

type Service struct {
	ShoppingList
	MQ
}

type ShoppingList interface {
	GetShoppingLists(userId uuid.UUID) ([]entity.ShoppingListInfo, error)
	CreateSharedShoppingList(userId uuid.UUID, shoppingListId *uuid.UUID, name *string) (uuid.UUID, error)
	GetShoppingList(shoppingListId *uuid.UUID, userId uuid.UUID) (entity.ShoppingList, error)
	SetShoppingListName(shoppingListId, userId uuid.UUID, name *string) error
	SetShoppingList(input entity.ShoppingListInput) (int32, error)
	AddPurchasesToShoppingList(input entity.ShoppingListInput) (int32, error)
	DeleteSharedShoppingList(shoppingListId uuid.UUID, userId uuid.UUID) error

	GetShoppingListUsers(shoppingListId, requesterId uuid.UUID) ([]uuid.UUID, error)
	GetShoppingListLink(shoppingListId, requesterId uuid.UUID, linkPattern string) (string, time.Time, error)
	JoinShoppingList(shoppingListId, userId, key uuid.UUID) error
	DeleteUserFromShoppingList(userId, shoppingListId, requesterId uuid.UUID) error
}

type MQ interface {
	CreatePersonalShoppingList(userId uuid.UUID, messageId uuid.UUID) error
	ImportFirebaseShoppingList(userId uuid.UUID, firebaseId string, messageId uuid.UUID) error
	DeletePersonalShoppingList(userId uuid.UUID, messageId uuid.UUID) error
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
		MQ:           mq.NewService(repo, client),
	}, nil
}
