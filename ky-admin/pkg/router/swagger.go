package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupSwaggerRoutes 配置Swagger路由
func SetupSwaggerRoutes(r *gin.Engine) {
	// 添加Swagger UI路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
