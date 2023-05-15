package postgres

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"time"
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

func (r *Repository) GetShoppingListKey(shoppingListId uuid.UUID) (uuid.UUID, time.Time, error) {
	tx, err := r.startTransaction()
	if err != nil {
		return uuid.UUID{}, time.Time{}, err
	}

	var key uuid.UUID
	var expiresAt time.Time

	createKeyQuery := fmt.Sprintf(`
			INSERT INTO %[1]v (shopping_list_id, expires_at)
			VALUES ($1, $2)
			WHERE NOT EXISTS
				(
					SELECT key, expires_at
					FROM %[1]v
					WHERE shopping_list_id=$1
				)
			RETURNING key, expires_at
		`, keysTable)

	row := tx.QueryRow(createKeyQuery, shoppingListId, time.Now().Add(r.keyTtl))
	if err := row.Scan(&key, &expiresAt); err != nil {
		log.Errorf("unable to create shopping list %s key", shoppingListId, err)
		return uuid.UUID{}, time.Time{}, errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}
	if expiresAt.Unix() < time.Now().Unix() {
		return r.updateShoppingListKey(tx, shoppingListId)
	}

	return key, expiresAt, commitTransaction(tx)
}

func (r *Repository) updateShoppingListKey(tx *sql.Tx, shoppingListId uuid.UUID) (uuid.UUID, time.Time, error) {
	key := uuid.New()
	expiresAt := time.Now().Add(r.keyTtl)

	updateKeyQuery := fmt.Sprintf(`
			UPDATE %s
			SET key=$1, expires_at=$2
			WHERE shopping_list_id=$3
		`, keysTable)

	if _, err := tx.Exec(updateKeyQuery, key, expiresAt, shoppingListId); err != nil {
		log.Errorf("unable to update shopping list %s key", shoppingListId, err)
		return uuid.UUID{}, time.Time{}, errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	return key, expiresAt, commitTransaction(tx)
}

func (r *Repository) IsShoppingListKeyValid(shoppingListId, key uuid.UUID) (bool, error) {
	valid := false

	query := fmt.Sprintf(`
			SELECT EXISTS
			(
				SELECT 1
				FROM %s
				WHERE shopping_list_id=$1 AND key=$2
			)
		`, keysTable)

	if err := r.db.Get(&valid, query, shoppingListId, key); err != nil {
		log.Errorf("unable to validate shopping list %s key: %s", shoppingListId, err)
		return false, fail.GrpcUnknown
	}
	return true, nil
}

func (r *Repository) AddUserToShoppingList(userId, shoppingListId uuid.UUID) error {
	query := fmt.Sprintf(`
			INSERT INTO %[1]v (shopping_list_id, user_id)
			VALUES ($1, $2)
			WHERE NOT EXISTS
				(
					SELECT shopping_list_id
					WHERE shopping_list_id=$1 AND user_id=$2
				)
		`, usersTable)

	if _, err := r.db.Exec(query, shoppingListId, userId); err != nil {
		log.Errorf("unable to add user %s to shopping list %s: %s", userId, shoppingListId, err)
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
