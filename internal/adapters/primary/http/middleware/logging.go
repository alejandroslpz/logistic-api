package middleware

import (
	"time"

	"logistics-api/internal/pkg/logger"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		statusCode := c.Writer.Status()

		log.Info("HTTP Request",
			logger.String("method", method),
			logger.String("path", path),
			logger.Int("status", statusCode),
			logger.String("latency", latency.String()),
			logger.String("ip", c.ClientIP()),
			logger.String("user_agent", c.Request.UserAgent()),
		)
	}
}
