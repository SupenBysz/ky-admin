package api

import (
	"github.com/SupenBysz/ky-admin/internal/common/errors"
	"github.com/SupenBysz/ky-admin/internal/dto"
	"github.com/SupenBysz/ky-admin/internal/service"
	"github.com/SupenBysz/ky-admin/pkg/api"
	"github.com/gin-gonic/gin"
)

// TokenResponse 令牌响应
type TokenResponse struct {
	Token        string      `json:"token"`         // 访问令牌
	RefreshToken string      `json:"refresh_token"` // 刷新令牌
	User         interface{} `json:"user"`          // 用户信息
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"` // 刷新令牌
}

// RefreshTokenResponse 刷新令牌响应
type RefreshTokenResponse struct {
	Token        string `json:"token"`         // 新访问令牌
	RefreshToken string `json:"refresh_token"` // 新刷新令牌
}

// AuthController 认证控制器
type AuthController struct {
	authService service.AuthService
}

// NewAuthController 创建认证控制器
func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login godoc
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param data body dto.LoginDTO true "登录信息"
// @Success 200 {object} api.Response{data=TokenResponse} "登录成功"
// @Failure 400 {object} api.Response "参数错误"
// @Failure 401 {object} api.Response "认证失败"
// @Router /api/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req dto.LoginDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的登录参数", err))
		return
	}

	// 调用认证服务登录
	token, refreshToken, user, err := c.authService.Login(&req)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	// 返回令牌和用户信息
	api.Success(ctx, TokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	})
}

// Register godoc
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param data body dto.RegisterDTO true "注册信息"
// @Success 200 {object} api.Response{data=dto.UserResponseDTO} "注册成功"
// @Failure 400 {object} api.Response "参数错误"
// @Router /api/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req dto.RegisterDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的注册参数", err))
		return
	}

	// 调用认证服务注册
	user, err := c.authService.Register(&req)
	if err != nil {
		api.Error(ctx, err)
		return
	}

	// 返回用户信息
	api.Success(ctx, user)
}

// RefreshToken godoc
// @Summary 刷新令牌
// @Description 刷新用户访问令牌
// @Tags 认证管理
// @Accept json
// @Produce json
// @Param data body RefreshTokenRequest true "刷新令牌信息"
// @Success 200 {object} api.Response{data=RefreshTokenResponse} "刷新成功"
// @Failure 400 {object} api.Response "参数错误"
// @Failure 401 {object} api.Response "令牌无效"
// @Router /api/auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		api.Error(ctx, errors.NewAppError(errors.CodeInvalidParams, "无效的刷新令牌参数", err))
		return
	}

	// 调用认证服务刷新令牌
	token, refreshToken, err := c.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		api.Error(ctx, errors.NewAuthError(err))
		return
	}

	// 返回新的令牌
	api.Success(ctx, RefreshTokenResponse{
		Token:        token,
		RefreshToken: refreshToken,
	})
}
