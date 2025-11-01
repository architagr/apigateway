package ratelimiter

import (
	"apigateway/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type LeakyBucket struct {
	capacity   int
	queue      chan struct{}
	leakRate   time.Duration
	lastLeaked time.Time
}

func newLeakyBucket(capacity, leakRatePerSec int, appLogger *logger.ZapLogger) *LeakyBucket {
	bucket := &LeakyBucket{
		capacity: capacity,
		queue:    make(chan struct{}, capacity),
		leakRate: time.Second / time.Duration(leakRatePerSec),
	}
	go bucket.startLeaking(appLogger)
	appLogger.Info("[leaky-bucket-rate-limiter] created.", zap.Int("capacity", capacity), zap.Int("refillRate (second)", leakRatePerSec))
	return bucket
}

func (lb *LeakyBucket) startLeaking(appLogger *logger.ZapLogger) {
	ticker := time.NewTicker(lb.leakRate)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-lb.queue:
			lb.lastLeaked = time.Now().UTC()
			appLogger.Info("[leaky-bucket-rate-limiter] Processed 1 request.", zap.Int("queue", len(lb.queue)))
		default:
			// Nothing to leak right now
		}
	}
}

func leakyBucketMiddleware(leakRatePerSec, capacity int, appLogger *logger.ZapLogger) gin.HandlerFunc {
	bucket := newLeakyBucket(capacity, leakRatePerSec, appLogger)

	return func(c *gin.Context) {
		rateLimitterHeader(c, bucket.capacity, bucket.capacity-len(bucket.queue))
		select {
		case bucket.queue <- struct{}{}:
			appLogger.Info("[leaky-bucket-rate-limiter] [allow] Request enqueued.", zap.Int("queue", len(bucket.queue)))
			c.Next()
		default:
			nextRefillTime := bucket.lastLeaked.Add(bucket.leakRate)
			appLogger.Info("[leaky-bucket-rate-limiter] [deny] Bucket full.", zap.String("uri", c.Request.URL.Path))
			requestThrottled(c, nextRefillTime)
		}
	}
}
