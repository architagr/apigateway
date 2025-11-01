package ratelimiter

import (
	"apigateway/config"
	"apigateway/pkg/logger"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetRateLimiterMiddlerware(c *config.Config, appLogger *logger.ZapLogger) gin.HandlerFunc {
	var f gin.HandlerFunc

	switch c.RateLimitter.LimiterType {
	case config.RateLimiterType_TokenBucket:
		f = rateLimiterMiddleware(c.RateLimitter.Rate, c.RateLimitter.Capacity, appLogger)
	case config.RateLimiterType_LeakyBucket:
		f = leakyBucketMiddleware(c.RateLimitter.Rate, c.RateLimitter.Capacity, appLogger)
	case config.RateLimiterType_FixedWindowBucket:
		f = fixedWindowCounterMiddleware(c.RateLimitter.Rate, c.RateLimitter.Capacity, appLogger)
	default:
		f = rateLimiterMiddleware(c.RateLimitter.Rate, c.RateLimitter.Capacity, appLogger)
	}
	return f
}

func rateLimitterHeader(c *gin.Context, bucketCapacity, remainingCapacity int) {
	c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", bucketCapacity))
	c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remainingCapacity))
}
func requestThrottled(c *gin.Context, nextRefillTime time.Time) {
	resetIn := time.Until(nextRefillTime).Seconds()
	c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", nextRefillTime.Unix()))
	c.Header("Retry-After", fmt.Sprintf("%.0f", resetIn))
	c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
		"message": "too many requests",
	})
}
