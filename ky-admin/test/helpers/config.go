package helpers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/SupenBysz/ky-admin/pkg/config"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	Charset         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
	LogLevel        int
}

// LoadTestConfig 加载测试配置
func LoadTestConfig() (*config.TestConfig, error) {
	return config.GetTestConfig()
}

// 查找项目根目录
func FindRootDir() (string, error) {
	// 获取当前目录
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 向上查找直到找到go.mod文件
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
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
	if _, err := os.Stat(filepath.Join(baseDir, "go.mod")); err == nil {
		return baseDir, nil
	}

	return "", fmt.Errorf("无法找到项目根目录")
}
