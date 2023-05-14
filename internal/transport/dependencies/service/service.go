package service

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/firebase"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/config"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/service/dependencies/repository"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/service/dependencies/services"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/service/mail"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/service/mq"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/service/shopping_list"
)

type Service struct {
	ShoppingList
	MQ
}

type ShoppingList interface {
	GetShoppingLists(userId uuid.UUID) ([]entity.ShoppingListInfo, error)
	CreateSharedShoppingList(userId uuid.UUID, shoppingListId *uuid.UUID, name *string) (uuid.UUID, error)
	GetShoppingList(shoppingListId *uuid.UUID, userId uuid.UUID) (entity.ShoppingList, error)
	SetShoppingListName(shoppingListId *uuid.UUID, name *string, requesterId uuid.UUID) error
	SetShoppingList(input entity.ShoppingListInput) (int32, error)
	AddPurchasesToShoppingList(input entity.ShoppingListInput) (int32, error)
	DeleteSharedShoppingList(shoppingListId uuid.UUID, userId uuid.UUID) error

	GetShoppingListInvites(userId uuid.UUID) ([]entity.ShoppingListInfo, error)
	GetShoppingListUsers(shoppingListId, requesterId uuid.UUID) ([]uuid.UUID, error)
	InviteShoppingListUser(userId, shoppingListId, requesterId uuid.UUID) error
	AcceptShoppingListInvite(userId, shoppingListId uuid.UUID) error
	DeclineShoppingListInvite(userId, shoppingListId uuid.UUID) error
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
	remoteServices *services.Remote,
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

	mailService, err := mail.NewService(cfg)
	if err != nil {
		return nil, err
	}

	return &Service{
		ShoppingList: shopping_list.NewService(repo, remoteServices.Auth, mailService),
		MQ:           mq.NewService(repo, client),
	}, nil
}
