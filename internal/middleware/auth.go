package middleware

import (
	"learning-platform/internal/handlers/response"
	"learning-platform/internal/models"
	"learning-platform/internal/platform/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

const CtxUserID = "user_id"
const CtxRole = "role"

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
		c.Set(CtxRole, claims.Roles)
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString(CtxRole)
		for _, role := range roles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.AbortWithError(c, models.ErrForbidden)
	}

}
