package repositories

import (
	"context"
	"errors"
	"learning-platform/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, password FROM users WHERE email = $1`, email,
	).Scan(&user.ID, &user.Email, &user.Password)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (email, password, full_name, roles, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		user.Email, user.Password, user.FullName, user.Roles, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return models.ErrEmailAlreadyExists
	}

	return err
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}
	err := r.db.QueryRow(ctx,
		`SELECT id, email, roles, full_name FROM users WHERE id = $1`, id,
	).Scan(&user.ID, &user.Email, &user.Roles, &user.FullName)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, models.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
