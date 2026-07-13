package middleware

import (
	"fmt"
	"learning-platform/internal/handlers/response"
	"learning-platform/internal/models"
	"learning-platform/internal/platform/ratelimit"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimit(limiter *ratelimit.RateLimiter, limit int, window time.Duration, fn func(ctx *gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := fn(c)

		allowed, remaining, err := limiter.Allow(c.Request.Context(), key, limit, int(window.Seconds()))
		if err != nil {
			c.Next()
			return
		}

		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		if !allowed {
			response.AbortWithError(c, models.ErrRateLimitExceeded)
			return
		}

		c.Next()
	}
}
