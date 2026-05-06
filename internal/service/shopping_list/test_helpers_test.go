package shopping_list

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	profileApi "github.com/mephistolie/chefbook-backend-profile/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	log.InitWithService("shopping-list", "", false)
	os.Exit(m.Run())
}

type fakeShoppingListRepo struct {
	lists        map[uuid.UUID]entity.ShoppingList
	personalIds  map[uuid.UUID]uuid.UUID
	owners       map[uuid.UUID]uuid.UUID
	users        map[uuid.UUID][]uuid.UUID
	listTypes    map[uuid.UUID]string
	keys         map[uuid.UUID]uuid.UUID
	setErr       error
	setVersion   int32
	setCalls     int
	lastSetInput entity.ShoppingListInput
	deletedUser  uuid.UUID
	deletedList  uuid.UUID
}

func (r *fakeShoppingListRepo) CreatePersonalShoppingList(context.Context, uuid.UUID, uuid.UUID) error {
	return errors.New("not implemented")
}

func (r *fakeShoppingListRepo) ImportFirebaseShoppingList(context.Context, uuid.UUID, []entity.Purchase, uuid.UUID) error {
	return errors.New("not implemented")
}

func (r *fakeShoppingListRepo) DeleteUserShoppingLists(context.Context, uuid.UUID, uuid.UUID) error {
	return errors.New("not implemented")
}

func (r *fakeShoppingListRepo) GetShoppingLists(context.Context, uuid.UUID) ([]entity.ShoppingListInfo, error) {
	return nil, errors.New("not implemented")
}

func (r *fakeShoppingListRepo) SetShoppingListName(context.Context, uuid.UUID, uuid.UUID, *string) error {
	return errors.New("not implemented")
}

func (r *fakeShoppingListRepo) CreateSharedShoppingList(context.Context, uuid.UUID, *uuid.UUID, *string) (uuid.UUID, error) {
	return uuid.Nil, errors.New("not implemented")
}

func (r *fakeShoppingListRepo) GetShoppingList(_ context.Context, shoppingListId, _ uuid.UUID) (entity.ShoppingList, error) {
	if list, ok := r.lists[shoppingListId]; ok {
		return list, nil
	}
	return entity.ShoppingList{}, errors.New("shopping list not found")
}

func (r *fakeShoppingListRepo) GetPersonalShoppingListId(_ context.Context, userId uuid.UUID) (uuid.UUID, error) {
	if id, ok := r.personalIds[userId]; ok {
		return id, nil
	}
	return uuid.Nil, errors.New("personal shopping list not found")
}

func (r *fakeShoppingListRepo) GetShoppingListType(_ context.Context, shoppingListId uuid.UUID) (string, error) {
	if listType, ok := r.listTypes[shoppingListId]; ok {
		return listType, nil
	}
	return "", errors.New("shopping list type not found")
}

func (r *fakeShoppingListRepo) GetShoppingListOwner(_ context.Context, shoppingListId uuid.UUID) (uuid.UUID, error) {
	if owner, ok := r.owners[shoppingListId]; ok {
		return owner, nil
	}
	if list, ok := r.lists[shoppingListId]; ok {
		return list.Owner.Id, nil
	}
	return uuid.Nil, errors.New("owner not found")
}

func (r *fakeShoppingListRepo) SetShoppingList(_ context.Context, input entity.ShoppingListInput) (int32, error) {
	r.setCalls += 1
	r.lastSetInput = input
	if r.setErr != nil {
		return 0, r.setErr
	}
	if r.setVersion == 0 {
		r.setVersion = 1
	}
	return r.setVersion, nil
}

func (r *fakeShoppingListRepo) DeleteSharedShoppingList(context.Context, uuid.UUID) error {
	return errors.New("not implemented")
}

func (r *fakeShoppingListRepo) GetShoppingListUsers(_ context.Context, shoppingListId uuid.UUID) ([]uuid.UUID, error) {
	if users, ok := r.users[shoppingListId]; ok {
		return users, nil
	}
	return nil, errors.New("shopping list users not found")
}

func (r *fakeShoppingListRepo) GetShoppingListKey(_ context.Context, shoppingListId uuid.UUID) (uuid.UUID, time.Time, error) {
	if key, ok := r.keys[shoppingListId]; ok {
		return key, time.Now().Add(time.Hour), nil
	}
	return uuid.Nil, time.Time{}, errors.New("shopping list key not found")
}

func (r *fakeShoppingListRepo) IsShoppingListKeyValid(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	return false, errors.New("not implemented")
}

func (r *fakeShoppingListRepo) AddUserToShoppingList(context.Context, uuid.UUID, uuid.UUID) error {
	return errors.New("not implemented")
}

func (r *fakeShoppingListRepo) DeleteUserFromShoppingList(_ context.Context, userId, shoppingListId uuid.UUID) error {
	r.deletedUser = userId
	r.deletedList = shoppingListId
	return nil
}

type fakeProfileClient struct {
	infos map[string]*profileApi.ProfileMinInfo
}

func (c *fakeProfileClient) GetProfilesMinInfo(
	context.Context,
	*profileApi.GetProfilesMinInfoRequest,
	...grpc.CallOption,
) (*profileApi.GetProfilesMinInfoResponse, error) {
	return &profileApi.GetProfilesMinInfoResponse{Infos: c.infos}, nil
}

func (c *fakeProfileClient) GetProfile(
	context.Context,
	*profileApi.GetProfileRequest,
	...grpc.CallOption,
) (*profileApi.GetProfileResponse, error) {
	return nil, errors.New("not implemented")
}
