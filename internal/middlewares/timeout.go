package middlewares

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeDuration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeDuration)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
		defer cancel()
	}
}
