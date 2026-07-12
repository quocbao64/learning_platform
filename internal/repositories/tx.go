package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txManager struct {
	db *pgxpool.Pool
}

func NewTxManager(db *pgxpool.Pool) *txManager {
	return &txManager{
		db: db,
	}
}

func (m *txManager) ExecTx(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(ctx, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
