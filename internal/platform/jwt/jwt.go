package jwt

import (
	"errors"
	"learning-platform/internal/configs"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/wire"
)

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewManager(secret string, ttl int) *Manager {
	return &Manager{
		secret: []byte(secret),
		ttl:    time.Duration(ttl) * time.Minute,
	}
}

func (m *Manager) GenerateToken(userID int64) (string, error) {
	now := time.Now()
	exp := now.Add(m.ttl)
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

func (m *Manager) ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return m.secret, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func Provide(cfg *configs.Config) *Manager {
	return NewManager(cfg.JWTConfig.Secret, cfg.JWTConfig.TTLMinutes)
}

var ProviderSet = wire.NewSet(Provide)
