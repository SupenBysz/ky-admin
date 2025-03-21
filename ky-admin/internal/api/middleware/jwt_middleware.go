package middleware

import (
	"strings"

	"github.com/SupenBysz/ky-admin/internal/pkg/response"
	"github.com/SupenBysz/ky-admin/internal/service"
	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware JWT认证中间件
func JWTAuthMiddleware(jwtService service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "请求未授权")
			c.Abort()
			return
		}

		// 检查Authorization格式
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "请求头Authorization格式错误")
			c.Abort()
			return
		}

		// 解析token
		token := parts[1]
		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			response.Unauthorized(c, "无效的令牌："+err.Error())
			c.Abort()
			return
		}

		// 将用户信息保存到上下文
		c.Set("userId", claims.UserID)
		c.Next()
	}
}

// CurrentUser 获取当前登录用户ID
func CurrentUser(c *gin.Context) uint {
	userID, exists := c.Get("userId")
	if !exists {
		return 0
	}
	return userID.(uint)
}
