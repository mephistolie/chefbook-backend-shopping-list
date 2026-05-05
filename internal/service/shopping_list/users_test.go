package shopping_list

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/uuid"
	profileApi "github.com/mephistolie/chefbook-backend-profile/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
	grpcRepo "github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/repository/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeleteUserFromShoppingListProtectsOwnership(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	memberId := uuid.New()
	listId := uuid.New()

	tests := []struct {
		name        string
		userId      uuid.UUID
		requesterId uuid.UUID
		wantCode    codes.Code
		wantDeleted bool
	}{
		{
			name:        "non-owner requester cannot delete users",
			userId:      memberId,
			requesterId: memberId,
			wantCode:    codes.PermissionDenied,
		},
		{
			name:        "owner cannot delete self",
			userId:      ownerId,
			requesterId: ownerId,
			wantCode:    codes.InvalidArgument,
		},
		{
			name:        "owner deletes member",
			userId:      memberId,
			requesterId: ownerId,
			wantDeleted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeShoppingListRepo{owners: map[uuid.UUID]uuid.UUID{listId: ownerId}}
			service := NewService(repo, nil)

			err := service.DeleteUserFromShoppingList(ctx, tt.userId, listId, tt.requesterId)
			if tt.wantCode != codes.OK {
				if status.Code(err) != tt.wantCode {
					t.Fatalf("expected code %s, got %s: %v", tt.wantCode, status.Code(err), err)
				}
				if repo.deletedUser != uuid.Nil {
					t.Fatalf("expected no user deletion, got %s", repo.deletedUser)
				}
				return
			}
			if err != nil {
				t.Fatalf("DeleteUserFromShoppingList returned error: %v", err)
			}
			if !tt.wantDeleted || repo.deletedUser != memberId || repo.deletedList != listId {
				t.Fatalf("expected member deletion, got user=%s list=%s", repo.deletedUser, repo.deletedList)
			}
		})
	}
}

func TestGetShoppingListUsersReturnsUsersWithProfileData(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	memberId := uuid.New()
	listId := uuid.New()
	ownerName := "Owner"
	memberName := "Member"
	ownerAvatar := "owner.png"

	repo := &fakeShoppingListRepo{
		owners: map[uuid.UUID]uuid.UUID{listId: ownerId},
		users:  map[uuid.UUID][]uuid.UUID{listId: {ownerId, memberId}},
	}
	grpc := &grpcRepo.Repository{
		Profile: &grpcRepo.Profile{ProfileServiceClient: &fakeProfileClient{
			infos: map[string]*profileApi.ProfileMinInfo{
				ownerId.String():  {VisibleName: &ownerName, Avatar: &ownerAvatar},
				memberId.String(): {VisibleName: &memberName},
			},
		}},
	}
	service := NewService(repo, grpc)

	users, err := service.GetShoppingListUsers(ctx, listId, ownerId)
	if err != nil {
		t.Fatalf("GetShoppingListUsers returned error: %v", err)
	}

	want := []entity.User{
		{Id: ownerId, Name: &ownerName, Avatar: &ownerAvatar},
		{Id: memberId, Name: &memberName},
	}
	if !reflect.DeepEqual(users, want) {
		t.Fatalf("unexpected users:\nwant: %#v\n got: %#v", want, users)
	}
}
