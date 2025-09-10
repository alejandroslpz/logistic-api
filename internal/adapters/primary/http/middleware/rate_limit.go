package middleware

import (
	"net/http"
	"sync"
	"time"

	"logistics-api/internal/adapters/primary/http/dto"
	"logistics-api/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
	logger   logger.Logger
}

func NewRateLimiter(rps rate.Limit, burst int, logger logger.Logger) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rps,
		burst:    burst,
		logger:   logger,
	}
}

func (rl *RateLimiter) getLimiter(ip string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[ip]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[ip] = limiter
		rl.mu.Unlock()
	}

	return limiter
}

func (rl *RateLimiter) RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := rl.getLimiter(ip)

		if !limiter.Allow() {
			rl.logger.Warn("Rate limit exceeded",
				logger.String("ip", ip),
				logger.String("path", c.Request.URL.Path))

			dto.ErrorResponse(c, http.StatusTooManyRequests,
				"rate_limit_exceeded", "Too many requests")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (rl *RateLimiter) StartCleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			rl.mu.Lock()
			for ip, limiter := range rl.limiters {
				if limiter.Tokens() == float64(rl.burst) {
					delete(rl.limiters, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
}
