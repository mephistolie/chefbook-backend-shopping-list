package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/mephistolie/chefbook-backend-common/log"
	"github.com/mephistolie/chefbook-backend-common/responses/fail"
	"github.com/mephistolie/chefbook-backend-shopping-list/internal/config"
)

const (
	shoppingListTable = "shopping_list"
	inboxTable        = "inbox"

	errUniqueViolation = "23505"
)

type Repository struct {
	db *sqlx.DB
}

func Connect(cfg config.Database) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx",
		fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=require",
			*cfg.Host, *cfg.Port, *cfg.User, *cfg.DBName, *cfg.Password))
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func errorWithTransactionRollback(tx *sql.Tx, err error) error {
	if err = tx.Rollback(); err != nil {
		log.Error("unable to rollback transaction: ", err)
	}
	return nil
}

func commitTransaction(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		log.Error("unable to commit transaction: ", err)
		if err := tx.Rollback(); err != nil {
			log.Error("unable to rollback transaction: ", err)
		}
		return fail.GrpcUnknown
	}
	return nil
}

func isUniqueViolationError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == errUniqueViolation
	}
	return false
}
