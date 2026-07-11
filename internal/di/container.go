package di

import (
	"learning-platform/internal/handlers"
	"learning-platform/internal/middleware"
	"learning-platform/internal/platform/jwt"

	"github.com/gin-gonic/gin"
)

type Container struct {
	AuthHandler *handlers.AuthHandler
	JWTManager  *jwt.Manager
}

func (c *Container) SetupRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api/v1")
	authMiddleware := middleware.Auth(c.JWTManager)
	c.AuthHandler.RegisterRoutes(api, authMiddleware)

	return r
}
