package middleware

import (
	"learning-platform/internal/handlers/response"
	"learning-platform/internal/models"
	"learning-platform/internal/platform/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

const CtxUserID = "user_id"

func UserID(c *gin.Context) int64 {
	return c.GetInt64(CtxUserID)
}

func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if header == "" {
			response.AbortWithError(c, models.ErrInvalidCredentials)
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.AbortWithError(c, models.ErrInvalidCredentials)
			return
		}

		claims, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			response.AbortWithError(c, models.ErrInvalidCredentials)
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Next()
	}
}
