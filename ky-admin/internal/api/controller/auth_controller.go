package controller

import (
	"github.com/SupenBysz/ky-admin/internal/common/errors"
	"github.com/SupenBysz/ky-admin/internal/dto"
	"github.com/SupenBysz/ky-admin/internal/model"
	"github.com/SupenBysz/ky-admin/internal/service"
	"github.com/SupenBysz/ky-admin/pkg/api"
	"github.com/gin-gonic/gin"
)

// AuthController 认证控制器
type AuthController struct {
	userService service.UserService
	jwtService  service.JWTService
	authService service.AuthService
}

// NewAuthController 创建认证控制器
func NewAuthController(userService service.UserService, jwtService service.JWTService) *AuthController {
	return &AuthController{
		userService: userService,
		jwtService:  jwtService,
	}
}

// Login 用户登录
func (c *AuthController) Login(ctx *gin.Context) {
	// 解析请求参数
	var loginReq dto.LoginDTO
	if err := ctx.ShouldBindJSON(&loginReq); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 用户登录
	user, err := c.userService.Login(loginReq.Username, loginReq.Password)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	// 生成token
	token, err := c.jwtService.GenerateToken(user.ID)
	if err != nil {
		api.Error(ctx, errors.NewSystemError(err))
		return
	}

	// 生成刷新令牌
	refreshToken, err := c.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		api.Error(ctx, errors.NewSystemError(err))
		return
	}

	// 获取用户信息
	userDTO, err := c.userService.GetUserWithRoles(user.ID)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	api.Success(ctx, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
		"user":          userDTO,
	})
}

// Register 用户注册
func (c *AuthController) Register(ctx *gin.Context) {
	// 解析请求参数
	var registerReq dto.RegisterDTO
	if err := ctx.ShouldBindJSON(&registerReq); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 创建用户
	// 哈希密码
	user := &model.User{
		Username: registerReq.Username,
		Password: registerReq.Password,
		Nickname: registerReq.Nickname,
		Email:    registerReq.Email,
		Phone:    registerReq.Phone,
		Status:   1, // 默认启用
	}

	if err := c.userService.CreateUser(user); err != nil {
		api.Error(ctx, err)
		return
	}

	// 分配角色
	if len(registerReq.RoleIDs) > 0 {
		if err := c.userService.AssignRoles(user.ID, registerReq.RoleIDs); err != nil {
			api.Error(ctx, err)
			return
		}
	}

	// 生成token
	token, err := c.jwtService.GenerateToken(user.ID)
	if err != nil {
		api.Error(ctx, errors.NewSystemError(err))
		return
	}

	// 生成刷新令牌
	refreshToken, err := c.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		api.Error(ctx, errors.NewSystemError(err))
		return
	}

	// 获取用户信息
	userDTO, err := c.userService.GetUserWithRoles(user.ID)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	api.Success(ctx, gin.H{
		"token":         token,
		"refresh_token": refreshToken,
		"user":          userDTO,
	})
}

// RefreshToken 刷新令牌
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	// 解析请求参数
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的请求参数", err))
		return
	}

	// 验证刷新令牌
	claims, err := c.jwtService.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		api.Error(ctx, errors.NewAuthError(err))
		return
	}

	// 生成新的访问令牌
	token, err := c.jwtService.GenerateToken(claims.UserID)
	if err != nil {
		api.Error(ctx, errors.NewSystemError(err))
		return
	}

	api.Success(ctx, gin.H{
		"token": token,
	})
}
