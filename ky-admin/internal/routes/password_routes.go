package routes

import (
	"github.com/SupenBysz/ky-admin/internal/api/controller"
	"github.com/gin-gonic/gin"
)

// RegisterPasswordRoutes 注册密码相关路由
func RegisterPasswordRoutes(r *gin.RouterGroup, passwordController *controller.PasswordController) {
	password := r.Group("/password")
	{
		// 忘记密码
		password.POST("/forgot", passwordController.ForgotPassword)
		// 验证重置令牌
		password.GET("/validate-token", passwordController.ValidateResetToken)
		// 重置密码
		password.POST("/reset", passwordController.ResetPassword)
	}
}
