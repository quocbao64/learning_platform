package ratelimit

import (
	"context"
	"time"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	client *redis.Client
}

func NewRateLimiter(client *redis.Client) *RateLimiter {
	return &RateLimiter{
		client: client,
	}
}

func (r *RateLimiter) Allow(c context.Context, key string, limit int, window int) (bool, int, error) {
	redisKey := "ratelimit:" + key

	count, err := r.client.Incr(c, redisKey).Result()
	if err != nil {
		return false, 0, err
	}

	if count == 1 {
		r.client.Expire(c, redisKey, time.Duration(window)*time.Second)
	}

	if count > int64(limit) {
		return false, 0, nil
	}

	return true, limit - int(count), nil
}

var ProviderSet = wire.NewSet(NewRateLimiter)
