package logger

import (
	"github.com/SupenBysz/ky-admin/pkg/config"
	"go.uber.org/zap"
)

// InitLoggerWithConfig 使用配置初始化日志系统
func InitLoggerWithConfig(conf *config.Config) error {
	// 从配置中获取日志相关设置
	logPath := conf.GetString("log.filepath")
	if logPath == "" {
		logPath = "./logs/app.log"
	}

	isDev := conf.GetString("app.env") == "development"

	logLevel := conf.GetString("log.level")
	if logLevel == "" {
		logLevel = "info"
	}

	// 调用基础日志初始化函数
	return InitLogger(logPath, logLevel, isDev)
}

// Sync 同步日志缓冲区
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

// Err 创建错误字段
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Str 创建字符串字段
func Str(key, val string) zap.Field {
	return zap.String(key, val)
}

// Int 创建整数字段
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Bool 创建布尔字段
func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

// Any 创建任意类型字段
func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}
