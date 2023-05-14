package postgres

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	shoppingListFail "github.com/mephistolie/chefbook-backend-shopping-list/internal/entity/fail"
)

func (r *Repository) GetShoppingListOwner(shoppingListId uuid.UUID) (uuid.UUID, error) {
	ownerId := uuid.UUID{}

	query := fmt.Sprintf(`
			SELECT owner_id
			FROM %s
			WHERE shopping_list_id=$1
		`, usersTable)

	err := r.db.Get(&ownerId, query, shoppingListId)
	if err != nil {
		log.Errorf("unable to get shopping list %s owner: %s", shoppingListId, err)
		return uuid.UUID{}, fail.GrpcUnknown
	}

	return ownerId, nil
}

func (r *Repository) GetShoppingListUsers(shoppingListId uuid.UUID) ([]uuid.UUID, error) {
	var users []uuid.UUID

	query := fmt.Sprintf(`
			SELECT user_id
			FROM %s
			WHERE shopping_list_id=$1 and accepted=true
		`, usersTable)

	rows, err := r.db.Query(query, shoppingListId)
	if err != nil {
		log.Errorf("unable to get shopping list %s users: %s", shoppingListId, err)
		return nil, fail.GrpcUnknown
	}

	for rows.Next() {
		user := uuid.UUID{}
		if err = rows.Scan(&user); err != nil {
			log.Warnf("unable to parse shopping list user id: ", err)
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *Repository) InviteUserToShoppingList(userId, shoppingListId uuid.UUID) error {
	tx, err := r.startTransaction()
	if err != nil {
		return err
	}

	var count int
	getShoppingListUsersCountQuery := fmt.Sprintf(`
			SELECT count(shopping_list_id)
			FROM %s
			WHERE shopping_list_id=$1
		`, usersTable)
	row := tx.QueryRow(getShoppingListUsersCountQuery)
	if err := row.Scan(&count); err != nil {
		log.Errorf("unable to get shopping list %s users count: %s", shoppingListId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}
	if count >= r.maxShoppingListUsersCount {
		log.Warnf("user %s tries to invite guest to shopping list over maximum users count %s", userId, r.maxShoppingListUsersCount)
		return errorWithTransactionRollback(tx, shoppingListFail.GrpcMaxShoppingListUsersCount(count))
	}

	inviteUserQuery := fmt.Sprintf(`
			INSERT INTO %[1]v (shopping_list_id, user_id)
			VALUES ($1, $2)
			WHERE
				NOT EXISTS (
					SELECT shopping_list_id FROM %[1]v WHERE shopping_list_id=$1 AND user_id=$2
				)
		`, usersTable)

	if _, err = tx.Exec(inviteUserQuery, shoppingListId); err != nil {
		log.Errorf("unable to add connection between shopping list %s and user %s: %s", shoppingListId, userId, err)
		return fail.GrpcUnknown
	}

	return nil
}

func (r *Repository) AcceptShoppingListInvite(userId, shoppingListId uuid.UUID) error {
	query := fmt.Sprintf(`
			UPDATE %s
			SET accepted=true
			WHERE shopping_list_id=$1 AND user_id=$1
		`, usersTable)

	if _, err := r.db.Exec(query, shoppingListId, userId); err != nil {
		log.Errorf("unable to accept shopping list %s invite for user %s: %s", shoppingListId, userId, err)
		return fail.GrpcUnknown
	}

	return nil
}

func (r *Repository) DeleteUserFromShoppingList(userId, shoppingListId uuid.UUID) error {
	query := fmt.Sprintf(`
			DELETE FROM %s
			WHERE shopping_list_id=$1 AND user_id=$2
		`, usersTable)

	if _, err := r.db.Exec(query, shoppingListId); err != nil {
		log.Errorf("unable to delete connection between shopping list %s and user %s: %s", shoppingListId, userId, err)
		return fail.GrpcUnknown
	}

	return nil
}
