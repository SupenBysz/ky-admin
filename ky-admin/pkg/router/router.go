package router

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SupenBysz/ky-admin/internal/api/controller"
	"github.com/SupenBysz/ky-admin/internal/config"
	"github.com/SupenBysz/ky-admin/internal/pkg/response"
	"github.com/SupenBysz/ky-admin/internal/repository"
	"github.com/SupenBysz/ky-admin/internal/routes"
	"github.com/SupenBysz/ky-admin/internal/service"
	pkgconfig "github.com/SupenBysz/ky-admin/pkg/config"
	"github.com/SupenBysz/ky-admin/pkg/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	// 导入Swagger文档
	_ "github.com/SupenBysz/ky-admin/swagger"
)

// SetupRouter 配置路由
func SetupRouter(cfg *pkgconfig.Config) *gin.Engine {
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

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ok"})
	})
	r.GET("/livez", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "alive"})
	})
	r.GET("/readyz", func(c *gin.Context) {
		response.Success(c, gin.H{"status": "ready"})
	})

	// API路由组
	apiGroup := r.Group("/api")
	{
		// v1版本API组
		v1Group := apiGroup.Group("/v1")
		{
			// 认证相关
			authGroup := v1Group.Group("/auth")
			{
				authGroup.POST("/login", func(c *gin.Context) {
					response.Success(c, gin.H{"message": "登录接口待实现"})
				})
			}

			// 用户相关API
			userGroup := v1Group.Group("/users")
			{
				userGroup.GET("", func(c *gin.Context) {
					response.Success(c, gin.H{"users": []interface{}{}})
				})
			}

			// 角色相关API
			roleGroup := v1Group.Group("/roles")
			{
				roleGroup.GET("", func(c *gin.Context) {
					response.Success(c, []interface{}{})
				})
			}

			// 权限相关API
			permGroup := v1Group.Group("/permissions")
			{
				permGroup.GET("", func(c *gin.Context) {
					response.Success(c, []interface{}{})
				})
			}

			// 密码相关API，无需认证
			passwordController := createPasswordController(cfg)
			routes.RegisterPasswordRoutes(v1Group, passwordController)
		}
	}

	// 404处理
	r.NoRoute(func(c *gin.Context) {
		response.NotFound(c, "请求的资源不存在")
	})

	return r
}

// createPasswordController 创建密码控制器
func createPasswordController(cfg *pkgconfig.Config) *controller.PasswordController {
	// 这里应该从依赖注入容器中获取，此处临时创建
	db := getDB() // 实际应用中应该从容器或上下文中获取

	userRepo := repository.NewUserRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)

	emailConfig := config.EmailConfig{
		SMTPHost:     cfg.GetString("email.smtp_host"),
		SMTPPort:     cfg.GetInt("email.smtp_port"),
		SMTPUsername: cfg.GetString("email.smtp_username"),
		SMTPPassword: cfg.GetString("email.smtp_password"),
		SenderEmail:  cfg.GetString("email.sender_email"),
		SenderName:   cfg.GetString("email.sender_name"),
		TemplateDir:  cfg.GetString("email.template_dir"),
		EnableTLS:    cfg.GetBool("email.enable_tls"),
	}

	emailSvc := service.NewEmailService(emailConfig)
	appURL := cfg.GetString("app.url")

	passwordSvc := service.NewPasswordService(userRepo, passwordResetRepo, emailSvc, appURL)
	return controller.NewPasswordController(passwordSvc)
}

// getDB 临时函数，获取DB实例，实际应从DI容器获取
func getDB() *gorm.DB {
	// 这里应该从全局上下文或依赖注入容器中获取DB实例
	// 实际应用中不应该这样做
	return nil // 真实应用需要返回实际的DB实例
}

// RunServer 启动服务器
func RunServer(cfg *pkgconfig.Config) error {
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
