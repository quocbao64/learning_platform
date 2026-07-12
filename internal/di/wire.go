//go:build wireinject
// +build wireinject

package di

import (
	"learning-platform/internal/configs"
	"learning-platform/internal/handlers"
	"learning-platform/internal/platform/db"
	"learning-platform/internal/platform/jwt"
	"learning-platform/internal/platform/ratelimit"
	"learning-platform/internal/platform/redis"
	"learning-platform/internal/repositories"
	"learning-platform/internal/services"

	"github.com/google/wire"
)

func Initialize() (*Container, error) {
	wire.Build(
		db.ProviderSet,
		repositories.ProviderSet,
		services.ProviderSet,
		handlers.ProviderSet,
		jwt.ProviderSet,
		configs.ProviderSet,
		redis.ProviderSet,
		ratelimit.ProviderSet,
		wire.Struct(new(Container), "*"),
	)
	return nil, nil
}
