package di

import (
	_ "learning-platform/docs"
	"learning-platform/internal/handlers"
	"learning-platform/internal/middleware"
	"learning-platform/internal/platform/jwt"
	"learning-platform/internal/platform/ratelimit"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	RateLimit         *ratelimit.RateLimiter
}

func (c *Container) SetupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	api := r.Group("/api/v1")
	authMiddleware := middleware.Auth(c.JWTManager)
	limitByIP := middleware.RateLimit(c.RateLimit, 5, time.Minute, keyByIP)
	c.AuthHandler.RegisterRoutes(api, limitByIP)
	c.UserHandler.RegisterRoutes(api, authMiddleware)
	c.CourseHandler.RegisterRoute(api, authMiddleware)
	c.LessonHandler.RegisterRoutes(api, authMiddleware)
	c.EnrollmentHandler.RegisterRoutes(api, authMiddleware)
	c.ProgressHandler.RegisterRoutes(api, authMiddleware)

	return r
}

func keyByIP(c *gin.Context) string {
	return "ip:" + c.ClientIP()
}
