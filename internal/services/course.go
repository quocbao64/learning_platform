package services

import (
	"context"
	"learning-platform/internal/models"
)

type CourseFilter struct {
	PageID  int
	PerPage int
	Status  string
	Keyword string
}

type CourseRepository interface {
	List(c context.Context, filter *CourseFilter) ([]*models.Course, error)
	Create(c context.Context, course *models.Course) error
}

type CourseService interface {
	ListCourses(c context.Context, filter *CourseFilter) ([]*models.Course, error)
	CreateCourse(c context.Context, course *models.Course) error
}

type courseService struct {
	repo CourseRepository
}

func NewCourseService(repo CourseRepository) *courseService {
	return &courseService{repo: repo}
}

func (s *courseService) ListCourses(c context.Context, filter *CourseFilter) ([]*models.Course, error) {
	return s.repo.List(c, filter)
}

func (s *courseService) CreateCourse(c context.Context, course *models.Course) error {
	return s.repo.Create(c, course)
}
