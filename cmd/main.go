package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	// TODO: 初始化配置

	// TODO: 初始化日志

	// TODO: 连接数据库

	// TODO: 初始化HTTP服务

	// 等待中断信号优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 启动一个临时HTTP服务器以便能正常退出
	fmt.Println("Server is starting...")

	// 模拟服务器运行
	fmt.Println("Server is running...")

	// 等待中断信号
	<-quit
	fmt.Println("Shutting down server...")

	// TODO: 关闭数据库连接

	// TODO: 关闭HTTP服务

	log.Println("Server exited properly")
}
