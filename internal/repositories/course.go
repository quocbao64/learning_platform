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

func (r *courseRepository) List(c context.Context, filter *services.CourseFilter) ([]*models.Course, error) {
	var courses []*models.Course
	query := `SELECT id, title, description, instructor_id, status, total_seats, created_at, updated_at 
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

	rows, err := r.db.Query(c, query, args...)

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

func (r *courseRepository) Create(c context.Context, course *models.Course) error {
	err := r.db.QueryRow(c,
		`INSERT INTO courses (title, description, instructor_id, status, total_seats, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		course.Title, course.Description, course.InstructorID, course.Status, course.TotalSeats, time.Now(), time.Now(),
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

func (r *courseRepository) GetByID(c context.Context, id int64) (*models.Course, error) {
	var course models.Course
	err := r.db.QueryRow(c,
		`SELECT id, title, description, instructor_id, status, total_seats, created_at, updated_at 
				FROM courses WHERE id = $1`, id,
	).Scan(
		&course.ID,
		&course.Title,
		&course.Description,
		&course.InstructorID,
		&course.Status,
		&course.TotalSeats,
		&course.CreatedAt,
		&course.UpdatedAt,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return nil, models.ErrCourseNotFound
	}

	if err != nil {
		return nil, err
	}

	lessonsRows, err := r.db.Query(c,
		`SELECT id, title, order_index, created_at, updated_at
				FROM lessons
    			WHERE course_id = $1 ORDER BY order_index`, id)
	if err != nil {
		return nil, err
	}

	defer lessonsRows.Close()

	var lessons []*models.Lesson
	for lessonsRows.Next() {
		lesson := models.Lesson{}
		err = lessonsRows.Scan(&lesson.ID, &lesson.Title, &lesson.OrderIndex, &lesson.CreatedAt, &lesson.UpdatedAt)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, &lesson)
	}

	course.Lessons = lessons

	return &course, nil
}

func (r *courseRepository) DecrementSeats(c context.Context, courseID int64) (bool, error) {
	cmd, err := r.db.Exec(c,
		`UPDATE courses SET total_seats = total_seats - 1 
               WHERE id = $1 AND total_seats > 0 RETURNING total_seats`, courseID)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, err
	}

	return cmd.RowsAffected() == 1, nil
}

func NewCourseRepository(db *pgxpool.Pool) *courseRepository {
	return &courseRepository{
		db: db,
	}
}
