package postgres

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/repository/postgres/dto"
)

func (r *Repository) AddUser(userId uuid.UUID) error {
	var shoppingListBSON, err = json.Marshal(dto.NewShoppingList(emptyShoppingList()))
	if err != nil {
		log.Errorf("unable to get marshal shopping list for user %s: %s", userId, err)
		return fail.GrpcUnknown
	}

	addUserQuery := fmt.Sprintf(`
			INSERT INTO %s (user_id, shopping_list)
			VALUES ($1, $2)
		`, shoppingListTable)

	if _, err := r.db.Exec(addUserQuery, userId, shoppingListBSON); err != nil {
		log.Errorf("unable to add user %s: %s", userId, err)
		return fail.GrpcUnknown
	}

	return nil
}

func (r *Repository) DeleteUser(userId uuid.UUID) error {
	deleteUserQuery := fmt.Sprintf(`
			DELETE FROM %s
			WHERE user_id=$1
		`, shoppingListTable)

	if _, err := r.db.Exec(deleteUserQuery, userId); err != nil {
		log.Errorf("unable to delete user %s: %s", userId, err)
		return fail.GrpcUnknown
	}

	return nil
}
