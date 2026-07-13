package repositories

import (
	"context"
	"fmt"
	"learning-platform/internal/models"
	"learning-platform/internal/services"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type lessonRepository struct {
	db *pgxpool.Pool
}

func (r *lessonRepository) List(c context.Context, filter *services.LessonFilter) ([]*models.Lesson, error) {
	query := `SELECT id, course_id, title, content, order_index, created_at, updated_at 
				FROM lessons 
				WHERE 1=1`
	args := make([]interface{}, 0)
	i := 1

	if filter.CourseID != 0 {
		query += fmt.Sprintf(` AND course_id = $%d`, i)
		args = append(args, filter.CourseID)
		i++
	}

	query += fmt.Sprintf(` ORDER BY order_index ASC LIMIT %d OFFSET %d`, filter.PerPage, (filter.PageID)*filter.PerPage)

	rows, err := r.db.Query(c, query, args...)
	if err != nil {
		return nil, models.ErrInternal.Wrap(err)
	}
	defer rows.Close()

	lessons, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Lesson])
	if err != nil {
		return nil, models.ErrInternal.Wrap(err)
	}

	return lessons, nil
}

func (r *lessonRepository) Create(c context.Context, lesson *models.Lesson) error {
	err := r.db.QueryRow(c,
		`INSERT INTO lessons (course_id, title, content, order_index, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		lesson.CourseID,
		lesson.Title,
		lesson.Content,
		lesson.OrderIndex,
		lesson.CreatedAt,
		lesson.UpdatedAt,
	).Scan(&lesson.ID)

	if err != nil {
		return err
	}

	return nil
}

func NewLessonRepository(db *pgxpool.Pool) *lessonRepository {
	return &lessonRepository{
		db: db,
	}
}
