package services

import (
	"context"
	"learning-platform/internal/models"
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
}

type EnrollmentService interface {
	CreateEnrollment(c context.Context, enrollment *models.Enrollment) (*models.Enrollment, error)
	ListEnrollment(c context.Context, filter *EnrollmentFilter) ([]*models.Enrollment, error)
	DeleteEnrollment(c context.Context, userID int64, courseID int64) error
}

type enrollmentService struct {
	repo       EnrollmentRepository
	courseRepo CourseRepository
}

func NewEnrollmentService(repo EnrollmentRepository, courseRepo CourseRepository) *enrollmentService {
	return &enrollmentService{
		repo:       repo,
		courseRepo: courseRepo,
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

	isDecrease, err := s.courseRepo.DecrementSeats(c, enrollment.CourseID)
	if err != nil {
		return nil, err
	}
	if !isDecrease {
		return nil, models.ErrCourseFull
	}

	return s.repo.Create(c, enrollment)
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

	return s.repo.Delete(c, enrollment.ID)
}
