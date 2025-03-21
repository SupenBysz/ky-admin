package api

import (
	"github.com/SupenBysz/ky-admin/pkg/api"
	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`       // 用户名
	Password string `json:"password" binding:"required" example:"password123"` // 密码
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token  string `json:"token"`  // JWT令牌
	UserID uint   `json:"userId"` // 用户ID
}

// UserController 用户控制器
type UserController struct{}

// NewUserController 创建用户控制器
func NewUserController() *UserController {
	return &UserController{}
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body LoginRequest true "登录信息"
// @Success 200 {object} api.Response{data=LoginResponse} "登录成功"
// @Failure 400 {object} api.Response "参数错误"
// @Failure 401 {object} api.Response "认证失败"
// @Router /api/auth/login [post]
func (c *UserController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.ParamError(ctx, "无效的登录参数")
		return
	}

	// 这里是登录逻辑，实际项目中需要实现真正的登录功能
	// 示例中返回模拟数据
	api.Success(ctx, LoginResponse{
		Token:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
		UserID: 1,
	})
}

// GetUsers godoc
// @Summary 获取用户列表
// @Description 获取系统用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} api.Response{data=[]map[string]interface{}} "用户列表"
// @Router /api/users [get]
func (c *UserController) GetUsers(ctx *gin.Context) {
	// 示例数据，实际项目中应该从数据库获取
	users := []map[string]interface{}{
		{"id": 1, "username": "admin", "nickname": "管理员", "email": "admin@example.com"},
		{"id": 2, "username": "user1", "nickname": "用户1", "email": "user1@example.com"},
	}
	api.Success(ctx, users)
}
