package logger

import (
	"path/filepath"

	"github.com/SupenBysz/ky-admin/pkg/config"
)

// InitFromConfig 从配置中初始化日志系统
func InitFromConfig() error {
	cfg := config.GetConfig()

	// 获取日志配置
	logPath := cfg.GetString("log.path")
	if logPath == "" {
		logPath = "logs/ky-admin.log"
	}

	// 确保日志路径是绝对路径
	if !filepath.IsAbs(logPath) {
		// 如果是相对路径，使用配置中的应用根目录作为基准
		appRoot := cfg.GetString("app.root")
		if appRoot != "" {
			logPath = filepath.Join(appRoot, logPath)
		}
	}

	logLevel := cfg.GetString("log.level")
	if logLevel == "" {
		logLevel = "info"
	}

	// 从环境变量获取开发模式标志
	isDevelopment := cfg.GetBool("app.development")

	// 使用通用初始化方法
	return InitLogger(logPath, logLevel, isDevelopment)
}
