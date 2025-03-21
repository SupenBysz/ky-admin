package middleware

import (
	"time"

	"github.com/SupenBysz/ky-admin/pkg/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 返回CORS中间件
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                                                     // 允许所有来源访问
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},                                 // 允许的HTTP方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"}, // 允许的HTTP头
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},                                        // 暴露的HTTP头
		AllowCredentials: true,                                                                              // 允许携带凭证
		MaxAge:           12 * time.Hour,                                                                    // 预检请求缓存时间
	})
}

// CORSWithConfig 根据配置返回CORS中间件
func CORSWithConfig(cfg *config.Config) gin.HandlerFunc {
	// 如果配置中定义了CORS设置，则使用配置值
	if cfg == nil {
		return CORS() // 使用默认设置
	}

	// 从配置获取CORS设置
	return cors.New(cors.Config{
		AllowOrigins:     cfg.GetStringSlice("cors.allow_origins"),
		AllowMethods:     cfg.GetStringSlice("cors.allow_methods"),
		AllowHeaders:     cfg.GetStringSlice("cors.allow_headers"),
		ExposeHeaders:    cfg.GetStringSlice("cors.expose_headers"),
		AllowCredentials: cfg.GetBool("cors.allow_credentials"),
		MaxAge:           time.Duration(cfg.GetInt("cors.max_age")) * time.Second,
	})
}
