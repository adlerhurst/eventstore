package memory

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type retriableTx struct {
	tx pgx.Tx
}

func (tx *retriableTx) Exec(ctx context.Context, stmt string, args ...interface{}) error {
	_, err := tx.tx.Exec(ctx, stmt, args...)
	return err
}

func (tx *retriableTx) Commit(ctx context.Context) error {
	return tx.tx.Commit(ctx)
}

func (tx *retriableTx) Rollback(ctx context.Context) error {
	return tx.tx.Rollback(ctx)
}
