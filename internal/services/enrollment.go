package services

import (
	"context"
	"errors"
	"learning-platform/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type EnrollmentFilter struct {
	UserID   int64
	CourseID int64
}

type EnrollmentRepository interface {
	Create(c context.Context, enrollment *models.Enrollment) (*models.Enrollment, error)
	List(c context.Context, filter *EnrollmentFilter) ([]*models.Enrollment, error)
	Delete(c context.Context, id int64) error
	FindByUserIDAndCourseID(c context.Context, userID int64, courseID int64) (*models.Enrollment, error)
	FindByID(c context.Context, id int64) (*models.Enrollment, error)
	CreateTx(c context.Context, tx pgx.Tx, enrollment *models.Enrollment) (*models.Enrollment, error)
	DeleteTx(c context.Context, tx pgx.Tx, id int64) error
}

type EnrollmentService interface {
	CreateEnrollment(c context.Context, enrollment *models.Enrollment) (*models.Enrollment, error)
	ListEnrollment(c context.Context, filter *EnrollmentFilter) ([]*models.Enrollment, error)
	DeleteEnrollment(c context.Context, userID int64, courseID int64) error
}

type enrollmentService struct {
	repo       EnrollmentRepository
	courseRepo CourseRepository
	txManager  TxManager
	cache      Cache
}

func NewEnrollmentService(
	repo EnrollmentRepository,
	courseRepo CourseRepository,
	txManager TxManager,
	cache Cache,
) *enrollmentService {
	return &enrollmentService{
		repo:       repo,
		courseRepo: courseRepo,
		txManager:  txManager,
		cache:      cache,
	}
}

func (s *enrollmentService) CreateEnrollment(c context.Context, enrollment *models.Enrollment) (*models.Enrollment, error) {
	course, err := s.courseRepo.GetByID(c, enrollment.CourseID)
	if err != nil {
		return nil, err
	}

	if course == nil {
		return nil, models.ErrCourseNotFound
	}

	if course.Status != models.CourseStatusPublished {
		return nil, models.ErrCourseNotPublished
	}

	existing, err := s.repo.FindByUserIDAndCourseID(c, enrollment.UserID, enrollment.CourseID)
	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, models.ErrEnrollmentAlreadyExists
	}

	err = s.txManager.ExecTx(c, func(c context.Context, tx pgx.Tx) error {
		isDecrease, err := s.courseRepo.DecrementSeatsTx(c, tx, enrollment.CourseID)
		if err != nil {
			return err
		}
		if !isDecrease {
			return models.ErrCourseFull
		}

		if _, err := s.repo.CreateTx(c, tx, enrollment); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return models.ErrEnrollmentAlreadyExists
			}

			return err
		}

		return nil
	})

	return enrollment, nil
}

func (s *enrollmentService) ListEnrollment(c context.Context, filter *EnrollmentFilter) ([]*models.Enrollment, error) {
	return s.repo.List(c, filter)
}

func (s *enrollmentService) DeleteEnrollment(c context.Context, userID int64, courseID int64) error {
	enrollment, err := s.repo.FindByUserIDAndCourseID(c, userID, courseID)
	if err != nil {
		return err
	}

	if enrollment == nil {
		return nil
	}

	if enrollment.UserID != userID {
		return models.ErrEnrollmentNotFound
	}

	err = s.txManager.ExecTx(c, func(c context.Context, tx pgx.Tx) error {
		if err := s.repo.DeleteTx(c, tx, enrollment.ID); err != nil {
			return err
		}

		if _, err := s.courseRepo.IncrementSeatsTx(c, tx, enrollment.CourseID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return s.repo.Delete(c, enrollment.ID)
}
