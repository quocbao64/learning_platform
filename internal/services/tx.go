package services

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type TxManager interface {
	ExecTx(ctx context.Context, fn func(ctx context.Context, tx pgx.Tx) error) error
}
