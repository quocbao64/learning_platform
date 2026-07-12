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
}

type CourseService interface {
	ListCourses(c context.Context, filter *CourseFilter) ([]*models.Course, error)
	CreateCourse(c context.Context, course *models.Course) error
	GetCourseByID(c context.Context, id int64) (*models.Course, error)
}

type courseService struct {
	repo  CourseRepository
	cache Cache
}

func NewCourseService(repo CourseRepository, cache Cache) *courseService {
	return &courseService{
		repo:  repo,
		cache: cache,
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
	return s.repo.Create(c, course)
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
