package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SupenBysz/ky-admin/pkg/config"
	"github.com/SupenBysz/ky-admin/pkg/database"
	"github.com/SupenBysz/ky-admin/pkg/logger"
	"github.com/SupenBysz/ky-admin/pkg/models"
	"github.com/SupenBysz/ky-admin/pkg/router"
	"go.uber.org/zap"
)

// 版本信息，通过编译时传入
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
)

func main() {
	// 显示版本信息
	fmt.Printf("KY-Admin Version: %s, Build: %s, Commit: %s\n", Version, BuildTime, GitCommit)

	// 初始化配置
	cfg := config.GetConfig()
	env := cfg.LoadEnvironment()
	fmt.Printf("当前环境: %s\n", env)

	// 初始化日志系统
	if err := logger.InitFromConfig(); err != nil {
		fmt.Printf("初始化日志系统失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// 记录启动日志
	logger.Info("KY-Admin服务启动",
		zap.String("version", Version),
		zap.String("build_time", BuildTime),
		zap.String("git_commit", GitCommit),
		zap.String("environment", env),
	)

	// 获取服务器配置
	serverHost := cfg.GetString("server.host")
	serverPort := cfg.GetInt("server.port")
	serverAddr := fmt.Sprintf("%s:%d", serverHost, serverPort)

	logger.Info("服务配置加载完成",
		zap.String("server_addr", serverAddr),
	)

	// 数据库连接信息
	dbHost := cfg.GetString("database.host")
	dbPort := cfg.GetInt("database.port")
	dbName := cfg.GetString("database.database")

	logger.Info("数据库配置",
		zap.String("db_host", dbHost),
		zap.Int("db_port", dbPort),
		zap.String("db_name", dbName),
	)

	// 连接数据库
	dbManager := database.NewDBManager(cfg)
	dbManager.SetLogger(logger.GetLogger())

	if err := dbManager.Connect(); err != nil {
		logger.Error("连接数据库失败", zap.Error(err))
		os.Exit(1)
	}
	defer dbManager.Close()

	// 注册数据库连接到模型
	models.RegisterDB(dbManager.GetDB())

	logger.Info("数据库连接成功")

	// 初始化HTTP服务
	r := router.SetupRouter(cfg)

	// 服务器配置
	server := &http.Server{
		Addr:           serverAddr,
		Handler:        r,
		ReadTimeout:    time.Duration(cfg.GetIntWithDefault("server.read_timeout", 10)) * time.Second,
		WriteTimeout:   time.Duration(cfg.GetIntWithDefault("server.write_timeout", 10)) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// 启动HTTP服务器（非阻塞）
	go func() {
		logger.Info("HTTP服务器正在启动", zap.String("addr", serverAddr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP服务器启动失败", zap.Error(err))
			os.Exit(1)
		}
	}()

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 等待中断信号
	<-quit
	logger.Info("收到退出信号，正在关闭服务器...")

	// 优雅关闭的超时时间
	timeout := cfg.GetInt("server.shutdown_timeout")
	if timeout <= 0 {
		timeout = 5
	}

	// 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	logger.Info("等待优雅关闭",
		zap.Int("timeout_seconds", timeout),
	)

	// 优雅关闭HTTP服务
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭出错", zap.Error(err))
	}

	logger.Info("服务已正常退出")
}
