package api

import (
	"net/http"

	"github.com/SupenBysz/ky-admin/internal/common/errors"
	"github.com/gin-gonic/gin"
)

// 响应状态码
const (
	CodeSuccess       = 0   // 成功
	CodeBadRequest    = 400 // 请求错误
	CodeUnauthorized  = 401 // 未授权
	CodeForbidden     = 403 // 禁止访问
	CodeNotFound      = 404 // 资源不存在
	CodeInternalError = 500 // 服务器内部错误
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 响应码
	Message string      `json:"message"` // 响应消息
	Data    interface{} `json:"data"`    // 响应数据
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "操作成功",
		Data:    data,
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// SuccessWithPage 带分页的成功响应
func SuccessWithPage(c *gin.Context, data interface{}, total int64, page, pageSize int) {
	c.JSON(http.StatusOK, gin.H{
		"code":    CodeSuccess,
		"message": "查询成功",
		"data":    data,
		"page":    page,
		"size":    pageSize,
		"total":   total,
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	var code int
	var message string
	var httpStatus int

	// 使用断言检查是否为AppError类型
	if appErr, ok := err.(*errors.AppError); ok {
		code = appErr.Code
		message = appErr.Message
		if appErr.Err != nil {
			message = message + ": " + appErr.Err.Error()
		}

		// 根据错误码设置HTTP状态码
		switch appErr.Code {
		case errors.CodeInvalidParams:
			httpStatus = http.StatusBadRequest
		case errors.CodeUnauthorized:
			httpStatus = http.StatusUnauthorized
		case errors.CodeForbidden:
			httpStatus = http.StatusForbidden
		case errors.CodeNotFound:
			httpStatus = http.StatusNotFound
		default:
			httpStatus = http.StatusInternalServerError
		}
	} else {
		// 默认为内部错误
		code = CodeInternalError
		message = err.Error()
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ParamError 参数错误响应
func ParamError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    CodeBadRequest,
		Message: message,
		Data:    nil,
	})
}

// Unauthorized 未授权响应
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    CodeUnauthorized,
		Message: message,
		Data:    nil,
	})
}

// Forbidden 禁止访问响应
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    CodeForbidden,
		Message: message,
		Data:    nil,
	})
}

// NotFound 资源不存在响应
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    CodeNotFound,
		Message: message,
		Data:    nil,
	})
}

// Failed 通用失败响应
func Failed(c *gin.Context, code int, message string, data ...interface{}) {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	}

	var httpStatus int
	switch code {
	case CodeBadRequest:
		httpStatus = http.StatusBadRequest
	case CodeUnauthorized:
		httpStatus = http.StatusUnauthorized
	case CodeForbidden:
		httpStatus = http.StatusForbidden
	case CodeNotFound:
		httpStatus = http.StatusNotFound
	default:
		httpStatus = http.StatusInternalServerError
	}

	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    responseData,
	})
}
