package controller

import (
	"github.com/SupenBysz/ky-admin/pkg/api"
	"github.com/gin-gonic/gin"
)

// HealthController 健康检查控制器
type HealthController struct{}

// NewHealthController 创建健康检查控制器
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Check 健康检查
func (c *HealthController) Check(ctx *gin.Context) {
	api.SuccessWithMessage(ctx, "系统运行正常", nil)
}

// Version 获取版本信息
func (c *HealthController) Version(ctx *gin.Context) {
	version := map[string]string{
		"version": "1.0.0",
		"build":   "2023-04-01",
	}
	api.Success(ctx, version)
}
