package repositories

import (
	"context"
	"errors"
	"fmt"
	"learning-platform/internal/models"
	"learning-platform/internal/services"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type enrollmentRepository struct {
	db *pgxpool.Pool
}

func NewEnrollmentRepository(db *pgxpool.Pool) *enrollmentRepository {
	return &enrollmentRepository{
		db: db,
	}
}

func (r *enrollmentRepository) Create(c context.Context, enrollment *models.Enrollment) (*models.Enrollment, error) {
	err := r.db.QueryRow(c,
		`INSERT INTO enrollments (user_id, course_id, status, enrolled_at)
				VALUES ($1, $2, $3, $4) RETURNING id`,
		enrollment.UserID,
		enrollment.CourseID,
		enrollment.Status,
		enrollment.EnrolledAt,
	).Scan(&enrollment.ID)

	if err != nil {
		return nil, models.ErrInternal.Wrap(err)
	}

	return enrollment, nil
}

func (r *enrollmentRepository) List(c context.Context, filter *services.EnrollmentFilter) ([]*models.Enrollment, error) {
	query := `SELECT id, user_id, course_id, status, enrolled_at
				FROM enrollments
				WHERE 1=1`
	args := make([]interface{}, 0)
	i := 1

	if filter.UserID != 0 {
		query += fmt.Sprintf(` AND user_id = $%d`, i)
		args = append(args, filter.UserID)
		i++
	}
	if filter.CourseID != 0 {
		query += fmt.Sprintf(` AND course_id = $%d`, i)
		args = append(args, filter.CourseID)
	}

	rows, err := r.db.Query(c, query, args...)
	if err != nil {
		return nil, models.ErrInternal.Wrap(err)
	}

	defer rows.Close()

	enrollments, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[models.Enrollment])
	if err != nil {
		return nil, models.ErrInternal.Wrap(err)
	}

	return enrollments, nil
}

func (r *enrollmentRepository) Delete(c context.Context, id int64) error {
	err := r.db.QueryRow(c, `DELETE FROM enrollments WHERE id = $1`, id).Scan()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	return nil
}

func (r *enrollmentRepository) FindByUserIDAndCourseID(c context.Context, userID, courseID int64) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.QueryRow(c,
		`SELECT id, user_id, course_id, status, enrolled_at 
				FROM enrollments 
				WHERE user_id = $1 AND course_id = $2`, userID, courseID,
	).Scan(
		&enrollment.ID,
		&enrollment.UserID,
		&enrollment.CourseID,
		&enrollment.Status,
		&enrollment.EnrolledAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, models.ErrInternal.Wrap(err)
	}

	return &enrollment, nil
}

func (r *enrollmentRepository) FindByID(c context.Context, id int64) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.QueryRow(c,
		`SELECT id, user_id, course_id, status, enrolled_at 
				FROM enrollments 
				WHERE id = $1`, id,
	).Scan(
		&enrollment.ID,
		&enrollment.UserID,
		&enrollment.CourseID,
		&enrollment.Status,
		&enrollment.EnrolledAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, models.ErrInternal.Wrap(err)
	}

	return &enrollment, nil
}

func (r *enrollmentRepository) CreateTx(c context.Context, tx pgx.Tx, enrollment *models.Enrollment) (*models.Enrollment, error) {
	err := tx.QueryRow(c,
		`INSERT INTO enrollments (user_id, course_id, status, enrolled_at)
				VALUES ($1, $2, $3, $4) RETURNING id`,
		enrollment.UserID,
		enrollment.CourseID,
		enrollment.Status,
		enrollment.EnrolledAt,
	).Scan(&enrollment.ID)

	if err != nil {
		return nil, models.ErrInternal.Wrap(err)
	}

	return enrollment, nil
}

func (r *enrollmentRepository) DeleteTx(c context.Context, tx pgx.Tx, id int64) error {
	err := tx.QueryRow(c, `DELETE FROM enrollments WHERE id = $1`, id).Scan()
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	return nil
}

func (r *enrollmentRepository) CountByCourseID(c context.Context, courseID int64) (int64, error) {
	var count int64
	err := r.db.QueryRow(c,
		`SELECT count(*) FROM enrollments WHERE course_id = $1`, courseID,
	).Scan(&count)

	if err != nil {
		return 0, models.ErrInternal.Wrap(err)
	}

	return count, nil
}
