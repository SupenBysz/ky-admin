package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 响应码
const (
	CodeSuccess       = 0    // 成功
	CodeParamError    = 1001 // 参数错误
	CodeAuthError     = 1002 // 认证错误
	CodeForbidden     = 1003 // 禁止访问
	CodeNotFound      = 1004 // 资源不存在
	CodeServerError   = 2001 // 服务器错误
	CodeDatabaseError = 2002 // 数据库错误
	CodeUnknownError  = 9999 // 未知错误
)

// Response 标准API响应结构
type Response struct {
	Code    int         `json:"code"`               // 业务码
	Message string      `json:"message"`            // 响应消息
	Data    interface{} `json:"data"`               // 响应数据
	TraceID string      `json:"trace_id,omitempty"` // 请求跟踪ID
}

// Success 返回成功响应
func Success(c *gin.Context, data interface{}) {
	traceID := c.GetString("trace_id")
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "操作成功",
		Data:    data,
		TraceID: traceID,
	})
}

// SuccessWithMsg 返回带自定义消息的成功响应
func SuccessWithMsg(c *gin.Context, message string, data interface{}) {
	traceID := c.GetString("trace_id")
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
		TraceID: traceID,
	})
}

// Fail 返回失败响应
func Fail(c *gin.Context, code int, message string) {
	traceID := c.GetString("trace_id")
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
		TraceID: traceID,
	})
}

// ParamError 返回参数错误响应
func ParamError(c *gin.Context, message string) {
	if message == "" {
		message = "参数错误"
	}
	Fail(c, CodeParamError, message)
}

// AuthError 返回认证错误响应
func AuthError(c *gin.Context, message string) {
	if message == "" {
		message = "认证失败"
	}
	c.JSON(http.StatusUnauthorized, Response{
		Code:    CodeAuthError,
		Message: message,
		Data:    nil,
		TraceID: c.GetString("trace_id"),
	})
}

// ForbiddenError 返回禁止访问响应
func ForbiddenError(c *gin.Context, message string) {
	if message == "" {
		message = "禁止访问"
	}
	c.JSON(http.StatusForbidden, Response{
		Code:    CodeForbidden,
		Message: message,
		Data:    nil,
		TraceID: c.GetString("trace_id"),
	})
}

// NotFoundError 返回资源不存在响应
func NotFoundError(c *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	c.JSON(http.StatusNotFound, Response{
		Code:    CodeNotFound,
		Message: message,
		Data:    nil,
		TraceID: c.GetString("trace_id"),
	})
}

// ServerError 返回服务器错误响应
func ServerError(c *gin.Context, message string) {
	if message == "" {
		message = "服务器内部错误"
	}
	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeServerError,
		Message: message,
		Data:    nil,
		TraceID: c.GetString("trace_id"),
	})
}

// DatabaseError 返回数据库错误响应
func DatabaseError(c *gin.Context, message string) {
	if message == "" {
		message = "数据库操作错误"
	}
	ServerError(c, message)
}
