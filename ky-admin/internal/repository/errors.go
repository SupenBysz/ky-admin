package repository

import (
	"errors"
)

// 定义错误常量
var (
	// 用户相关错误
	ErrUserNotFound         = errors.New("用户不存在")
	ErrUserAlreadyExists    = errors.New("用户已存在")
	ErrInvalidCredentials   = errors.New("用户名或密码不正确")
	ErrWeakPassword         = errors.New("密码强度不足")
	ErrIncorrectOldPassword = errors.New("旧密码不正确")

	// 角色相关错误
	ErrRoleNotFound      = errors.New("角色不存在")
	ErrRoleAlreadyExists = errors.New("角色已存在")

	// 权限相关错误
	ErrPermissionDenied = errors.New("权限不足")
)
