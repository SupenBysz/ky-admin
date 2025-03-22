package controller

import (
	"net/http"

	"github.com/SupenBysz/ky-admin/internal/dto"
	"github.com/SupenBysz/ky-admin/internal/pkg/response"
	"github.com/SupenBysz/ky-admin/internal/service"
	"github.com/gin-gonic/gin"
)

// PasswordController 密码控制器
type PasswordController struct {
	passwordService service.PasswordService
}

// NewPasswordController 创建密码控制器
func NewPasswordController(passwordService service.PasswordService) *PasswordController {
	return &PasswordController{
		passwordService: passwordService,
	}
}

// ForgotPassword 忘记密码
// @Summary 请求密码重置
// @Description 通过邮箱请求密码重置
// @Tags 密码
// @Accept json
// @Produce json
// @Param request body dto.ForgotPasswordDTO true "忘记密码请求"
// @Success 200 {object} response.Response{data=bool} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/password/forgot [post]
func (c *PasswordController) ForgotPassword(ctx *gin.Context) {
	var req dto.ForgotPasswordDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	if err := c.passwordService.ForgotPassword(&req); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "请求密码重置失败", err.Error())
		return
	}

	response.SuccessWithMessage(ctx, true, "密码重置邮件已发送")
}

// ValidateResetToken 验证重置令牌
// @Summary 验证重置令牌
// @Description 验证密码重置令牌是否有效
// @Tags 密码
// @Accept json
// @Produce json
// @Param token query string true "重置令牌"
// @Success 200 {object} response.Response{data=bool} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/password/validate-token [get]
func (c *PasswordController) ValidateResetToken(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		response.Fail(ctx, http.StatusBadRequest, "令牌不能为空", nil)
		return
	}

	valid, err := c.passwordService.ValidateResetToken(token)
	if err != nil {
		response.Fail(ctx, http.StatusBadRequest, "令牌验证失败", err.Error())
		return
	}

	response.SuccessWithMessage(ctx, valid, "令牌验证成功")
}

// ResetPassword 重置密码
// @Summary 重置密码
// @Description 使用令牌重置密码
// @Tags 密码
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordDTO true "重置密码请求"
// @Success 200 {object} response.Response{data=bool} "成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "内部服务器错误"
// @Router /api/v1/password/reset [post]
func (c *PasswordController) ResetPassword(ctx *gin.Context) {
	var req dto.ResetPasswordDTO
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "请求参数错误", err.Error())
		return
	}

	if err := c.passwordService.ResetPassword(&req); err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "密码重置失败", err.Error())
		return
	}

	response.SuccessWithMessage(ctx, true, "密码已重置")
}
