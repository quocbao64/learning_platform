package repositories

import (
	"learning-platform/internal/services"

	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserRepository, wire.Bind(new(services.UserRepository), new(*userRepository)),
	NewCourseRepository, wire.Bind(new(services.CourseRepository), new(*courseRepository)),
	NewLessonRepository, wire.Bind(new(services.LessonRepository), new(*lessonRepository)),
	NewEnrollmentRepository, wire.Bind(new(services.EnrollmentRepository), new(*enrollmentRepository)),
	NewProgressRepository, wire.Bind(new(services.ProgressRepository), new(*progressRepository)),
	NewRedisCache, wire.Bind(new(services.Cache), new(*redisCache)),
)
