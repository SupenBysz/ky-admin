// Package errors 提供统一的错误处理机制
package errors

import (
	"errors"
	"fmt"
)

// 业务错误码
const (
	CodeSuccess         = 0    // 成功
	CodeInvalidParams   = 400  // 参数错误
	CodeUnauthorized    = 401  // 未授权
	CodeForbidden       = 403  // 禁止访问
	CodeNotFound        = 404  // 资源不存在
	CodeUserError       = 1000 // 用户相关错误
	CodeRoleError       = 1100 // 角色相关错误
	CodePermissionError = 1200 // 权限相关错误
	CodeSystemError     = 5000 // 系统错误
)

// 业务错误
var (
	// 用户相关错误
	ErrUserNotFound       = errors.New("用户不存在")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserAlreadyExists  = errors.New("用户已存在")
	ErrWeakPassword       = errors.New("密码强度不足")
	ErrOldPasswordWrong   = errors.New("原密码错误")
	ErrUserDisabled       = errors.New("用户已被禁用")

	// 角色相关错误
	ErrRoleNotFound      = errors.New("角色不存在")
	ErrRoleAlreadyExists = errors.New("角色已存在")

	// 权限相关错误
	ErrPermissionDenied = errors.New("权限不足")
	ErrNoPermission     = errors.New("没有操作权限")

	// 系统错误
	ErrInternalServer    = errors.New("服务器内部错误")
	ErrDatabaseOperation = errors.New("数据库操作失败")
	ErrInvalidToken      = errors.New("无效的令牌")
	ErrTokenExpired      = errors.New("令牌已过期")
)

// AppError 应用错误
type AppError struct {
	Code    int    // 错误码
	Message string // 错误消息
	Err     error  // 原始错误
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("错误码: %d, 消息: %s, 原因: %s", e.Code, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("错误码: %d, 消息: %s", e.Code, e.Message)
}

// Unwrap 返回原始错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// NewAppError 创建新的应用错误
func NewAppError(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewUserError 创建用户相关错误
func NewUserError(err error) *AppError {
	return &AppError{
		Code:    CodeUserError,
		Message: "用户操作失败",
		Err:     err,
	}
}

// NewAuthError 创建认证相关错误
func NewAuthError(err error) *AppError {
	return &AppError{
		Code:    CodeUnauthorized,
		Message: "认证失败",
		Err:     err,
	}
}

// NewPermissionError 创建权限相关错误
func NewPermissionError(err error) *AppError {
	return &AppError{
		Code:    CodeForbidden,
		Message: "权限不足",
		Err:     err,
	}
}

// NewSystemError 创建系统相关错误
func NewSystemError(err error) *AppError {
	return &AppError{
		Code:    CodeSystemError,
		Message: "系统错误",
		Err:     err,
	}
}

// Is 判断错误是否为指定错误
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As 将错误转换为指定类型
func As(err error, target interface{}) bool {
	return errors.As(err, target)
}
