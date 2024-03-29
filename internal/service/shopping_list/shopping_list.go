package shopping_list

import (
	"context"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/utils/slices"
	api "github.com/mephistolie/chefbook-backend-recipe/api/proto/implementation/v1"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
	shoppingListFail "github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity/fail"
	"sync"
	"time"
)

func (s *Service) GetShoppingLists(userId uuid.UUID) ([]entity.ShoppingListInfo, error) {
	shoppingLists, err := s.repo.GetShoppingLists(userId)
	if err != nil {
		return nil, err
	}

	var rawIds []string
	for _, shoppingList := range shoppingLists {
		rawIds = append(rawIds, shoppingList.Owner.Id.String())
	}

	profiles := s.getProfilesInfo(rawIds)
	for i := range shoppingLists {
		if profile, ok := profiles[shoppingLists[i].Owner.Id.String()]; ok {
			shoppingLists[i].Owner.Name = profile.VisibleName
			shoppingLists[i].Owner.Avatar = profile.Avatar
		}
	}

	return shoppingLists, nil
}

func (s *Service) CreateSharedShoppingList(userId uuid.UUID, shoppingListId *uuid.UUID, name *string) (uuid.UUID, error) {
	return s.repo.CreateSharedShoppingList(userId, shoppingListId, name)
}

func (s *Service) GetShoppingList(shoppingListId *uuid.UUID, userId uuid.UUID) (entity.ShoppingList, error) {
	if shoppingListId == nil {
		return s.getPersonalShoppingList(userId)
	}

	shoppingList, err := s.repo.GetShoppingList(*shoppingListId)
	if err != nil {
		return entity.ShoppingList{}, err
	}

	if userId != shoppingList.Owner.Id {
		if err = s.checkUserHasAccessToShoppingList(userId, *shoppingListId, false); err != nil {
			return entity.ShoppingList{}, err
		}
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		profiles := s.getProfilesInfo([]string{shoppingList.Owner.Id.String()})
		if profile, ok := profiles[shoppingList.Owner.Id.String()]; ok {
			shoppingList.Owner.Name = profile.VisibleName
			shoppingList.Owner.Avatar = profile.Avatar
		}
		wg.Done()
	}()
	go func() {
		shoppingList.RecipeNames = s.getRecipeNames(shoppingList.Purchases, userId)
		wg.Done()
	}()
	wg.Wait()

	return shoppingList, nil
}

func (s *Service) getPersonalShoppingList(userId uuid.UUID) (entity.ShoppingList, error) {
	id, err := s.repo.GetPersonalShoppingListId(userId)
	if err != nil {
		return entity.ShoppingList{}, err
	}

	shoppingList, err := s.repo.GetShoppingList(id)
	if err != nil {
		return entity.ShoppingList{}, err
	}

	shoppingList.RecipeNames = s.getRecipeNames(shoppingList.Purchases, userId)

	return shoppingList, nil
}

func (s *Service) SetShoppingListName(shoppingListId, userId uuid.UUID, name *string) error {
	return s.repo.SetShoppingListName(shoppingListId, userId, name)
}

func (s *Service) SetShoppingList(input entity.ShoppingListInput) (int32, error) {
	if input.ShoppingListId == nil {
		return s.SetPersonalShoppingList(input)
	}

	if err := s.checkUserHasAccessToShoppingList(input.EditorId, *input.ShoppingListId, true); err != nil {
		return 0, err
	}

	return s.repo.SetShoppingList(input)
}

func (s *Service) SetPersonalShoppingList(input entity.ShoppingListInput) (int32, error) {
	id, err := s.repo.GetPersonalShoppingListId(input.EditorId)
	if err != nil {
		return 0, err
	}
	input.ShoppingListId = &id
	return s.repo.SetShoppingList(input)
}

func (s *Service) AddPurchasesToShoppingList(input entity.ShoppingListInput) (int32, error) {
	var shoppingList entity.ShoppingList
	var err error = nil
	if input.ShoppingListId == nil {
		shoppingList, err = s.getPersonalShoppingList(input.EditorId)
	} else {
		if shoppingList, err = s.repo.GetShoppingList(*input.ShoppingListId); err == nil && input.EditorId != shoppingList.Owner.Id {
			err = s.checkUserHasAccessToShoppingList(input.EditorId, *input.ShoppingListId, false)
		}
	}
	if err != nil {
		return 0, err
	}

	if input.LastVersion != nil && shoppingList.Version != *input.LastVersion {
		return 0, shoppingListFail.GrpcOutdatedVersion
	}

	purchases := concatenateShoppingLists(shoppingList.Purchases, input.Purchases)
	concatenatedInput := entity.ShoppingListInput{
		ShoppingListId: input.ShoppingListId,
		EditorId:       input.EditorId,
		Purchases:      purchases,
		LastVersion:    input.LastVersion,
	}

	return s.repo.SetShoppingList(concatenatedInput)
}

func concatenateShoppingLists(
	oldPurchases []entity.Purchase,
	newPurchases []entity.Purchase,
) []entity.Purchase {
	var purchases []entity.Purchase

	oldPurchasesByIds := make(map[uuid.UUID]*entity.Purchase)
	oldPurchasesByName := make(map[string]*entity.Purchase)

	for i, oldPurchase := range oldPurchases {
		purchases = append(purchases, oldPurchase)
		oldPurchasesByIds[purchases[i].Id] = &purchases[i]
		oldPurchasesByName[purchases[i].Name] = &purchases[i]
	}

	for _, newPurchase := range newPurchases {
		id := newPurchase.Id
		name := newPurchase.Name
		multiplier := newPurchase.Multiplier
		amount := newPurchase.Amount

		oldPurchase := oldPurchasesByIds[id]
		if oldPurchase == nil {
			oldPurchase = oldPurchasesByName[name]
			if oldPurchase == nil || !isSamePurchase(*oldPurchase, newPurchase) {
				purchases = append(purchases, newPurchase)
				log.Debugf("purchase %s not found in shopping list; adding..", newPurchase.Id)
				continue
			}
		}
		log.Debugf("found existing purchase %s on add purchase to shopping list...", oldPurchase.Id)

		if multiplier != nil && *multiplier > 0 {
			log.Debugf("sum multipliers for purchase %s...", oldPurchase.Id)
			totalMultiplier := *multiplier
			log.Debugf("new multiplier for purchase %s is %s", oldPurchase.Id, totalMultiplier)
			if oldPurchase.Multiplier != nil {
				log.Debugf("old multiplier for purchase %s is %s", oldPurchase.Id, *oldPurchase.Multiplier)
				totalMultiplier += *oldPurchase.Multiplier
			}
			(*oldPurchase).Multiplier = &totalMultiplier
		}
		if amount != nil && *amount > 0 {
			log.Debugf("sum amounts for purchase %s...", oldPurchase.Id)
			totalAmount := *amount
			log.Debugf("new amount for purchase %s is %s", oldPurchase.Id, totalAmount)
			if oldPurchase.Amount != nil {
				log.Debugf("old amount for purchase %s is %s", oldPurchase.Id, *oldPurchase.Amount)
				totalAmount += *oldPurchase.Amount
			}
			(*oldPurchase).Amount = &totalAmount
		}
	}

	return purchases
}

func isSamePurchase(first, second entity.Purchase) bool {
	firstMeasureUnit, secondMeasureUnit := "", ""
	if first.MeasureUnit != nil {
		firstMeasureUnit = *first.MeasureUnit
	}
	if second.MeasureUnit != nil {
		secondMeasureUnit = *second.MeasureUnit
	}
	if firstMeasureUnit != secondMeasureUnit {
		return false
	}

	firstRecipeId, secondRecipeId := "", ""
	if first.RecipeId != nil {
		firstRecipeId = first.RecipeId.String()
	}
	if second.RecipeId != nil {
		secondRecipeId = second.RecipeId.String()
	}
	if firstRecipeId != secondRecipeId {
		return false
	}

	return true
}

func (s *Service) DeleteSharedShoppingList(shoppingListId uuid.UUID, userId uuid.UUID) error {
	if err := s.checkUserIsShoppingListOwner(userId, shoppingListId); err != nil {
		return err
	}
	return s.repo.DeleteSharedShoppingList(shoppingListId)
}

func (s *Service) getRecipeNames(purchases []entity.Purchase, userId uuid.UUID) map[string]string {

	var recipeIds []string
	for _, purchase := range purchases {
		if purchase.RecipeId != nil {
			recipeIds = append(recipeIds, purchase.RecipeId.String())
		}
	}
	recipeIds = slices.RemoveDuplicates(recipeIds)

	ctx, cancelCtx := context.WithTimeout(context.Background(), 3*time.Second)
	res, err := s.grpc.Recipe.GetRecipeNames(ctx, &api.GetRecipeNamesRequest{
		RecipeIds: recipeIds,
		UserId:    userId.String(),
	})
	cancelCtx()

	if err != nil {
		log.Debugf("unable to get recipe names: %s", err)
		return map[string]string{}
	}
	log.Debugf("got recipe names: %s", res.RecipeNames)
	return res.RecipeNames
}
