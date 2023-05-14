package postgres

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/entity"
	shoppingListFail "github.com/mephistolie/chefbook-backend-shopping-list/internal/entity/fail"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/repository/postgres/dto"
)

func (r *Repository) GetShoppingLists(userId uuid.UUID, onlyPending bool) ([]entity.ShoppingListInfo, error) {
	var shoppingLists []entity.ShoppingListInfo

	query := fmt.Sprintf(`
			SELECT %[1]v.shopping_list_id, %[2]v.name, %[2]v.type, %[2]v.owner_id
			FROM %[1]v
			LEFT JOIN %[2]v ON %[1]v.shopping_list_id=%[2]v.shopping_list_id
			WHERE %[1]v.user_id=$1
		`, usersTable, shoppingListsTable)
	if onlyPending {
		query = query + fmt.Sprintf(" AND %s.accepted=false", usersTable)
	}

	rows, err := r.db.Query(query, userId)
	if err != nil {
		log.Errorf("unable to get shopping lists for user %s: %s", userId, err)
		return nil, fail.GrpcUnknown
	}

	for rows.Next() {
		shoppingList := entity.ShoppingListInfo{}
		if err = rows.Scan(&shoppingList.Id, &shoppingList.Name, &shoppingList.Type, &shoppingList.OwnerId); err != nil {
			log.Warnf("unable to parse shopping list info: ", err)
			continue
		}
		shoppingLists = append(shoppingLists, shoppingList)
	}

	return shoppingLists, nil
}

func (r *Repository) CreateSharedShoppingList(userId uuid.UUID, shoppingListId *uuid.UUID, name *string) (uuid.UUID, error) {
	var id uuid.UUID
	if shoppingListId != nil {
		id = *shoppingListId
	} else {
		id = uuid.New()
	}

	tx, err := r.startTransaction()
	if err != nil {
		return uuid.UUID{}, err
	}

	if err = r.ensureShoppingListsLimit(tx, userId); err != nil {
		return uuid.UUID{}, err
	}

	bsonShoppingList, err := json.Marshal([]dto.Purchase{})
	if err != nil {
		log.Errorf("unable to marshal purchases for user %s: %s", userId, err)
		return uuid.UUID{}, errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	bsonRecipeNames, err := json.Marshal(entity.RecipeNames{})
	if err != nil {
		log.Errorf("unable to marshal recipe names for user %s: %s", userId, err)
		return uuid.UUID{}, errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	query := fmt.Sprintf(`
			INSERT INTO %s (shopping_list_id, name, type, owner_id, purchases, recipe_names)
			VALUES ($1, $2, 'shared', $3, $4, $5)
		`, shoppingListsTable)

	if _, err := tx.Exec(query, id, name, userId, bsonShoppingList, bsonRecipeNames); err != nil {
		log.Errorf("unable to create shared shopping list for user %s: %s", userId, err)
		return uuid.UUID{}, errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	return id, commitTransaction(tx)
}

func (r *Repository) ensureShoppingListsLimit(tx *sql.Tx, userId uuid.UUID) error {
	var count int
	getShoppingListsCountQuery := fmt.Sprintf(`
			SELECT count(shopping_list_id)
			FROM %s
			WHERE owner_id=$1
		`, shoppingListsTable)

	row := tx.QueryRow(getShoppingListsCountQuery)
	if err := row.Scan(&count); err != nil {
		log.Errorf("unable to get shopping lists count for user %s: %s", userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}
	if count >= r.maxShoppingListsCount {
		log.Warnf("user %s tries to create new shopping list over maximum count %s", userId, r.maxShoppingListsCount)
		return errorWithTransactionRollback(tx, shoppingListFail.GrpcMaxShoppingListsCount(count))
	}

	return nil
}

func (r *Repository) GetShoppingList(shoppingListId uuid.UUID) (entity.ShoppingList, error) {
	shoppingList := entity.ShoppingList{Id: shoppingListId}
	var bsonPurchases []byte
	var bsonRecipeNames []byte
	var purchases []dto.Purchase

	query := fmt.Sprintf(`
			SELECT name, type, purchases, recipe_names, owner_id, version
			FROM %s
			WHERE shopping_list_id=$1
		`, shoppingListsTable)

	row := r.db.QueryRow(query, shoppingListId)
	if err := row.Scan(&shoppingList.Name, &shoppingList.Type, &bsonPurchases, &bsonRecipeNames, &shoppingList.OwnerId,
		&shoppingList.Version); err != nil {
		log.Errorf("unable to get shopping list %s: %s", shoppingListId, err)
		return entity.ShoppingList{}, shoppingListFail.GrpcShoppingListNotFound
	}

	if err := json.Unmarshal(bsonPurchases, &purchases); err != nil {
		log.Warnf("unable to unmarshal shopping list %s purchases: %s", shoppingListId, err)
		shoppingList.Version, err = r.SetShoppingList(entity.ShoppingListInput{ShoppingListId: &shoppingListId})
		return shoppingList, err
	}
	shoppingList.Purchases = dto.NewPurchasesEntity(purchases)

	if err := json.Unmarshal(bsonRecipeNames, &shoppingList.RecipeNames); err != nil {
		log.Warnf("unable to unmarshal shopping list %s recipe names: %s", shoppingListId, err)
		shoppingList.Version, err = r.SetShoppingList(entity.ShoppingListInput{ShoppingListId: &shoppingListId})
		return shoppingList, err
	}

	return shoppingList, nil
}

func (r *Repository) GetPersonalShoppingListId(userId uuid.UUID) (uuid.UUID, error) {
	var id uuid.UUID

	query := fmt.Sprintf(`
			SELECT shopping_list_id
			FROM %s
			WHERE owner_id=$1 AND type='personal'
		`, shoppingListsTable)

	if err := r.db.Get(&id, query, userId); err != nil {
		log.Errorf("unable to get personal shopping list id for user %s: %s", userId, err)
		return uuid.UUID{}, shoppingListFail.GrpcShoppingListNotFound
	}

	return id, nil
}

func (r *Repository) SetShoppingListName(shoppingListId uuid.UUID, name *string) error {
	query := fmt.Sprintf(`
			UPDATE %s
			SET name=$1
			WHERE shopping_list_id=$2
		`, shoppingListsTable)

	if _, err := r.db.Exec(query, name, shoppingListId); err != nil {
		log.Errorf("unable to set shopping list %s name: %s", shoppingListId, err)
		return fail.GrpcUnknown
	}

	return nil
}

func (r *Repository) SetShoppingList(input entity.ShoppingListInput) (int32, error) {
	query, bsonShoppingList, bsonRecipeNames, err := getSetShoppingListBaseQuery(input.Purchases, input.RecipeNames)
	if err != nil {
		return 0, err
	}

	var version int32
	if input.LastVersion != nil {
		query = query + " AND version=$4 RETURNING version"
		if err = r.db.Get(&version, query, bsonShoppingList, bsonRecipeNames, *input.ShoppingListId, *input.LastVersion); err != nil {
			log.Warnf("try to update shopping list %s with outdated version %s: %s", *input.ShoppingListId, *input.LastVersion, err)
			return 0, shoppingListFail.GrpcOutdatedVersion
		}
	} else {
		query = query + " RETURNING version"
		if err = r.db.Get(&version, query, bsonShoppingList, bsonRecipeNames, *input.ShoppingListId); err != nil {
			log.Errorf("unable to set shopping list %s: %s", *input.ShoppingListId, err)
			return 0, shoppingListFail.GrpcShoppingListNotFound
		}
	}

	return version, nil
}

func (r *Repository) DeleteSharedShoppingList(shoppingListId uuid.UUID) error {
	query := fmt.Sprintf(`
			DELETE FROM %s
			WHERE shopping_list_id=$1 AND type='shared' 
		`, shoppingListsTable)

	if _, err := r.db.Exec(query, shoppingListId); err != nil {
		log.Errorf("unable to delete shared shopping list %s: %s", shoppingListId, err)
		return fail.GrpcUnknown
	}

	return nil
}

func getSetShoppingListBaseQuery(purchases []entity.Purchase, recipeNames entity.RecipeNames) (string, []byte, []byte, error) {
	bsonPurchases, err := json.Marshal(dto.NewPurchasesDto(purchases))
	if err != nil {
		log.Errorf("unable to marshal shopping list: %s", err)
		return "", nil, nil, shoppingListFail.GrpcShoppingListNotFound
	}

	bsonRecipeNames, err := json.Marshal(recipeNames)
	if err != nil {
		log.Errorf("unable to marshal shopping list: %s", err)
		return "", nil, nil, shoppingListFail.GrpcShoppingListNotFound
	}

	setShoppingListQuery := fmt.Sprintf(`
			UPDATE %s
			SET purchases=$1, recipe_names=$2, version=version+1
			WHERE shopping_list_id=$3
		`, shoppingListsTable)

	return setShoppingListQuery, bsonPurchases, bsonRecipeNames, nil
}
