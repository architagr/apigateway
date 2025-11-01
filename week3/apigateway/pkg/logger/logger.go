package logger

import (
	"apigateway/pkg/constants"
	"context"
	"time"

	"go.uber.org/zap"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewLogger() (*ZapLogger, error) {
	z, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	return &ZapLogger{logger: z}, nil
}


func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func Now() time.Time {
	return time.Now()
}

func Since(t time.Time) time.Duration {
	return time.Since(t)
}

func GetTraceID(ctx context.Context) string {
	if v := ctx.Value(constants.TraceIDKey); v != nil {
		if tid, ok := v.(string); ok {
			return tid
		}
	}
	return ""
}

func (l *ZapLogger) WithTrace(ctx context.Context) *zap.Logger {
	traceID := GetTraceID(ctx)
	return l.logger.With(zap.String("trace_id", traceID))
}
