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

func (r *Repository) AddUser(userId uuid.UUID, messageId uuid.UUID) error {
	tx, err := r.handleMessageIdempotently(messageId)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil
		} else {
			return fail.GrpcUnknown
		}
	}

	shoppingListBSON, err := json.Marshal([]dto.Purchase{})
	if err != nil {
		log.Errorf("unable to get marshal shopping list for user %s: %s", userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	addUserQuery := fmt.Sprintf(`
			INSERT INTO %s (user_id, purchases)
			VALUES ($1, $2)
		`, shoppingListTable)

	if _, err := tx.Exec(addUserQuery, userId, shoppingListBSON); err != nil {
		log.Errorf("unable to add user %s: %s", userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	return commitTransaction(tx)
}

func (r *Repository) ImportFirebaseProfile(userId uuid.UUID, purchases []entity.Purchase, messageId uuid.UUID) error {
	tx, err := r.handleMessageIdempotently(messageId)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil
		} else {
			return fail.GrpcUnknown
		}
	}

	setShoppingListQuery, shoppingListBSON, err := getSetShoppingListBaseQuery(purchases)
	if err != nil {
		return errorWithTransactionRollback(tx, err)
	}

	if _, err = tx.Exec(setShoppingListQuery, shoppingListBSON, userId); err != nil {
		log.Errorf("unable to set shopping list for user %s: %s", userId, err)
		return errorWithTransactionRollback(tx, shoppingListFail.GrpcShoppingListNotFound)
	}

	return commitTransaction(tx)
}

func (r *Repository) DeleteUser(userId uuid.UUID, messageId uuid.UUID) error {
	tx, err := r.handleMessageIdempotently(messageId)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil
		} else {
			return fail.GrpcUnknown
		}
	}

	deleteUserQuery := fmt.Sprintf(`
			DELETE FROM %s
			WHERE user_id=$1
		`, shoppingListTable)

	if _, err := tx.Exec(deleteUserQuery, userId); err != nil {
		log.Errorf("unable to delete user %s: %s", userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	return commitTransaction(tx)
}

func (r *Repository) handleMessageIdempotently(messageId uuid.UUID) (*sql.Tx, error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Error("unable to begin transaction: ", err)
		return nil, err
	}

	query := fmt.Sprintf(`
			INSERT INTO %s (message_id)
			VALUES ($1)
		`, inboxTable)

	if _, err = tx.Exec(query, messageId); err != nil {
		if !isUniqueViolationError(err) {
			log.Error("unable to add message to inbox: ", err)
		}
		return nil, errorWithTransactionRollback(tx, err)
	}

	return tx, nil
}
