package services

import (
	"context"
	"learning-platform/internal/models"
	"learning-platform/internal/platform/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(c context.Context, email, password string) (string, error)
}

type authService struct {
	users      UserRepository
	jwtManager *jwt.Manager
}

func NewAuthService(users UserRepository, jwtManager *jwt.Manager) *authService {
	return &authService{
		users:      users,
		jwtManager: jwtManager,
	}
}

func (s *authService) Login(c context.Context, email, password string) (string, error) {
	user, err := s.users.GetByEmail(c, email)
	if err != nil || user == nil {
		return "", models.ErrInvalidCredentials
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", models.ErrInvalidCredentials
	}

	return s.jwtManager.GenerateToken(user.ID, user.Roles)
}
