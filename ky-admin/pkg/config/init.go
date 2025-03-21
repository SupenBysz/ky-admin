package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 全局配置变量，仅供Init函数使用
var configInstance *viper.Viper

// Init 初始化配置
// 这是一个简化版的初始化函数，用于迁移和命令行工具
func Init(configPath string) error {
	// 创建viper实例
	v := viper.New()

	// 设置配置文件类型
	v.SetConfigType("yaml")

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取配置文件
	v.SetConfigFile(configPath)
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %s", err)
	}

	// 设置环境变量前缀，允许环境变量覆盖配置文件
	v.SetEnvPrefix("KYADMIN")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 添加配置热更新
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件已变更: %s\n", e.Name)
		// 这里可以加入更多的配置变更处理逻辑
	})

	// 设置全局实例
	configInstance = v

	return nil
}

// GetViper 获取Viper实例，供命令行工具和迁移使用
func GetViper() *viper.Viper {
	return configInstance
}

// LoadConfigFromEnv 从环境变量加载配置
func LoadConfigFromEnv() {
	if configInstance == nil {
		configInstance = viper.New()
		configInstance.SetEnvPrefix("KYADMIN")
		configInstance.AutomaticEnv()
		configInstance.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	}
}

// GetViperString 获取字符串配置，支持默认值
func GetViperString(key string, defaultValue string) string {
	if configInstance == nil {
		return defaultValue
	}

	value := configInstance.GetString(key)
	if value == "" {
		return defaultValue
	}

	return value
}

// GetViperInt 获取整数配置，支持默认值
func GetViperInt(key string, defaultValue int) int {
	if configInstance == nil {
		return defaultValue
	}

	if !configInstance.IsSet(key) {
		return defaultValue
	}

	return configInstance.GetInt(key)
}

// GetViperBool 获取布尔配置，支持默认值
func GetViperBool(key string, defaultValue bool) bool {
	if configInstance == nil {
		return defaultValue
	}

	if !configInstance.IsSet(key) {
		return defaultValue
	}

	return configInstance.GetBool(key)
}
