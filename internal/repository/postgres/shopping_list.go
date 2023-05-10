package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
	shoppingListFail "github.com/mephistolie/chefbook-backend-shopping-list/internal/entity/fail"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/repository/postgres/dto"
	"time"
)

func (r *Repository) GetShoppingList(userId uuid.UUID) (entity.ShoppingList, error) {
	var shoppingList dto.ShoppingList
	var shoppingListBSON []byte

	getShoppingListQuery := fmt.Sprintf(`
			SELECT shopping_list
			FROM %s
			WHERE user_id=$1
		`, shoppingListTable)

	if err := r.db.Get(&shoppingListBSON, getShoppingListQuery, userId); err != nil {
		log.Errorf("unable to get shopping list for user %s: %s", userId, err)
		return entity.ShoppingList{}, shoppingListFail.GrpcShoppingListNotFound
	}

	if err := json.Unmarshal(shoppingListBSON, &shoppingList); err != nil {
		log.Warnf("unable to unmarshal shopping list for user %s: %s", userId, err)
		emptyShoppingList := emptyShoppingList()
		_ = r.SetShoppingList(userId, emptyShoppingList)
		return emptyShoppingList, nil
	}

	return shoppingList.Entity(), nil
}

func (r *Repository) SetShoppingList(userId uuid.UUID, shoppingList entity.ShoppingList) error {
	var shoppingListBSON, err = json.Marshal(dto.NewShoppingList(shoppingList))
	if err != nil {
		log.Errorf("unable to marshal shopping list for user %s: %s", userId, err)
		return shoppingListFail.GrpcShoppingListNotFound
	}

	setShoppingListQuery := fmt.Sprintf(`
			UPDATE %s
			SET shopping_list=$1
			WHERE user_id=$2
		`, shoppingListTable)

	if _, err = r.db.Exec(setShoppingListQuery, shoppingListBSON, userId); err != nil {
		log.Errorf("unable to set shopping list for user %s: %s", userId, err)
		return shoppingListFail.GrpcShoppingListNotFound
	}

	return nil
}

func emptyShoppingList() entity.ShoppingList {
	return entity.ShoppingList{
		Timestamp: time.Now(),
	}
}
