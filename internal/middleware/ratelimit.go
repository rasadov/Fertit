package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rasadov/MailManagerApp/internal/services"
)

func RateLimitMiddleware(limiter services.RateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		ctx := c.Request.Context()

		// Check if IP is already blocked
		blocked, err := limiter.IsBlocked(ctx, ip)
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error":   err.Error(),
				"message": "Rate limiter service unavailable",
			})
			c.Abort()
			return
		}

		if blocked {
			resetTime, err := limiter.GetTimeUntilReset(ctx, ip)
			if err != nil {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"error":   err.Error(),
					"message": "Rate limiter service unavailable",
				})
				c.Abort()
				return
			}

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":               "Too many failed login attempts",
				"message":             "Account temporarily locked",
				"retry_after_seconds": int(resetTime.Seconds()),
				"retry_after_minutes": int(resetTime.Minutes()) + 1,
			})
			c.Abort()
			return
		}

		// Continue with the request
		c.Next()
	}
}
