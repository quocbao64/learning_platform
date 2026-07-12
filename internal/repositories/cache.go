package repositories

import (
	"context"
	"errors"
	"learning-platform/internal/configs"
	"learning-platform/internal/services"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisCache struct {
	client *redis.Client
	cfg    *configs.Config
}

func NewRedisCache(client *redis.Client, redisCfg *configs.Config) *redisCache {
	return &redisCache{
		client: client,
		cfg:    redisCfg,
	}
}

func (c *redisCache) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", services.ErrCacheMiss
	}

	return value, err
}

func (c *redisCache) Set(ctx context.Context, key, value string) error {
	expiration := time.Duration(c.cfg.RedisConfig.TTL) * time.Minute
	return c.client.Set(ctx, key, value, expiration).Err()
}

func (c *redisCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}
