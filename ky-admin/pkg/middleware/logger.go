package middleware

import (
	"fmt"
	"time"

	"github.com/SupenBysz/ky-admin/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 返回一个日志中间件，记录请求的详细信息
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 获取跟踪ID
		traceID := c.GetString("trace_id")
		if traceID == "" {
			traceID = logger.NewTraceID()
			c.Set("trace_id", traceID)
		}

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 请求状态
		status := c.Writer.Status()
		method := c.Request.Method
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 错误信息
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 构建完整路径
		fullPath := path
		if query != "" {
			fullPath = fmt.Sprintf("%s?%s", path, query)
		}

		// 根据状态码选择日志级别
		log := logger.GetLogger()
		logFields := []zap.Field{
			zap.String("trace_id", traceID),
			zap.String("method", method),
			zap.String("path", fullPath),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.String("ip", ip),
			zap.String("user-agent", userAgent),
		}

		if errorMessage != "" {
			logFields = append(logFields, zap.String("error", errorMessage))
		}

		// 根据状态码选择日志级别
		if status >= 500 {
			log.Error("Server Error", logFields...)
		} else if status >= 400 {
			log.Warn("Client Error", logFields...)
		} else {
			log.Info("Request", logFields...)
		}
	}
}
