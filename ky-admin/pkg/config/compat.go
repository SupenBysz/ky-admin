package config

import (
	"time"
)

// LoadConfig 加载配置文件并返回配置管理器
func LoadConfig(configPath string) (*Config, error) {
	return New(WithConfigPath(configPath))
}

// GetStringWithDefault 获取字符串配置，支持默认值
func (c *Config) GetStringWithDefault(key string, defaultValue string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetString(key)
}

// GetIntWithDefault 获取整数配置，支持默认值
func (c *Config) GetIntWithDefault(key string, defaultValue int) int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetInt(key)
}

// GetBoolWithDefault 获取布尔配置，支持默认值
func (c *Config) GetBoolWithDefault(key string, defaultValue bool) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetBool(key)
}

// GetDurationWithDefault 获取时间间隔配置，支持默认值
func (c *Config) GetDurationWithDefault(key string, defaultValue time.Duration) time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetDuration(key)
}
