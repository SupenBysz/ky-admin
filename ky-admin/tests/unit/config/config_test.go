package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/SupenBysz/ky-admin/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigLoading(t *testing.T) {
	// 保存原始环境变量
	originalEnv := os.Getenv("KY_ADMIN_ENV")
	defer os.Setenv("KY_ADMIN_ENV", originalEnv)

	// 设置测试环境
	os.Setenv("KY_ADMIN_ENV", "test")

	// 获取项目根目录路径
	rootDir, err := findRootDir()
	require.NoError(t, err, "获取项目根目录失败")

	// 创建配置实例
	cfg, err := config.New(config.WithConfigPath(filepath.Join(rootDir, "configs")))
	require.NoError(t, err, "加载配置失败")
	assert.NotNil(t, cfg, "配置应该不为空")

	// 测试基本配置获取方法
	assert.Equal(t, "ky-admin", cfg.GetString("app.name"), "应用名称应为ky-admin")
	assert.Equal(t, "1.0.0", cfg.GetString("app.version"), "应用版本应为1.0.0")
	assert.Equal(t, 8081, cfg.GetInt("app.port"), "测试环境端口应为8081")
	assert.Equal(t, "test", cfg.GetString("app.mode"), "应用模式应为test")
	assert.True(t, cfg.GetBool("app.debug"), "调试模式应为开启状态")

	// 测试数组配置
	origins := cfg.GetStringSlice("app.allow_origins")
	assert.Contains(t, origins, "http://localhost:8080", "允许源应包含localhost:8080")
	assert.Contains(t, origins, "http://127.0.0.1:8080", "允许源应包含127.0.0.1:8080")

	// 测试嵌套配置
	assert.True(t, cfg.GetBool("app.api.health_check.enabled"), "健康检查应启用")
	assert.Equal(t, "/livez", cfg.GetString("app.api.health_check.livez_path"), "存活检查路径应为/livez")

	// 测试数据库配置
	assert.Equal(t, "localhost", cfg.GetString("database.host"), "数据库主机应为localhost")
	assert.Equal(t, 3306, cfg.GetInt("database.port"), "数据库端口应为3306")
	assert.Equal(t, "ky_admin_test", cfg.GetString("database.name"), "测试数据库名应为ky_admin_test")

	// 测试持续时间配置
	assert.Equal(t, 30*time.Second, cfg.GetDuration("app.timeout"), "应用超时应为30秒")
	assert.Equal(t, 1*time.Hour, cfg.GetDuration("database.conn_max_lifetime"), "连接最大生命周期应为1小时")
}

func TestGetWithDefault(t *testing.T) {
	// 获取项目根目录路径
	rootDir, err := findRootDir()
	require.NoError(t, err, "获取项目根目录失败")

	// 创建配置实例
	cfg, err := config.New(config.WithConfigPath(filepath.Join(rootDir, "configs")))
	require.NoError(t, err, "加载配置失败")
	assert.NotNil(t, cfg, "配置应该不为空")

	// 测试带默认值的获取方法
	assert.Equal(t, "default_value", cfg.GetStringWithDefault("non_existent_key", "default_value"), "非存在键应返回默认值")
	assert.Equal(t, 42, cfg.GetIntWithDefault("non_existent_key", 42), "非存在键应返回默认值")
	assert.True(t, cfg.GetBoolWithDefault("non_existent_key", true), "非存在键应返回默认值")
	assert.Equal(t, 3.14, cfg.GetFloat64WithDefault("non_existent_key", 3.14), "非存在键应返回默认值")
	assert.Equal(t, 5*time.Minute, cfg.GetDurationWithDefault("non_existent_key", 5*time.Minute), "非存在键应返回默认值")
}

func TestEnvironmentDetection(t *testing.T) {
	// 保存原始环境变量
	originalEnv := os.Getenv("KY_ADMIN_ENV")
	defer os.Setenv("KY_ADMIN_ENV", originalEnv)

	// 获取项目根目录路径
	rootDir, err := findRootDir()
	require.NoError(t, err, "获取项目根目录失败")

	// 测试开发环境
	os.Setenv("KY_ADMIN_ENV", "development")
	cfg, err := config.New(config.WithConfigPath(filepath.Join(rootDir, "configs")))
	require.NoError(t, err, "加载配置失败")
	assert.True(t, cfg.IsDevelopment(), "应该识别为开发环境")
	assert.False(t, cfg.IsProduction(), "不应该识别为生产环境")

	// 测试生产环境
	os.Setenv("KY_ADMIN_ENV", "production")
	cfg, err = config.New(config.WithConfigPath(filepath.Join(rootDir, "configs")))
	require.NoError(t, err, "加载配置失败")
	assert.True(t, cfg.IsProduction(), "应该识别为生产环境")
	assert.False(t, cfg.IsDevelopment(), "不应该识别为开发环境")
}

// 查找项目根目录
func findRootDir() (string, error) {
	// 获取当前目录
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 向上查找直到找到configs目录
	for {
		if _, err := os.Stat(filepath.Join(dir, "configs")); err == nil {
			return dir, nil
		}

		// 到达根目录，无法继续向上查找
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// 如果当前不在项目目录中，尝试使用预设的路径
	baseDir := "/Volumes/DataDocument/CodeSpace/AI Solution/kysion.com/backend/ky-admin"
	if _, err := os.Stat(filepath.Join(baseDir, "configs")); err == nil {
		return baseDir, nil
	}

	return "", os.ErrNotExist
}
