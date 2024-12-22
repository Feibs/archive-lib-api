package repo

import (
	"context"
	"database/sql"
)

type TransactionRepo interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type transactionRepoImpl struct {
	db *sql.DB
}

func NewTransactionRepo(db *sql.DB) transactionRepoImpl {
	return transactionRepoImpl{
		db: db,
	}
}

type TxKey struct{}

func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, TxKey{}, tx)
}

func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(TxKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

func (t transactionRepoImpl) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = fn(injectTx(ctx, tx))
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
