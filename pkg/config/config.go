package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 全局配置实例
var (
	globalConfig *Config
	once         sync.Once
)

// Config 配置管理器
type Config struct {
	viper      *viper.Viper
	configPath string
	configName string
	configType string
	mu         sync.RWMutex
}

// ConfigOption 配置选项
type ConfigOption func(*Config)

// WithConfigPath 设置配置文件路径
func WithConfigPath(path string) ConfigOption {
	return func(c *Config) {
		c.configPath = path
	}
}

// WithConfigName 设置配置文件名（不包含扩展名）
func WithConfigName(name string) ConfigOption {
	return func(c *Config) {
		c.configName = name
	}
}

// WithConfigType 设置配置文件类型
func WithConfigType(cfgType string) ConfigOption {
	return func(c *Config) {
		c.configType = cfgType
	}
}

// LoadConfig 加载配置文件
func (c *Config) LoadConfig() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	v := c.viper

	// 设置配置文件路径
	v.AddConfigPath(c.configPath)
	v.SetConfigName(c.configName)
	v.SetConfigType(c.configType)

	// 支持环境变量覆盖
	v.AutomaticEnv()
	v.SetEnvPrefix("KY_ADMIN")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 加载配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("配置文件未找到: %s", err)
		}
		return fmt.Errorf("读取配置文件错误: %s", err)
	}

	// 监听配置文件变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("配置文件已变更: %s\n", e.Name)
	})

	return nil
}

// New 创建配置管理器实例
func New(opts ...ConfigOption) (*Config, error) {
	c := &Config{
		viper:      viper.New(),
		configPath: "configs",
		configName: "config",
		configType: "yaml",
	}

	// 应用配置选项
	for _, opt := range opts {
		opt(c)
	}

	// 确保配置目录存在
	if _, err := os.Stat(c.configPath); os.IsNotExist(err) {
		if err := os.MkdirAll(c.configPath, 0755); err != nil {
			return nil, fmt.Errorf("创建配置目录失败: %s", err)
		}
	}

	// 加载配置
	if err := c.LoadConfig(); err != nil {
		return nil, err
	}

	return c, nil
}

// GetConfig 获取全局配置实例
func GetConfig() *Config {
	once.Do(func() {
		// 获取工作目录
		workDir, err := os.Getwd()
		if err != nil {
			panic(fmt.Errorf("获取工作目录失败: %s", err))
		}

		// 尝试查找配置文件目录
		configPath := filepath.Join(workDir, "configs")
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			// 回退到上一级目录查找
			parentDir := filepath.Dir(workDir)
			configPath = filepath.Join(parentDir, "configs")
		}

		cfg, err := New(WithConfigPath(configPath))
		if err != nil {
			panic(fmt.Errorf("初始化配置失败: %s", err))
		}
		globalConfig = cfg
	})

	return globalConfig
}

// Get 获取配置值
func (c *Config) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.viper.Get(key)
}

// GetString 获取字符串配置
func (c *Config) GetString(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.viper.GetString(key)
}

// GetInt 获取整数配置
func (c *Config) GetInt(key string) int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.viper.GetInt(key)
}

// GetBool 获取布尔配置
func (c *Config) GetBool(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.viper.GetBool(key)
}

// GetFloat64 获取浮点数配置
func (c *Config) GetFloat64(key string) float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.viper.GetFloat64(key)
}

// GetDuration 获取时间间隔配置
func (c *Config) GetDuration(key string) time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.viper.GetDuration(key)
}

// GetStringSlice 获取字符串切片配置
func (c *Config) GetStringSlice(key string) []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.viper.GetStringSlice(key)
}

// GetStringMap 获取字符串映射配置
func (c *Config) GetStringMap(key string) map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.viper.GetStringMap(key)
}

// SetDefault 设置默认值
func (c *Config) SetDefault(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.viper.SetDefault(key, value)
}

// Set 设置配置值
func (c *Config) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.viper.Set(key, value)
}

// WriteConfig 写入配置到文件
func (c *Config) WriteConfig() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.viper.WriteConfig()
}

// IsSet 检查配置是否已设置
func (c *Config) IsSet(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.viper.IsSet(key)
}

// AllSettings 获取所有配置
func (c *Config) AllSettings() map[string]interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.viper.AllSettings()
}

// LoadEnvironment 获取当前环境
func (c *Config) LoadEnvironment() string {
	env := os.Getenv("KY_ADMIN_ENV")
	if env == "" {
		env = "development"
	}
	return env
}

// IsProduction 检查是否是生产环境
func (c *Config) IsProduction() bool {
	return c.LoadEnvironment() == "production"
}

// IsDevelopment 检查是否是开发环境
func (c *Config) IsDevelopment() bool {
	return c.LoadEnvironment() == "development"
}

// IsTest 检查是否是测试环境
func (c *Config) IsTest() bool {
	return c.LoadEnvironment() == "test"
}
