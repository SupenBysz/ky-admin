package logger

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// GinTraceMiddleware 为每个HTTP请求生成和附加跟踪ID的Gin中间件
func GinTraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 为请求生成跟踪ID
		traceID := c.GetHeader("X-Trace-ID")
		if traceID == "" {
			traceID = NewTraceID()
		}

		// 将跟踪ID添加到上下文
		ctx := WithTraceID(c.Request.Context(), traceID)
		c.Request = c.Request.WithContext(ctx)

		// 在响应头中添加跟踪ID
		c.Header("X-Trace-ID", traceID)

		// 继续请求处理
		c.Next()
	}
}

// GinLoggerMiddleware Gin框架的日志中间件，记录请求和响应信息
func GinLoggerMiddleware(skipPaths ...string) gin.HandlerFunc {
	skipPathMap := make(map[string]bool)
	for _, path := range skipPaths {
		skipPathMap[path] = true
	}

	return func(c *gin.Context) {
		// 检查是否跳过该路径
		path := c.Request.URL.Path
		if skipPathMap[path] {
			c.Next()
			return
		}

		start := time.Now()
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 获取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 重置请求体，以便后续中间件和处理函数可以读取
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 包装响应Writer以捕获响应数据
		responseWriter := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = responseWriter

		// 继续处理请求
		c.Next()

		// 请求处理完成后记录日志
		latency := time.Since(start)
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 获取上下文中的跟踪ID
		traceID := GetTraceID(c.Request.Context())

		// 日志字段
		fields := []zap.Field{
			zap.String(TraceIDField, traceID),
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.String("client_ip", clientIP),
			zap.String("user_agent", userAgent),
			zap.Duration("latency", latency),
		}

		// 根据需要添加请求体和响应体（可能需要根据内容类型和大小进行过滤）
		if len(requestBody) > 0 && shouldLogBody(c.ContentType()) {
			// 限制请求体大小，避免日志过大
			if len(requestBody) > 1024 {
				fields = append(fields, zap.String("request_body", string(requestBody[:1024])+"..."))
			} else {
				fields = append(fields, zap.String("request_body", string(requestBody)))
			}
		}

		responseBody := responseWriter.body.String()
		if responseBody != "" && shouldLogBody(c.Writer.Header().Get("Content-Type")) {
			// 限制响应体大小
			if len(responseBody) > 1024 {
				fields = append(fields, zap.String("response_body", responseBody[:1024]+"..."))
			} else {
				fields = append(fields, zap.String("response_body", responseBody))
			}
		}

		if errorMessage != "" {
			fields = append(fields, zap.String("error", errorMessage))
		}

		// 根据状态码确定日志级别
		if statusCode >= 500 {
			Error("HTTP请求处理失败", fields...)
		} else if statusCode >= 400 {
			Warn("HTTP请求处理警告", fields...)
		} else {
			Info("HTTP请求处理完成", fields...)
		}
	}
}

// 响应体写入器，用于捕获响应内容
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写Write方法以捕获响应内容
func (w *responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString 重写WriteString方法以捕获响应内容
func (w *responseBodyWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// 判断是否应该记录请求/响应体
func shouldLogBody(contentType string) bool {
	// 根据Content-Type决定是否记录请求/响应体
	switch contentType {
	case "application/json",
		"application/xml",
		"application/x-www-form-urlencoded",
		"text/plain",
		"text/html",
		"text/xml":
		return true
	}
	return false
}
