package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`    // 业务代码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 数据
}

// Success 成功响应
func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带消息的成功响应
func SuccessWithMessage(ctx *gin.Context, data interface{}, message string) {
	ctx.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Fail 失败响应
func Fail(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// BadRequest 请求参数错误
func BadRequest(ctx *gin.Context, message string) {
	Fail(ctx, http.StatusBadRequest, message, nil)
}

// Unauthorized 未授权
func Unauthorized(ctx *gin.Context, message string) {
	if message == "" {
		message = "未授权或登录已过期"
	}
	Fail(ctx, http.StatusUnauthorized, message, nil)
}

// Forbidden 无权访问
func Forbidden(ctx *gin.Context, message string) {
	if message == "" {
		message = "无权访问"
	}
	Fail(ctx, http.StatusForbidden, message, nil)
}

// NotFound 资源不存在
func NotFound(ctx *gin.Context, message string) {
	if message == "" {
		message = "资源不存在"
	}
	Fail(ctx, http.StatusNotFound, message, nil)
}

// ServerError 服务器内部错误
func ServerError(ctx *gin.Context, message string) {
	if message == "" {
		message = "服务器内部错误"
	}
	Fail(ctx, http.StatusInternalServerError, message, nil)
}

// Page 分页响应
type Page struct {
	Total    int64       `json:"total"`    // 总数
	Page     int         `json:"page"`     // 页码
	PageSize int         `json:"pageSize"` // 每页条数
	List     interface{} `json:"list"`     // 数据列表
}

// PageSuccess 分页成功响应
func PageSuccess(ctx *gin.Context, total int64, page int, pageSize int, list interface{}) {
	Success(ctx, Page{
		Total:    total,
		Page:     page,
		PageSize: pageSize,
		List:     list,
	})
}

// SuccessWithPage 分页成功响应（别名）
func SuccessWithPage(ctx *gin.Context, list interface{}, total int64, page int, pageSize int) {
	PageSuccess(ctx, total, page, pageSize, list)
}
