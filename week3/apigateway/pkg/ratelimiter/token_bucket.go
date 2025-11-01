package ratelimiter

import (
	"apigateway/pkg/logger"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type tokenBucketRateLimiter struct {
	bucket         []int64
	rate           int
	bucketCapacity int
	nextRefillTime time.Time
	mu             sync.Mutex
}

func (rl *tokenBucketRateLimiter) refillBucket(appLogger *logger.ZapLogger) {
	d := time.Second * time.Duration(rl.rate)
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	for range ticker.C {
		func() {
			rl.mu.Lock()
			defer rl.mu.Unlock()
			refilled := 0
			for i := 0; i < rl.rate && len(rl.bucket) < rl.bucketCapacity; i++ {
				rl.bucket = append(rl.bucket, time.Now().Unix())
				refilled++
			}
			rl.nextRefillTime = time.Now().UTC().Add(d)
			appLogger.Info("[token-bucket-rate-limiter] refilled with new tokens.", zap.Time("nextRefill", rl.nextRefillTime), zap.Int("refilled", refilled))
		}()
	}
}

func rateLimiterMiddleware(r, cap int, appLogger *logger.ZapLogger) gin.HandlerFunc {
	rl := &tokenBucketRateLimiter{
		bucketCapacity: cap,
		rate:           r,
		bucket:         make([]int64, cap),
		nextRefillTime: time.Now().UTC().Add(time.Second * (time.Duration(r))),
	}
	for i := 0; i < rl.bucketCapacity; i++ {
		rl.bucket[i] = time.Now().Unix()
	}

	appLogger.Info(
		"[token-bucket-rate-limiter] initialized",
		zap.Int("capacity", rl.bucketCapacity),
		zap.Int("refillRate (second)", rl.rate),
	)

	go rl.refillBucket(appLogger)

	return func(c *gin.Context) {
		rl.mu.Lock()
		defer rl.mu.Unlock()

		appLogger.Info(
			"[token-bucket-rate-limiter] Incoming request",
			zap.String("method", c.Request.Method),
			zap.String("uri", c.Request.URL.Path),
		)

		rateLimitterHeader(c, rl.bucketCapacity, len(rl.bucket))

		if len(rl.bucket) > 0 {
			rl.bucket = rl.bucket[1:]
			appLogger.Info("[token-bucket-rate-limiter] [allow] Request allowed.", zap.Int("Remaining tokens", len(rl.bucket)))
			c.Next()
		} else {
			appLogger.Error("[token-bucket-rate-limiter] [deny] Too many requests.", zap.Time("nextRefill", rl.nextRefillTime))
			requestThrottled(c, rl.nextRefillTime)
			return
		}
	}
}
