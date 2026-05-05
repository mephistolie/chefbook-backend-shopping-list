package shopping_list

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
	shoppingListFail "github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity/fail"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAddPurchasesToShoppingListMergesExistingPurchases(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	listId := uuid.New()
	recipeId := uuid.New()
	existingPurchaseId := uuid.New()
	differentMeasurePurchaseId := uuid.New()
	measureUnit := "g"
	otherMeasureUnit := "ml"
	oldMultiplier := int32(2)
	newMultiplier := int32(3)
	oldAmount := float32(100)
	newAmount := float32(50)
	ignoredMultiplier := int32(0)
	ignoredAmount := float32(-25)
	version := int32(7)

	repo := &fakeShoppingListRepo{
		lists: map[uuid.UUID]entity.ShoppingList{
			listId: {
				Id:      listId,
				Owner:   entity.User{Id: ownerId},
				Version: version,
				Purchases: []entity.Purchase{
					{
						Id:          existingPurchaseId,
						Name:        "Flour",
						Multiplier:  &oldMultiplier,
						Amount:      &oldAmount,
						MeasureUnit: &measureUnit,
						RecipeId:    &recipeId,
					},
				},
			},
		},
	}
	service := NewService(repo, nil)

	nextVersion, err := service.AddPurchasesToShoppingList(ctx, entity.ShoppingListInput{
		ShoppingListId: &listId,
		EditorId:       ownerId,
		LastVersion:    &version,
		Purchases: []entity.Purchase{
			{
				Id:          uuid.New(),
				Name:        "Flour",
				Multiplier:  &newMultiplier,
				Amount:      &newAmount,
				MeasureUnit: &measureUnit,
				RecipeId:    &recipeId,
			},
			{
				Id:          differentMeasurePurchaseId,
				Name:        "Flour",
				MeasureUnit: &otherMeasureUnit,
			},
			{
				Id:         existingPurchaseId,
				Name:       "Flour",
				Multiplier: &ignoredMultiplier,
				Amount:     &ignoredAmount,
			},
		},
	})
	if err != nil {
		t.Fatalf("AddPurchasesToShoppingList returned error: %v", err)
	}
	if nextVersion != repo.setVersion {
		t.Fatalf("expected next version %d, got %d", repo.setVersion, nextVersion)
	}
	if repo.setCalls != 1 {
		t.Fatalf("expected SetShoppingList to be called once, got %d", repo.setCalls)
	}

	got := repo.lastSetInput.Purchases
	if len(got) != 2 {
		t.Fatalf("expected 2 purchases after merge, got %d: %#v", len(got), got)
	}
	if got[0].Id != existingPurchaseId {
		t.Fatalf("expected first purchase to keep original id %s, got %s", existingPurchaseId, got[0].Id)
	}
	if got[0].Multiplier == nil || *got[0].Multiplier != oldMultiplier+newMultiplier {
		t.Fatalf("expected multiplier to be summed to %d, got %#v", oldMultiplier+newMultiplier, got[0].Multiplier)
	}
	if got[0].Amount == nil || *got[0].Amount != oldAmount+newAmount {
		t.Fatalf("expected amount to be summed to %f, got %#v", oldAmount+newAmount, got[0].Amount)
	}
	if got[1].Id != differentMeasurePurchaseId {
		t.Fatalf("expected same name with different measure unit to stay separate")
	}
}

func TestAddPurchasesToShoppingListRejectsOutdatedVersion(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	listId := uuid.New()
	inputVersion := int32(3)

	repo := &fakeShoppingListRepo{
		lists: map[uuid.UUID]entity.ShoppingList{
			listId: {
				Id:      listId,
				Owner:   entity.User{Id: ownerId},
				Version: 4,
				Purchases: []entity.Purchase{
					{Id: uuid.New(), Name: "Milk"},
				},
			},
		},
	}
	service := NewService(repo, nil)

	_, err := service.AddPurchasesToShoppingList(ctx, entity.ShoppingListInput{
		ShoppingListId: &listId,
		EditorId:       ownerId,
		LastVersion:    &inputVersion,
		Purchases:      []entity.Purchase{{Id: uuid.New(), Name: "Eggs"}},
	})
	if status.Code(err) != codes.FailedPrecondition {
		t.Fatalf("expected outdated version error code %s, got %s: %v", codes.FailedPrecondition, status.Code(err), err)
	}
	if status.Convert(err).Message() != status.Convert(shoppingListFail.GrpcOutdatedVersion).Message() {
		t.Fatalf("expected outdated version message, got %q", status.Convert(err).Message())
	}
	if repo.setCalls != 0 {
		t.Fatalf("expected SetShoppingList not to be called, got %d calls", repo.setCalls)
	}
}

func TestSetShoppingListRequiresSharedListAccess(t *testing.T) {
	ctx := context.Background()
	editorId := uuid.New()
	listId := uuid.New()
	repo := &fakeShoppingListRepo{
		owners: map[uuid.UUID]uuid.UUID{listId: uuid.New()},
		users:  map[uuid.UUID][]uuid.UUID{listId: {uuid.New()}},
	}
	service := NewService(repo, nil)

	_, err := service.SetShoppingList(ctx, entity.ShoppingListInput{
		ShoppingListId: &listId,
		EditorId:       editorId,
		Purchases:      []entity.Purchase{{Id: uuid.New(), Name: "Bread"}},
	})
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("expected access denied, got %s: %v", status.Code(err), err)
	}
	if repo.setCalls != 0 {
		t.Fatalf("expected SetShoppingList not to be called, got %d calls", repo.setCalls)
	}
}
