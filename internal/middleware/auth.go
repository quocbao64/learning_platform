package middleware

import (
	"learning-platform/internal/platform/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const CtxUserID = "userID"

func UserID(c *gin.Context) string {
	return c.GetString(CtxUserID)
}

func Auth(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
			})
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
			})
			return
		}

		claims, err := jwtManager.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": http.StatusUnauthorized,
			})
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Next()
	}
}
