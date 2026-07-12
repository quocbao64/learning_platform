package di

import (
	"learning-platform/internal/handlers"
	"learning-platform/internal/middleware"
	"learning-platform/internal/platform/jwt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Container struct {
	AuthHandler       *handlers.AuthHandler
	JWTManager        *jwt.Manager
	UserHandler       *handlers.UserHandler
	CourseHandler     *handlers.CourseHandler
	LessonHandler     *handlers.LessonHandler
	EnrollmentHandler *handlers.EnrollmentHandler
	ProgressHandler   *handlers.ProgressHandler
	RedisClient       *redis.Client
}

func (c *Container) SetupRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/v1")
	authMiddleware := middleware.Auth(c.JWTManager)
	c.AuthHandler.RegisterRoutes(api, authMiddleware)
	c.UserHandler.RegisterRoutes(api, authMiddleware)
	c.CourseHandler.RegisterRoute(api, authMiddleware)
	c.LessonHandler.RegisterRoutes(api, authMiddleware)
	c.EnrollmentHandler.RegisterRoutes(api, authMiddleware)
	c.ProgressHandler.RegisterRoutes(api, authMiddleware)

	return r
}
