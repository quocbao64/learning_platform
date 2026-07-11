package jwt

import (
	"learning-platform/internal/configs"

	"github.com/google/wire"
)

func Provide(cfg *configs.Config) *Manager {
	return NewManager(cfg.JWTConfig.Secret, cfg.JWTConfig.TTLMinutes)
}

var ProviderSet = wire.NewSet(Provide)
