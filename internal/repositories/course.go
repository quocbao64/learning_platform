package repositories

import (
	"context"
	"errors"
	"fmt"
	"learning-platform/internal/models"
	"learning-platform/internal/services"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type courseRepository struct {
	db *pgxpool.Pool
}

func (c2 courseRepository) List(c context.Context, filter *services.CourseFilter) ([]*models.Course, error) {
	var courses []*models.Course
	query := `SELECT id, title, description, instructor_id, status, created_at, updated_at 
				FROM courses
				WHERE 1=1`
	i := 1
	args := make([]interface{}, 0)

	if filter.Status != "" {
		query += fmt.Sprintf(` AND status = $%d`, i)
		i++
		args = append(args, filter.Status)
	}

	if filter.Keyword != "" {
		query += fmt.Sprintf(` AND (title ILIKE $%d OR description ILIKE $%d)`, i, i)
		args = append(args, "%"+filter.Keyword+"%")
		i++
	}

	query += fmt.Sprintf(` ORDER BY created_at DESC LIMIT %d OFFSET %d`, filter.PerPage, (filter.PageID)*filter.PerPage)

	rows, err := c2.db.Query(c, query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	courses, err = pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Course])
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (c2 courseRepository) Create(c context.Context, course *models.Course) error {
	err := c2.db.QueryRow(c,
		`INSERT INTO courses (title, description, instructor_id, status, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		course.Title, course.Description, course.InstructorID, course.Status, time.Now(), time.Now(),
	).Scan(&course.ID)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return models.ErrCourseAlreadyExists
	}

	if err != nil {
		return err
	}

	return nil
}

func NewCourseRepository(db *pgxpool.Pool) *courseRepository {
	return &courseRepository{
		db: db,
	}
}
