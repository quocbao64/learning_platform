package services

import (
	"context"
	"encoding/json"
	"fmt"
	"learning-platform/internal/models"

	"github.com/jackc/pgx/v5"
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
	GetByID(c context.Context, id int64) (*models.Course, error)
	DecrementSeats(c context.Context, courseID int64) (bool, error)
	GetSeatsByIDs(c context.Context, courseIDs []int64) (map[int64]int, error)
	DecrementSeatsTx(c context.Context, tx pgx.Tx, courseID int64) (bool, error)
	IncrementSeatsTx(c context.Context, tx pgx.Tx, courseID int64) (bool, error)
	Update(c context.Context, course *models.Course) error
	Delete(c context.Context, id int64) error
}

type CourseService interface {
	ListCourses(c context.Context, filter *CourseFilter) ([]*models.Course, error)
	CreateCourse(c context.Context, course *models.Course) error
	GetCourseByID(c context.Context, id int64) (*models.Course, error)
	UpdateCourse(c context.Context, userID, courseID int64, courseInput *models.UpdateCourse) error
	DeleteCourse(c context.Context, userID, courseID int64) error
}

type courseService struct {
	repo                  CourseRepository
	cache                 Cache
	enrollmentsRepository EnrollmentRepository
}

func NewCourseService(repo CourseRepository, cache Cache, enrollmentRepository EnrollmentRepository) *courseService {
	return &courseService{
		repo:                  repo,
		cache:                 cache,
		enrollmentsRepository: enrollmentRepository,
	}
}

func (s *courseService) ListCourses(c context.Context, filter *CourseFilter) ([]*models.Course, error) {
	cacheKey := fmt.Sprintf("course:list:status=%s:keyword=%s:page=%d:per_page=%d",
		filter.Status, filter.Keyword, filter.PageID, filter.PerPage)

	if cachedData, err := s.cache.Get(c, cacheKey); err == nil {
		var cachedCourses []*models.CachedCourse
		if err := json.Unmarshal([]byte(cachedData), &cachedCourses); err == nil {
			courses := make([]*models.Course, len(cachedCourses))
			for i, course := range cachedCourses {
				courses[i] = &models.Course{
					ID:           course.ID,
					InstructorID: course.InstructorID,
					Title:        course.Title,
					Description:  course.Description,
					Status:       course.Status,
					CreatedAt:    course.CreatedAt,
					UpdatedAt:    course.UpdatedAt,
				}
			}

			if err := s.attachSeatsToCourses(c, courses); err != nil {
				return nil, err
			}

			return courses, nil
		}
	}

	courses, err := s.repo.List(c, filter)
	if err != nil {
		return courses, err
	}

	cachedCourses := make([]models.CachedCourse, len(courses))
	for i, course := range courses {
		cachedCourses[i] = models.CachedCourse{
			ID:           course.ID,
			InstructorID: course.InstructorID,
			Title:        course.Title,
			Description:  course.Description,
			Status:       course.Status,
			CreatedAt:    course.CreatedAt,
			UpdatedAt:    course.UpdatedAt,
		}
	}

	if data, err := json.Marshal(cachedCourses); err == nil {
		_ = s.cache.Set(c, cacheKey, string(data))

	}

	if err := s.attachSeatsToCourses(c, courses); err != nil {
		return nil, err
	}

	return courses, nil
}

func (s *courseService) CreateCourse(c context.Context, course *models.Course) error {
	if err := s.repo.Create(c, course); err != nil {
		return err
	}

	_ = s.cache.DeleteByPattern(c, "course:list:*")

	return nil
}

func (s *courseService) GetCourseByID(c context.Context, id int64) (*models.Course, error) {
	cacheKey := fmt.Sprintf("course:get:id=%d", id)
	if cachedCourses, err := s.cache.Get(c, cacheKey); err == nil {
		var course *models.Course
		if err := json.Unmarshal([]byte(cachedCourses), &course); err == nil {
			return course, err
		}
	}

	course, err := s.repo.GetByID(c, id)
	if err != nil {
		return nil, err
	}

	if data, err := json.Marshal(course); err == nil {
		err = s.cache.Set(c, cacheKey, string(data))
		if err != nil {
			return nil, err
		}
	}
	return course, nil
}

func (s *courseService) attachSeatsToCourses(c context.Context, courses []*models.Course) error {
	if len(courses) == 0 {
		return nil
	}

	courseIDs := make([]int64, len(courses))
	for i, course := range courses {
		courseIDs[i] = course.ID
	}
	seats, err := s.repo.GetSeatsByIDs(c, courseIDs)
	if err != nil {
		return err
	}
	for _, course := range courses {
		course.TotalSeats = seats[course.ID]
	}

	return nil
}

func (s *courseService) UpdateCourse(c context.Context, userID, courseID int64, courseInput *models.UpdateCourse) error {
	course, err := s.repo.GetByID(c, courseID)
	if err != nil {
		return err
	}

	if course.InstructorID != userID {
		return models.ErrForbidden
	}

	if courseInput.Title != nil {
		course.Title = *courseInput.Title
	}
	if courseInput.Description != nil {
		course.Description = *courseInput.Description
	}
	if courseInput.TotalSeats != nil {
		course.TotalSeats = *courseInput.TotalSeats
	}
	if courseInput.Status != nil {
		if !isValidStatus(*courseInput.Status) {
			return models.ErrCourseInvalidStatus
		}
		course.Status = *courseInput.Status
	}

	if err := s.repo.Update(c, course); err != nil {
		return err
	}

	_ = s.cache.Delete(c, fmt.Sprintf("course:get:id=%d", course.ID))
	_ = s.cache.DeleteByPattern(c, "course:list:*")

	return nil
}

func (s *courseService) DeleteCourse(c context.Context, userID, courseID int64) error {
	course, err := s.repo.GetByID(c, courseID)
	if err != nil {
		return err
	}

	if course.InstructorID != userID {
		return models.ErrForbidden
	}

	count, err := s.enrollmentsRepository.CountByCourseID(c, courseID)
	if err != nil {
		return err
	}

	if count > 0 {
		return models.ErrCourseHasEnrollments
	}

	if err := s.repo.Delete(c, courseID); err != nil {
		return err
	}

	_ = s.cache.Delete(c, fmt.Sprintf("course:get:id=%d", courseID))
	_ = s.cache.DeleteByPattern(c, "course:list:*")

	return nil
}

func isValidStatus(status string) bool {
	switch status {
	case models.CourseStatusDraft,
		models.CourseStatusPublished,
		models.CourseStatusArchived:
		return true
	}
	return false
}
