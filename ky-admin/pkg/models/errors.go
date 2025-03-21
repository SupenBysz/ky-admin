package models

import "errors"

// 定义各种模型操作的错误
var (
	// 通用错误
	ErrRecordNotFound       = errors.New("记录不存在")
	ErrCreateFailed         = errors.New("创建记录失败")
	ErrUpdateFailed         = errors.New("更新记录失败")
	ErrDeleteFailed         = errors.New("删除记录失败")
	ErrDuplicateKey         = errors.New("记录已存在")
	ErrInvalidParam         = errors.New("无效的参数")
	ErrForeignKeyConstraint = errors.New("外键约束错误")

	// 用户相关错误
	ErrUserNotFound          = errors.New("用户不存在")
	ErrUserAlreadyExists     = errors.New("用户已存在")
	ErrUserDisabled          = errors.New("用户已禁用")
	ErrInvalidPassword       = errors.New("密码错误")
	ErrOldPasswordIncorrect  = errors.New("旧密码不正确")
	ErrNoPermission          = errors.New("没有权限")
	ErrUserLocked            = errors.New("用户已锁定")
	ErrPasswordLocked        = errors.New("密码已锁定，无法修改")
	ErrOriginalPasswordWrong = errors.New("原密码错误")
	ErrUsernameExists        = errors.New("用户名已存在")
	ErrEmailExists           = errors.New("邮箱已存在")
	ErrMobileExists          = errors.New("手机号已存在")
	ErrUserArchived          = errors.New("用户已归档")
	ErrInvalidUserStatus     = errors.New("无效的用户状态")
	ErrPasswordIncorrect     = errors.New("密码不正确")

	// 角色相关错误
	ErrRoleNotFound           = errors.New("角色不存在")
	ErrRoleAlreadyExists      = errors.New("角色已存在")
	ErrRoleDisabled           = errors.New("角色已禁用")
	ErrSystemRoleNotDeletable = errors.New("系统角色不可删除")

	// 权限相关错误
	ErrPermissionNotFound           = errors.New("权限不存在")
	ErrPermissionAlreadyExists      = errors.New("权限已存在")
	ErrPermissionDisabled           = errors.New("权限已禁用")
	ErrSystemPermissionNotDeletable = errors.New("系统权限不可删除")
	ErrHasChildPermissions          = errors.New("存在子权限，无法删除")

	// 数据库相关错误
	ErrDBNotInitialized = errors.New("数据库未初始化")

	// 其他通用错误
	ErrInvalidRequest = errors.New("无效的请求")
	ErrUnauthorized   = errors.New("未授权")
	ErrForbidden      = errors.New("禁止访问")
	ErrInternalServer = errors.New("内部服务器错误")

	// 新增的错误
	ErrWrongPassword                 = errors.New("密码错误")
	ErrPasswordLock                  = errors.New("密码已锁定")
	ErrAccountDisabled               = errors.New("账号已禁用")
	ErrRoleCodeExists                = errors.New("角色代码已存在")
	ErrSystemRoleNotModifiable       = errors.New("系统角色不允许修改")
	ErrPermissionCodeExists          = errors.New("权限代码已存在")
	ErrSystemPermissionNotModifiable = errors.New("系统权限不允许修改")
	ErrInvalidToken                  = errors.New("无效的令牌")
	ErrTokenExpired                  = errors.New("令牌已过期")
	ErrNotFound                      = errors.New("资源不存在")
	ErrInvalidParameter              = errors.New("无效的参数")
	ErrServiceUnavailable            = errors.New("服务不可用")
	ErrTooManyRequests               = errors.New("请求过多")
	ErrBadRequest                    = errors.New("错误的请求")
	ErrMethodNotAllowed              = errors.New("方法不允许")
	ErrRequestTimeout                = errors.New("请求超时")
)
