package db

import (
	"context"
	"errors"
	"learning-platform/internal/configs"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ProviderSet = wire.NewSet(
	ConnectDB,
)

func ConnectDB(cfg *configs.Config) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(cfg.DbConfig.DatabaseURL)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func Migration(cfg *configs.Config) error {
	sourceURL := "file://" + cfg.MigrationConfig.MigrationPath
	migration, err := migrate.New(sourceURL, cfg.DbConfig.DatabaseURL)
	if err != nil {
		return err
	}

	err = migration.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
