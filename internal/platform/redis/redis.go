package redis

import (
	"context"
	"learning-platform/internal/configs"
	"time"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

func Connect(cfg *configs.Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisConfig.Host + ":" + cfg.RedisConfig.Port,
		Password: cfg.RedisConfig.Password,
		DB:       cfg.RedisConfig.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

var ProviderSet = wire.NewSet(Connect)
