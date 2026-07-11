package services

import (
	"context"
	"learning-platform/internal/models"
)

type LessonFilter struct {
	PageID   int
	PerPage  int
	CourseID int64
}

type LessonRepository interface {
	List(c context.Context, filter *LessonFilter) ([]*models.Lesson, error)
	Create(c context.Context, lesson *models.Lesson) error
}

type LessonService interface {
	List(c context.Context, filter *LessonFilter) ([]*models.Lesson, error)
	CreateLesson(c context.Context, lesson *models.Lesson) error
}

type lessonService struct {
	repo LessonRepository
}

func NewLessonService(repo LessonRepository) *lessonService {
	return &lessonService{
		repo: repo,
	}
}

func (s *lessonService) List(c context.Context, filter *LessonFilter) ([]*models.Lesson, error) {
	return s.repo.List(c, filter)
}

func (s *lessonService) CreateLesson(c context.Context, lesson *models.Lesson) error {
	return s.repo.Create(c, lesson)
}
