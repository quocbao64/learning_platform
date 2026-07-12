package services

import (
	"context"
	"errors"
)

var ErrCacheMiss = errors.New("cache miss")

type Cache interface {
	Get(c context.Context, key string) (string, error)
	Set(c context.Context, key string, value string) error
	Delete(c context.Context, key string) error
}
