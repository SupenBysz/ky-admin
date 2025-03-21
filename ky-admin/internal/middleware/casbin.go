package middleware

import (
	"strconv"
	"strings"

	appError "github.com/SupenBysz/ky-admin/internal/common/errors"
	"github.com/SupenBysz/ky-admin/internal/service"
	"github.com/SupenBysz/ky-admin/pkg/api"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

// Authorize 授权中间件
func Authorize(enforcer *casbin.Enforcer, userService service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		userID := CurrentUser(c)
		if userID == 0 {
			api.Error(c, appError.NewAuthError(appError.ErrInvalidToken))
			c.Abort()
			return
		}

		// 获取请求路径和方法
		obj := c.Request.URL.Path
		act := c.Request.Method

		// 获取用户角色信息
		userDTO, err := userService.GetUserWithRoles(userID)
		if err != nil {
			api.Error(c, appError.NewSystemError(err))
			c.Abort()
			return
		}

		if len(userDTO.Roles) == 0 {
			api.Error(c, appError.NewPermissionError(appError.ErrNoPermission))
			c.Abort()
			return
		}

		// 检查角色是否有权限访问此路径
		hasPermission := false
		sub := strconv.FormatUint(uint64(userID), 10)

		// 检查用户是否有权限访问
		for _, role := range userDTO.Roles {
			// 超级管理员拥有所有权限
			if role.Code == "superadmin" {
				hasPermission = true
				break
			}
		}

		// 如果不是超级管理员，检查权限
		if !hasPermission {
			ok, err := enforcer.Enforce(sub, obj, act)
			if err != nil {
				api.Error(c, appError.NewSystemError(err))
				c.Abort()
				return
			}

			hasPermission = ok
		}

		if !hasPermission {
			api.Error(c, appError.NewPermissionError(appError.ErrPermissionDenied))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthorizeByRoles 基于角色的授权中间件
func AuthorizeByRoles(userService service.UserService, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		userID := CurrentUser(c)
		if userID == 0 {
			api.Error(c, appError.NewAuthError(appError.ErrInvalidToken))
			c.Abort()
			return
		}

		// 获取用户角色信息
		userDTO, err := userService.GetUserWithRoles(userID)
		if err != nil {
			api.Error(c, appError.NewSystemError(err))
			c.Abort()
			return
		}

		if len(userDTO.Roles) == 0 {
			api.Error(c, appError.NewPermissionError(appError.ErrNoPermission))
			c.Abort()
			return
		}

		// 检查用户是否有允许的角色
		hasRole := false
		for _, userRole := range userDTO.Roles {
			// 超级管理员拥有所有权限
			if userRole.Code == "superadmin" {
				hasRole = true
				break
			}

			for _, allowedRole := range allowedRoles {
				if userRole.Code == allowedRole {
					hasRole = true
					break
				}
			}

			if hasRole {
				break
			}
		}

		if !hasRole {
			api.Error(c, appError.NewPermissionError(appError.ErrPermissionDenied))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AuthorizeByPermission 基于权限的授权中间件
func AuthorizeByPermission(enforcer *casbin.Enforcer, userService service.UserService, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		userID := CurrentUser(c)
		if userID == 0 {
			api.Error(c, appError.NewAuthError(appError.ErrInvalidToken))
			c.Abort()
			return
		}

		// 解析权限字符串 "module:action"
		parts := strings.Split(permission, ":")
		if len(parts) != 2 {
			api.Error(c, appError.NewSystemError(appError.ErrInternalServer))
			c.Abort()
			return
		}

		module, action := parts[0], parts[1]

		// 获取用户角色信息
		userDTO, err := userService.GetUserWithRoles(userID)
		if err != nil {
			api.Error(c, appError.NewSystemError(err))
			c.Abort()
			return
		}

		if len(userDTO.Roles) == 0 {
			api.Error(c, appError.NewPermissionError(appError.ErrNoPermission))
			c.Abort()
			return
		}

		// 检查角色是否有此权限
		hasPermission := false
		sub := strconv.FormatUint(uint64(userID), 10)

		for _, role := range userDTO.Roles {
			// 超级管理员拥有所有权限
			if role.Code == "superadmin" {
				hasPermission = true
				break
			}
		}

		// 如果不是超级管理员，检查权限
		if !hasPermission {
			ok, err := enforcer.Enforce(sub, module, action)
			if err != nil {
				api.Error(c, appError.NewSystemError(err))
				c.Abort()
				return
			}

			hasPermission = ok
		}

		if !hasPermission {
			api.Error(c, appError.NewPermissionError(appError.ErrPermissionDenied))
			c.Abort()
			return
		}

		c.Next()
	}
}
