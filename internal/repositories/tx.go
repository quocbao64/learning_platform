package repositories

import (
	"context"
	"learning-platform/internal/models"

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
		return models.ErrInternal.Wrap(err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	if err := fn(ctx, tx); err != nil {
		return models.ErrInternal.Wrap(err)
	}

	return tx.Commit(ctx)
}
