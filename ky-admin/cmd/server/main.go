package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SupenBysz/ky-admin/internal/api/router"
	"github.com/SupenBysz/ky-admin/internal/config"
	"github.com/SupenBysz/ky-admin/internal/repository"
	"github.com/SupenBysz/ky-admin/internal/service"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AppConfig 应用配置
type AppConfig struct {
	Env  string
	HTTP struct {
		Host string
		Port string
	}
}

func main() {
	// 初始化配置（简化，实际应该从文件或环境变量加载）
	appConfig := &AppConfig{
		Env: "dev",
		HTTP: struct {
			Host string
			Port string
		}{
			Host: "0.0.0.0",
			Port: "8080",
		},
	}

	// 初始化Gin
	gin.SetMode(gin.DebugMode)
	if appConfig.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化数据库连接（简化，实际应该使用配置参数）
	var db *gorm.DB = nil // 实际应该初始化真实的数据库连接

	// 初始化Casbin（简化，实际应该从持久化存储加载模型和策略）
	var enforcer *casbin.Enforcer = nil // 实际应该初始化真实的Casbin enforcer

	// 创建仓库
	userRepo := repository.NewUserRepository(db)

	// 创建JWT配置
	jwtConfig := config.JWTConfig{
		SecretKey:     "your-secret-key",
		TokenExpire:   24,  // 24小时
		RefreshExpire: 168, // 7天
		Issuer:        "ky-admin-server",
	}

	// 创建服务
	jwtService := service.NewJWTService(jwtConfig)
	userService := service.NewUserService(userRepo)
	permissionService := service.NewPermissionService(enforcer)

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	router.RegisterRoutes(r, userService, jwtService, permissionService, enforcer)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    appConfig.HTTP.Host + ":" + appConfig.HTTP.Port,
		Handler: r,
	}

	// 启动HTTP服务器（非阻塞）
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号（优雅关闭）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 设置关闭超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
