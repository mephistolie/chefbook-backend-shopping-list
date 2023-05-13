package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
	shoppingListFail "github.com/mephistolie/chefbook-backend-shopping-list/internal/entity/fail"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/repository/postgres/dto"
)

func (r *Repository) GetShoppingList(userId uuid.UUID) (entity.ShoppingList, error) {
	var shoppingList entity.ShoppingList
	var bsonPurchases []byte
	var purchases []dto.Purchase

	query := fmt.Sprintf(`
			SELECT purchases, version
			FROM %s
			WHERE user_id=$1
		`, shoppingListTable)

	row := r.db.QueryRow(query, userId)
	if err := row.Scan(&bsonPurchases, &shoppingList.Version); err != nil {
		log.Errorf("unable to get shopping list for user %s: %s", userId, err)
		return entity.ShoppingList{}, shoppingListFail.GrpcShoppingListNotFound
	}

	if err := json.Unmarshal(bsonPurchases, &purchases); err != nil {
		log.Warnf("unable to unmarshal shopping list for user %s: %s", userId, err)
		version, err := r.SetShoppingList(userId, []entity.Purchase{}, nil)
		return entity.ShoppingList{Version: version}, err
	}
	shoppingList.Purchases = dto.NewPurchasesEntity(purchases)

	return shoppingList, nil
}

func (r *Repository) SetShoppingList(userId uuid.UUID, purchases []entity.Purchase, lastVersion *int32) (int32, error) {
	query, shoppingListBSON, err := getSetShoppingListBaseQuery(purchases)
	if err != nil {
		return 0, err
	}

	var version int32
	if lastVersion != nil {
		query = query + " AND version=$3 RETURNING version"
		if err = r.db.Get(&version, query, shoppingListBSON, userId, *lastVersion); err != nil {
			log.Warnf("try to update shopping list with outdated version %s for user %s: %s", *lastVersion, userId, err)
			return 0, shoppingListFail.GrpcOutdatedVersion
		}
	} else {
		query = query + " RETURNING version"
		if err = r.db.Get(&version, query, shoppingListBSON, userId); err != nil {
			log.Errorf("unable to set shopping list for user %s: %s", userId, err)
			return 0, shoppingListFail.GrpcShoppingListNotFound
		}
	}

	return version, nil
}

func getSetShoppingListBaseQuery(purchases []entity.Purchase) (string, []byte, error) {
	shoppingListBSON, err := json.Marshal(dto.NewPurchasesDto(purchases))
	if err != nil {
		log.Errorf("unable to marshal shopping list: %s", err)
		return "", nil, shoppingListFail.GrpcShoppingListNotFound
	}

	setShoppingListQuery := fmt.Sprintf(`
			UPDATE %s
			SET purchases=$1, version=version+1
			WHERE user_id=$2
		`, shoppingListTable)

	return setShoppingListQuery, shoppingListBSON, nil
}
