package shopping_list

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
	shoppingListFail "github.com/mephistolie/chefbook-backend-shopping-list/internal/entity/fail"
)

func (s *Service) GetShoppingLists(userId uuid.UUID) ([]entity.ShoppingListInfo, error) {
	return s.repo.GetShoppingLists(userId, false)
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

	if userId != shoppingList.OwnerId {
		if err = s.checkUserHasAccessToShoppingList(userId, *shoppingListId, false); err != nil {
			return entity.ShoppingList{}, err
		}
	}

	return shoppingList, nil
}

func (s *Service) getPersonalShoppingList(userId uuid.UUID) (entity.ShoppingList, error) {
	id, err := s.repo.GetPersonalShoppingListId(userId)
	if err != nil {
		return entity.ShoppingList{}, err
	}
	return s.repo.GetShoppingList(id)
}

func (s *Service) SetShoppingListName(shoppingListId *uuid.UUID, name *string, requesterId uuid.UUID) error {
	if shoppingListId == nil {
		id, err := s.repo.GetPersonalShoppingListId(requesterId)
		if err != nil {
			return err
		}
		shoppingListId = &id
	}

	if err := s.checkUserIsShoppingListOwner(requesterId, *shoppingListId); err != nil {
		return err
	}
	return s.repo.SetShoppingListName(*shoppingListId, name)
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
		if shoppingList, err = s.repo.GetShoppingList(*input.ShoppingListId); err == nil && input.EditorId != shoppingList.OwnerId {
			err = s.checkUserHasAccessToShoppingList(input.EditorId, *input.ShoppingListId, false)
		}
	}
	if err != nil {
		return 0, err
	}

	if input.LastVersion != nil && shoppingList.Version != *input.LastVersion {
		return 0, shoppingListFail.GrpcOutdatedVersion
	}

	purchases, recipeNames := concatenateShoppingLists(shoppingList.Purchases, input.Purchases, shoppingList.RecipeNames, input.RecipeNames)
	concatenatedInput := entity.ShoppingListInput{
		ShoppingListId: input.ShoppingListId,
		EditorId:       input.EditorId,
		Purchases:      purchases,
		RecipeNames:    recipeNames,
		LastVersion:    input.LastVersion,
	}

	return s.repo.SetShoppingList(concatenatedInput)
}

func concatenateShoppingLists(
	oldPurchases []entity.Purchase,
	newPurchases []entity.Purchase,
	oldRecipeNames entity.RecipeNames,
	newRecipeNames entity.RecipeNames,
) ([]entity.Purchase, entity.RecipeNames) {
	var purchases []entity.Purchase
	recipeNames := oldRecipeNames

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
		amount := newPurchase.Amount
		multiplier := newPurchase.Multiplier

		oldPurchase := oldPurchasesByIds[id]
		if oldPurchase == nil {
			oldPurchase = oldPurchasesByName[name]
			if oldPurchase == nil {
				purchases = append(purchases, newPurchase)
			}
		}

		if amount != nil && *amount > 0 {
			totalAmount := *amount
			if oldPurchase.Amount != nil {
				totalAmount += *oldPurchase.Amount
			}
			oldPurchase.Amount = &totalAmount
		}
		if multiplier != nil && *multiplier > 0 {
			totalMultiplier := *multiplier
			if oldPurchase.Multiplier != nil {
				totalMultiplier += *oldPurchase.Multiplier
			}
			oldPurchase.Amount = &totalMultiplier
		}
	}

	for recipeId, recipeName := range newRecipeNames {
		recipeNames[recipeId] = recipeName
	}

	return purchases, recipeNames
}

func (s *Service) DeleteSharedShoppingList(shoppingListId uuid.UUID, userId uuid.UUID) error {
	if err := s.checkUserIsShoppingListOwner(userId, shoppingListId); err != nil {
		return err
	}
	return s.repo.DeleteSharedShoppingList(shoppingListId)
}