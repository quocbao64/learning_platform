package services

import (
	"context"
	"errors"
	"learning-platform/internal/models"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type UserService interface {
	Register(ctx context.Context, name, email, password string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *userService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, name, email, password string) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, models.ErrUserNotFound) {
		return nil, err
	}
	if user != nil {
		return nil, models.ErrEmailAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user = &models.User{
		FullName: name,
		Email:    email,
		Password: string(hash),
		Roles:    "user",
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByID(ctx context.Context, id int64) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}
