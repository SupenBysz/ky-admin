package router

import (
	"fmt"
	"net/http"
	"time"

	internalAPI "github.com/SupenBysz/ky-admin/internal/api"
	"github.com/SupenBysz/ky-admin/pkg/api"
	"github.com/SupenBysz/ky-admin/pkg/config"
	"github.com/SupenBysz/ky-admin/pkg/middleware"
	"github.com/gin-gonic/gin"

	// 导入Swagger文档
	_ "github.com/SupenBysz/ky-admin/swagger"
)

// SetupRouter 配置路由
func SetupRouter(cfg *config.Config) *gin.Engine {
	// 根据配置设置模式
	if cfg.GetString("environment") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	// 中间件
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(middleware.DefaultRateLimit())

	// 配置Swagger路由
	SetupSwaggerRoutes(r)

	// 初始化控制器
	healthController := internalAPI.NewHealthController()
	userController := internalAPI.NewUserController()

	// 健康检查
	r.GET("/health", healthController.Health)
	r.GET("/livez", healthController.LivenessProbe)
	r.GET("/readyz", healthController.ReadinessProbe)

	// API路由组
	apiGroup := r.Group("/api")
	{
		// 认证相关
		authGroup := apiGroup.Group("/auth")
		{
			authGroup.POST("/login", userController.Login)
		}

		// 用户相关API
		userGroup := apiGroup.Group("/users")
		{
			userGroup.GET("", userController.GetUsers)
		}

		// 角色相关API
		roleGroup := apiGroup.Group("/roles")
		{
			roleGroup.GET("", func(c *gin.Context) {
				api.Success(c, []interface{}{})
			})
		}

		// 权限相关API
		permGroup := apiGroup.Group("/permissions")
		{
			permGroup.GET("", func(c *gin.Context) {
				api.Success(c, []interface{}{})
			})
		}
	}

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		api.NotFoundError(c, "请求的资源不存在")
	})

	return r
}

// RunServer 启动服务器
func RunServer(cfg *config.Config) error {
	router := SetupRouter(cfg)

	// 获取服务器配置
	host := cfg.GetString("server.host")
	port := cfg.GetInt("server.port")
	addr := fmt.Sprintf("%s:%d", host, port)
	readTimeout := cfg.GetIntWithDefault("server.read_timeout", 10)
	writeTimeout := cfg.GetIntWithDefault("server.write_timeout", 10)

	server := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    time.Duration(readTimeout) * time.Second,
		WriteTimeout:   time.Duration(writeTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return server.ListenAndServe()
}
