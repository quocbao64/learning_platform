package repositories

import (
	"learning-platform/internal/services"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserRepository, wire.Bind(new(services.UserRepository), new(*userRepository)),
)
