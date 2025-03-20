package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SupenBysz/ky-admin/pkg/config"
	"github.com/SupenBysz/ky-admin/pkg/logger"
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
		log.Fatalf("初始化日志系统失败: %v", err)
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

	// TODO: 连接数据库

	// TODO: 初始化HTTP服务

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动一个临时HTTP服务器以便能正常退出
	logger.Info("服务器正在启动...")

	// 模拟服务器运行
	logger.Info("服务器正在运行...")

	// 等待中断信号
	<-quit
	logger.Info("收到退出信号，正在关闭服务器...")

	// 优雅关闭的超时时间
	timeout := cfg.GetInt("server.shutdown_timeout")
	if timeout <= 0 {
		timeout = 5
	}

	logger.Info("等待优雅关闭",
		zap.Int("timeout_seconds", timeout),
	)
	time.Sleep(time.Duration(timeout) * time.Second)

	// TODO: 关闭数据库连接

	// TODO: 关闭HTTP服务

	logger.Info("服务已正常退出")
}
