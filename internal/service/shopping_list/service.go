package shopping_list

import (
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/firebase"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/service/dependencies/repository"
	"time"
)

type Service struct {
	repo     repository.ShoppingList
	firebase *firebase.Client
}

func NewService(
	repo repository.ShoppingList,
) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetShoppingList(userId uuid.UUID) (entity.ShoppingList, error) {
	return s.repo.GetShoppingList(userId)
}

func (s *Service) SetShoppingList(userId uuid.UUID, purchases []entity.Purchase) error {
	shoppingList := entity.ShoppingList{
		Purchases: purchases,
		Timestamp: time.Now().UTC(),
	}
	return s.repo.SetShoppingList(userId, shoppingList)
}

func (s *Service) AddToShoppingList(userId uuid.UUID, purchases []entity.Purchase) error {
	shoppingList, err := s.repo.GetShoppingList(userId)
	if err != nil {
		return err
	}

	purchasesByIds := make(map[uuid.UUID]*entity.Purchase)
	purchasesByName := make(map[string]*entity.Purchase)

	for i := range shoppingList.Purchases {
		purchasesByIds[shoppingList.Purchases[i].Id] = &shoppingList.Purchases[i]
		purchasesByName[shoppingList.Purchases[i].Name] = &shoppingList.Purchases[i]
	}

	for i := range purchases {
		id := purchases[i].Id
		name := purchases[i].Name
		amount := purchases[i].Amount
		multiplier := purchases[i].Multiplier
		if amount != nil && *amount > 0 && purchasesByIds[id] != nil {
			totalAmount := *amount
			if (*purchasesByIds[id]).Amount != nil {
				totalAmount += *((*purchasesByIds[id]).Amount)
			}
			(*purchasesByIds[id]).Amount = &totalAmount
		} else if purchasesByName[name] != nil && multiplier > 0 {
			(*purchasesByIds[id]).Multiplier += multiplier
		} else {
			shoppingList.Purchases = append(shoppingList.Purchases, purchases[i])
		}
	}
	shoppingList.Timestamp = time.Now().UTC()

	return s.repo.SetShoppingList(userId, shoppingList)
}
