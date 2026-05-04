package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity"
	shoppingListFail "github.com/mephistolie/chefbook-backend-shopping-list/v2/internal/entity/fail"
)

func (r *Repository) CreatePersonalShoppingList(ctx context.Context, userId uuid.UUID, messageId uuid.UUID) error {
	tx, err := r.handleMessageIdempotently(ctx, messageId)
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

	if _, err = tx.ExecContext(ctx, createShoppingListQuery, shoppingListId, userId); err != nil {
		log.Errorf("unable to create shopping list for user %s: %s", userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	createConnectionQuery := fmt.Sprintf(`
			INSERT INTO %s (shopping_list_id, user_id)
			VALUES ($1, $2)
		`, usersTable)

	if _, err = tx.ExecContext(ctx, createConnectionQuery, shoppingListId, userId); err != nil {
		log.Errorf("unable to create connection between shopping list %s and user %s: %s", shoppingListId, userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	return commitTransaction(tx)
}

func (r *Repository) ImportFirebaseShoppingList(ctx context.Context, shoppingListId uuid.UUID, purchases []entity.Purchase, messageId uuid.UUID) error {
	tx, err := r.handleMessageIdempotently(ctx, messageId)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil
		} else {
			return fail.GrpcUnknown
		}
	}

	query, bsonShoppingList, err := getSetShoppingListBaseQuery(purchases)
	if err != nil {
		return errorWithTransactionRollback(tx, err)
	}

	if _, err = tx.ExecContext(ctx, query, bsonShoppingList, shoppingListId); err != nil {
		log.Errorf("unable to set shopping list %s: %s", shoppingListId, err)
		return errorWithTransactionRollback(tx, shoppingListFail.GrpcShoppingListNotFound)
	}

	return commitTransaction(tx)
}

func (r *Repository) DeleteUserShoppingLists(ctx context.Context, userId uuid.UUID, messageId uuid.UUID) error {
	tx, err := r.handleMessageIdempotently(ctx, messageId)
	if err != nil {
		if isUniqueViolationError(err) {
			return nil
		} else {
			return fail.GrpcUnknown
		}
	}

	query := fmt.Sprintf(`
			DELETE FROM %s
			WHERE owner_id=$1
		`, shoppingListsTable)

	if _, err := tx.ExecContext(ctx, query, userId); err != nil {
		log.Errorf("unable to delete user %s: %s", userId, err)
		return errorWithTransactionRollback(tx, fail.GrpcUnknown)
	}

	return commitTransaction(tx)
}

func (r *Repository) handleMessageIdempotently(ctx context.Context, messageId uuid.UUID) (*sql.Tx, error) {
	tx, err := r.startTransaction(ctx)
	if err != nil {
		return nil, err
	}

	addMessageQuery := fmt.Sprintf(`
			INSERT INTO %s (message_id)
			VALUES ($1)
		`, inboxTable)

	if _, err = tx.ExecContext(ctx, addMessageQuery, messageId); err != nil {
		if !isUniqueViolationError(err) {
			log.Error("unable to add message to inbox: ", err)
		}
		return nil, errorWithTransactionRollback(tx, err)
	}

	deleteOutdatedMessagesQuery := fmt.Sprintf(`
			DELETE FROM %[1]v
			WHERE ctid IN
			(
				SELECT ctid
				FROM %[1]v
				ORDER BY timestamp DESC
				OFFSET 1000
			)
		`, inboxTable)

	if _, err = tx.ExecContext(ctx, deleteOutdatedMessagesQuery); err != nil {
		return nil, errorWithTransactionRollback(tx, err)
	}

	return tx, nil
}
