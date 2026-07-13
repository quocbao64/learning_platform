package repositories

import (
	"context"
	"learning-platform/internal/models"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type progressRepository struct {
	db *pgxpool.Pool
}

func NewProgressRepository(db *pgxpool.Pool) *progressRepository {
	return &progressRepository{
		db: db,
	}
}

func (r *progressRepository) Upsert(c context.Context, progress *models.Progress) (*models.Progress, error) {
	var completedAt time.Time
	if progress.Status == models.ProgressStatusCompleted {
		completedAt = time.Now()
	}

	err := r.db.QueryRow(c,
		`INSERT INTO progress (enrollment_id, lesson_id, status, completed_at)
				VALUES ($1, $2, $3, $4)
				ON CONFLICT (enrollment_id, lesson_id) 
				DO UPDATE SET status = $3, completed_at = $4
				RETURNING id`,
		progress.EnrollmentID,
		progress.LessonID,
		progress.Status,
		completedAt,
	).Scan(&progress.ID)

	if err != nil {
		return nil, models.ErrInternal.Wrap(err)
	}

	return progress, nil
}

func (r *progressRepository) CountCompletedLessons(c context.Context, userID int64, courseID int64) (int, error) {
	var count int
	err := r.db.QueryRow(c,
		`SELECT COUNT(*) 
				FROM progress p
				JOIN enrollments e ON p.enrollment_id = e.id
				WHERE e.user_id = $1 AND e.course_id = $2 AND p.status = 'completed'`,
		userID,
		courseID,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *progressRepository) CountLessonsByEnrollment(c context.Context, userID int64, courseID int64) (int, error) {
	var count int
	err := r.db.QueryRow(c,
		`SELECT COUNT(*) 
				FROM lessons l
				JOIN enrollments e ON e.course_id = l.course_id
				WHERE e.user_id = $1 AND e.course_id = $2`,
		userID,
		courseID,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
