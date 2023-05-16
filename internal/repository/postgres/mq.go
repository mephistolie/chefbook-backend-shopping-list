package postgres

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
	shoppingListFail "github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity/fail"
)

func (r *Repository) CreatePersonalShoppingList(userId uuid.UUID, messageId uuid.UUID) error {
	tx, err := r.handleMessageIdempotently(messageId)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil
		} else {
			return fail.GrpcUnknown
		}
	}

	shoppingListId := uuid.New()
	createShoppingListQuery := fmt.Sprintf(`
			INSERT INTO %s (shopping_list_id, owner_id)
			VALUES ($1, $2)
		`, shoppingListsTable)

	if _, err = tx.Exec(createShoppingListQuery, shoppingListId, userId); err != nil {
		log.Errorf("unable to create shopping list for user %s: %s", userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	createConnectionQuery := fmt.Sprintf(`
			INSERT INTO %s (shopping_list_id, user_id)
			VALUES ($1, $2)
		`, usersTable)

	if _, err = tx.Exec(createConnectionQuery, shoppingListId, userId); err != nil {
		log.Errorf("unable to create connection between shopping list %s and user %s: %s", shoppingListId, userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	return commitTransaction(tx)
}

func (r *Repository) ImportFirebaseShoppingList(shoppingListId uuid.UUID, purchases []entity.Purchase, messageId uuid.UUID) error {
	tx, err := r.handleMessageIdempotently(messageId)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil
		} else {
			return fail.GrpcUnknown
		}
	}

	query, bsonShoppingList, bsonRecipeNames, err := getSetShoppingListBaseQuery(purchases, entity.RecipeNames{})
	if err != nil {
		return errorWithTransactionRollback(tx, err)
	}

	if _, err = tx.Exec(query, bsonShoppingList, bsonRecipeNames, shoppingListId); err != nil {
		log.Errorf("unable to set shopping list %s: %s", shoppingListId, err)
		return errorWithTransactionRollback(tx, shoppingListFail.GrpcShoppingListNotFound)
	}

	return commitTransaction(tx)
}

func (r *Repository) DeletePersonalShoppingList(userId uuid.UUID, messageId uuid.UUID) error {
	tx, err := r.handleMessageIdempotently(messageId)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil
		} else {
			return fail.GrpcUnknown
		}
	}

	query := fmt.Sprintf(`
			DELETE FROM %s
			WHERE type='personal' AND user_id=$1
		`, shoppingListsTable)

	if _, err := tx.Exec(query, userId); err != nil {
		log.Errorf("unable to delete user %s: %s", userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	return commitTransaction(tx)
}

func (r *Repository) handleMessageIdempotently(messageId uuid.UUID) (*sql.Tx, error) {
	tx, err := r.startTransaction()
	if err != nil {
		return nil, err
	}

	addMessageQuery := fmt.Sprintf(`
			INSERT INTO %s (message_id)
			VALUES ($1)
		`, inboxTable)

	if _, err = tx.Exec(addMessageQuery, messageId); err != nil {
		if !isUniqueViolationError(err) {
			log.Error("unable to add message to inbox: ", err)
		}
		return nil, errorWithTransactionRollback(tx, err)
	}

	deleteOutdatedMessagesQuery := fmt.Sprintf(`
			DELETE FROM %[1]v
			WHERE ctid IN
			(
				SELECT ctid IN
				FROM %[1]v
				ORDER BY timestamp DESC
				OFFSET 1000
			)
		`, inboxTable)

	if _, err = tx.Exec(deleteOutdatedMessagesQuery); err != nil {
		return nil, errorWithTransactionRollback(tx, err)
	}

	return tx, nil
}
