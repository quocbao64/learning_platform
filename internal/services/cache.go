package services

import (
	"context"
)

type Cache interface {
	Get(c context.Context, key string) (string, error)
	Set(c context.Context, key, value string) error
	Delete(c context.Context, key string) error
	DeleteByPattern(ctx context.Context, pattern string) error
}
