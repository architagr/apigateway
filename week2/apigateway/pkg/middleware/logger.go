package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"apigateway/pkg/constants"
	"apigateway/pkg/logger"
)

//requestLoggerMiddleware logs requests and injects a trace ID
func RequestLoggerMiddleware(l *logger.ZapLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		//generate a trace ID
		traceID := uuid.New().String()

		//add trace ID to request context
		ctx := context.WithValue(c.Request.Context(), constants.TraceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)

		l.Info("request received",
			zap.String("trace_id", traceID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("client_ip", c.ClientIP()),
		)

		start := logger.Now()
		c.Next()
		duration := logger.Since(start)

		l.Info("request completed",
			zap.String("trace_id", traceID),
			zap.Int("status", c.Writer.Status()),
			zap.Int64("duration_ms", duration.Milliseconds()),
		)
	}
}
