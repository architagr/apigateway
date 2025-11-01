package ratelimiter

import (
	"apigateway/pkg/logger"
	"log"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FixedWindowCounter struct {
	capacity  int
	counter   int
	window    time.Duration
	lastReset time.Time
	mu        sync.Mutex
}

func newFixedWindowCounter(capacity, rate int, appLogger *logger.ZapLogger) *FixedWindowCounter {
	bucket := &FixedWindowCounter{
		capacity:  capacity,
		counter:   0,
		window:    time.Duration(rate) * time.Second,
		lastReset: time.Now().UTC(),
	}
	go bucket.startWindow(appLogger)
	appLogger.Info("[fixed-window-rate-limiter] 🪟 Initialized", zap.Int("capacity", capacity), zap.Duration("window", bucket.window))
	return bucket
}

// Periodically resets counter at window boundary
func (f *FixedWindowCounter) startWindow(appLogger *logger.ZapLogger) {
	ticker := time.NewTicker(f.window)
	defer ticker.Stop()

	for t := range ticker.C {
		f.mu.Lock()
		f.lastReset = t.UTC()
		f.counter = 0
		appLogger.Info("[fixed-window-rate-limiter] 🔄 Window reset", zap.Time("lastReset", f.lastReset))
		f.mu.Unlock()
	}
}

// Gin middleware
func fixedWindowCounterMiddleware(rate, capacity int, appLogger *logger.ZapLogger) gin.HandlerFunc {
	limiter := newFixedWindowCounter(capacity, rate, appLogger)

	return func(c *gin.Context) {
		limiter.mu.Lock()
		defer limiter.mu.Unlock()

		remaining := limiter.capacity - limiter.counter
		if remaining < 0 {
			remaining = 0
		}
		rateLimitterHeader(c, limiter.capacity, remaining)

		if limiter.counter < limiter.capacity {
			limiter.counter++
			appLogger.Info(
				"[fixed-window-rate-limiter] ✅ Allowed",
				zap.String("uri", c.Request.URL.Path),
				zap.Int("count", limiter.counter/limiter.capacity),
			)
			c.Next()
		} else {
			log.Printf("[fixed-window-rate-limiter] ⛔ Rate limit exceeded for %s | IP: %s", c.Request.URL.Path, c.ClientIP())
			requestThrottled(c, limiter.lastReset)
		}
	}
}
