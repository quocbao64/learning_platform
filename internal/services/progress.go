package services

import (
	"context"
	"learning-platform/internal/models"
)

type ProgressRepository interface {
	Upsert(c context.Context, progress *models.Progress) (*models.Progress, error)
	CountCompletedLessons(c context.Context, userID int64, courseID int64) (int, error)
	CountLessonsByEnrollment(c context.Context, userID int64, courseID int64) (int, error)
}

type ProgressService interface {
	UpdateProgress(c context.Context, progress *models.Progress) (*models.Progress, error)
	GetProgress(c context.Context, userID int64, enrollmentID int64) (*models.CourseProgress, error)
}

type progressService struct {
	repo           ProgressRepository
	enrollmentRepo EnrollmentRepository
}

func NewProgressService(repo ProgressRepository, enrollmentRepo EnrollmentRepository) *progressService {
	return &progressService{
		repo:           repo,
		enrollmentRepo: enrollmentRepo,
	}
}

func (s *progressService) UpdateProgress(c context.Context, progress *models.Progress) (*models.Progress, error) {
	return s.repo.Upsert(c, progress)
}

func (s *progressService) GetProgress(c context.Context, userID int64, enrollmentID int64) (*models.CourseProgress, error) {
	enrollment, err := s.enrollmentRepo.FindByID(c, enrollmentID)
	if err != nil {
		return nil, err
	}

	if enrollment == nil {
		return nil, models.ErrEnrollmentNotFound
	}

	totalLessons, err := s.repo.CountLessonsByEnrollment(c, userID, enrollment.CourseID)
	if err != nil {
		return nil, err
	}

	totalCompletedLessons, err := s.repo.CountCompletedLessons(c, userID, enrollment.CourseID)
	if err != nil {
		return nil, err
	}

	var percentage float64
	if totalLessons > 0 {
		percentage = float64(totalCompletedLessons) / float64(totalLessons) * 100
	}

	return &models.CourseProgress{
		EnrollmentID:     enrollment.ID,
		TotalLessons:     totalLessons,
		CompletedCount:   totalCompletedLessons,
		PercentCompleted: percentage,
	}, err
}
