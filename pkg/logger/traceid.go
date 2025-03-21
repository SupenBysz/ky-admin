package logger

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// TraceID上下文键
type traceIDKey struct{}

// TraceIDField 用于日志的跟踪ID字段名
const TraceIDField = "trace_id"

// NewTraceID 生成新的跟踪ID
func NewTraceID() string {
	return uuid.New().String()
}

// WithTraceID 在上下文中添加跟踪ID
func WithTraceID(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		traceID = NewTraceID()
	}
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// GetTraceID 从上下文中获取跟踪ID
func GetTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	traceID, ok := ctx.Value(traceIDKey{}).(string)
	if !ok {
		return ""
	}
	return traceID
}

// TraceIDField 创建一个带有跟踪ID的zap.Field
func TraceField(ctx context.Context) zap.Field {
	return zap.String(TraceIDField, GetTraceID(ctx))
}

// WithContext 使用上下文创建带有跟踪ID的Logger
func WithContext(ctx context.Context) *zap.Logger {
	traceID := GetTraceID(ctx)
	if traceID == "" {
		return GetLogger()
	}
	return GetLogger().With(zap.String(TraceIDField, traceID))
}

// DebugWithContext 使用上下文输出Debug级别日志，自动包含跟踪ID
func DebugWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	logger := WithContext(ctx)
	logger.Debug(msg, fields...)
}

// InfoWithContext 使用上下文输出Info级别日志，自动包含跟踪ID
func InfoWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	logger := WithContext(ctx)
	logger.Info(msg, fields...)
}

// WarnWithContext 使用上下文输出Warn级别日志，自动包含跟踪ID
func WarnWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	logger := WithContext(ctx)
	logger.Warn(msg, fields...)
}

// ErrorWithContext 使用上下文输出Error级别日志，自动包含跟踪ID
func ErrorWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	logger := WithContext(ctx)
	logger.Error(msg, fields...)
}

// FatalWithContext 使用上下文输出Fatal级别日志，自动包含跟踪ID
func FatalWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	logger := WithContext(ctx)
	logger.Fatal(msg, fields...)
}
