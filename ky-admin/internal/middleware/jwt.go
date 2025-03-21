package middleware

import (
	"strings"

	appError "github.com/SupenBysz/ky-admin/internal/common/errors"
	"github.com/SupenBysz/ky-admin/internal/service"
	"github.com/SupenBysz/ky-admin/pkg/api"
	"github.com/gin-gonic/gin"
)

// JWT JWT认证中间件
func JWT(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			api.Error(c, appError.NewAuthError(appError.ErrInvalidToken))
			c.Abort()
			return
		}

		// 检查Token格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			api.Error(c, appError.NewAuthError(appError.ErrInvalidToken))
			c.Abort()
			return
		}

		// 解析并验证Token
		claims, err := jwtService.ValidateToken(parts[1])
		if err != nil {
			api.Error(c, appError.NewAuthError(err))
			c.Abort()
			return
		}

		// 将用户ID存储在上下文中
		c.Set("userID", claims.UserID)

		c.Next()
	}
}

// CurrentUser 获取当前用户ID
func CurrentUser(c *gin.Context) uint {
	userID, exists := c.Get("userID")
	if !exists {
		return 0
	}
	return userID.(uint)
}
