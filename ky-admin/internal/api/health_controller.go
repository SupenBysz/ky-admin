package api

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

// Health godoc
// @Summary 健康检查接口
// @Description 返回服务健康状态
// @Tags 健康检查
// @Accept json
// @Produce json
// @Success 200 {object} api.Response{data=map[string]string} "服务健康"
// @Router /health [get]
func (h *HealthController) Health(c *gin.Context) {
	api.Success(c, gin.H{"status": "ok"})
}

// LivenessProbe godoc
// @Summary 存活探针
// @Description 用于Kubernetes存活探测
// @Tags 健康检查
// @Accept json
// @Produce json
// @Success 200 {object} api.Response{data=map[string]string} "服务存活"
// @Router /livez [get]
func (h *HealthController) LivenessProbe(c *gin.Context) {
	api.Success(c, gin.H{"status": "alive"})
}

// ReadinessProbe godoc
// @Summary 就绪探针
// @Description 用于Kubernetes就绪探测
// @Tags 健康检查
// @Accept json
// @Produce json
// @Success 200 {object} api.Response{data=map[string]string} "服务就绪"
// @Router /readyz [get]
func (h *HealthController) ReadinessProbe(c *gin.Context) {
	api.Success(c, gin.H{"status": "ready"})
}
