package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/firebase"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/config"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/repository/grpc"
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
	GetShoppingLists(ctx context.Context, userId uuid.UUID) ([]entity.ShoppingListInfo, error)
	CreateSharedShoppingList(ctx context.Context, userId uuid.UUID, shoppingListId *uuid.UUID, name *string) (uuid.UUID, error)
	GetShoppingList(ctx context.Context, shoppingListId *uuid.UUID, userId uuid.UUID) (entity.ShoppingList, error)
	SetShoppingListName(ctx context.Context, shoppingListId, userId uuid.UUID, name *string) error
	SetShoppingList(ctx context.Context, input entity.ShoppingListInput) (int32, error)
	AddPurchasesToShoppingList(ctx context.Context, input entity.ShoppingListInput) (int32, error)
	DeleteSharedShoppingList(ctx context.Context, shoppingListId uuid.UUID, userId uuid.UUID) error

	GetShoppingListUsers(ctx context.Context, shoppingListId, requesterId uuid.UUID) ([]entity.User, error)
	GetShoppingListLink(ctx context.Context, shoppingListId, requesterId uuid.UUID, linkPattern string) (string, time.Time, error)
	JoinShoppingList(ctx context.Context, shoppingListId, userId, key uuid.UUID) error
	DeleteUserFromShoppingList(ctx context.Context, userId, shoppingListId, requesterId uuid.UUID) error
}

type MQ interface {
	CreatePersonalShoppingList(ctx context.Context, userId uuid.UUID, messageId uuid.UUID) error
	ImportFirebaseShoppingList(ctx context.Context, userId uuid.UUID, firebaseId string, messageId uuid.UUID) error
	DeleteUserShoppingLists(ctx context.Context, userId uuid.UUID, messageId uuid.UUID) error
}

func New(
	ctx context.Context,
	cfg *config.Config,
	repo repository.ShoppingList,
	grpc *grpc.Repository,
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
		ShoppingList: shopping_list.NewService(repo, grpc),
		MQ:           mq.NewService(repo, client),
	}, nil
}
