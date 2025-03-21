package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/SupenBysz/ky-admin/pkg/config"
	"github.com/SupenBysz/ky-admin/pkg/database"
	"github.com/SupenBysz/ky-admin/pkg/logger"
	"github.com/SupenBysz/ky-admin/pkg/models"
	"github.com/SupenBysz/ky-admin/pkg/router"
	"go.uber.org/zap"
)

func main() {
	// 获取工作目录
	workDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("获取工作目录失败: %s\n", err)
		os.Exit(1)
	}

	// 初始化配置
	configPath := filepath.Join(workDir, "configs", "config.yaml")
	fmt.Printf("使用配置文件: %s\n", configPath)
	// 这里可以加载配置文件，例如 config.Init(configPath)

	cfg := config.GetConfig()

	// 初始化日志
	log := logger.GetLogger()

	// 创建数据库管理器
	dbManager := database.NewDBManager(cfg)
	dbManager.SetLogger(log)

	// 连接数据库
	if err := dbManager.Connect(); err != nil {
		log.Error("连接数据库失败", zap.Error(err))
		os.Exit(1)
	}
	defer dbManager.Close()

	// 注册数据库到模型包
	models.RegisterDB(dbManager.GetDB())

	// 检查数据库健康状态
	if err := dbManager.CheckHealth(); err != nil {
		log.Error("数据库健康检查失败", zap.Error(err))
		os.Exit(1)
	}

	log.Info("KY-Admin API服务启动中...")

	// 初始化HTTP服务器
	r := router.SetupRouter(cfg)

	// 获取服务器配置
	serverHost := cfg.GetString("server.host")
	serverPort := cfg.GetInt("server.port")
	serverAddr := fmt.Sprintf("%s:%d", serverHost, serverPort)

	// 服务器配置
	server := &http.Server{
		Addr:           serverAddr,
		Handler:        r,
		ReadTimeout:    time.Duration(cfg.GetIntWithDefault("server.read_timeout", 10)) * time.Second,
		WriteTimeout:   time.Duration(cfg.GetIntWithDefault("server.write_timeout", 10)) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// 设置关闭信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器（非阻塞）
	go func() {
		log.Info("HTTP服务器启动", zap.String("addr", serverAddr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP服务器启动失败", zap.Error(err))
			os.Exit(1)
		}
	}()

	// 等待中断信号
	<-quit
	log.Info("收到退出信号，正在关闭服务器...")

	// 优雅关闭的超时时间
	timeout := cfg.GetIntWithDefault("server.shutdown_timeout", 5)

	// 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// 优雅关闭HTTP服务
	if err := server.Shutdown(ctx); err != nil {
		log.Error("服务器关闭出错", zap.Error(err))
	}

	log.Info("KY-Admin API服务已停止")
}
